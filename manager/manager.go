package manager

import (
	"fmt"
	"log"
	"net/url"
	"strconv"

	"github.com/JanaSabuj/concurrent-file-downloader/cli"
	"github.com/JanaSabuj/concurrent-file-downloader/greenhttp"
	"github.com/JanaSabuj/concurrent-file-downloader/models"
	"github.com/JanaSabuj/concurrent-file-downloader/util"
)

func Init() {
	fmt.Println(util.InitMotd)
}

func End() {
	// fmt.Println(util.EndMotd)
}

func Run(urlPtr *url.URL) {
	// init http client
	client := greenhttp.NewHTTPClient()

	// make HEAD call
	url := urlPtr.String()
	method := "HEAD"
	headers := map[string]string{
		"User-Agent": "CFD Downloader",
	}
	resp, err := client.Do(method, url, headers)
	if err != nil {
		log.Fatal(err)
	}

	// get Content-Length
	contentLength := resp.Header.Get("Content-Length")
	contentLengthInBytes, err := strconv.Atoi(contentLength)
	if err != nil {
		log.Fatal("Unsupported file download type.... Empty size sent by server", err)
	}
	log.Println("Content-Length:", contentLengthInBytes)

	// get file name
	fname, err := cli.ExtractFileName(url)
	if err != nil {
		log.Fatal("Error extracting filename...")
	}
	log.Println("Filename extracted: ", fname)

	// set concurrent workers
	chunks := util.WORKER_ROUTINES
	log.Println(fmt.Sprintf("Set %v parallel workers/connections", chunks))

	// calculate chunk size
	chunksize := contentLengthInBytes / chunks
	log.Println("Each chunk size: ", chunksize)

	// create the download object
	downReq := models.DownloadRequest{
		Url:      url,
		FileName: fname,
		Chunks:   chunks,
	}

	// chunk it up
	// downReq.SplitIntoChunks()

	// get each chunk concurrently

	// merge
}
