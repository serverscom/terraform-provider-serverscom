package serverscom

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"runtime"
	"strings"
)

type fakeRequest struct {
	RequestMethod string
	RequestPath   string
	RequestParams string
	RequestBody   string

	ResponseCode    int
	ResponseHeaders map[string]string
	ResponseBody    []byte
}

type fakeServer struct {
	Server         *httptest.Server
	Requests       []*fakeRequest
	CurrentRequest *fakeRequest
}

func newFakeServer() *fakeServer {
	request := &fakeRequest{}

	return &fakeServer{
		Requests:       []*fakeRequest{request},
		CurrentRequest: request,
	}
}

func (fs *fakeServer) WithRequestMethod(method string) *fakeServer {
	fs.CurrentRequest.RequestMethod = method
	return fs
}

func (fs *fakeServer) WithRequestBody(body string) *fakeServer {
	fs.CurrentRequest.RequestBody = body
	return fs
}

func (fs *fakeServer) WithRequestPath(path string) *fakeServer {
	fs.CurrentRequest.RequestPath = path
	return fs
}

func (fs *fakeServer) WithRequestParams(params string) *fakeServer {
	fs.CurrentRequest.RequestParams = params
	return fs
}

func (fs *fakeServer) WithResponseCode(code int) *fakeServer {
	fs.CurrentRequest.ResponseCode = code
	return fs
}

func (fs *fakeServer) WithResponseHeaders(headers map[string]string) *fakeServer {
	fs.CurrentRequest.ResponseHeaders = headers
	return fs
}

func (fs *fakeServer) Next() *fakeServer {
	request := &fakeRequest{}

	fs.CurrentRequest = request
	fs.Requests = append(fs.Requests, request)

	return fs
}

func (fs *fakeServer) EnsureScenarioWasFinished() error {
	if len(fs.Requests) == 0 {
		return nil
	}

	var unfinished []string

	for _, r := range fs.Requests {
		unfinished = append(unfinished, fmt.Sprintf("%s %s", r.RequestMethod, r.RequestPath))
	}

	return errors.New(strings.Join(unfinished, "\n"))
}

func (fs *fakeServer) Close() {
	if fs.Server == nil {
		return
	}

	fs.Server.Close()
}

func (fs *fakeServer) WithResponseBodyStub(filename string) *fakeServer {
	_, currentFile, _, _ := runtime.Caller(1)
	stubFilePath := path.Join(path.Dir(currentFile), "..", filename)

	stub, err := ioutil.ReadFile(stubFilePath)
	if err != nil {
		panic(fmt.Sprintf("Stub error: %q", err))
	}

	fs.CurrentRequest.ResponseBody = []byte(stub)

	return fs
}

func (fs *fakeServer) WithResponseBodyStubInline(body string) *fakeServer {
	fs.CurrentRequest.ResponseBody = []byte(body)

	return fs
}

func (fs *fakeServer) WithResponseBodyStubFile(filePath string) *fakeServer {
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}

	fs.CurrentRequest.ResponseBody = b

	return fs
}

func (fs *fakeServer) Build() (*fakeServer, *Client) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if len(fs.Requests) == 0 {
			w.WriteHeader(http.StatusTeapot)
			_, err := w.Write([]byte(fmt.Sprintf("End of scenario was reached")))
			crashIfErrorPresent(err)
			return
		}

		currentRequest := fs.Requests[:1][0]
		fs.Requests = fs.Requests[1:len(fs.Requests)]

		finalURL := ""
		if strings.HasPrefix(currentRequest.RequestPath, "/") {
			finalURL = fmt.Sprintf("/v1%s", currentRequest.RequestPath)
		} else {
			panic(`Each path should be started from "/"`)
		}

		if currentRequest.RequestPath == "" || r.URL.Path == finalURL {
			if currentRequest.RequestMethod != "" && currentRequest.RequestMethod != r.Method {
				w.WriteHeader(http.StatusTeapot)
				_, err := w.Write([]byte(fmt.Sprintf(
					"Unexpected request method, expected: %s %s, but got: %s %s",
					currentRequest.RequestMethod,
					finalURL,
					r.Method,
					r.URL.Path,
				)))
				crashIfErrorPresent(err)
				return
			}

			if currentRequest.RequestParams != "" && currentRequest.RequestParams != r.URL.Query().Encode() {
				w.WriteHeader(http.StatusTeapot)
				_, err := w.Write([]byte(fmt.Sprintf("Unexpected query params, expected: %s, but got: %s", currentRequest.RequestParams, r.URL.Query().Encode())))
				crashIfErrorPresent(err)
				return
			}

			if currentRequest.RequestBody != "" {
				b, err := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					http.Error(w, err.Error(), 500)
					return
				}

				if currentRequest.RequestBody != string(b) {
					w.WriteHeader(http.StatusTeapot)
					_, err := w.Write([]byte(fmt.Sprintf("Unexpected request body: %s", string(b))))
					crashIfErrorPresent(err)
					return
				}
			}

			w.Header().Set("Content-Type", "application/json")
			for k, v := range currentRequest.ResponseHeaders {
				w.Header().Set(k, v)
			}

			if currentRequest.ResponseCode != 0 {
				w.WriteHeader(currentRequest.ResponseCode)
			} else {
				w.WriteHeader(http.StatusOK)
			}

			_, err := w.Write(currentRequest.ResponseBody)
			crashIfErrorPresent(err)
		} else {
			w.WriteHeader(http.StatusTeapot)
			_, err := w.Write([]byte(fmt.Sprintf("Unhandled route: %s %s, expected: %s %s", r.Method, r.URL.String(), currentRequest.RequestMethod, finalURL)))
			crashIfErrorPresent(err)
		}
	}))

	fs.Server = ts

	client := NewClientWithEndpoint("testing_token", fmt.Sprintf("%s/v1", ts.URL))

	return fs, client
}

func crashIfErrorPresent(err error) {
	if err != nil {
		panic(err)
	}
}
