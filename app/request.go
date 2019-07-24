package app

import (
    "fmt"
    "errors"
    "net/http"
    "github.com/mercadolibre/go-meli-toolkit/restful/rest"
    "github.com/mercadolibre/go-meli-toolkit/goutils/apierrors"
)

type Request struct {
    Method   string
    BaseURL  string
    Endpoint string
    Payload  map[string]interface{}
}

type Response struct {
    StatusCode int
    Headers    http.Header
    Payload    []byte
}

func (request *Request) Do() (*Response, apierrors.ApiError) {
    const (
        apiFormat           = "%s%s"
        errNilResponse      = "nil response received from %s"
        errInvalidMethod    = "invalid method"
        errExecutingRequest = "error executing request"
    )

    full := fmt.Sprintf(
        apiFormat,
        request.BaseURL,
        request.Endpoint,
    )

    var response *rest.Response

    switch request.Method {
    case http.MethodGet:
        response = rest.Get(full)
    case http.MethodHead:
        response = rest.Head(full)
    case http.MethodPost:
        response = rest.Post(full, request.Payload)
    case http.MethodPut:
        response = rest.Put(full, request.Payload)
    case http.MethodPatch:
        response = rest.Patch(full, request.Payload)
    case http.MethodDelete:
        response = rest.Delete(full)
    case http.MethodOptions:
        response = rest.Options(full)
    default:
        return nil, apierrors.NewBadRequestApiError(errInvalidMethod)
    }

    if response == nil {
        err := errors.New(fmt.Sprintf(errNilResponse, full))
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
