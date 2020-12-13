package downloader

import "diezone/domain"

type Downloader interface {
	Download(url string) (*domain.WebResource, error)
}
