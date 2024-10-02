package serverscom

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	serverCollector *ServerCollector

	// when timer expires the collector triggers the ExecuteRequests method
	serverCollectorTimer = 5 * time.Second
)

// ServerCollector represents server collector abstraction.
// It accept requests from create events for 'serverscom_dedicated_server' resources
// groups it by checksum based on all server fields except 'hosts' and create these servers in one batch request
type ServerCollector struct {
	Client   *scgo.Client
	Requests map[string][]*Request
	Mutex    sync.Mutex
	Timer    *time.Timer
}

// Request represents a request with create server input and result channel
type Request struct {
	Input      *scgo.DedicatedServerCreateInput
	ResultChan chan Result
}

// Result represents the api response and error.
// It's used for sending back the result to the create event through the ResultChan
type Result struct {
	Servers []scgo.DedicatedServer
	Error   error
}

// NewServerCollector creates new ServerCollector
func NewServerCollector(client *scgo.Client) *ServerCollector {
	return &ServerCollector{
		Client:   client,
		Requests: make(map[string][]*Request),
		Timer:    time.NewTimer(5 * time.Second),
	}
}

// AddRequest adds request to server collector
// Each request resets the serverCollectorTimer
func (sc *ServerCollector) AddRequest(ctx context.Context, model string, request *scgo.DedicatedServerCreateInput) (<-chan Result, error) {
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	resultChan := make(chan Result, 1)
	checksum, err := calculateServerChecksum(*request)
	if err != nil {
		return nil, err
	}
	sc.Requests[checksum] = append(sc.Requests[checksum], &Request{Input: request, ResultChan: resultChan})

	if !sc.Timer.Stop() {
		select {
		case <-sc.Timer.C:
		default:
		}
	}
	sc.Timer.Reset(serverCollectorTimer)

	return resultChan, nil
}

// ExecuteRequests triggers when timer expires and runs CreateServersBatch for each requests checksum group
func (sc *ServerCollector) ExecuteRequests() {
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	for checksum, requests := range sc.Requests {
		CreateServersBatch(sc.Client, requests)
		sc.Requests[checksum] = nil
	}
}

// Run runs the collector to listen for requests
func (sc *ServerCollector) Run() {
	go func() {
		for {
			<-sc.Timer.C
			sc.ExecuteRequests()
		}
	}()
}

// CreateServersBatch aggregates hostnames from all requests in one input and creates these servers in one api request
func CreateServersBatch(client *scgo.Client, requests []*Request) {
	if len(requests) == 0 {
		return
	}
	for _, req := range requests {
		defer close(req.ResultChan)
	}

	// combine hosts input
	// for any duplicate hostname return error
	uniqueHostnames := make(map[string]bool)
	createInput := *requests[0].Input
	createInput.Hosts = nil
	for _, req := range requests {
		for _, host := range req.Input.Hosts {
			if _, ok := uniqueHostnames[host.Hostname]; ok {
				result := Result{
					Error: fmt.Errorf("duplicate hostname found: %s", host.Hostname),
				}
				req.ResultChan <- result
				continue
			}
			uniqueHostnames[host.Hostname] = true
			createInput.Hosts = append(createInput.Hosts, host)
		}
	}

	// create servers
	dedicatedServers, err := client.Hosts.CreateDedicatedServers(context.TODO(), createInput)
	result := Result{
		Servers: dedicatedServers,
		Error:   err,
	}

	for _, req := range requests {
		req.ResultChan <- result
	}
}

// calculateServerChecksum generate checksum for server create input excepting the Hosts field
func calculateServerChecksum(input scgo.DedicatedServerCreateInput) (string, error) {
	input.Hosts = []scgo.DedicatedServerHostInput{}

	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
