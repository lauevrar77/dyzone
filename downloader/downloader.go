package downloader

import "github.com/lauevrar77/dyzone/domain"

type Downloader interface {
	Download(url string) (*domain.WebResource, error)
}
