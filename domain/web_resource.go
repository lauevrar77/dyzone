package domain

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

type WebResource struct {
	domain      string
	uri         string
	contentType string
	rawContent  []byte
	htmlContent *string
}

func NewWebResource(domain string, uri string, contentType string, rawContent []byte) WebResource {
	return WebResource{
		domain:      domain,
		uri:         uri,
		contentType: contentType,
		rawContent:  rawContent,
		htmlContent: nil,
	}
}

func (resource WebResource) IsWebPage() bool {
	return strings.Contains(resource.contentType, "text/html")
}

func (resource *WebResource) StyleSheetsUrls() ([]string, error) {
	resource.parseHtml()

	urls := make([]string, 0)
	doc, err := htmlquery.Parse(strings.NewReader(*resource.htmlContent))

	if err != nil {
		return urls, err
	}

	list := htmlquery.Find(doc, "//link[@rel='stylesheet']/@href")
	for _, n := range list {
		url := htmlquery.SelectAttr(n, "href")
		urls = append(urls, url)
	}

	return urls, nil
}

func (resource *WebResource) ImagesUrls() ([]string, error) {
	resource.parseHtml()

	urls := make([]string, 0)
	doc, err := htmlquery.Parse(strings.NewReader(*resource.htmlContent))

	if err != nil {
		return urls, err
	}

	list := htmlquery.Find(doc, "//img/@src")
	for _, n := range list {
		url := htmlquery.SelectAttr(n, "src")
		urls = append(urls, url)
	}

	return urls, nil
}

func (resource *WebResource) JavascriptUrls() ([]string, error) {
	resource.parseHtml()

	urls := make([]string, 0)
	doc, err := htmlquery.Parse(strings.NewReader(*resource.htmlContent))

	if err != nil {
		return urls, err
	}

	list := htmlquery.Find(doc, "//script/@src")
	for _, n := range list {
		url := htmlquery.SelectAttr(n, "src")
		urls = append(urls, url)
	}

	return urls, nil
}

func (resource *WebResource) InternalLinksUrls() ([]string, error) {
	resource.parseHtml()

	urls := make([]string, 0)
	doc, err := htmlquery.Parse(strings.NewReader(*resource.htmlContent))

	if err != nil {
		return urls, err
	}

	list := htmlquery.Find(doc, "//a/@href")
	for _, n := range list {
		urlString := htmlquery.SelectAttr(n, "href")

		url, _ := url.Parse(urlString)

		if url.Host == resource.domain || len(url.Host) == 0 {
			if len(url.Host) == 0 {
				scheme := "http"
				if len(url.Scheme) > 0 {
					scheme = url.Scheme
				}
				urlString = fmt.Sprintf("%s://%s/%s", scheme, resource.domain, urlString)
			}
			urls = append(urls, urlString)
		}
	}

	return urls, nil
}

func (resource *WebResource) parseHtml() {
	if !resource.IsWebPage() {
		panic("Only web pages can be parsed to html")
	}

	if resource.htmlContent == nil {
		stringContent := string(resource.rawContent)
		resource.htmlContent = &stringContent
	}
}
