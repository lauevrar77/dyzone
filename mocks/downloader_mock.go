package mocks

import (
	"github.com/lauevrar77/dyzone/domain"
)

type DownloaderMock struct {
	managementFunc func(url string) (*domain.WebResource, error)
}

func NewDownloaderMock(manageWebResourceFunc func(url string) (*domain.WebResource, error)) DownloaderMock {
	return DownloaderMock{
		managementFunc: manageWebResourceFunc,
	}
}

func (downloader DownloaderMock) Download(url string) (*domain.WebResource, error) {
	return downloader.managementFunc(url)
}
