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
	cloudComputingImageListPath = "/cloud_computing/regions/%d/images"
)

// CloudComputingImagesCollection is an interface for interfacing with the collection of CloudComputingImage
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ListCloudImages
type CloudComputingImagesCollection interface {
	IsClean() bool

	HasPreviousPage() bool
	HasNextPage() bool
	HasFirstPage() bool
	HasLastPage() bool

	NextPage(ctx context.Context) ([]CloudComputingImage, error)
	PreviousPage(ctx context.Context) ([]CloudComputingImage, error)
	FirstPage(ctx context.Context) ([]CloudComputingImage, error)
	LastPage(ctx context.Context) ([]CloudComputingImage, error)

	Collect(ctx context.Context) ([]CloudComputingImage, error)
	List(ctx context.Context) ([]CloudComputingImage, error)

	SetPage(page int) CloudComputingImagesCollection
	SetPerPage(perPage int) CloudComputingImagesCollection

	Refresh(ctx context.Context) error
}

// CloudComputingImagesCollectionHandler handles opertations aroud collection
type CloudComputingImagesCollectionHandler struct {
	client *Client

	regionID int64

	params map[string]string

	clean bool

	rels       map[string]string
	collection []CloudComputingImage
}

// NewCloudComputingImagesCollection produces a new CloudComputingImagesCollectionHandler and represents this as an interface of CloudComputingImagesCollection
func NewCloudComputingImagesCollection(client *Client, regionID int64) CloudComputingImagesCollection {
	return &CloudComputingImagesCollectionHandler{
		client: client,

		regionID: regionID,

		params:     make(map[string]string),
		rels:       make(map[string]string),
		clean:      true,
		collection: make([]CloudComputingImage, 0),
	}
}

// IsClean returns a bool value where true is means, this collection not used yet and doesn't contain any state.
func (col *CloudComputingImagesCollectionHandler) IsClean() bool {
	return col.clean
}

// HasPreviousPage returns a bool value where truth is means collection has a previous page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingImagesCollectionHandler) HasPreviousPage() bool {
	return col.hasRel("prev")
}

// HasNextPage returns a bool value where truth is means collection has a next page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingImagesCollectionHandler) HasNextPage() bool {
	return col.hasRel("next")
}

// HasFirstPage returns a bool value where truth is means collection has a first page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingImagesCollectionHandler) HasFirstPage() bool {
	return col.hasRel("first")
}

// HasLastPage returns a bool value where truth is means collection has a last page.
//
// In case when IsClean returns true, this method will return false, which means no request(s)
// were made and collection doesn't have metadata to know about pagination.
//
// First metadata will come with the first called methods such: NextPage, PreviousPage, LastPage, FirstPage, List, Refresh, Collect.
func (col *CloudComputingImagesCollectionHandler) HasLastPage() bool {
	return col.hasRel("last")
}

// NextPage navigates to the next page returns a []CloudComputingImage, produces an error, when a
// collection has no next page.
//
// Before using this method please ensure IsClean returns false and HasNextPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingImagesCollectionHandler) NextPage(ctx context.Context) ([]CloudComputingImage, error) {
	return col.navigate(ctx, "next")
}

// PreviousPage navigates to the previous page returns a []CloudComputingImage, produces an error, when a
// collection has no previous page.
//
// Before using this method please ensure IsClean returns false and HasPreviousPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingImagesCollectionHandler) PreviousPage(ctx context.Context) ([]CloudComputingImage, error) {
	return col.navigate(ctx, "prev")
}

// FirstPage navigates to the first page returns a []CloudComputingImage, produces an error, when a
// collection has no first page.
//
// Before using this method please ensure IsClean returns false and HasFirstPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingImagesCollectionHandler) FirstPage(ctx context.Context) ([]CloudComputingImage, error) {
	return col.navigate(ctx, "first")
}

// LastPage navigates to the last page returns a []CloudComputingImage, produces an error, when a
// collection has no last page.
//
// Before using this method please ensure IsClean returns false and HasLastPage returns true.
// You can force to load pagination metadata by calling Refresh or List methods.
func (col *CloudComputingImagesCollectionHandler) LastPage(ctx context.Context) ([]CloudComputingImage, error) {
	return col.navigate(ctx, "last")
}

// Collect navigates by pages until the last page is reached will be reached and returns accumulated data between pages.
//
// This method uses NextPage.
func (col *CloudComputingImagesCollectionHandler) Collect(ctx context.Context) ([]CloudComputingImage, error) {
	var accumulatedCollectionElements []CloudComputingImage

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
func (col *CloudComputingImagesCollectionHandler) List(ctx context.Context) ([]CloudComputingImage, error) {
	if col.IsClean() {
		if err := col.Refresh(ctx); err != nil {
			return nil, err
		}
	}

	return col.collection, nil
}

// SetPage sets current page param.
func (col *CloudComputingImagesCollectionHandler) SetPage(page int) CloudComputingImagesCollection {
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
func (col *CloudComputingImagesCollectionHandler) SetPerPage(perPage int) CloudComputingImagesCollection {
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
func (col *CloudComputingImagesCollectionHandler) Refresh(ctx context.Context) error {
	if err := col.fireHTTPRequest(ctx); err != nil {
		return err
	}

	return nil
}

func (col *CloudComputingImagesCollectionHandler) fireHTTPRequest(ctx context.Context) error {
	var accumulatedCollectionElements []CloudComputingImage

	initialURL := col.client.buildURL(cloudComputingImageListPath, col.regionID)
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

func (col *CloudComputingImagesCollectionHandler) navigate(ctx context.Context, name string) ([]CloudComputingImage, error) {
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

func (col *CloudComputingImagesCollectionHandler) applyParam(name, value string) {
	if value == "" {
		delete(col.params, name)
	} else {
		col.params[name] = value
	}
}

func (col *CloudComputingImagesCollectionHandler) applyRel(name string) error {
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

func (col *CloudComputingImagesCollectionHandler) hasRel(name string) bool {
	if _, ok := col.rels[name]; ok {
		return true
	}

	return false
}
