package network

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

type apiClient struct {
	httpClient  *retryablehttp.Client
	baseURL     string
	accessToken string
}

type restoreResponse struct {
	URL string `json:"url"`
}

func newApiClient(client *retryablehttp.Client, baseURL string, accessToken string) apiClient {
	return apiClient{
		httpClient:  client,
		baseURL:     baseURL,
		accessToken: accessToken,
	}
}

func (c apiClient) restore(cacheKeys []string) (string, error) {
	url := fmt.Sprintf("%s/restore?cache_keys=%s", c.baseURL, strings.Join(cacheKeys, ","))

	req, err := retryablehttp.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode == http.StatusNotFound {
		return "", ErrCacheNotFound
	}
	if resp.StatusCode != http.StatusOK {
		errorResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp)
	}

	var response restoreResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return "", err
	}

	return response.URL, nil
}

func (c apiClient) downloadArchive(url string) (io.ReadCloser, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusOK {
		return resp.Body, nil
	} else {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				panic(err)
			}
		}(resp.Body)
		errorResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, errorResp)
	}
}
