package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type jsonCacheClient struct {
	baseURL    string
	httpClient *http.Client
}

type jsonCacheClientResponse[T any] struct {
	http *http.Response
	json T
}

type RestoreCacheRequest struct {
	AppID     string
	CacheKeys []string
}

type RestoreCacheResponse struct {
	URL string
}

type RestoreCacheJSONResponse jsonCacheClientResponse[RestoreCacheResponse]

func (client *jsonCacheClient) restoreCache(
	request RestoreCacheRequest,
) (response RestoreCacheJSONResponse, err error) {
	encodedCacheKeys := strings.Join(request.CacheKeys, ",")
	path := fmt.Sprintf("/restore?app_id=%s&cache_keys=%s", request.AppID, encodedCacheKeys)
	httpResponse, err := client.do(
		http.MethodGet,
		path,
		nil,
	)
	if err != nil {
		return response, err
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			// TODO: Log a warning
		}
	}(httpResponse.Body)

	var jsonResponse RestoreCacheResponse
	if err = json.NewDecoder(httpResponse.Body).Decode(&jsonResponse); err != nil {
		return response, err
	}

	response.http = httpResponse
	response.json = jsonResponse

	return response, nil
}

func (client *jsonCacheClient) do(
	method string,
	path string,
	request any,
) (response *http.Response, err error) {
	var body io.Reader

	if request != nil {
		jsonData, err := json.Marshal(request)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(jsonData)
	}

	url := fmt.Sprintf("%s%s", client.baseURL, path)
	httpRequest, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	httpRequest.Header.Set("content-type", "application/json")

	response, err = client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}

	return response, nil
}
