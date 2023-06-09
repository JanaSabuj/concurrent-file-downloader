package manager

import (
	"fmt"
	"net/url"

	"github.com/JanaSabuj/concurrent-file-downloader/util"
)

func Init() {
	fmt.Println(util.InitMotd)
}

func End() {
	fmt.Println(util.EndMotd)
}

func Run(url *url.URL) {
	// make HEAD call and get size

	// chunk it up

	// get each chunk concurrently

	// merge
}
