package dyzone

import (
	"errors"
	"testing"

	"github.com/lauevrar77/dyzone/domain"
	"github.com/lauevrar77/dyzone/mocks"
)

func workingDownloader(url string) (*domain.WebResource, error) {
	return domain.NewWebResource(url, "text/html", []byte("Hello, World!"))
}
func failingDownloader(url string) (*domain.WebResource, error) {
	return nil, errors.New("error")
}

func workingSpiderFunc(resource *domain.WebResource) ([]string, *domain.WebResource, error) {
	followingUrls := []string{}

	return followingUrls, resource, nil
}

func workingFollowingSpiderFunc(resource *domain.WebResource) ([]string, *domain.WebResource, error) {
	followingUrls := []string{}

	if resource.URI() != "/example.html" {
		followingUrls = append(followingUrls, "https://example.com/example.html")
	}

	return followingUrls, resource, nil
}

func workingIgnoringSpiderFunc(resource *domain.WebResource) ([]string, *domain.WebResource, error) {
	followingUrls := []string{}

	return followingUrls, nil, nil
}

func failingSpiderFunc(resource *domain.WebResource) ([]string, *domain.WebResource, error) {
	followingUrls := []string{}

	return followingUrls, nil, errors.New("error")
}

func workingPipelineFunc(resource *domain.WebResource) (*domain.WebResource, error) {
	return resource, nil
}

func workingPipelineModifyingFunc(resource *domain.WebResource) (*domain.WebResource, error) {
	resource.ChangeRawContent([]byte("Hello modified"))
	return resource, nil
}
func failingPipelineFunc(resource *domain.WebResource) (*domain.WebResource, error) {
	return nil, errors.New("error")
}

func TestSpiderRunner(t *testing.T) {
	downloader := mocks.NewDownloaderMock(workingDownloader)
	spider := mocks.NewSpiderMock(workingSpiderFunc)
	pipeline := mocks.NewPipelineMock(workingPipelineFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	resources, err := runner.Run("http://example.com")

	if err != nil {
		t.Log("Error in spider")
		t.FailNow()
	}

	if len(resources) != 1 {
		t.Log("Wrong number of result resources")
		t.Fail()
	}

	resource := resources[0]

	if resource.Domain() != "example.com" {
		t.Log("Wrong domain")
		t.Fail()
	}

	if resource.URI() != "/" {
		t.Logf("Wrong uri : %s", resource.URI())
		t.Fail()
	}

	if resource.ContentType() != "text/html" {
		t.Log("Wrong content type")
		t.Fail()
	}
}

func TestSpiderRunnerFollow(t *testing.T) {
	downloader := mocks.NewDownloaderMock(workingDownloader)
	spider := mocks.NewSpiderMock(workingFollowingSpiderFunc)
	pipeline := mocks.NewPipelineMock(workingPipelineFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	resources, err := runner.Run("http://example.com")

	if err != nil {
		t.Log("Error in spider")
		t.FailNow()
	}

	if len(resources) != 2 {
		t.Log("Wrong number of result resources")
		t.Fail()
	}

	resource := resources[1]

	if resource.Domain() != "example.com" {
		t.Log("Wrong domain")
		t.Fail()
	}

	if resource.URI() != "/example.html" {
		t.Logf("Wrong uri : %s", resource.URI())
		t.Fail()
	}

	if resource.ContentType() != "text/html" {
		t.Log("Wrong content type")
		t.Fail()
	}
}

func TestSpiderRunnerPipeline(t *testing.T) {
	downloader := mocks.NewDownloaderMock(workingDownloader)
	spider := mocks.NewSpiderMock(workingSpiderFunc)
	pipeline := mocks.NewPipelineMock(workingPipelineModifyingFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	resources, err := runner.Run("http://example.com")

	if err != nil {
		t.Log("Error in spider")
		t.FailNow()
	}

	if len(resources) != 1 {
		t.Log("Wrong number of result resources")
		t.Fail()
	}

	resource := resources[0]

	if resource.Domain() != "example.com" {
		t.Log("Wrong domain")
		t.Fail()
	}

	if resource.URI() != "/" {
		t.Logf("Wrong uri : %s", resource.URI())
		t.Fail()
	}

	if resource.ContentType() != "text/html" {
		t.Log("Wrong content type")
		t.Fail()
	}

	expectedContent := []byte("Hello modified")
	if !contentMatch(expectedContent, resource.RawContent()) {
		t.Log("Wrong modified content")
		t.Fail()
	}
}

func TestSpiderRunnerIgnoringScraper(t *testing.T) {
	downloader := mocks.NewDownloaderMock(workingDownloader)
	spider := mocks.NewSpiderMock(workingIgnoringSpiderFunc)
	pipeline := mocks.NewPipelineMock(workingPipelineModifyingFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	resources, err := runner.Run("http://example.com")

	if err != nil {
		t.Log("Error in spider")
		t.FailNow()
	}

	if len(resources) != 0 {
		t.Log("Wrong number of result resources")
		t.Fail()
	}
}

func TestSpiderRunnerFailingDownloader(t *testing.T) {
	downloader := mocks.NewDownloaderMock(failingDownloader)
	spider := mocks.NewSpiderMock(workingIgnoringSpiderFunc)
	pipeline := mocks.NewPipelineMock(workingPipelineModifyingFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	_, err := runner.Run("http://example.com")

	if err == nil {
		t.Log("Downloader should fail")
		t.FailNow()
	}
}

func TestSpiderRunnerFailingSpider(t *testing.T) {
	downloader := mocks.NewDownloaderMock(workingDownloader)
	spider := mocks.NewSpiderMock(failingSpiderFunc)
	pipeline := mocks.NewPipelineMock(workingPipelineModifyingFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	_, err := runner.Run("http://example.com")

	if err == nil {
		t.Log("Spider should fail")
		t.FailNow()
	}
}

func TestSpiderRunnerFailingPipeline(t *testing.T) {
	downloader := mocks.NewDownloaderMock(workingDownloader)
	spider := mocks.NewSpiderMock(workingSpiderFunc)
	pipeline := mocks.NewPipelineMock(failingPipelineFunc)

	runner := NewSpiderRunner(downloader, spider, pipeline)
	_, err := runner.Run("http://example.com")

	if err == nil {
		t.Log("Spider should fail")
		t.FailNow()
	}
}
func contentMatch(expected []byte, received []byte) bool {
	if len(expected) != len(received) {
		return false
	}

	for index, expect := range expected {
		if expect != received[index] {
			return false
		}
	}
	return true
}
