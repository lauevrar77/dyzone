package downloader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/lauevrar77/dyzone/domain"
)

type httpDownloader struct {
	httpClient http.Client
}

func NewHttpDownloader() httpDownloader {
	return httpDownloader{
		httpClient: http.Client{},
	}
}

func (downloader httpDownloader) Download(url string) (*domain.WebResource, error) {
	response, err := downloader.httpClient.Get(url)

	if err != nil {
		return nil, err
	}

	if downloader.requestFailed(response) {
		errorMessage := downloader.formatRequestError(url, response)
		return nil, errors.New(errorMessage)
	}

	webResource, err := downloader.responseToWebResource(url, response)

	defer response.Body.Close()
	return webResource, nil
}

func (downloader httpDownloader) requestFailed(response *http.Response) bool {
	return response.StatusCode >= 400
}

func (downloader httpDownloader) formatRequestError(url string, response *http.Response) string {
	return fmt.Sprintf(
		"Error performing HTTP Get to %s. HTTP code is %s.",
		url,
		response.Status,
	)
}

func (downloader httpDownloader) responseToWebResource(requestUrl string, response *http.Response) (*domain.WebResource, error) {
	responseUrl, err := response.Location()

	if err != nil {
		responseUrl, err = url.Parse(requestUrl)
	}

	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	webResource := domain.NewWebResource(
		responseUrl.Host,
		responseUrl.Path,
		response.Header.Get("Content-Type"),
		body,
	)

	return &webResource, nil
}
