package main

import (
	"log"

	"github.com/JanaSabuj/concurrent-file-downloader/cli"
	"github.com/JanaSabuj/concurrent-file-downloader/manager"
)

func main() {
	// init motd
	manager.Init()
	defer manager.End()

	// get URL as input from user to download
	url, err := cli.GetURLFromUser()
	if err != nil {
		log.Fatal("Invalid URL:", err)
	}

	// call Manager to download the file
	manager.Run(url)
}
