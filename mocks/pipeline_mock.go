package mocks

import (
	"github.com/lauevrar77/dyzone/domain"
)

type PipelineMock struct {
	managementFunc func(resource *domain.WebResource) (*domain.WebResource, error)
}

func NewPipelineMock(manageWebResourceFunc func(resource *domain.WebResource) (*domain.WebResource, error)) PipelineMock {
	return PipelineMock{
		managementFunc: manageWebResourceFunc,
	}
}

func (pipeline PipelineMock) ManageWebResource(resource *domain.WebResource) (*domain.WebResource, error) {
	return pipeline.managementFunc(resource)
}
