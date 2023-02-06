package utils

import (
	"regexp"
	"strings"
)

// URLToFilename removes the 'http://' or 'https://' and 'www.' prefixes from the URL.
// Also remove the last "/" character if it exists
func URLToFilename(url string) string {
	// Compile the regular expression to remove the prefixes.
	re := regexp.MustCompile(`https?://(www\.)?`)

	filename := re.ReplaceAllString(url, "")

	// Remove the last "/" character if it exists
	if len(filename) > 0 && string(filename[len(filename)-1]) == "/" {
		filename = filename[:len(filename)-1]
	}

	return filename
}

// AssetURLToFilename converts the URL to a filename by replacing slashes with underscores.
func AssetURLToFilename(url string) string {
	// Remove the 'http://' or 'https://' and 'www.' prefixes.
	filename := strings.ReplaceAll(url, "https://", "")
	filename = strings.ReplaceAll(filename, "http://", "")
	filename = strings.ReplaceAll(filename, "www.", "")
	// Replace slashes with underscores.
	filename = strings.ReplaceAll(filename, "/", "_")

	return filename
}

// WrapAssetDir will wrap URL with dir or dot (.)
// Sometimes HTML page link can't work with ./dir format and should have /dir format
func WrapAssetDir(dot, dir, asset string) string {
	return dot + "/" + dir + "/" + asset
}

// WrapURL adds the path to the end of the URL, ensuring that there is a single slash between them.
func WrapURL(url, path string) string {
	if strings.Contains(path, "http") {
		return path
	}

	if len(path) > 0 && string(path[0]) != "/" {
		url += "/"
	}

	return url + path
}
