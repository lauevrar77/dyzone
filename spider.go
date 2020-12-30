package dyzone

import (
	"github.com/lauevrar77/dyzone/domain"
	"github.com/lauevrar77/dyzone/downloader"
)

type Spider interface {
	OnWebResourceFetched(webResource *domain.WebResource) ([]string, *domain.WebResource, error)
}

type WebResourcePipeline interface {
	ManageWebResource(webResource *domain.WebResource) (*domain.WebResource, error)
}

type SpiderRunner struct {
	downloader downloader.Downloader
	spider     Spider
	pipeline   WebResourcePipeline
}

func (runner SpiderRunner) Run(startUrl string) ([]*domain.WebResource, error) {
	resources := make([]*domain.WebResource, 0)

	// Run Downloader
	webResource, err := runner.downloader.Download(startUrl)

	if err != nil {
		return nil, err
	}

	// Give result to spider to generate following requests and result
	newRequests, webResource, err := runner.spider.OnWebResourceFetched(webResource)

	if err != nil {
		return nil, err
	}

	// Send found resource into the management pipeline
	if webResource != nil {
		webResource, err = runner.pipeline.ManageWebResource(webResource)

		if err != nil {
			return nil, err
		}

		if webResource != nil {
			resources = append(resources, webResource)
		}
	}

	// Perform following requests recusively
	for _, request := range newRequests {
		childrenResources, err := runner.Run(request)

		if err != nil {
			return nil, err
		}

		for _, childResource := range childrenResources {
			resources = append(resources, childResource)
		}
	}

	return resources, nil
}

func NewSpiderRunner(downloader downloader.Downloader, spider Spider, pipeline WebResourcePipeline) SpiderRunner {
	return SpiderRunner{
		downloader: downloader,
		spider:     spider,
		pipeline:   pipeline,
	}
}
