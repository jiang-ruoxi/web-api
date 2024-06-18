package utils

import (
	"bytes"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func New() Http {
	return newBaseHttp(&http.Client{})
}

func Post(url string) Request {
	return New().Post(url)
}
func Delete(url string) Request {
	return New().Delete(url)
}
func Put(url string) Request {
	return New().Put(url)
}
func Get(url string) Request {
	return New().Get(url)
}

type Http interface {
	Timeout(timeout time.Duration) Http
	Post(url string) Request
	Delete(url string) Request
	Put(url string) Request
	Get(url string) Request
}

type Request interface {
	AddHeader(key, value string) Request
	AddHeaders(headers map[string]string) Request
	AddQuery(key, value string) Request
	AddQueries(queries map[string]string) Request
	AddBody(body any) Request
	Do() (response Response, err error)
}

type Response interface {
	Status() string
	StatusCode() int
	Url() string
	Body() io.ReadCloser
	BodyBytes() ([]byte, error)
	BodyString() (string, error)
	UnmarshalBody(data any) (err error)
}

type BaseHttp struct {
	client *http.Client
}

func (b *BaseHttp) Post(url string) Request {
	return b.request(http.MethodPost, url)
}

func (b *BaseHttp) Delete(url string) Request {
	return b.request(http.MethodDelete, url)
}

func (b *BaseHttp) Put(url string) Request {
	return b.request(http.MethodPut, url)
}

func (b *BaseHttp) Get(url string) Request {
	return b.request(http.MethodGet, url)
}

func (b *BaseHttp) Timeout(timeout time.Duration) Http {
	b.client.Timeout = timeout
	return b
}

func (b *BaseHttp) request(method string, url string) Request {
	return newBaseRequest(b.client, method, url)
}

func newBaseHttp(client *http.Client) *BaseHttp {
	return &BaseHttp{
		client: client,
	}
}

type BaseRequest struct {
	method  string
	url     string
	headers map[string]string
	queries map[string]string
	body    any
	client  *http.Client
}

func newBaseRequest(client *http.Client, method string, url string) *BaseRequest {
	return &BaseRequest{
		client:  client,
		method:  method,
		url:     url,
		headers: map[string]string{},
		queries: map[string]string{},
	}
}
func (b *BaseRequest) AddHeaders(headers map[string]string) Request {
	if headers == nil {
		return b
	}
	for key, value := range headers {
		b.AddHeader(key, value)
	}
	return b
}
func (b *BaseRequest) AddHeader(key, value string) Request {
	b.headers[key] = value
	return b
}
func (b *BaseRequest) AddQuery(key, value string) Request {
	b.queries[key] = value
	return b
}
func (b *BaseRequest) AddQueries(queries map[string]string) Request {
	for key, value := range queries {
		b.queries[key] = value
	}
	return b
}
func (b *BaseRequest) AddBody(body any) Request {
	b.body = body
	return b
}
func (b *BaseRequest) Do() (response Response, err error) {
	var errPath = "BaseRequest.Do error"
	var httpResponse *http.Response
	for key, value := range b.queries {
		value = url.QueryEscape(value)
		var urlLen = len(b.url)
		var last = b.url[urlLen-1 : urlLen]
		if last == "?" || last == "&" {
			b.url += fmt.Sprintf("%s=%s", key, value)
			continue
		}
		if strings.Contains(b.url, "?") {
			b.url += fmt.Sprintf("&%s=%s", key, value)
			continue
		}
		b.url += fmt.Sprintf("?%s=%s", key, value)
		continue
	}
	var bodyBytes []byte
	if bodyBytes, err = jsoniter.Marshal(b.body); err != nil {
		err = errors.Wrap(err, errPath)
		return
	}
	var request *http.Request
	if request, err = http.NewRequest(b.method, b.url, bytes.NewReader(bodyBytes)); err != nil {
		err = errors.Wrap(err, errPath)
		return
	}
	for key, value := range b.headers {
		request.Header.Set(key, value)
	}
	if httpResponse, err = b.client.Do(request); err != nil {
		err = errors.Wrap(err, errPath)
		return
	}
	if response, err = newResponse(httpResponse); err != nil {
		err = errors.Wrap(err, errPath)
		return
	}
	return
}

type BaseResponse struct {
	body         any
	httpResponse *http.Response
}

func newResponse(httpResponse *http.Response) (response Response, err error) {
	if httpResponse == nil {
		err = errors.Wrap(errors.New("http response is nil"), "http.newResponse error")
		return
	}
	response = &BaseResponse{
		httpResponse: httpResponse,
	}
	return
}
func (b *BaseResponse) UnmarshalBody(body any) (err error) {
	var httpResponseBodyBytes []byte
	if httpResponseBodyBytes, err = io.ReadAll(b.httpResponse.Body); err != nil {
		err = errors.Wrap(err, "BaseResponse.Body error")
		return
	}
	if err = jsoniter.Unmarshal(httpResponseBodyBytes, body); err != nil {
		err = errors.Wrap(err, "BaseResponse.Body error")
		return
	}
	return
}
func (b *BaseResponse) Status() string {
	return b.httpResponse.Status
}
func (b *BaseResponse) StatusCode() int {
	return b.httpResponse.StatusCode
}
func (b *BaseResponse) Body() io.ReadCloser {
	return b.httpResponse.Body
}
func (b *BaseResponse) BodyBytes() ([]byte, error) {
	return io.ReadAll(b.Body())
}
func (b *BaseResponse) BodyString() (string, error) {
	bytes, err := b.BodyBytes()
	if err != nil {
		return "", err
	}
	s := string(bytes)
	return s, err
}
func (b *BaseResponse) Url() (url string) {
	if b.httpResponse == nil {
		return
	}
	var request = b.httpResponse.Request
	if request == nil {
		return
	}
	var uri = request.URL
	if uri == nil {
		return
	}
	url = uri.String()
	return
}
func GetFileNameByUrl(url string) (fileName string) {
	var urlSplits = strings.Split(url, "/")
	if len(urlSplits) < 1 {
		return
	}
	fileName = urlSplits[len(urlSplits)-1]
	var fileNameSplits = strings.Split(fileName, ".")
	if len(fileNameSplits) < 1 {
		return
	}
	fileName = fileNameSplits[0]
	return
}
