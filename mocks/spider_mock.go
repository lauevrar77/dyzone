package mocks

import (
	"github.com/lauevrar77/dyzone/domain"
)

type SpiderMock struct {
	managementFunc func(resource *domain.WebResource) ([]string, *domain.WebResource, error)
}

func NewSpiderMock(onResourceFetchFunc func(resource *domain.WebResource) ([]string, *domain.WebResource, error)) SpiderMock {
	return SpiderMock{
		managementFunc: onResourceFetchFunc,
	}
}

func (spider SpiderMock) OnWebResourceFetched(resource *domain.WebResource) ([]string, *domain.WebResource, error) {
	return spider.managementFunc(resource)
}
