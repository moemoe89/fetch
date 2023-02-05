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
	}
	
	// walk through all the files in the directory
	for _, dir := range dirs {
		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}

			// create a new file header for the current file
			header, err := zip.FileInfoHeader(info)
			if err != nil {
				return err
			}

			// set the name of the file in the zip archive
			header.Name = path

			// create a new file in the zip archive
			f, err := zipWriter.CreateHeader(header)
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
			_, err = io.Copy(f, file)
			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
