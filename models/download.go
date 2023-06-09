package models

type DownloadRequest struct {
	Url      string
	FileName string
	Chunks   int
}
