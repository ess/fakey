package fakey

import (
	"errors"
	"net/url"

	"github.com/engineyard/eycore/core"
)

type Client struct {
	requests  map[string][]string
	responses map[string][]string
}

func (api *Client) setup() {
	if api.responses == nil {
		api.responses = make(map[string][]string)
	}

	if api.requests == nil {
		api.requests = make(map[string][]string)
	}
}

func (api *Client) Requests(method string) []string {
	var requests []string

	api.setup()

	requests = append(requests, api.requests[method]...)

	return requests
}

func (api *Client) AddResponse(method string, response string) {
	api.setup()

	api.responses[method] = append(api.responses[method], response)
}

func (api *Client) consume(method string) (string, error) {
	var response string

	api.setup()

	switch len(api.responses[method]) {
	case 0:
		return response, errors.New("No response")
	case 1:
		response = api.responses[method][0]
		api.responses[method] = nil
	default:
		response, api.responses[method] = api.responses[method][0], api.responses[method][1:]
	}

	return response, nil
}

func (api *Client) handle(method string, path string) ([]byte, error) {
	api.setup()

	api.requests[method] = append(api.requests[method], path)

	response, err := api.consume(method)
	if err != nil {
		return nil, err
	}

	return []byte(response), nil
}

func (api *Client) Get(path string, params url.Values) ([]byte, error) {
	return api.handle("get", path)
}

func (api *Client) Post(path string, params url.Values, body core.Body) ([]byte, error) {
	return api.handle("post", path)
}

func (api *Client) Put(path string, params url.Values, body core.Body) ([]byte, error) {
	return api.handle("put", path)
}

func (api *Client) Delete(path string, params url.Values) ([]byte, error) {
	return api.handle("delete", path)
}
