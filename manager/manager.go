package manager

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"sync"

	"github.com/JanaSabuj/concurrent-file-downloader/cli"
	"github.com/JanaSabuj/concurrent-file-downloader/greenhttp"
	"github.com/JanaSabuj/concurrent-file-downloader/models"
	"github.com/JanaSabuj/concurrent-file-downloader/util"
)

func Init() {
	fmt.Println(util.InitMotd)
}

func End() {
	fmt.Println(util.EndMotd)
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
		log.Fatal("Unsupported file download type.... Empty size sent by server ", err)
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

	// create the downloadRequest object
	downReq := &models.DownloadRequest{
		Url:        url,
		FileName:   fname,
		Chunks:     chunks,
		Chunksize:  chunksize,
		TotalSize:  contentLengthInBytes,
		HttpClient: client,
	}

	// chunk it up
	byteRangeArray := make([][2]int, chunks)
	byteRangeArray = downReq.SplitIntoChunks()
	fmt.Println(byteRangeArray)

	// download each chunk concurrently
	var wg sync.WaitGroup
	for idx, byteChunk := range byteRangeArray {
		wg.Add(1) // add wait before goroutine invocation

		go func(idx int, byteChunk [2]int) {
			defer wg.Done() // defer done at start of goroutine
			err := downReq.Download(idx, byteChunk)
			if err != nil {
				log.Fatal(fmt.Sprintf("Failed to download chunk %v", idx), err)
			}
		}(idx, byteChunk)
	}
	wg.Wait()

	// merge
	err = downReq.MergeDownloads()
	if err != nil {
		log.Fatal("Failed merging tmp downloaded files...", err)
	}

	// cleanup
	err = downReq.CleanupTmpFiles()
	if err != nil {
		log.Fatal("Failed cleaning up tmp downloaded files...", err)
	}

	// final file generated
	log.Println(fmt.Sprintf("File generated: %v\n\n", downReq.FileName))
}
