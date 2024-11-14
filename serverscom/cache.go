package serverscom

import (
	"context"
	"fmt"
	"sync"

	lru "github.com/hashicorp/golang-lru"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var cache *Cache

func NewCache(cli *scgo.Client) *Cache {
	newLru, err := lru.New(100)
	if err != nil {
		panic(err)
	}

	return &Cache{
		client: cli,
		lru:    newLru,
		ctx:    context.TODO(),
	}
}

type Cache struct {
	client *scgo.Client
	lru    *lru.Cache
	ctx    context.Context

	sync.Mutex
}

func (c *Cache) Locations() ([]scgo.Location, error) {
	c.Lock()
	defer c.Unlock()

	val, ok := c.lru.Get("locations")
	if ok {
		return val.([]scgo.Location), nil
	}

	locations, err := c.client.Locations.Collection().Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add("locations", locations)

	return locations, nil
}

func (c *Cache) ServerModels(locationID int64) ([]scgo.ServerModelOption, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("locations/%d/server_models", locationID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.ServerModelOption), nil
	}

	serverModels, err := c.client.Locations.ServerModelOptions(locationID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, serverModels)

	return serverModels, nil
}

func (c *Cache) DriveModels(locationID int64, serverModelID int64) ([]scgo.DriveModel, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("locations/%d/server_models/%d/drive_models", locationID, serverModelID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.DriveModel), nil
	}

	driveModels, err := c.client.Locations.DriveModelOptions(locationID, serverModelID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, driveModels)

	return driveModels, nil
}

func (c *Cache) OperatingSystems(locationID int64, serverModelID int64) ([]scgo.OperatingSystemOption, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("locations/%d/server_models/%d/operating_systems", locationID, serverModelID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.OperatingSystemOption), nil
	}

	operatingSystems, err := c.client.Locations.OperatingSystemOptions(locationID, serverModelID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, operatingSystems)

	return operatingSystems, nil
}

func (c *Cache) Uplinks(locationID int64, serverModelID int64) ([]scgo.UplinkOption, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("locations/%d/server_models/%d/uplinks", locationID, serverModelID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.UplinkOption), nil
	}

	uplinks, err := c.client.Locations.UplinkOptions(locationID, serverModelID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, uplinks)

	return uplinks, nil
}

func (c *Cache) Bandwidth(locationID int64, serverModelID int64, uplinkModelID int64) ([]scgo.BandwidthOption, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("locations/%d/server_models/%d/uplinks/%d/bandwidth", locationID, serverModelID, uplinkModelID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.BandwidthOption), nil
	}

	bandwidthList, err := c.client.Locations.BandwidthOptions(locationID, serverModelID, uplinkModelID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, bandwidthList)

	return bandwidthList, nil
}

func (c *Cache) LocationGroups() ([]scgo.L2LocationGroup, error) {
	c.Lock()
	defer c.Unlock()

	val, ok := c.lru.Get("l2/location_group")
	if ok {
		return val.([]scgo.L2LocationGroup), nil
	}

	locationGroups, err := c.client.L2Segments.LocationGroups().Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add("l2/location_group", locationGroups)

	return locationGroups, nil
}

func (c *Cache) CloudComputingRegions() ([]scgo.CloudComputingRegion, error) {
	c.Lock()
	defer c.Unlock()

	val, ok := c.lru.Get("cloud_computing/regions")
	if ok {
		return val.([]scgo.CloudComputingRegion), nil
	}

	cloudRegions, err := c.client.CloudComputingRegions.Collection().Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add("cloud_computing/regions", cloudRegions)

	return cloudRegions, nil
}

func (c *Cache) CloudComputingImages(regionID int64) ([]scgo.CloudComputingImage, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("cloud_computing/regions/%d/images", regionID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.CloudComputingImage), nil
	}

	cloudImages, err := c.client.CloudComputingRegions.Images(regionID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, cloudImages)

	return cloudImages, nil
}

func (c *Cache) CloudComputingFlavors(regionID int64) ([]scgo.CloudComputingFlavor, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("cloud_computing/regions/%d/flavors", regionID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.CloudComputingFlavor), nil
	}

	cloudFlavors, err := c.client.CloudComputingRegions.Flavors(regionID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, cloudFlavors)

	return cloudFlavors, nil
}

func (c *Cache) SBMOperatingSystems(locationID int64, sbmFlavorModelID int64) ([]scgo.OperatingSystemOption, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("locations/%d/sbm_flavor_models/%d/operating_systems", locationID, sbmFlavorModelID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.OperatingSystemOption), nil
	}

	operatingSystems, err := c.client.Locations.SBMOperatingSystemOptions(locationID, sbmFlavorModelID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, operatingSystems)

	return operatingSystems, nil
}

func (c *Cache) SBMFlavors(regionID int64) ([]scgo.SBMFlavor, error) {
	c.Lock()
	defer c.Unlock()

	key := fmt.Sprintf("sbm/regions/%d/flavors", regionID)

	val, ok := c.lru.Get(key)
	if ok {
		return val.([]scgo.SBMFlavor), nil
	}

	sbmFlavors, err := c.client.Locations.SBMFlavorOptions(regionID).Collect(c.ctx)
	if err != nil {
		return nil, err
	}

	c.lru.Add(key, sbmFlavors)

	return sbmFlavors, nil
}
