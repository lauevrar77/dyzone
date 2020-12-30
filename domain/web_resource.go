package domain

import (
	"log"
	"net/url"
	"strings"

	"github.com/antchfx/htmlquery"
)

type WebResource struct {
	url         *url.URL
	contentType string
	rawContent  []byte
	htmlContent *string
}

func NewWebResource(webUrl string, contentType string, rawContent []byte) (*WebResource, error) {
	parsedUrl, err := url.Parse(webUrl)

	if err != nil {
		return nil, err
	}

	return &WebResource{
		url:         parsedUrl,
		contentType: contentType,
		rawContent:  rawContent,
		htmlContent: nil,
	}, nil
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
		cannonicalUrl, err := MakeUrlCannonical(url, resource.url)
		if err != nil {
			log.Printf("Could not parse url %s\n", url)
		} else {
			urls = append(urls, cannonicalUrl)
		}
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
		cannonicalUrl, err := MakeUrlCannonical(url, resource.url)
		if err != nil {
			log.Printf("Could not parse url %s\n", url)
		} else {
			urls = append(urls, cannonicalUrl)
		}
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
		cannonicalUrl, err := MakeUrlCannonical(url, resource.url)
		if err != nil {
			log.Printf("Could not parse url %s\n", url)
		} else {
			urls = append(urls, cannonicalUrl)
		}
	}

	return urls, nil
}

func (resource *WebResource) LinksUrls() ([]string, error) {
	resource.parseHtml()

	urls := make([]string, 0)
	doc, err := htmlquery.Parse(strings.NewReader(*resource.htmlContent))

	if err != nil {
		return urls, err
	}

	list := htmlquery.Find(doc, "//a/@href")
	for _, n := range list {
		urlString := htmlquery.SelectAttr(n, "href")
		cannonicalUrl, err := MakeUrlCannonical(urlString, resource.url)
		if err != nil {
			log.Printf("Could not parse url %s\n", urlString)
		} else {
			urls = append(urls, cannonicalUrl)
		}
	}

	return urls, nil
}

func (resource *WebResource) InternalLinksUrls() ([]string, error) {
	urls, err := resource.LinksUrls()

	if err != nil {
		return urls, err
	}

	filteredUrls := make([]string, 0)
	for _, url := range urls {
		if isInternal, err := isInternalUrl(url, resource.url); err == nil && isInternal {
			filteredUrls = append(filteredUrls, url)
		}
	}

	return filteredUrls, nil
}

func (resource WebResource) ContentType() string {
	return resource.contentType
}

func (resource WebResource) RawContent() []byte {
	return resource.rawContent
}

func (resource WebResource) Domain() string {
	return resource.url.Host
}

func (resource WebResource) URI() string {
	return resource.url.RequestURI()
}

func (resource *WebResource) ChangeRawContent(content []byte) {
	resource.rawContent = content
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

func MakeUrlCannonical(rawUrl string, parentUrl *url.URL) (string, error) {
	parsedUrl, err := url.Parse(rawUrl)

	if err != nil {
		return "", err
	}

	if parsedUrl.Scheme == "" {
		slashSplits := strings.Split(parsedUrl.Path, "/")
		containsMultipleSlashes := len(slashSplits) > 1

		if containsMultipleSlashes {
			// We are sure enough we have a domain.
			// So we just have to add scheme

			if strings.HasPrefix(parsedUrl.Path, parentUrl.Host) { // Same Domain as parent ?
				parsedUrl.Scheme = parentUrl.Scheme // Use parent scheme
			} else {
				parsedUrl.Scheme = "http" // Otherwise suppose http/https redirection exists
			}
		} else {
			// May be just the domain without trailing slash e.g : google.com
			// Or a relative link to a page e.g : index.html
			// Trying to determine based on "extension"

			dotSplits := strings.Split(slashSplits[0], ".")
			extension := dotSplits[len(dotSplits)-1]

			if isWebResourceExtension(extension) {
				parsedUrl.Scheme = parentUrl.Scheme
				parsedUrl.Host = parentUrl.Host
			} else {
				parsedUrl.Scheme = "http" // Suppose that HTTPS redirection exists
			}

		}

		return parsedUrl.String(), nil
	}

	return rawUrl, nil

}

func isWebResourceExtension(extension string) bool {
	extension = strings.ToLower(extension)

	switch extension {
	case "html":
		return true
	case "php":
		return true
	case "jpg":
		return true
	case "jpeg":
		return true
	case "png":
		return true
	case "css":
		return true
	case "js":
		return true
	}

	return false
}

func isInternalUrl(rawUrl string, parentUrl *url.URL) (bool, error) {
	parsedUrl, err := url.Parse(rawUrl)

	if err != nil {
		return false, err
	}

	return parsedUrl.Host == parentUrl.Host, nil
}
