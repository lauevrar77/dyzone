package domain

import (
	"net/url"
	"testing"
)

func TestNewWebResource(t *testing.T) {
	content := make([]byte, 10)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	if webResource.Domain() != "example.com" {
		t.Logf("Wrong domain %s", webResource.Domain())
		t.Fail()
	}

	if webResource.URI() != "/home" {
		t.Logf("Wrong URI %s", webResource.URI())
		t.Fail()
	}

	if webResource.ContentType() != "text/html" {
		t.Logf("Wrong content type %s", webResource.ContentType())
		t.Fail()
	}

	contentEquals := true
	pageContent := webResource.RawContent()
	for index, contentbyte := range pageContent {
		contentEquals = content[index] == contentbyte

		if !contentEquals {
			t.Log("Contents are different")
			t.FailNow()
		}

	}
}

func TestIsWebPageHTML(t *testing.T) {
	content := make([]byte, 10)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	if !webResource.IsWebPage() {
		t.Log("Not detected as web page")
		t.Fail()
	}
}

func TestIsWebPageHTM(t *testing.T) {
	content := make([]byte, 10)
	webResource, err := NewWebResource("https://example.com/home", "text/htm", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	if !webResource.IsWebPage() {
		t.Log("Not detected as web page")
		t.Fail()
	}
}

func TestIsNotWebPage(t *testing.T) {
	content := make([]byte, 10)
	webResource, err := NewWebResource("https://example.com/home", "image/jpg", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	if webResource.IsWebPage() {
		t.Log("Not detected as web page")
		t.Fail()
	}
}

func TestInternalWebLinks(t *testing.T) {
	strContent := `
	<html>
		<body>
			<a href="hello1.html"></a>
			<a href="http://example.com/hello2.html"></a>
			<a href="google.com/hello3.html"></a>
		</body>

	</html>
	`
	content := []byte(strContent)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	links, err := webResource.InternalLinksUrls()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := []string{"https://example.com/hello1.html", "http://example.com/hello2.html"}
	if !linksMatch(expected, links) {
		t.Logf("Extracted urls %s does not match expected ones %s", links, expected)
		t.Fail()
	}
}

func TestWebLinks(t *testing.T) {
	strContent := `
	<html>
		<body>
			<a href="hello1.html"></a>
			<a href="http://example.com/hello2.html"></a>
			<a href="google.com/hello3.html"></a>
		</body>

	</html>
	`
	content := []byte(strContent)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	links, err := webResource.LinksUrls()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := []string{"https://example.com/hello1.html", "http://example.com/hello2.html", "http://google.com/hello3.html"}
	if !linksMatch(expected, links) {
		t.Logf("Extracted urls %s does not match expected ones %s", links, expected)
		t.Fail()
	}
}

func TestImages(t *testing.T) {
	strContent := `
	<html>
		<body>
			<a href="hello1.html"></a>
			<a href="http://example.com/hello2.html"></a>
			<a href="google.com/hello3.html"></a>
			<img src="https://example.com/assets/test.jpg"></a>
			<img src="image.png"></a>
			<img src="google.com/hello3.png"></a>
		</body>

	</html>
	`
	content := []byte(strContent)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	links, err := webResource.ImagesUrls()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := []string{"https://example.com/assets/test.jpg", "https://example.com/image.png", "http://google.com/hello3.png"}
	if !linksMatch(expected, links) {
		t.Logf("Extracted urls %s does not match expected ones %s", links, expected)
		t.Fail()
	}
}

func TestStyleSheets(t *testing.T) {
	strContent := `
	<html>
		<body>
			<a href="hello1.html"></a>
			<a href="http://example.com/hello2.html"></a>
			<a href="google.com/hello3.html"></a>
			<img src="https://example.com/assets/test.jpg"></a>
			<img src="image.png"></a>
			<img src="google.com/hello3.png"></a>
			<link rel="stylesheet" href="style.css"/>
			<link rel="stylesheet" href="https://google.com/style.css"/>
		</body>

	</html>
	`
	content := []byte(strContent)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	links, err := webResource.StyleSheetsUrls()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := []string{"https://example.com/style.css", "https://google.com/style.css"}
	if !linksMatch(expected, links) {
		t.Logf("Extracted urls %s does not match expected ones %s", links, expected)
		t.Fail()
	}
}

func TestJavascript(t *testing.T) {
	strContent := `
	<html>
		<body>
			<script src="https://google.com/script.js"></script>
			<script src="script.js"></script>
		</body>

	</html>
	`
	content := []byte(strContent)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)

	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	links, err := webResource.JavascriptUrls()

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	expected := []string{"https://example.com/script.js", "https://google.com/script.js"}
	if !linksMatch(expected, links) {
		t.Logf("Extracted urls %s does not match expected ones %s", links, expected)
		t.Fail()
	}
}

func TestChangeWebContent(t *testing.T) {
	strContent := `
	<html>
		<body>
			<script src="https://google.com/script.js"></script>
			<script src="script.js"></script>
		</body>

	</html>
	`
	content := []byte(strContent)
	webResource, err := NewWebResource("https://example.com/home", "text/html", content)
	if err != nil {
		t.Log("Could not create WebResource")
		t.FailNow()
	}

	newContent := []byte("Hello, World!")
	webResource.ChangeRawContent(newContent)

	if !contentMatch(newContent, webResource.RawContent()) {
		t.Logf("Wrong modified content : %s", webResource.RawContent())
		t.Fail()
	}
}

func TestMakeUrlCannonical(t *testing.T) {
	parsedUrl, _ := url.Parse("https://example.com")

	cannonicalUrl, _ := MakeUrlCannonical("no_domain_scheme.html", parsedUrl)
	expected := "https://example.com/no_domain_scheme.html"
	if cannonicalUrl != expected {
		t.Logf("%s different of %s", cannonicalUrl, expected)
		t.Fail()

	}

	cannonicalUrl, _ = MakeUrlCannonical("example.com/no_scheme.html", parsedUrl)
	expected = "https://example.com/no_scheme.html"
	if cannonicalUrl != expected {
		t.Logf("%s different of %s", cannonicalUrl, expected)
		t.Fail()

	}

	cannonicalUrl, _ = MakeUrlCannonical("https://google.com/other_complete.html", parsedUrl)
	expected = "https://google.com/other_complete.html"
	if cannonicalUrl != expected {
		t.Logf("%s different of %s", cannonicalUrl, expected)
		t.Fail()

	}

	cannonicalUrl, _ = MakeUrlCannonical("google.com/other_no_scheme.html", parsedUrl)
	expected = "http://google.com/other_no_scheme.html"
	if cannonicalUrl != expected {
		t.Logf("%s different of %s", cannonicalUrl, expected)
		t.Fail()
	}
}

func linksMatch(expected []string, received []string) bool {
	if len(expected) != len(received) {
		return false
	}

	expectedMap := make(map[string]int, len(expected))
	receivedMap := make(map[string]int, len(received))

	for index, expect := range expected {
		expectedMap[expect] += 1
		receivedMap[received[index]] += 1
	}

	for key, expect := range expectedMap {
		if expect != receivedMap[key] {
			return false
		}
	}

	return true
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
