package fetcher

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func (c *client) Zip(filename string, filePaths, dirs []string) error {
	// Create the zip file
	zipFile, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer func() { _ = zipFile.Close() }()

	// Create a new zip writer
	zipWriter := zip.NewWriter(zipFile)

	defer func() { _ = zipWriter.Close() }()

	for _, filePath := range filePaths {
		if err := c.addFile(zipWriter, filePath); err != nil {
			return err
		}
	}

	// walk through all the files in the directory
	for _, dir := range dirs {
		if err := c.addDir(zipWriter, dir); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) addFile(zipWriter *zip.Writer, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	fileWriter, err := zipWriter.Create(filePath)
	if err != nil {
		return err
	}

	if _, err := io.Copy(fileWriter, file); err != nil {
		return err
	}

	return nil
}

func (c *client) addDir(zipWriter *zip.Writer, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error { //nolint:staticcheck
		if info.IsDir() {
			return nil
		}

		// create a new file header for the current file
		header, err := zip.FileInfoHeader(info) //nolint:staticcheck
		if err != nil {
			return err
		}

		// set the name of the file in the zip archive
		header.Name = path

		// create a new file in the zip archive
		targetFile, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		// open the current file
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		defer func() { _ = file.Close() }()

		// copy the contents of the current file to the new file in the zip archive
		_, err = io.Copy(targetFile, file)
		if err != nil {
			return err
		}

		return nil
	})
}
