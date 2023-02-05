package utils

import (
	"regexp"
	"strings"
)

// URLToFilename removes the 'http://' or 'https://' and 'www.' prefixes from the URL.
func URLToFilename(url string) string {
	// Compile the regular expression to remove the prefixes.
	re := regexp.MustCompile(`https?://(www\.)?`)

	return re.ReplaceAllString(url, "")
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
