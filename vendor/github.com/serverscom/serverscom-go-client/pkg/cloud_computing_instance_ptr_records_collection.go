// This code generated automatically

package serverscom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

const (
	cloudInstancePTRsListPath = "/cloud_computing/instances/%s/ptr_records"
)

// CloudComputingInstancePTRRecordsCollection is an interface for interfacing with the collection of PTRRecord
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ReturnsInstancePtrRecords
type CloudComputingInstancePTRRecordsCollection interface {
	IsClean() bool

	HasPreviousPage() bool
	HasNextPage() bool
	HasFirstPage() bool
	HasLastPage() bool

	NextPage(ctx context.Context) ([]PTRRecord, error)
	PreviousPage(ctx context.Context) ([]PTRRecord, error)
	FirstPage(ctx context.Context) ([]PTRRecord, error)
	LastPage(ctx context.Context) ([]PTRRecord, error)

	Collect(ctx context.Context) ([]PTRRecord, error)
	List(ctx context.Context) ([]PTRRecord, error)

	SetPage(page int) CloudComputingInstancePTRRecordsCollection
	SetPerPage(perPage int) CloudComputingInstancePTRRecordsCollection

	Refresh(ctx context.Context) error
}

// CloudComputingInstancePTRRecordsCollectionHandler handles opertations aroud collection
type CloudComputingInstancePTRRecordsCollectionHandler struct {
	client *Client

	cloudInstanceID string

	params map[string]string

	clean bool

	rels       map[string]string
	collection []PTRRecord
}

// NewCloudComputingInstancePTRRecordsCollection produces a new CloudComputingInstancePTRRecordsCollectionHandler and represents this as an interface of CloudComputingInstancePTRRecordsCollection
func NewCloudComputingInstancePTRRecordsCollection(client *Client, cloudInstanceID string) CloudComputingInstancePTRRecordsCollection {
	return &CloudComputingInstancePTRRecordsCollectionHandler{
		client: client,

		cloudInstanceID: cloudInstanceID,

		params:     make(map[string]string),
		rels:       make(map[string]string),
		clean:      true,
		collection: make([]PTRRecord, 0),
	}
}

// IsClean returns a bool value where true is means, this collection not used yet and doesn't contain any state.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) IsClean() bool {
	return col.clean
}

// HasPreviousPage returns a bool value where truth is means collection has a previous page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) HasPreviousPage() bool {
	return col.hasRel("prev")
}

// HasNextPage returns a bool value where truth is means collection has a next page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) HasNextPage() bool {
	return col.hasRel("next")
}

// HasFirstPage returns a bool value where truth is means collection has a first page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) HasFirstPage() bool {
	return col.hasRel("first")
}

// HasLastPage returns a bool value where truth is means collection has a last page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) HasLastPage() bool {
	return col.hasRel("last")
}

// NextPage navigates to the next page returns a []PTRRecord, produces an error, when a
// collection has no next page.
//
// Before using this method please ensure IsClean returns false and HasNextPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) NextPage(ctx context.Context) ([]PTRRecord, error) {
	return col.navigate(ctx, "next")
}

// PreviousPage navigates to the previous page returns a []PTRRecord, produces an error, when a
// collection has no previous page.
//
// Before using this method please ensure IsClean returns false and HasPreviousPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) PreviousPage(ctx context.Context) ([]PTRRecord, error) {
	return col.navigate(ctx, "prev")
}

// FirstPage navigates to the first page returns a []PTRRecord, produces an error, when a
// collection has no first page.
//
// Before using this method please ensure IsClean returns false and HasFirstPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) FirstPage(ctx context.Context) ([]PTRRecord, error) {
	return col.navigate(ctx, "first")
}

// LastPage navigates to the last page returns a []PTRRecord, produces an error, when a
// collection has no last page.
//
// Before using this method please ensure IsClean returns false and HasLastPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) LastPage(ctx context.Context) ([]PTRRecord, error) {
	return col.navigate(ctx, "last")
}

// Collect navigates by pages until the last page is reached will be reached and returns accumulated data between pages.
//
// This method uses NextPage.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) Collect(ctx context.Context) ([]PTRRecord, error) {
	var accumulatedCollectionElements []PTRRecord

	currentCollectionElements, err := col.List(ctx)

	if err != nil {
		return nil, err
	}

	for _, element := range currentCollectionElements {
		accumulatedCollectionElements = append(accumulatedCollectionElements, element)
	}

	for col.HasNextPage() {
		nextCollectionElements, err := col.NextPage(ctx)

		if err != nil {
			return nil, err
		}

		for _, element := range nextCollectionElements {
			accumulatedCollectionElements = append(accumulatedCollectionElements, element)
		}
	}

	return accumulatedCollectionElements, nil
}

// List returns a []BandwidthOption limited by pagination.
//
// This method performs request only once when IsClean returns false, also this request doesn't
// perform request when such methods were called before: NextPage, PreviousPage, LastPage, FirstPage, Refresh, Collect.
//
// In the case when previously called method is Collect, this method returns data from the last page.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) List(ctx context.Context) ([]PTRRecord, error) {
	if col.IsClean() {
		if err := col.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	return col.collection, nil
}

// SetPage sets current page param.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) SetPage(page int) CloudComputingInstancePTRRecordsCollection {
	var currentPage string

	if page > 1 {
		currentPage = strconv.Itoa(page)
	} else {
		currentPage = ""
	}

	col.applyParam("page", currentPage)

	return col
}

// SetPerPage sets current per page param.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) SetPerPage(perPage int) CloudComputingInstancePTRRecordsCollection {
	var currentPerPage string

	if perPage > 0 {
		currentPerPage = strconv.Itoa(perPage)
	} else {
		currentPerPage = ""
	}

	col.applyParam("per_page", currentPerPage)

	return col
}

// Refresh performs the request and then updates accumulated data limited by pagination.
//
// After calling this method accumulated data can be extracted by List method.
func (col *CloudComputingInstancePTRRecordsCollectionHandler) Refresh(ctx context.Context) error {
	if err := col.fireHTTPRequest(ctx); err != nil {
		return err
	}

	return nil
}

func (col *CloudComputingInstancePTRRecordsCollectionHandler) fireHTTPRequest(ctx context.Context) error {
	var accumulatedCollectionElements []PTRRecord

	initialURL := col.client.buildURL(cloudInstancePTRsListPath, col.cloudInstanceID)
	url := col.client.applyParams(
		initialURL,
		col.params,
	)

	response, body, err := col.client.buildAndExecRequestWithResponse(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, &accumulatedCollectionElements); err != nil {
		return err
	}

	col.clean = false
	col.collection = accumulatedCollectionElements
	col.rels = hyperHeaderParser(response.Header)

	return nil
}

func (col *CloudComputingInstancePTRRecordsCollectionHandler) navigate(ctx context.Context, name string) ([]PTRRecord, error) {
	if col.IsClean() {
		if err := col.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	if err := col.applyRel(name); err != nil {
		return nil, err
	}

	if err := col.Refresh(ctx); err != nil {
		return nil, err
	}

	return col.collection, nil
}

func (col *CloudComputingInstancePTRRecordsCollectionHandler) applyParam(name, value string) {
	if value == "" {
		delete(col.params, name)
	} else {
		col.params[name] = value
	}
}

func (col *CloudComputingInstancePTRRecordsCollectionHandler) applyRel(name string) error {
	if !col.hasRel(name) {
		return fmt.Errorf("No rel for: %s", name)
	}

	url, err := url.Parse(col.rels[name])

	if err != nil {
		return err
	}

	col.applyParam("page", url.Query().Get("page"))
	col.applyParam("per_page", url.Query().Get("per_page"))

	return nil
}

func (col *CloudComputingInstancePTRRecordsCollectionHandler) hasRel(name string) bool {
	if _, ok := col.rels[name]; ok {
		return true
	}

	return false
}
