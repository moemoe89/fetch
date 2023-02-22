package fetcher

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// Zip zips the given `filePaths` and `dirs` into a single archive file specified by `filename`.
func (c *client) Zip(filename string, filePaths, dirs []string) error {
	// Create the zip file.
	zipFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create zip: %w", err)
	}

	defer func() { _ = zipFile.Close() }()

	// Create a new zip writer.
	zipWriter := zip.NewWriter(zipFile)

	defer func() { _ = zipWriter.Close() }()

	// Channel to wait for all the goroutines to complete.
	var wg sync.WaitGroup

	c.errChan = make(chan error, len(filePaths)+len(dirs))

	// Zip the file with concurrency.
	for _, filePath := range filePaths {
		wg.Add(1)

		go func(filePath string) {
			defer wg.Done()

			if err := c.addFile(zipWriter, filePath); err != nil {
				c.errChan <- fmt.Errorf("failed to add file to zip: %w", err)
				return
			}
		}(filePath)
	}

	// Zip the dir with concurrency.
	for _, dir := range dirs {
		wg.Add(1)

		go func(dir string) {
			defer wg.Done()

			if err := c.addDir(zipWriter, dir); err != nil {
				c.errChan <- fmt.Errorf("failed to add dir to zip: %w", err)
				return
			}
		}(dir)
	}

	// Wait for all the goroutines to complete.
	wg.Wait()

	// Check the error channels for any errors.
	select {
	case err := <-c.errChan:
		return err
	default:
		close(c.errChan)
		return nil
	}
}

// addFile adds filepath to zip.
func (c *client) addFile(zipWriter *zip.Writer, filePath string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed open file to zip: %w", err)
	}

	fileWriter, err := zipWriter.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed create file to zip: %w", err)
	}

	if _, err := io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("failed copy file to zip: %w", err)
	}

	return nil
}

// addDir adds directory to zip.
func (c *client) addDir(zipWriter *zip.Writer, dir string) error {
	// Walk through all the files in the directory.
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { //nolint:staticcheck
		if info.IsDir() {
			return nil
		}

		// Create a new file header for the current file.
		header, err := zip.FileInfoHeader(info) //nolint:staticcheck
		if err != nil {
			return fmt.Errorf("failed get file header info: %w", err)
		}

		// Set the name of the file in the zip archive.
		header.Name = path

		// Create a new file in the zip archive.
		targetFile, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed create file header: %w", err)
		}

		// Open the current file.
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("failed open file from dir to zip: %w", err)
		}

		defer func() { _ = file.Close() }()

		// Copy the contents of the current file to the new file in the zip archive.
		_, err = io.Copy(targetFile, file)
		if err != nil {
			return fmt.Errorf("failed copy file from dir to zip: %w", err)
		}

		return nil
	})
}
