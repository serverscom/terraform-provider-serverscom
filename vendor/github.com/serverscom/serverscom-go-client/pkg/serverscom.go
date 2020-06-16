package serverscom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const defaultAPIEndpoint string = "https://api.servers.com/v1"

// Client manages communication with Servers.com API V1
type Client struct {
	baseURL   string
	UserAgent string

	CloudComputingInstances CloudComputingInstancesService
	Hosts                   HostsService
	L2Segments              L2SegmentsService

	Locations             LocationsService
	CloudComputingRegions CloudComputingRegionsService
	SSHKeys               SSHKeysService
	SSLCertificates       SSLCertificatesService

	token string

	client *http.Client
}

// NewClient builds a new client with token
func NewClient(token string) *Client {
	return NewClientWithEndpoint(token, defaultAPIEndpoint)
}

// NewClientWithEndpoint builds a new client with token and api endpoint
func NewClientWithEndpoint(token, baseURL string) *Client {
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
	}

	client := &http.Client{Transport: tr}

	scClient := &Client{
		baseURL:   baseURL,
		UserAgent: "go-serverscom-client",
		token:     token,
		client:    client,
	}

	scClient.configureResources()

	return scClient
}

// SetupUserAgent setups custom User-Agent header, by default: go-serverscom-client
func (cli *Client) SetupUserAgent(userAgent string) {
	cli.UserAgent = userAgent
}

func (cli *Client) configureResources() {
	cli.CloudComputingInstances = &CloudComputingInstancesHandler{cli}
	cli.Hosts = &HostsHandler{cli}
	cli.L2Segments = &L2SegmentsHandler{cli}
	cli.Locations = &LocationsHandler{cli}
	cli.CloudComputingRegions = &CloudComputingRegionsHandler{cli}
	cli.SSHKeys = &SSHKeysHandler{cli}
	cli.SSLCertificates = &SSLCertificatesHandler{cli}
}

func (cli *Client) buildURL(path string, values ...interface{}) string {
	return fmt.Sprintf(
		"%s%s",
		cli.baseURL,
		cli.buildPath(path, values...),
	)
}

func (cli *Client) buildPath(path string, values ...interface{}) string {
	return fmt.Sprintf(path, values...)
}

func (cli *Client) applyParams(endpointURL string, params map[string]string) string {
	if len(params) == 0 {
		return endpointURL
	}

	queryParams := url.Values{}

	for key, val := range params {
		queryParams.Set(key, val)
	}

	return fmt.Sprintf(
		"%s?%s",
		endpointURL,
		queryParams.Encode(),
	)
}

func (cli *Client) buildAndExecRequestWithResponse(ctx context.Context, method, endpointURL string, body []byte) (*http.Response, []byte, error) {
	var req *http.Request
	var err error

	if body != nil {
		reader := bytes.NewReader(body)
		req, err = http.NewRequest(method, endpointURL, reader)
	} else {
		req, err = http.NewRequest(method, endpointURL, nil)
	}

	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", cli.token))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Set("User-Agent", cli.UserAgent)

	resp, err := cli.client.Do(req.WithContext(ctx))

	if err != nil {
		return nil, nil, fmt.Errorf("Client request error: %q", err)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, nil, fmt.Errorf("Client can't read body: %q", err)
	}

	if resp.StatusCode < 400 {
		return resp, contents, nil
	}

	contentType := resp.Header.Get("Content-Type")

	var responseError = responseErrorWrapper{}

	if strings.HasPrefix(contentType, "application/json") {
		if err := json.Unmarshal(contents, &responseError); err != nil {
			return nil, nil, newParsingError(
				resp.StatusCode,
				string(contents),
				err,
			)
		}
	} else {
		responseError.Code = "UNKNOWN"
		responseError.Message = string(contents)
	}

	switch resp.StatusCode {
	case 400:
		return nil, nil, newBadRequestError(resp.StatusCode, responseError.Code, responseError.Message)
	case 401:
		return nil, nil, newUnauthorizedError(resp.StatusCode, responseError.Code, responseError.Message)
	case 403:
		return nil, nil, newForbiddenError(resp.StatusCode, responseError.Code, responseError.Message)
	case 404:
		return nil, nil, newNotFoundError(resp.StatusCode, responseError.Code, responseError.Message)
	case 409:
		return nil, nil, newConflictError(resp.StatusCode, responseError.Code, responseError.Message)
	case 422:
		return nil, nil, newUnprocessableEntityError(resp.StatusCode, responseError.Code, responseError.Message, responseError.Errors)
	case 500:
		return nil, nil, newInternalServerError(resp.StatusCode, responseError.Code, responseError.Message)
	default:
		return nil, nil, fmt.Errorf("Unexpected response code: %d, with body: %s", resp.StatusCode, string(contents))
	}
}

func (cli *Client) buildAndExecRequest(ctx context.Context, method, endpointURL string, body []byte) ([]byte, error) {
	_, body, err := cli.buildAndExecRequestWithResponse(ctx, method, endpointURL, body)

	return body, err
}

func hyperHeaderParser(header http.Header) map[string]string {
	var rels = make(map[string]string)

	link := header.Get("Link")
	if len(link) == 0 {
		return rels
	}

	for _, l := range strings.Split(link, ",") {
		l = strings.TrimSpace(l)
		segments := strings.Split(l, ";")

		if len(segments) < 2 {
			continue
		}

		if !strings.HasPrefix(segments[0], "<") || !strings.HasSuffix(segments[0], ">") {
			continue
		}

		url, err := url.Parse(segments[0][1 : len(segments[0])-1])
		if err != nil {
			continue
		}

		link := url.String()

		for _, segment := range segments[1:] {
			switch strings.TrimSpace(segment) {
			case `rel="next"`:
				rels["next"] = link
			case `rel="prev"`:
				rels["prev"] = link
			case `rel="first"`:
				rels["first"] = link
			case `rel="last"`:
				rels["last"] = link
			}
		}
	}

	return rels
}
