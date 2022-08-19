package network

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
)

const maxKeyLength = 512
const maxKeyCount = 8

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
	keysInQuery, err := validateKeys(cacheKeys)
	if err != nil {
		return "", err
	}
	apiURL := fmt.Sprintf("%s/restore?cache_keys=%s", c.baseURL, keysInQuery)

	req, err := retryablehttp.NewRequest("GET", apiURL, nil)
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

func validateKeys(keys []string) (string, error) {
	if len(keys) > maxKeyCount {
		return "", fmt.Errorf("maximum number of keys is %d", maxKeyCount)
	}
	var truncatedKeys []string
	for _, key := range keys {
		if strings.Contains(key, ",") {
			return "", fmt.Errorf("commas are not allowed in keys")
		}
		if len(key) > maxKeyLength {
			truncatedKeys = append(truncatedKeys, key[:maxKeyLength])
		} else {
			truncatedKeys = append(truncatedKeys, key)
		}
	}

	return url.QueryEscape(strings.Join(truncatedKeys, ",")), nil
}
