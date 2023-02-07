//nolint:lll
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/moemoe89/fetch/pkg/fetcher"
	"github.com/moemoe89/fetch/pkg/utils"
)

// usageText is a message to describe how to use the CLI.
const usageText = `fetch is a command line interface (CLI) program that can be used to fetch web pages and save their contents to disk for later retrieval and browsing. The program takes one or more URLs as input and downloads the HTML contents of each URL. The contents are saved to disk as HTML files with a name derived from the URL.

In addition to fetching and saving the HTML contents, the program also has an optional --metadata flag that can be used to print metadata about the fetched pages. The metadata includes the date and time of last fetch, the number of links on the page, and the number of images on the page.

This program provides a convenient way to retrieve and store web pages for offline viewing and is useful for developers and users who need to save and view web pages at a later time.

Example:
	fetch https://www.google.com
	fetch --metadata https://www.google.com
	fetch --metadata https://www.google.com https://www.github.com

`

var (
	// metadata is a flag to fetch the page with metadata (site name, number of link & image, last fetch) or not.
	metadata = flag.Bool("metadata", false, "Print metadata about the fetched pages such as site name, number of links, number of images and last fetch time")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 1 {
		usage()
		log.Fatal("Expected minimum one argument")
	}

	urls := flag.Args()

	// Initialize fetcher.
	client, err := fetcher.New()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	// Fetch the URLs with concurrency.
	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			err = fetchPage(client, url)
			if err != nil {
				// If something wrong happen, print the error.
				_, _ = io.WriteString(os.Stderr, err.Error()+"\n\n")
			}

			wg.Done()
		}(url)
	}

	wg.Wait()

	os.Exit(0)
}

func usage() {
	_, _ = io.WriteString(os.Stderr, usageText)

	flag.PrintDefaults()
}

func fetchPage(client fetcher.Fetcher, url string) error {
	body, err := client.FetchPage(context.Background(), url)
	if err != nil {
		return fmt.Errorf("failed to fetch page: %s: %w", url, err)
	}

	filename := utils.URLToFilename(url)

	htmlFile := filename + ".html"

	// Stop here if the argument doesn't includes metadata.
	if !*metadata {
		err = client.SavePage(htmlFile, body)
		if err != nil {
			return fmt.Errorf("failed to save page: %s: %w", url, err)
		}

		os.Exit(0)
	}

	dir := filename
	zipFile := filename + ".zip"
	jsonFile := filename + ".json"

	// Creates assets directory.
	err = os.Mkdir(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create dir: %s: %w", dir, err)
	}

	// Extract metadata.
	metadata, err := client.ExtractMetadata(url, dir+"/"+jsonFile, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to extract metadata: %s: %w", url, err)
	}

	newBody := string(body)

	newBody, err = fetchAssets(client, metadata, dir, newBody)
	if err != nil {
		return err
	}

	// Save HTML page.
	err = client.SavePage(htmlFile, []byte(newBody))
	if err != nil {
		return fmt.Errorf("failed to save page: %s: %w", url, err)
	}

	// Zip assets and HTML file.
	err = client.Zip(zipFile, []string{htmlFile}, []string{dir})
	if err != nil {
		return fmt.Errorf("failed to zip page: %s: %w", url, err)
	}

	// Print the metadata string.
	stringMetadata := client.StringMetadata(metadata)
	_, _ = io.WriteString(os.Stderr, stringMetadata)

	return nil
}

func fetchAssets(
	client fetcher.Fetcher,
	metadata *fetcher.Metadata,
	dir, newBody string,
) (string, error) {
	var wg sync.WaitGroup

	var mutex sync.Mutex

	// Error channel for file paths.
	errChan := make(chan error, len(metadata.Assets))

	// Fetch the assets with concurrency.
	for _, asset := range metadata.Assets {
		wg.Add(1)

		go func(asset string) {
			defer wg.Done()

			// Sometimes URL not contains full URL e.g. /dir/image.png
			// To download the assets, assume it falls under target URL.
			// e.g. www.example.com/dir/image.png
			wrapAsset := utils.WrapURL(metadata.Site, asset)

			body, err := client.FetchPage(context.Background(), wrapAsset)
			if err != nil {
				errChan <- fmt.Errorf("failed to fetch page: %s: %w", wrapAsset, err)

				return
			}

			wrapAssetDir, assetFile := buildAssetDirFile(dir, asset, wrapAsset)

			mutex.Lock()
			newBody = strings.ReplaceAll(newBody, asset, wrapAssetDir)
			mutex.Unlock()

			// Save assets file.
			err = client.SavePage(dir+"/"+assetFile, body)
			if err != nil {
				errChan <- fmt.Errorf("failed to save page: %s: %w", wrapAsset, err)

				return
			}
		}(asset)
	}

	wg.Wait()

	// Handle error channel from downloading assets.
	select {
	case err := <-errChan:
		return "", err
	default:
	}

	return newBody, nil
}

func buildAssetDirFile(dir, asset, wrapAsset string) (string, string) {
	// Removes unnecessary characters.
	assetFile := utils.AssetURLToFilename(wrapAsset)

	// Sometimes HTML page doesn't work well if the link contains dot (.)
	// e.g ./dir/image.png and just need /dir/image.png
	dot := "."
	if len(asset) > 0 && string(asset[0]) == "/" {
		dot = ""
	}

	wrapAssetDir := utils.WrapAssetDir(dot, dir, assetFile)

	return wrapAssetDir, assetFile
}
