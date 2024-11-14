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
	Requests map[string]map[string][]*Request
	Mutex    sync.Mutex
	Timer    *time.Timer
}

// ServerCreateInput represents server (sbm, dedicated, ...) create input interface
type ServerCreateInput interface {
	GetHosts() []interface{}
	SetHosts([]interface{})
}

// ServersResponse represents server (sbm, dedicated, ...) response interface
type ServersResponse interface {
	GetIdByHostname(hoostname string) string
	Count() int
}

// Request represents a request with create server input and result channel
type Request struct {
	Input      ServerCreateInput
	ResultChan chan Result
}

// Result represents the api response and error.
// It's used for sending back the result to the create event through the ResultChan
type Result struct {
	Servers ServersResponse
	Error   error
}

// NewServerCollector creates new ServerCollector
func NewServerCollector(client *scgo.Client) *ServerCollector {
	return &ServerCollector{
		Client:   client,
		Requests: make(map[string]map[string][]*Request),
		Timer:    time.NewTimer(serverCollectorTimer),
	}
}

// AddRequest adds request to server collector
// Each request resets the serverCollectorTimer
func (sc *ServerCollector) AddRequest(ctx context.Context, resourceType string, input ServerCreateInput) (<-chan Result, error) {
	sc.Mutex.Lock()
	defer sc.Mutex.Unlock()

	resultChan := make(chan Result, 1)
	checksum, err := calculateServerChecksum(input)
	if err != nil {
		return nil, err
	}
	if sc.Requests[resourceType] == nil {
		sc.Requests[resourceType] = make(map[string][]*Request)
	}
	sc.Requests[resourceType][checksum] = append(
		sc.Requests[resourceType][checksum],
		&Request{
			Input:      input,
			ResultChan: resultChan,
		},
	)

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

	ctx := context.Background()
	for resourceType, checksums := range sc.Requests {
		for checksum, requests := range checksums {
			CreateServersBatch(ctx, sc.Client, requests)
			sc.Requests[resourceType][checksum] = nil
		}
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
func CreateServersBatch(ctx context.Context, client *scgo.Client, requests []*Request) {
	if len(requests) == 0 {
		return
	}
	for _, req := range requests {
		defer close(req.ResultChan)
	}

	// combine hosts input
	// for any duplicate hostname return error
	uniqueHostnames := make(map[string]struct{})
	var combinedHosts []interface{}
	result := Result{}
	for _, req := range requests {
		for _, host := range req.Input.GetHosts() {
			hostname, err := getHostHostname(host)
			if err != nil {
				result.Error = err
				req.ResultChan <- result
				continue
			}
			if _, ok := uniqueHostnames[hostname]; ok {
				result.Error = fmt.Errorf("duplicate hostname found: %s", hostname)
				req.ResultChan <- result
				continue
			}
			uniqueHostnames[hostname] = struct{}{}
			combinedHosts = append(combinedHosts, host)
		}
	}

	firstReq := requests[0]
	firstReq.Input.SetHosts(combinedHosts)

	switch v := firstReq.Input.(type) {
	case *DedicatedServerCreateInput:
		resp, err := client.Hosts.CreateDedicatedServers(ctx, v.DedicatedServerCreateInput)
		servers := &DedicatedServerResponse{servers: resp}
		result.Servers = servers
		result.Error = err
	case *SBMServerCreateInput:
		resp, err := client.Hosts.CreateSBMServers(ctx, v.SBMServerCreateInput)
		servers := &SBMServerResponse{servers: resp}
		result.Servers = servers
		result.Error = err
	default:
		result.Error = fmt.Errorf("Unknown resource type: %T\n", v)
	}

	for _, req := range requests {
		req.ResultChan <- result
	}
}

// calculateServerChecksum generate checksum for server create input excepting the Hosts field
func calculateServerChecksum(input ServerCreateInput) (string, error) {
	originalHosts := input.GetHosts()
	input.SetHosts(nil)

	data, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	input.SetHosts(originalHosts)

	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}

// getHostHostname returns host hostname from host input
func getHostHostname(hostInput interface{}) (string, error) {
	switch v := hostInput.(type) {
	case scgo.DedicatedServerHostInput:
		return v.Hostname, nil
	case scgo.SBMServerHostInput:
		return v.Hostname, nil
	default:
		return "", fmt.Errorf("Unknown host input type: %T\n", v)
	}
}
