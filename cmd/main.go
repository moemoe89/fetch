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

	"github.com/moemoe89/fetch/pkg/fetcher"
	"github.com/moemoe89/fetch/pkg/utils"
)

const usageText = `fetch is a command line interface (CLI) program that can be used to fetch web pages and save their contents to disk for later retrieval and browsing. The program takes one or more URLs as input and downloads the HTML contents of each URL. The contents are saved to disk as HTML files with a name derived from the URL.

In addition to fetching and saving the HTML contents, the program also has an optional --metadata flag that can be used to print metadata about the fetched pages. The metadata includes the date and time of last fetch, the number of links on the page, and the number of images on the page.

This program provides a convenient way to retrieve and store web pages for offline viewing and is useful for developers and users who need to save and view web pages at a later time.

Example:
	fetch https://www.google.com
	fetch --metadata https://www.google.com
	fetch --metadata https://www.google.com https://www.github.com

`

var (
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

	client, err := fetcher.New()
	if err != nil {
		log.Fatal(err)
	}

	for _, url := range urls {
		body, err := client.FetchPage(context.Background(), url)
		if err != nil {
			log.Fatal(err)
		}

		filename := utils.URLToFilename(url)

		dir := filename
		htmlFile := filename + ".html"
		zipFile := filename + ".zip"

		err = os.Mkdir(dir, 0755)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		assets, err := client.ExtractMetadata(bytes.NewReader(body))

		newBody := string(body)

		for _, asset := range assets {
			assetNew := asset

			dot := "."
			if len(asset) > 0 && string(asset[0]) == "/" {
				dot = ""
			}

			assetNew = utils.WrapURL(url, asset)

			body, err := client.FetchPage(context.Background(), assetNew)
			if err != nil {
				log.Fatal(err)
			}

			assetFile := utils.AssetURLToFilename(assetNew)

			assetFile = dir + "/" + assetFile

			newBody = strings.ReplaceAll(newBody, asset, dot+"/"+assetFile)

			err = client.SavePage(assetFile, body)
			if err != nil {
				log.Fatal(err)
			}
		}

		err = client.SavePage(htmlFile, []byte(newBody))
		if err != nil {
			log.Fatal(err)
		}

		err = client.Zip(zipFile, []string{htmlFile}, []string{dir})
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(*metadata)
}

func usage() {
	_, _ = io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}
