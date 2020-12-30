package mocks

import "net/http"

type MockHttpClient struct {
	getFunc func() (*http.Response, error)
}

func NewHttpMockClient(getFunc func() (*http.Response, error)) MockHttpClient {
	return MockHttpClient{
		getFunc: getFunc,
	}
}

func (client MockHttpClient) Get(url string) (*http.Response, error) {
	return client.getFunc()
}
