package main

import (
	"io"
	"log"
	"net/http"

	utils "github.com/bifrostcloud/protoc-gen-httpclient/pkg/utils"
	stacktrace "github.com/palantir/stacktrace"

	// example "github.com/bifrostcloud/protoc-gen-httpclient/example"
	basic "github.com/bifrostcloud/protoc-gen-httpclient/pkg/client/basic"
)

// StopResponse -
type StopResponse struct {
	Status string
}

// ServiceWithBasicAuth -
type ServiceWithBasicAuth struct {
	Agent    string
	Endpoint string
	Username string
	Password string
}

func test() {
	c := basic.New()
	var body io.Reader
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodPost, `http://localhost:8080/host`, body)
	if err != nil {
		log.Fatal(stacktrace.Propagate(err, "[POST] request creation failed for DaemonService.PostConfigure with input"))
	}

	request.Header.Set("User-Agent", "Sia-Agent")
	basicauth := utils.BasicAuth("", "somepassword")
	request.Header.Set("Authorization", "Basic "+basicauth)

	_, err = c.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	// log.Printf("%#v", response)
	// // s := NewServiceWithBasicAuth("", "pass")
	// // s.Call()
}

// NewServiceWithBasicAuth - sets up a new client with basic authentication scheme
func NewServiceWithBasicAuth(username, password string) *ServiceWithBasicAuth {
	return &ServiceWithBasicAuth{
		Username: username,
		Password: password,
	}
}

// RequestManipulator -
type RequestManipulator func(*http.Request) error

// Call - returns the result of this rpc call
func (s *ServiceWithBasicAuth) Call(rms []RequestManipulator, opts []basic.Option) {
	// c := basic.New(opts...)
	// target := ""
	// method := "GET"
	// request, err := http.NewRequest(method, s.Endpoint+target, nil)
	// if err != nil {
	// 	// return response, stacktrace.Propagate(err, "[GET] request creation failed")
	// }
	// for _, r := range rms {
	// 	err := r(request)
	// 	if err != nil {
	// 		// return response, stacktrace.Propagate(err, "[GET] request creation failed")
	// 	}
	// }
	// request.Header.Set("User-Agent", s.Agent)
	// basicauth := utils.BasicAuth(s.Username, s.Password)
	// request.Header.Set("Authorization", "Basic "+basicauth)
	// response, err := c.Do(request)
	// if err != nil {
	// 	// return nil, err
	// }
	// result := &StopResponse{}
	// json.NewDecoder(response.Body).Decode(result)
	// result.Status = fmt.Sprintf("%d", response.StatusCode)

}
