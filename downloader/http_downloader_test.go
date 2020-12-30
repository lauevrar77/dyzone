package downloader

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/lauevrar77/dyzone/mocks"
)

func makeGetFunction(contentType string, location string, body string, statusCode int) func() (*http.Response, error) {
	return func() (*http.Response, error) {

		headers := make(http.Header, 0)
		headers.Set("Content-Type", contentType)
		headers.Set("Location", location)
		return &http.Response{
			Status:        fmt.Sprintf("%d", statusCode),
			StatusCode:    statusCode,
			Proto:         "HTTP/1.1",
			ProtoMajor:    1,
			ProtoMinor:    1,
			Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
			ContentLength: int64(len(body)),
			Header:        headers,
		}, nil
	}
}

func makeGetExceptionFunction() func() (*http.Response, error) {
	return func() (*http.Response, error) {

		return nil, errors.New("Wow error!")
	}
}
func TestChangeClient(t *testing.T) {
	mockClient := mocks.NewHttpMockClient(func() (*http.Response, error) {
		return &http.Response{}, nil
	})
	downloader := NewHttpDownloader()
	downloader.ChangeHttpClient(mockClient)
}

func TestDownload(t *testing.T) {
	mockClient := mocks.NewHttpMockClient(
		makeGetFunction("text/html", "https://example.com", "Hello, World!", 200),
	)
	downloader := NewHttpDownloader()
	downloader.ChangeHttpClient(mockClient)

	webResponse, err := downloader.Download("https://example.com")

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	if webResponse == nil {
		t.Log("Return WebResource should not be null")
		t.FailNow()
	}

	if webResponse.ContentType() != "text/html" {
		t.Log("Wrong content type")
		t.Fail()
	}

	if webResponse.Domain() != "example.com" {
		t.Log("Wrong domain")
		t.Fail()
	}

	if !sameContent(webResponse.RawContent(), []byte("Hello, World!")) {
		t.Log("Wrong content")
		t.Fail()
	}
}

func TestDownloadHttpError(t *testing.T) {
	mockClient := mocks.NewHttpMockClient(
		makeGetFunction("text/html", "https://example.com", "Hello, World!", 500),
	)
	downloader := NewHttpDownloader()
	downloader.ChangeHttpClient(mockClient)

	_, err := downloader.Download("https://example.com")

	if err == nil {
		t.Log(err)
		t.FailNow()
	}
}

func TestDownloadExceptionError(t *testing.T) {
	mockClient := mocks.NewHttpMockClient(
		makeGetExceptionFunction(),
	)
	downloader := NewHttpDownloader()
	downloader.ChangeHttpClient(mockClient)

	_, err := downloader.Download("https://example.com")

	if err == nil {
		t.Log(err)
		t.FailNow()
	}
}
func sameContent(content1 []byte, content2 []byte) bool {
	if len(content1) != len(content2) {
		return false
	}

	for index, content := range content1 {
		if content != content2[index] {
			return false
		}
	}

	return true
}
