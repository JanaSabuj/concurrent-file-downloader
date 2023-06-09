package util

import (
	"fmt"
	"net/url"
	"path"
)

func ExtractFileName(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	fileName := path.Base(parsedURL.Path)
	if fileName == "/" || fileName == "." {
		return "", fmt.Errorf("unable to extract file name from URL: %s", urlStr)
	}

	return fileName, nil
}
