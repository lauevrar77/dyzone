package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/lauevrar77/dyzone"
	"github.com/lauevrar77/dyzone/domain"
	"github.com/lauevrar77/dyzone/downloader"
)

var acceptedContentTypes = [...]string{"image/jpeg", "image/png"}

type ImageSpider struct {
}

func (spider ImageSpider) OnWebResourceFetched(webResource *domain.WebResource) ([]string, *domain.WebResource, error) {
	fmt.Println("Spider got result")
	if webResource.IsWebPage() {
		fmt.Println("Spider result is a web page")
		imageUrls, err := webResource.ImagesUrls()

		if err != nil {
			fmt.Println("Could not extract images")
			return imageUrls, nil, err
		}
		fmt.Printf("Web page contains %d images \n", len(imageUrls))

		return imageUrls, nil, nil
	}

	contentType := webResource.ContentType()
	if isAcceptedContentType(contentType) {
		return []string{}, webResource, nil
	}

	return []string{}, nil, nil
}

func isAcceptedContentType(contentType string) bool {
	for _, acceptedContentType := range acceptedContentTypes {
		if contentType == acceptedContentType {
			return true
		}
	}

	return false
}

type ImagePipeline struct{}

func (pipeline ImagePipeline) ManageWebResource(webResource *domain.WebResource) (*domain.WebResource, error) {
	fmt.Println("Saving image to 'images' directory")

	rawContent := webResource.RawContent()

	if exist, err := exists("images"); !exist || err != nil {
		os.Mkdir("images", 0755)
	}

	imageNameParts := strings.Split(webResource.URI(), "/")
	imageName := imageNameParts[len(imageNameParts)-1]
	file, err := os.Create(fmt.Sprintf("images/%s", imageName))
	if err != nil {
		fmt.Println("Cannot open file")
		fmt.Println(err)
		return nil, err
	}
	writer := bufio.NewWriter(file)
	writer.Write(rawContent)
	writer.Flush()

	return nil, nil
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func main() {
	var pipeline dyzone.WebResourcePipeline
	pipe := ImagePipeline{}
	pipeline = pipe

	downloader := downloader.NewHttpDownloader()

	var spider dyzone.Spider
	imageSpider := ImageSpider{}
	spider = imageSpider

	runner := dyzone.NewSpiderRunner(downloader, spider, pipeline)

	runner.Run("https://korben.info")
}
