package serverscom

// CloudComputingRegionsService is an interface to interfacing with the cloud computing regions endpoints
// API documentation:
// https://developers.servers.com/api-documentation/v1/#tag/Cloud-Region
type CloudComputingRegionsService interface {
	// Primary collection
	Collection() CloudComputingRegionsCollection

	// Additional collections
	Images(regionID int64) CloudComputingImagesCollection
	Flavors(regionID int64) CloudComputingFlavorsCollection
}

// CloudComputingRegionsHandler handles operations around cloud computing regions
type CloudComputingRegionsHandler struct {
	client *Client
}

// Collection builds a new CloudComputingInstancesCollection interface
func (cr *CloudComputingRegionsHandler) Collection() CloudComputingRegionsCollection {
	return NewCloudComputingRegionsCollection(cr.client)
}

// Images builds a new CloudComputingImagesCollection interface
func (cr *CloudComputingRegionsHandler) Images(regionID int64) CloudComputingImagesCollection {
	return NewCloudComputingImagesCollection(cr.client, regionID)
}

// Flavors builds a new CloudComputingFlavorsCollection interface
func (cr *CloudComputingRegionsHandler) Flavors(regionID int64) CloudComputingFlavorsCollection {
	return NewCloudComputingFlavorsCollection(cr.client, regionID)
}
