package serverscom

// LocationsService is an interface to interfacing with the Location and Order options endpoints
// API documentation:
// https://developers.servers.com/api-documentation/v1/#tag/Location
// https://developers.servers.com/api-documentation/v1/#tag/Server-Model-Option
// https://developers.servers.com/api-documentation/v1/#tag/Drive-Model-Option
// https://developers.servers.com/api-documentation/v1/#tag/Ram-Option
// https://developers.servers.com/api-documentation/v1/#tag/Operating-System-Option
// https://developers.servers.com/api-documentation/v1/#tag/Uplink-Model-Option
// https://developers.servers.com/api-documentation/v1/#tag/Bandwidth-Option
type LocationsService interface {
	// Primary collection
	Collection() LocationsCollection

	// Generic operations
	ServerModelOptions(LocationID int64) ServerModelOptionsCollection
	RAMOptions(LocationID, ServerModelID int64) RAMOptionsCollection
	OperatingSystemOptions(LocationID, ServerModelID int64) OperatingSystemOptionsCollection
	DriveModelOptions(LocationID, ServerModelID int64) DriveModelOptionsCollection
	UplinkOptions(LocationID, ServerModelID int64) UplinkOptionsCollection
	BandwidthOptions(LocationID, ServerModelID, uplinkID int64) BandwidthOptionsCollection
}

// LocationsHandler handles operations around cloud instances
type LocationsHandler struct {
	client *Client
}

// Collection builds a new LocationsCollection interface
func (resource *LocationsHandler) Collection() LocationsCollection {
	return NewLocationsCollection(resource.client)
}

// ServerModelOptions builds a new ServerModelOptionsCollection interface
func (resource *LocationsHandler) ServerModelOptions(LocationID int64) ServerModelOptionsCollection {
	return NewServerModelOptionsCollection(resource.client, LocationID)
}

// RAMOptions builds a new RAMOptionsCollection interface
func (resource *LocationsHandler) RAMOptions(LocationID, ServerModelID int64) RAMOptionsCollection {
	return NewRAMOptionsCollection(resource.client, LocationID, ServerModelID)
}

// OperatingSystemOptions builds a new OperatingSystemOptionsCollection interface
func (resource *LocationsHandler) OperatingSystemOptions(LocationID, ServerModelID int64) OperatingSystemOptionsCollection {
	return NewOperatingSystemOptionsCollection(resource.client, LocationID, ServerModelID)
}

// DriveModelOptions builds a new DriveModelOptionsCollection interface
func (resource *LocationsHandler) DriveModelOptions(LocationID, ServerModelID int64) DriveModelOptionsCollection {
	return NewDriveModelOptionsCollection(resource.client, LocationID, ServerModelID)
}

// UplinkOptions builds a new UplinkOptionsCollection interface
func (resource *LocationsHandler) UplinkOptions(LocationID, ServerModelID int64) UplinkOptionsCollection {
	return NewUplinkOptionsCollection(resource.client, LocationID, ServerModelID)
}

// BandwidthOptions builds a new BandwidthOptionsCollection interface
func (resource *LocationsHandler) BandwidthOptions(LocationID, ServerModelID, uplinkID int64) BandwidthOptionsCollection {
	return NewBandwidthOptionsCollection(resource.client, LocationID, ServerModelID, uplinkID)
}
