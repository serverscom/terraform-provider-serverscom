package serverscom

import (
	"context"
	"encoding/json"
)

const (
	l2SegmentPath               = "/l2_segments/%s"
	l2SegmentCreatePath         = "/l2_segments"
	l2SegmentUpdatePath         = "/l2_segments/%s"
	l2SegemntDeletePath         = "/l2_segments/%s"
	l2SegmentChangeNetworksPath = "/l2_segments/%s/networks"
)

// L2SegmentsService is an interface for interfacing with Host, Dedicated Server endpoints
// API documentation:
// https://developers.servers.com/api-documentation/v1/#tag/L2-Segment
type L2SegmentsService interface {
	// Primary collection
	Collection() L2SegmentsCollection

	// Extra collections
	LocationGroups() L2LocationGroupsCollection

	// Generic operations
	Get(ctx context.Context, segmentID string) (*L2Segment, error)
	Create(ctx context.Context, input L2SegmentCreateInput) (*L2Segment, error)
	Update(ctx context.Context, segmentID string, input L2SegmentUpdateInput) (*L2Segment, error)
	Delete(ctx context.Context, segmentID string) error

	// Additional operations
	ChangeNetworks(ctx context.Context, segmentID string, input L2SegmentChangeNetworksInput) (*L2Segment, error)

	// Additional collections
	Members(segmentID string) L2MembersCollection
	Networks(segmentID string) L2NetworksCollection
}

// L2SegmentsHandler handles  operatings around l2 segments
type L2SegmentsHandler struct {
	client *Client
}

// Collection builds a new L2SegmentsCollection interface
func (l2 *L2SegmentsHandler) Collection() L2SegmentsCollection {
	return NewL2SegmentsCollection(l2.client)
}

// Get l2 segment
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/RetrieveAnExistingL2Segment
func (l2 *L2SegmentsHandler) Get(ctx context.Context, segmentID string) (*L2Segment, error) {
	url := l2.client.buildURL(l2SegmentPath, []interface{}{segmentID}...)

	body, err := l2.client.buildAndExecRequest(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	l2Segment := new(L2Segment)

	if err := json.Unmarshal(body, &l2Segment); err != nil {
		return nil, err
	}

	return l2Segment, nil
}

// Create l2 segment
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/CreateANewL2Segment
func (l2 *L2SegmentsHandler) Create(ctx context.Context, input L2SegmentCreateInput) (*L2Segment, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := l2.client.buildURL(l2SegmentCreatePath)

	body, err := l2.client.buildAndExecRequest(ctx, "POST", url, payload)

	if err != nil {
		return nil, err
	}

	l2Segment := new(L2Segment)

	if err := json.Unmarshal(body, &l2Segment); err != nil {
		return nil, err
	}

	return l2Segment, nil
}

// Update l2 segment
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/UpdateAnExistingL2Segment
func (l2 *L2SegmentsHandler) Update(ctx context.Context, segmentID string, input L2SegmentUpdateInput) (*L2Segment, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := l2.client.buildURL(l2SegmentUpdatePath, []interface{}{segmentID}...)

	body, err := l2.client.buildAndExecRequest(ctx, "PUT", url, payload)

	if err != nil {
		return nil, err
	}

	l2Segment := new(L2Segment)

	if err := json.Unmarshal(body, &l2Segment); err != nil {
		return nil, err
	}

	return l2Segment, nil
}

// Delete l2 segment
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/DeleteAnExistingL2Segment
func (l2 *L2SegmentsHandler) Delete(ctx context.Context, segmentID string) error {
	url := l2.client.buildURL(l2SegemntDeletePath, []interface{}{segmentID}...)

	_, err := l2.client.buildAndExecRequest(ctx, "DELETE", url, nil)

	return err
}

// LocationGroups builds a new L2LocationGroupsCollection interface
func (l2 *L2SegmentsHandler) LocationGroups() L2LocationGroupsCollection {
	return NewL2LocationGroupsCollection(l2.client)
}

// Members builds a new L2MembersCollection interface
func (l2 *L2SegmentsHandler) Members(segmentID string) L2MembersCollection {
	return NewL2MembersCollection(l2.client, segmentID)
}

// Networks builds a new L2NetworksCollection interface
func (l2 *L2SegmentsHandler) Networks(segmentID string) L2NetworksCollection {
	return NewL2NetworksCollection(l2.client, segmentID)
}

// ChangeNetworks changes networks set
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/UpdateAnExistingL2SegmentNetworks
func (l2 *L2SegmentsHandler) ChangeNetworks(ctx context.Context, segmentID string, input L2SegmentChangeNetworksInput) (*L2Segment, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := l2.client.buildURL(l2SegmentChangeNetworksPath, []interface{}{segmentID}...)

	body, err := l2.client.buildAndExecRequest(ctx, "PUT", url, payload)

	if err != nil {
		return nil, err
	}

	l2Segment := new(L2Segment)

	if err := json.Unmarshal(body, &l2Segment); err != nil {
		return nil, err
	}

	return l2Segment, nil
}
