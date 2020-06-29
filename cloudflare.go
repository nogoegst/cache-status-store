package cachestatusstore

import (
	"net/http"
	"strings"
)

type CloudFlareCache struct {
	httpClient *http.Client
	url        string
}

func NewCloudFlareCache(u string) *CloudFlareCache {
	ret := &CloudFlareCache{
		httpClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse
			},
		},
		url: strings.TrimRight(u, "/"),
	}
	return ret
}

func (cf *CloudFlareCache) Touch(key string) (bool, error) {
	u := cf.url + "/" + key
	resp, err := cf.httpClient.Head(u)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	cacheStatus := resp.Header.Get("cf-cache-status")
	hit := cacheStatus == "HIT" || cacheStatus == "EXPIRED"
	return hit, nil
}
