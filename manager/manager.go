package manager

import (
	"fmt"
	"log"
	"net/url"

	"github.com/JanaSabuj/concurrent-file-downloader/greenhttp"
	"github.com/JanaSabuj/concurrent-file-downloader/util"
)

func Init() {
	fmt.Println(util.InitMotd)
}

func End() {
	// fmt.Println(util.EndMotd)
}

func Run(urlPtr *url.URL) {
	// init client
	client := greenhttp.NewHTTPClient()

	// make HEAD call and get size
	url := urlPtr.String()
	method := "HEAD"
	headers := map[string]string{
		"User-Agent": "CFD Downloader",
	}
	req, err := client.NewRequest(method, url, headers, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.DoRequest(req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.Header)
	// get len
	contentLength := resp.Header.Get("Content-Length")
	fmt.Println("Content-Length:", contentLength)

	// chunk it up

	// get each chunk concurrently

	// merge
}
