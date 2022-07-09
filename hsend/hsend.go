package hsend

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	goquery "github.com/google/go-querystring/query"
)

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}

type HSend struct {
	httpClient      Doer
	method          string
	rawURL          string
	header          http.Header
	queryStructs    []interface{}
	bodyProvider    BodyProvider
	responseDecoder ResponseDecoder
}

func New() *HSend {
	return &HSend{
		httpClient:      http.DefaultClient,
		method:          "GET",
		header:          make(http.Header),
		queryStructs:    make([]interface{}, 0),
		responseDecoder: jsonDecoder{},
	}
}

// New returns a copy of an HSend client. This child HSend client will have properties
// from the Parent HSend client.
// Note that query and body values are copied to avoid pointers shenanigans :)
func (h *HSend) New() *HSend {
	headerCopy := make(http.Header)
	for k, v := range h.header {
		headerCopy[k] = v
	}

	return &HSend{
		httpClient:      h.httpClient,
		method:          h.method,
		rawURL:          h.rawURL,
		header:          headerCopy,
		queryStructs:    append([]interface{}{}, h.queryStructs...),
		bodyProvider:    h.bodyProvider,
		responseDecoder: h.responseDecoder,
	}
}

// Http Client

func (h *HSend) Client(httpClient *http.Client) *HSend {
	return h.Doer(httpClient)
}

func (h *HSend) Doer(doer Doer) *HSend {
	if doer == nil {
		h.httpClient = http.DefaultClient
	} else {
		h.httpClient = doer
	}

	return h
}

// Method

// Get sets the HSend method to GET and sets the given pathURL.
func (h *HSend) Get(pathURL string) *HSend {
	h.method = "GET"
	return h.Path(pathURL)
}

// Post sets the HSend method to POST and sets the given pathURL.
func (h *HSend) Post(pathURL string) *HSend {
	h.method = "POST"
	return h.Path(pathURL)
}

// Delete sets the HSend method to DELETE and sets the given pathURL.
func (h *HSend) Delete(pathURL string) *HSend {
	h.method = "DELETE"
	return h.Path(pathURL)
}

// URL

// Base sets the rawURL. If you intend to extend the url with Path,
// baseURL should be specified with a trailing slash.
func (h *HSend) Base(rawURL string) *HSend {
	h.rawURL = rawURL
	return h
}

// Path extends the rawURL with the given path by resolving the reference to
// an absolute URL. If parsing errors occur, the rawURL is left unmodified.
func (h *HSend) Path(path string) *HSend {
	baseURL, baseErr := url.Parse(h.rawURL)
	pathURL, pathErr := url.Parse(path)
	if baseErr == nil && pathErr == nil {
		h.rawURL = baseURL.ResolveReference(pathURL).String()
		return h
	}
	return h
}

func (h *HSend) QueryStruct(queryStruct interface{}) *HSend {
	if queryStruct != nil {
		h.queryStructs = append(h.queryStructs, queryStruct)
	}
	return h
}

// Header

// Add adds the (key, value) pair in Headers, appending values
// associated with key. Header keys are canonicalized.
func (h *HSend) Add(key, value string) *HSend {
	h.header.Add(key, value)
	return h
}

// Set sets the key, value pair in Headers, replacing existing values
// associated with key. Header keys are canonicalized.
func (h *HSend) Set(key, value string) *HSend {
	h.header.Set(key, value)
	return h
}

// Body

// BodyProvider sets the Sling's body provider.
func (h *HSend) BodyProvider(body BodyProvider) *HSend {
	if body == nil {
		return h
	}
	h.bodyProvider = body

	ct := body.ContentType()
	if ct != "" {
		h.Set("Content-Type", ct)
	}

	return h
}

func (h *HSend) BodyMultiPartForm(bodyForm interface{}, contentType string) *HSend {
	if bodyForm == nil {
		return h
	}

	return h.BodyProvider(multipartFormProvider{payload: bodyForm, contentType: contentType})
}

// Requests

// Request returns a new http.Request created with the Sling properties.
// Returns any errors parsing the rawURL, encoding query structs, encoding
// the body, or creating the http.Request.
func (h *HSend) Request() (*http.Request, error) {
	reqURL, err := url.Parse(h.rawURL)
	if err != nil {
		return nil, err
	}

	fmt.Println("hmm", h.queryStructs)
	err = addQueryStructs(reqURL, h.queryStructs)
	if err != nil {
		return nil, err
	}

	var body io.Reader
	if h.bodyProvider != nil {
		body, err = h.bodyProvider.Body()
		if err != nil {
			return nil, err
		}
	}
	req, err := http.NewRequest(h.method, reqURL.String(), body)
	if err != nil {
		return nil, err
	}
	addHeaders(req, h.header)

	return req, err
}

// addQueryStructs parses url tagged query structs using go-querystring to
// encode them to url.Values and format them onto the url.RawQuery. Any
// query parsing or encoding errors are returned.
func addQueryStructs(reqURL *url.URL, queryStructs []interface{}) error {
	urlValues, err := url.ParseQuery(reqURL.RawQuery)
	if err != nil {
		return err
	}
	// encodes query structs into a url.Values map and merges maps
	for _, queryStruct := range queryStructs {
		fmt.Println("struct? ", queryStruct)
		val := reflect.ValueOf(queryStruct)
		queryValues, err := goquery.Values(queryStruct)
		if err != nil {
			fmt.Println("query struct error", queryStruct, val.Kind())
			return err
		}
		fmt.Println("no query struct error")
		for key, values := range queryValues {
			for _, value := range values {
				urlValues.Add(key, value)
			}
		}
	}
	// url.Values format to a sorted "url encoded" string, e.g. "key=val&foo=bar"
	reqURL.RawQuery = urlValues.Encode()
	return nil
}

// addHeaders adds the key, value pairs from the given http.Header to the
// request. Values for existing keys are appended to the keys values.
func addHeaders(req *http.Request, header http.Header) {
	for key, values := range header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
}

// Sending

// Receive creates a new HTTP request and returns the response. Success
// responses (2XX) are JSON decoded into the value pointed to by successV and
// other responses are JSON decoded into the value pointed to by failureV.
// If the status code of response is 204(no content) or the Content-Lenght is 0,
// decoding is skipped. Any error creating the request, sending it, or decoding
// the response is returned.
// Receive is shorthand for calling Request and Do.
func (h *HSend) Receive(successV, failureV interface{}) (*http.Response, error) {
	req, err := h.Request()
	if err != nil {
		return nil, err
	}
	return h.Do(req, successV, failureV)
}

func (h *HSend) Do(req *http.Request, successV, failureV interface{}) (*http.Response, error) {
	fmt.Println("Got to do", req.URL, req.Header)
	resp, err := h.httpClient.Do(req)
	if err != nil {
		return resp, err
	}

	fmt.Println("Got to resp")

	// when err is nil, resp contains a non-nil resp.Body which must be closed
	defer resp.Body.Close()

	// The default HTTP client's Transport may not
	// reuse HTTP/1.x "keep-alive" TCP connections if the Body is
	// not read to completion and closed.
	// See: https://golang.org/pkg/net/http/#Response
	fmt.Println("Got to copy")
	defer io.Copy(ioutil.Discard, resp.Body)

	// Don't try to decode on 204s or Content-Length is 0
	if resp.StatusCode == http.StatusNoContent || resp.ContentLength == 0 {
		return resp, nil
	}

	// Decode from json
	if successV != nil || failureV != nil {
		fmt.Println("Got to decode")
		err = decodeResponse(resp, h.responseDecoder, successV, failureV)
	}
	return resp, err
}

// decodeResponse decodes response Body into the value pointed to by successV
// if the response is a success (2XX) or into the value pointed to by failureV
// otherwise. If the successV or failureV argument to decode into is nil,
// decoding is skipped.
// Caller is responsible for closing the resp.Body.
func decodeResponse(resp *http.Response, decoder ResponseDecoder, successV, failureV interface{}) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		if successV != nil {
			return decoder.Decode(resp, successV)
		}
	} else {
		if failureV != nil {
			return decoder.Decode(resp, failureV)
		}
	}
	return nil
}
