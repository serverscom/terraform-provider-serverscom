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
	l2MemberListPath = "/l2_segments/%s/members"
)

// L2MembersCollection is an interface for interfacing with the collection of L2Member
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ListAllL2SegmentMembers
type L2MembersCollection interface {
	IsClean() bool

	HasPreviousPage() bool
	HasNextPage() bool
	HasFirstPage() bool
	HasLastPage() bool

	NextPage(ctx context.Context) ([]L2Member, error)
	PreviousPage(ctx context.Context) ([]L2Member, error)
	FirstPage(ctx context.Context) ([]L2Member, error)
	LastPage(ctx context.Context) ([]L2Member, error)

	Collect(ctx context.Context) ([]L2Member, error)
	List(ctx context.Context) ([]L2Member, error)

	SetPage(page int) L2MembersCollection
	SetPerPage(perPage int) L2MembersCollection

	Refresh(ctx context.Context) error
}

// L2MembersCollectionHandler handles opertations aroud collection
type L2MembersCollectionHandler struct {
	client *Client

	segmentID string

	params map[string]string

	clean bool

	rels       map[string]string
	collection []L2Member
}

// NewL2MembersCollection produces a new L2MembersCollectionHandler and represents this as an interface of L2MembersCollection
func NewL2MembersCollection(client *Client, segmentID string) L2MembersCollection {
	return &L2MembersCollectionHandler{
		client: client,

		segmentID: segmentID,

		params:     make(map[string]string),
		rels:       make(map[string]string),
		clean:      true,
		collection: make([]L2Member, 0),
	}
}

// IsClean returns a bool value where true is means, this collection not used yet and doesn't contain any state.
func (col *L2MembersCollectionHandler) IsClean() bool {
	return col.clean
}

// HasPreviousPage returns a bool value where truth is means collection has a previous page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *L2MembersCollectionHandler) HasPreviousPage() bool {
	return col.hasRel("prev")
}

// HasNextPage returns a bool value where truth is means collection has a next page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *L2MembersCollectionHandler) HasNextPage() bool {
	return col.hasRel("next")
}

// HasFirstPage returns a bool value where truth is means collection has a first page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *L2MembersCollectionHandler) HasFirstPage() bool {
	return col.hasRel("first")
}

// HasLastPage returns a bool value where truth is means collection has a last page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *L2MembersCollectionHandler) HasLastPage() bool {
	return col.hasRel("last")
}

// NextPage navigates to the next page returns a []L2Member, produces an error, when a
// collection has no next page.
//
// Before using this method please ensure IsClean returns false and HasNextPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *L2MembersCollectionHandler) NextPage(ctx context.Context) ([]L2Member, error) {
	return col.navigate(ctx, "next")
}

// PreviousPage navigates to the previous page returns a []L2Member, produces an error, when a
// collection has no previous page.
//
// Before using this method please ensure IsClean returns false and HasPreviousPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *L2MembersCollectionHandler) PreviousPage(ctx context.Context) ([]L2Member, error) {
	return col.navigate(ctx, "prev")
}

// FirstPage navigates to the first page returns a []L2Member, produces an error, when a
// collection has no first page.
//
// Before using this method please ensure IsClean returns false and HasFirstPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *L2MembersCollectionHandler) FirstPage(ctx context.Context) ([]L2Member, error) {
	return col.navigate(ctx, "first")
}

// LastPage navigates to the last page returns a []L2Member, produces an error, when a
// collection has no last page.
//
// Before using this method please ensure IsClean returns false and HasLastPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *L2MembersCollectionHandler) LastPage(ctx context.Context) ([]L2Member, error) {
	return col.navigate(ctx, "last")
}

// Collect navigates by pages until the last page is reached will be reached and returns accumulated data between pages.
//
// This method uses NextPage.
func (col *L2MembersCollectionHandler) Collect(ctx context.Context) ([]L2Member, error) {
	var accumulatedCollectionElements []L2Member

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
func (col *L2MembersCollectionHandler) List(ctx context.Context) ([]L2Member, error) {
	if col.IsClean() {
		if err := col.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	return col.collection, nil
}

// SetPage sets current page param.
func (col *L2MembersCollectionHandler) SetPage(page int) L2MembersCollection {
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
func (col *L2MembersCollectionHandler) SetPerPage(perPage int) L2MembersCollection {
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
func (col *L2MembersCollectionHandler) Refresh(ctx context.Context) error {
	if err := col.fireHTTPRequest(ctx); err != nil {
		return err
	}

	return nil
}

func (col *L2MembersCollectionHandler) fireHTTPRequest(ctx context.Context) error {
	var accumulatedCollectionElements []L2Member

	initialURL := col.client.buildURL(l2MemberListPath, col.segmentID)
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

func (col *L2MembersCollectionHandler) navigate(ctx context.Context, name string) ([]L2Member, error) {
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

func (col *L2MembersCollectionHandler) applyParam(name, value string) {
	if value == "" {
		delete(col.params, name)
	} else {
		col.params[name] = value
	}
}

func (col *L2MembersCollectionHandler) applyRel(name string) error {
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

func (col *L2MembersCollectionHandler) hasRel(name string) bool {
	if _, ok := col.rels[name]; ok {
		return true
	}

	return false
}
