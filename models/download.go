package models

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/JanaSabuj/concurrent-file-downloader/greenhttp"
	"github.com/JanaSabuj/concurrent-file-downloader/util"
)

type DownloadRequest struct {
	Url        string // eg: https://cdn.videvo.net/videvo_files/video/premium/video0042/large_watermarked/900-2_900-6334-PD2_preview.mp4
	FileName   string
	Chunks     int
	Chunksize  int
	TotalSize  int
	HttpClient *greenhttp.HTTPClient
}

func (d *DownloadRequest) SplitIntoChunks() [][2]int {
	arr := make([][2]int, d.Chunks)
	for i := 0; i < d.Chunks; i++ {
		if i == 0 {
			arr[i][0] = 0
			arr[i][1] = d.Chunksize
		} else if i == d.Chunks-1 {
			arr[i][0] = arr[i-1][1] + 1
			arr[i][1] = d.TotalSize - 1
		} else {
			arr[i][0] = arr[i-1][1] + 1
			arr[i][1] = arr[i][0] + d.Chunksize
		}
	}

	return arr
}

func (d *DownloadRequest) Download(idx int, byteChunk [2]int) error {
	log.Println(fmt.Sprintf("Downloading chunk %v", idx))
	// make GET request with range
	method := "GET"
	headers := map[string]string{
		"User-Agent": "CFD Downloader",
		"Range":      fmt.Sprintf("bytes=%v-%v", byteChunk[0], byteChunk[1]),
	}
	resp, err := d.HttpClient.Do(method, d.Url, headers)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Chunk fail %v", resp.StatusCode))
	}

	if resp.StatusCode > 299 {
		return fmt.Errorf(fmt.Sprintf("Can't process, response is %v", resp.StatusCode))
	}

	// open a file to write the body to
	fname := (fmt.Sprintf("%v-%v.tmp", util.TMP_FILE_PREFIX, idx))
	file, err := os.Create(fname)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Can't create a file %v", fname))
	}
	defer file.Close()

	// write to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("Failed to write to file: %v", err)
	}
	log.Println(fmt.Sprintf("Wrote chunk %v to file", idx))

	return nil
}

func (d *DownloadRequest) MergeDownloads() error {
	// create output file
	out, err := os.Create(d.FileName)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer out.Close()

	// append each chunk to final file
	for idx := 0; idx < d.Chunks; idx++ {
		fname := fmt.Sprintf("%v-%v.tmp", util.TMP_FILE_PREFIX, idx)
		in, err := os.Open(fname)
		if err != nil {
			return fmt.Errorf("Failed to open chunk file %s: %v", fname, err)
		}
		defer in.Close()

		_, err = io.Copy(out, in)
		if err != nil {
			return fmt.Errorf("Failed to merge chunk file %s: %v", fname, err)
		}
	}

	fmt.Println("File chunks merged successfully...")
	return nil
}

func (d *DownloadRequest) CleanupTmpFiles() error {
	log.Println("Starting to clean tmp downloaded files...")

	// delete each chunk file
	for idx := 0; idx < d.Chunks; idx++ {
		fname := fmt.Sprintf("%v-%v.tmp", util.TMP_FILE_PREFIX, idx)
		err := os.Remove(fname)
		if err != nil {
			return fmt.Errorf("Failed to remove chunk file %s: %v", fname, err)
		}
	}

	return nil
}
