package main

import (
	"diezone/downloader"
	"fmt"
)

func main() {
	// This is actually a "toy" main application as this will eventually become a library
	var downloadClient downloader.Downloader
	httpDownloadClient := downloader.NewHttpDownloader()
	downloadClient = &httpDownloadClient

	page, _ := downloadClient.Download("https://wikipedia.com")
	fmt.Println(page)

	cssUrls, _ := page.InternalLinksUrls()
	for _, url := range cssUrls {
		fmt.Println(url)
	}
}
