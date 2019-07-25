package app

import (
    "fmt"
    "errors"
    "net/http"
    "github.com/mercadolibre/go-meli-toolkit/restful/rest"
    "github.com/mercadolibre/go-meli-toolkit/goutils/apierrors"
)

type Request struct {
    Method  string
    URL     string
    Payload map[string]interface{}
}

type Response struct {
    StatusCode int
    Headers    http.Header
    Payload    []byte
}

func (request *Request) Do() (*Response, apierrors.ApiError) {
    const (
        errNilResponse      = "nil response received from %s"
        errInvalidMethod    = "invalid method"
        errExecutingRequest = "error executing request"
    )

    var response *rest.Response

    // Change timeout

    switch request.Method {
    case http.MethodGet:
        response = rest.Get(request.URL)
    case http.MethodHead:
        response = rest.Head(request.URL)
    case http.MethodPost:
        response = rest.Post(request.URL, request.Payload)
    case http.MethodPut:
        response = rest.Put(request.URL, request.Payload)
    case http.MethodPatch:
        response = rest.Patch(request.URL, request.Payload)
    case http.MethodDelete:
        response = rest.Delete(request.URL)
    case http.MethodOptions:
        response = rest.Options(request.URL)
    default:
        return nil, apierrors.NewBadRequestApiError(errInvalidMethod)
    }

    if response == nil {
        err := errors.New(fmt.Sprintf(errNilResponse, request.URL))
        return nil, apierrors.NewInternalServerApiError(errExecutingRequest, err)
    }

    if response.Err != nil {
        return nil, apierrors.NewInternalServerApiError(errExecutingRequest, response.Err)
    }

    return &Response{
        StatusCode: response.StatusCode,
        Headers:    response.Header,
        Payload:    response.Bytes(),
    }, nil
}
