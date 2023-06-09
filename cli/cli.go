package cli

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"path"
)

func GetURLFromUser() (*url.URL, error) {
	// Create a new scanner to read from standard input
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the file URL to download: ")
	scanner.Scan()

	// Retrieve the user input
	userInput := scanner.Text()

	// Validate if the input is a valid URL (lenient parsing)
	parsedURL, err := url.Parse(userInput)
	if err != nil {
		return nil, err
	}

	return parsedURL, nil
}

// func GetFileNameFromUser() (string, error) {

// }

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
