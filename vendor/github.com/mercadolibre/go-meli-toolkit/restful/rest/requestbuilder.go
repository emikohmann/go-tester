package rest

import (
	"net/http"
	"sync"
	"time"

	"github.com/mercadolibre/go-meli-toolkit/golimiter"
	"github.com/mercadolibre/go-meli-toolkit/goutils"
	"github.com/mercadolibre/go-meli-toolkit/restful/rest/breaker"
	"github.com/mercadolibre/go-meli-toolkit/restful/rest/retry"
)

// Retry rate limiter
var retryLimiter *golimiter.Limiter

// DefaultTimeout is the default timeout for all clients.
// DefaultConnectTimeout is the time it takes to make a connection
// Type: time.Duration
var DefaultTimeout = 500 * time.Millisecond

var DefaultConnectTimeout = 1500 * time.Millisecond

// DefaultMaxIdleConnsPerHost is the default maxium idle connections to have
// per Host for all clients, that use *any* RequestBuilder that don't set
// a CustomPool
var DefaultMaxIdleConnsPerHost = 2

// ContentType represents the Content Type for the Body of HTTP Verbs like
// POST, PUT, and PATCH
type ContentType int

const (
	// JSON represents a JSON Content Type
	JSON ContentType = iota

	// XML represents an XML Content Type
	XML

	// BYTES represents a plain Content Type
	BYTES
)

// RequestBuilder is the baseline for creating requests
// There's a Default Builder that you may use for simple requests
// RequestBuilder si thread-safe, and you should store it for later re-used.
type RequestBuilder struct {

	// Headers to be send in the request
	Headers    http.Header
	headersMtx sync.RWMutex

	// Complete request time out.
	Timeout time.Duration

	// Connection timeout, it bounds the time spent obtaining a successful connection
	ConnectTimeout time.Duration

	// Base URL to be used for each Request. The final URL will be BaseURL + URL.
	BaseURL string

	// ContentType
	ContentType ContentType

	// Enable internal caching of Responses
	EnableCache bool

	// Disable timeout and default timeout = no timeout
	DisableTimeout bool

	// Set the http client to follow a redirect if we get a 3xx response
	FollowRedirect bool

	// Create a CustomPool if you don't want to share the *transport*, with others
	// RequestBuilder
	CustomPool *CustomPool

	// Set Basic Auth for this RequestBuilder
	BasicAuth *BasicAuth

	// Set an specific User Agent for this RequestBuilder
	UserAgent string

	// Public for custom fine tuning
	Client *http.Client

	clientMtxOnce sync.Once

	// Optional retry strategy
	RetryStrategy retry.RetryStrategy

	// If true, automatically uncompress the response body for supported formats
	UncompressResponse bool

	//Metrics report config
	MetricsConfig MetricsReportConfig

	// Optional circuit breaker
	circuitBreaker breaker.CircuitBreaker
}

type MetricsReportConfig struct {
	// Every metric will report a tag called target_id with this value. It can be used to filter metrics
	TargetId string
	// True to avoid sending http connections info (connections requests, connections new, connections request result)
	DisableHttpConnectionsMetrics bool
	// True to avoid sending api call metrics (api call requests, api call time, api call result)
	DisableApiCallMetrics bool
}

// CustomPool defines a separate internal *transport* and connection pooling.
type CustomPool struct {
	MaxIdleConnsPerHost int
	Proxy               string

	// Public for custom fine tuning
	Transport http.RoundTripper
}

// BasicAuth gives the possibility to set UserName and Password for a given
// RequestBuilder. Basic Auth is used by some APIs
type BasicAuth struct {
	UserName string
	Password string
}

// Add a circuit breaker to the RequestBuilder
func (rb *RequestBuilder) AddCircuitBreaker(cb breaker.CircuitBreaker) {
	rb.circuitBreaker = cb
	rb.circuitBreaker.SetTarget(rb.MetricsConfig.TargetId)
}

// Get issues a GET HTTP verb to the specified URL.
//
// In Restful, GET is used for "reading" or retrieving a resource.
// Client should expect a response status code of 200(OK) if resource exists,
// 404(Not Found) if it doesn't, or 400(Bad Request).
func (rb *RequestBuilder) Get(url string) *Response {
	return rb.DoRequest(http.MethodGet, url, nil)
}

// Post issues a POST HTTP verb to the specified URL.
//
// In Restful, POST is used for "creating" a resource.
// Client should expect a response status code of 201(Created), 400(Bad Request),
// 404(Not Found), or 409(Conflict) if resource already exist.
//
// Body could be any of the form: string, []byte, struct & map.
func (rb *RequestBuilder) Post(url string, body interface{}) *Response {
	return rb.DoRequest(http.MethodPost, url, body)
}

// Put issues a PUT HTTP verb to the specified URL.
//
// In Restful, PUT is used for "updating" a resource.
// Client should expect a response status code of of 200(OK), 404(Not Found),
// or 400(Bad Request). 200(OK) could be also 204(No Content)
//
// Body could be any of the form: string, []byte, struct & map.
func (rb *RequestBuilder) Put(url string, body interface{}) *Response {
	return rb.DoRequest(http.MethodPut, url, body)
}

// Patch issues a PATCH HTTP verb to the specified URL.
//
// In Restful, PATCH is used for "partially updating" a resource.
// Client should expect a response status code of of 200(OK), 404(Not Found),
// or 400(Bad Request). 200(OK) could be also 204(No Content)
//
// Body could be any of the form: string, []byte, struct & map.
func (rb *RequestBuilder) Patch(url string, body interface{}) *Response {
	return rb.DoRequest(http.MethodPatch, url, body)
}

// Delete issues a DELETE HTTP verb to the specified URL
//
// In Restful, DELETE is used to "delete" a resource.
// Client should expect a response status code of of 200(OK), 404(Not Found),
// or 400(Bad Request).
func (rb *RequestBuilder) Delete(url string) *Response {
	return rb.DoRequest(http.MethodDelete, url, nil)
}

// Head issues a HEAD HTTP verb to the specified URL
//
// In Restful, HEAD is used to "read" a resource headers only.
// Client should expect a response status code of 200(OK) if resource exists,
// 404(Not Found) if it doesn't, or 400(Bad Request).
func (rb *RequestBuilder) Head(url string) *Response {
	return rb.DoRequest(http.MethodHead, url, nil)
}

// Options issues a OPTIONS HTTP verb to the specified URL
//
// In Restful, OPTIONS is used to get information about the resource
// and supported HTTP verbs.
// Client should expect a response status code of 200(OK) if resource exists,
// 404(Not Found) if it doesn't, or 400(Bad Request).
func (rb *RequestBuilder) Options(url string) *Response {
	return rb.DoRequest(http.MethodOptions, url, nil)
}

// AsyncGet is the *asynchronous* option for GET.
// The go routine calling AsyncGet(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncGet(url string, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Get(url))
	}()
}

// AsyncPost is the *asynchronous* option for POST.
// The go routine calling AsyncPost(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncPost(url string, body interface{}, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Post(url, body))
	}()
}

// AsyncPut is the *asynchronous* option for PUT.
// The go routine calling AsyncPut(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncPut(url string, body interface{}, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Put(url, body))
	}()
}

// AsyncPatch is the *asynchronous* option for PATCH.
// The go routine calling AsyncPatch(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncPatch(url string, body interface{}, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Patch(url, body))
	}()
}

// AsyncDelete is the *asynchronous* option for DELETE.
// The go routine calling AsyncDelete(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncDelete(url string, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Delete(url))
	}()
}

// AsyncHead is the *asynchronous* option for HEAD.
// The go routine calling AsyncHead(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncHead(url string, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Head(url))
	}()
}

// AsyncOptions is the *asynchronous* option for OPTIONS.
// The go routine calling AsyncOptions(), will not be blocked.
//
// Whenever the Response is ready, the *f* function will be called back.
func (rb *RequestBuilder) AsyncOptions(url string, f func(*Response)) {
	go func() {
		defer goutils.Recover()
		f(rb.Options(url))
	}()
}

// ForkJoin let you *fork* requests, and *wait* until all of them have return.
//
// Concurrent has methods for Get, Post, Put, Patch, Delete, Head & Options,
// with the almost the same API as the synchronous methods.
// The difference is that these methods return a FutureResponse, which holds a pointer to
// Response. Response inside FutureResponse is nil until the request has finished.
//
// 	var futureA, futureB *rest.FutureResponse
//
// 	rest.ForkJoin(func(c *rest.Concurrent){
//		futureA = c.Get("/url/1")
//		futureB = c.Get("/url/2")
//	})
//
//	fmt.Println(futureA.Response())
//	fmt.Println(futureB.Response())
//
func (rb *RequestBuilder) ForkJoin(f func(*Concurrent)) {

	c := new(Concurrent)
	c.reqBuilder = rb

	f(c)

	c.wg.Add(c.list.Len())

	for e := c.list.Front(); e != nil; e = e.Next() {
		go e.Value.(func())()
	}

	c.wg.Wait()
}

func init() {
	retryLimiter = golimiter.New(1000, time.Second)

	defaultCheckRedirectFunc = http.Client{}.CheckRedirect
}
