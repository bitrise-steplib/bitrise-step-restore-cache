package service

import (
	"io"
	"net/http"
	"strings"
)

type gcsClient struct {
	baseURL    string
	httpClient *http.Client
}

func (client *gcsClient) download(url string) (archive []byte, err error) {
	request, err := http.NewRequest(http.MethodGet, client.replaceHost(url), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Set("content-type", "application/octet-stream")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			// TODO: Log a warning
		}
	}(response.Body)

	archive, err = io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return archive, nil
}

func (client *gcsClient) replaceHost(url string) string {
	return strings.Replace(url, "https://storage.googleapis.com", client.baseURL, 1)
}
