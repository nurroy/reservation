package requestapi

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"bitbucket.org/pinang-och-go/pkg/v1/utils/logger"

	"bitbucket.org/pinang-och-go/pkg/v1/cert"
	"bitbucket.org/pinang-och-go/pkg/v1/converter"
	"bitbucket.org/pinang-och-go/pkg/v1/utils/constants"
	"bitbucket.org/pinang-och-go/pkg/v1/utils/errors"
)

// Info is the http req info
type ReqInfo struct {
	URL        string
	HeaderInfo HeaderInfoSchema
	Body       []byte
}

type HeaderInfoSchema struct {
	ContentType   string
	Auth          string
	XBRISignature string
	XBRITimestamp string
}

type ResInfo struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// RequestAPI this is the interface for creating new external request
type RequestAPI interface {
	GenerateHeaderInfoSchema(requestBytes []byte, opt Options) HeaderInfoSchema
	GetBaseURL() string
	DO(ctx context.Context, req Request, opt Options) (*ResInfo, error)
}

type requestAPI struct {
	keyHMAC string
	baseURL string
	client  http.Client
}

// Options is the struct for creating header schema (default)
// Types is equal to content Type ex. "application/json", Method is equal to method used ex. POST
// Auth is used if there is authentication needed (cookie / token)
// NeedSignature is to append Signature to the request (xBRISignature)
type Options struct {
	Types         string
	Method        string
	Auth          string
	Timeout       time.Duration
	NeedSignature bool
}

// Request is requestAPI request format
// URLPath is the url path after the base url, ex. if you have complete URL "http://xyz?a=123"
// baseURL defined as "http://xyz" the URLPath should be "?a=123"
type Request struct {
	URLPath   string
	Body      []byte
	AddHeader map[string]string
}

const (
	//JSONType request api options type for json request
	JSONType = "application/json"
	//XMLType request api options type for xml request
	XMLType = "application/xml"

	//MethodPOST request api options method for POST
	MethodPOST = "POST"
	//MethodGET request api options method for GET
	MethodGET = "GET"
)

// NewRequest is function to create request api details
func NewRequest(HMACKey, baseURL string, crt *cert.Cert) RequestAPI {
	// define the default http client
	defaultRoundTripper := http.DefaultTransport
	defaultTransportPtr, _ := defaultRoundTripper.(*http.Transport)
	tr := defaultTransportPtr.Clone()

	// need cert check
	if crt != nil {
		tr.TLSClientConfig = &tls.Config{
			RootCAs: crt.Pool,
		}
		if crt.AllowSkip {
			tr.TLSClientConfig = &tls.Config{
				InsecureSkipVerify: true,
			}
		}
	}

	// define the max connection idle
	tr.MaxIdleConns = 100
	tr.MaxIdleConnsPerHost = 20

	c := http.Client{
		Timeout:   5 * time.Minute,
		Transport: tr,
	}

	return &requestAPI{
		keyHMAC: HMACKey,
		baseURL: baseURL,
		client:  c,
	}
}

// GET API request
func GET(reqinf *ReqInfo, timeout time.Duration) ([]byte, error) {
	req, err := http.NewRequest("GET", reqinf.URL, bytes.NewReader(reqinf.Body))
	if err != nil {
		return nil, err
	}

	//req.Header.Add("X-BRI-Signature", reqinf.HeaderInfo.XBRISignature)
	//req.Header.Add("X-BRI-Timestamp", reqinf.HeaderInfo.XBRITimestamp)
	req.Header.Add("Authorization", reqinf.HeaderInfo.Auth)
	req.Header.Add("Content-Type", reqinf.HeaderInfo.ContentType)

	logger.ServiceRequestLog(reqinf.URL, req, string(reqinf.Body))

	// execute
	cl := &http.Client{
		Timeout: timeout,
	}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%v", res.StatusCode))
	}

	// read body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	logger.ServiceResponseLog(reqinf.URL, string(buf))

	return buf, nil
}

// POST API request
func POST(reqinf *ReqInfo, timeout time.Duration) (*ResInfo, error) {
	req, err := http.NewRequest("POST", reqinf.URL, bytes.NewReader(reqinf.Body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", reqinf.HeaderInfo.ContentType)
	req.Header.Add("Authorization", reqinf.HeaderInfo.Auth)
	req.Header.Add("X-BRI-Signature", reqinf.HeaderInfo.XBRISignature)
	req.Header.Add("X-BRI-Timestamp", reqinf.HeaderInfo.XBRITimestamp)

	logger.ServiceRequestLog(reqinf.URL, req, string(reqinf.Body))

	// execute
	cl := &http.Client{
		Timeout: timeout,
	}

	res, err := cl.Do(req)

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%v", res.StatusCode))
	}

	// read body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	logger.ServiceResponseLog(reqinf.URL, string(buf))

	return &ResInfo{
		StatusCode: res.StatusCode,
		Body:       buf,
	}, nil
}

// GenerateHeaderInfoSchema general function to generate header info schema
func (r requestAPI) GenerateHeaderInfoSchema(requestBytes []byte, opt Options) HeaderInfoSchema {
	xBRITimestamp := time.Now().Format(constants.XBRITimeLayout)
	// generate message signature
	messageSignature := fmt.Sprintf("path=%s&verb=%s&token=%s&timestamp=%s&body=%s", r.baseURL, opt.Method, r.keyHMAC, xBRITimestamp, requestBytes)
	header := HeaderInfoSchema{
		ContentType:   opt.Types,
		Auth:          opt.Auth,
		XBRITimestamp: xBRITimestamp,
	}
	if opt.NeedSignature {
		header.XBRISignature = converter.GenerateHMAC256(r.keyHMAC, messageSignature)
	}
	return header
}

// GetBaseURL general function to get the request endpoint
func (r requestAPI) GetBaseURL() string {
	return r.baseURL
}

func (r requestAPI) DO(ctx context.Context, req Request, opt Options) (*ResInfo, error) {
	httpReq, err := r.GenerateRequest(ctx, req, opt)
	if err != nil {
		return nil, err
	}

	logger.ServiceRequestHttpLog(httpReq.URL, httpReq.Header, string(req.Body))

	res, err := r.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%v", res.StatusCode))
	}

	// read body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	logger.ServiceResponseLog(fmt.Sprintf("%s", httpReq.URL), string(buf))

	return &ResInfo{
		StatusCode: res.StatusCode,
		Body:       buf,
	}, nil
}

func (r requestAPI) GenerateRequest(ctx context.Context, reqData Request, opt Options) (*http.Request, error) {
	url := r.baseURL + reqData.URLPath
	req, err := http.NewRequestWithContext(ctx, opt.Method, url, bytes.NewReader(reqData.Body))
	if err != nil {
		return nil, err
	}

	xBRITimestamp := time.Now().Format(constants.XBRITimeLayout)
	// generate message signature for header and set the default header
	messageSignature := fmt.Sprintf("path=%s&verb=%s&token=%s&timestamp=%s&body=%s", url, opt.Method, r.keyHMAC, xBRITimestamp, reqData.Body)
	req.Header.Add("Content-Type", opt.Types)
	req.Header.Add("Authorization", opt.Auth)
	req.Header.Add("X-BRI-Timestamp", xBRITimestamp)

	if opt.NeedSignature {
		req.Header.Add("X-BRI-Signature", converter.GenerateHMAC256(r.keyHMAC, messageSignature))
	}

	for key, value := range reqData.AddHeader {
		req.Header.Add(key, value)
	}

	return req, nil
}

// POSTwithAdditionalHeader to create the new POST request with the additional header beside the standard header
func POSTwithAdditionalHeader(reqinf *ReqInfo, timeout time.Duration, headerMap map[string]string) (*ResInfo, error) {
	req, err := http.NewRequest("POST", reqinf.URL, bytes.NewReader(reqinf.Body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", reqinf.HeaderInfo.ContentType)
	req.Header.Add("Authorization", reqinf.HeaderInfo.Auth)
	req.Header.Add("X-BRI-Signature", reqinf.HeaderInfo.XBRISignature)
	req.Header.Add("X-BRI-Timestamp", reqinf.HeaderInfo.XBRITimestamp)

	for key, value := range headerMap {
		req.Header.Add(key, value)
	}

	logger.ServiceRequestLog(reqinf.URL, req, string(reqinf.Body))

	// execute
	cl := &http.Client{
		Timeout: timeout,
	}
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("%v", res.StatusCode))
	}

	// read body
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	logger.ServiceResponseLog(reqinf.URL, string(buf))

	return &ResInfo{
		StatusCode: res.StatusCode,
		Body:       buf,
	}, nil
}
