package oktools

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/okgotool/okgoweb/okmonitor"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	restyClientCache map[string]*resty.Client = map[string]*resty.Client{}
)

func GetHttpClient(timeout time.Duration) *http.Client {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	return client
}

func GetRestyRetryClient(timeout time.Duration) *resty.Client {
	client := GetHttpClient(timeout)

	restyClient := resty.NewWithClient(client)

	restyClient.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			if r == nil || r.StatusCode() == 0 || (r.StatusCode() >= http.StatusLocked && r.StatusCode() < http.StatusNotExtended) {
				return true
			} else {
				return false
			}
		},
	)

	return restyClient
}

func AddNotSuccessRetryCondition(r *resty.Client) *resty.Client {
	r.AddRetryCondition(
		func(r *resty.Response, err error) bool {
			if r == nil || r.StatusCode() == 0 || (r.StatusCode() >= http.StatusLocked && r.StatusCode() < http.StatusNotExtended) {
				return true
			} else if r.StatusCode() < http.StatusBadRequest { // StatusBadRequest=400
				return true
			} else {
				return false
			}
		},
	)

	return r
}

// 先从缓存查，查不到创建新的
// retryCount: 重试次数;
// retryWaitMills: 每次请求之间等待的毫秒数;
// totalTimeoutSeconds 包括重试，总超时时间；
// 每个请求超时时间 ~= totalTimeoutSeconds / retryCount
func GetRestyClient(retryCount int, retryWaitMills int, totalTimeoutSeconds int) *resty.Client {
	key := getRestyClientCacheKey(retryCount, retryWaitMills, totalTimeoutSeconds)
	if c, ok := restyClientCache[key]; ok {
		return c
	} else {
		c := GetNewRestyClient(retryCount, retryWaitMills, totalTimeoutSeconds)
		restyClientCache[key] = c
		return c
	}
}

// retryCount: 重试次数;
// retryWaitMills: 每次请求之间等待的毫秒数;
// totalTimeoutSeconds 包括重试，总超时时间；
// 每个请求超时时间 ~= totalTimeoutSeconds / retryCount
func GetNewRestyClient(retryCount int, retryWaitMills int, totalTimeoutSeconds int) *resty.Client {
	if retryCount < 1 {
		retryCount = 1
	}

	timeOut := 60 * time.Second
	if totalTimeoutSeconds > 0 {
		timeOutPerRequest := totalTimeoutSeconds / retryCount
		timeOut = time.Duration(timeOutPerRequest) * time.Second
	}

	restyClient := GetRestyRetryClient(timeOut)
	restyClient.SetRetryCount(retryCount)

	if retryWaitMills > 0 {
		waitTime := time.Duration(retryWaitMills) * time.Millisecond
		// maxWaitTime := time.Duration(retryWaitMills) * time.Millisecond
		restyClient.SetRetryWaitTime(waitTime)
		restyClient.SetRetryMaxWaitTime(waitTime)
	}

	return restyClient
}

func getRestyClientCacheKey(retryCount int, retryWaitMills int, totalTimeoutSeconds int) string {
	return fmt.Sprintf("%d-%d-%d", retryCount, retryWaitMills, totalTimeoutSeconds)
}

// retryCount: 重试次数;
// retryWaitMills: 每次请求之间等待的毫秒数;
// totalTimeoutSeconds 包括重试，总超时时间；
// 每个请求超时时间 ~= totalTimeoutSeconds / retryCount
func GetMonitorRestyClient(retryCount int, retryWaitMills int, totalTimeoutSeconds int) *resty.Client {
	restyClient := GetRestyClient(retryCount, retryWaitMills, totalTimeoutSeconds)

	if okmonitor.ApiCallMonitorEnabled {
		// restyClient.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// 	// Now you have access to Client and current Request object
		// 	// manipulate it as per your need
		// 	//host, path := GetHostFromUrl(req.URL)
		// 	//monitor.ApiCallNumTotal.With(prometheus.Labels{
		// 	//	"host":   host,
		// 	//	"url":    path,
		// 	//	"method": req.Method,
		// 	//}).Inc()

		// 	return nil // if its success otherwise return error
		// })

		restyClient.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
			// Now you have access to Client and current Response object
			// manipulate it as per your need
			latency := resp.Time().Seconds()
			statusCode := resp.StatusCode()
			host, path := GetHostFromUrl(resp.Request.URL)
			okmonitor.ApiCallRequestDuration.With(prometheus.Labels{
				"host":   host,
				"url":    path,
				"method": resp.Request.Method,
				"code":   fmt.Sprintf("%d", statusCode),
			}).Observe(latency)

			return nil // if its success otherwise return error
		})

		//restyClient.OnError(func(req *resty.Request, err error) {
		//	host, path := GetHostFromUrl(req.URL)
		//	monitor.ApiCallErrorNumTotal.With(prometheus.Labels{
		//		"host":   host,
		//		"url":    path,
		//		"method": req.Method,
		//	}).Inc()
		//})
	}

	return restyClient
}

func GetAuthRequest(authorization string) *resty.Request {
	r := GetRestyClient(3, 100, 90).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authorization)

	return r
}

func GetShortAuthRequest(authorization string) *resty.Request {
	r := GetRestyClient(2, 100, 20).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authorization)

	return r
}

// retryCount: 重试次数;
// retryWaitMills: 每次请求之间等待的毫秒数;
// totalTimeoutSeconds 包括重试，总超时时间；
// authorization: header Authorization
func GetCusAuthRequest(retryCount int, retryWaitMills int, totalTimeoutSeconds int, authorization string) *resty.Request {
	r := GetRestyClient(retryCount, retryWaitMills, totalTimeoutSeconds).R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", authorization)

	return r
}

// return host, path
func GetHostFromUrl(urlStr string) (string, string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return "", ""
	}

	return u.Host, u.Path
}
