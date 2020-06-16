package serverscom

import (
	"context"
	"encoding/json"
)

const (
	cloudInstanceCreatePath          = "/cloud_computing/instances"
	cloudInstancePath                = "/cloud_computing/instances/%s"
	cloudInstanceUpdatePath          = "/cloud_computing/instances/%s"
	cloudInstanceDeletePath          = "/cloud_computing/instances/%s"
	cloudInstanceReinstallPath       = "/cloud_computing/instances/%s/reinstall"
	cloudInstanceRescuePath          = "/cloud_computing/instances/%s/rescue"
	cloudInstanceUnrescuePath        = "/cloud_computing/instances/%s/unrescue"
	cloudInstanceUpgradePath         = "/cloud_computing/instances/%s/upgrade"
	cloudInstanceRevertUpgradePath   = "/cloud_computing/instances/%s/revert_upgrade"
	cloudInstanceApproveUpgradePath  = "/cloud_computing/instances/%s/approve_upgrade"
	cloudInstancePowerOnPath         = "/cloud_computing/instances/%s/switch_power_on"
	cloudInstancePowerOffPath        = "/cloud_computing/instances/%s/switch_power_off"
	cloudInstanceCreatePTRRecordPath = "/cloud_computing/instances/%s/ptr_records"
	cloudInstanceDeletePTRRecordPath = "/cloud_computing/instances/%s/ptr_records/%s"
)

// CloudComputingInstancesService is an interface to interfacing with the Cloud Instance endpoints
// API documentation: https://developers.servers.com/api-documentation/v1/#tag/Cloud-Instance
type CloudComputingInstancesService interface {
	// Primary collection
	Collection() CloudComputingInstancesCollection

	// Generic operations
	Get(ctx context.Context, id string) (*CloudComputingInstance, error)
	Create(ctx context.Context, input CloudComputingInstanceCreateInput) (*CloudComputingInstance, error)
	Update(ctx context.Context, id string, input CloudComputingInstanceUpdateInput) (*CloudComputingInstance, error)
	Delete(ctx context.Context, id string) error

	// Additional operations
	Reinstall(ctx context.Context, id string, input CloudComputingInstanceReinstallInput) (*CloudComputingInstance, error)
	Rescue(ctx context.Context, id string) (*CloudComputingInstance, error)
	Unrescue(ctx context.Context, id string) (*CloudComputingInstance, error)
	Upgrade(ctx context.Context, id string, input CloudComputingInstanceUpgradeInput) (*CloudComputingInstance, error)
	RevertUpgrade(ctx context.Context, id string) (*CloudComputingInstance, error)
	ApproveUpgrade(ctx context.Context, id string) (*CloudComputingInstance, error)
	PowerOn(ctx context.Context, id string) (*CloudComputingInstance, error)
	PowerOff(ctx context.Context, id string) (*CloudComputingInstance, error)
	CreatePTRRecord(ctx context.Context, cloudInstanceID string, input PTRRecordCreateInput) (*PTRRecord, error)
	DeletePTRRecord(ctx context.Context, cloudInstanceID string, ptrRecordID string) error

	// Additional collections
	PTRRecords(id string) CloudComputingInstancePTRRecordsCollection
}

// CloudComputingInstancesHandler handles operations around cloud instances
type CloudComputingInstancesHandler struct {
	client *Client
}

// Collection builds a new CloudComputingInstancesCollection interface
func (ci *CloudComputingInstancesHandler) Collection() CloudComputingInstancesCollection {
	return NewCloudComputingInstancesCollection(ci.client)
}

// Get cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ShowCloudComputingInstance
func (ci *CloudComputingInstancesHandler) Get(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstancePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	cloudInstance := new(CloudComputingInstance)

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// Create cloud instace
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/CreateANewCloudComputingInstance
func (ci *CloudComputingInstancesHandler) Create(ctx context.Context, input CloudComputingInstanceCreateInput) (*CloudComputingInstance, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := ci.client.buildURL(cloudInstanceCreatePath)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, payload)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// Update cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/UpdateCloudComputingInstance
func (ci *CloudComputingInstancesHandler) Update(ctx context.Context, id string, input CloudComputingInstanceUpdateInput) (*CloudComputingInstance, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := ci.client.buildURL(cloudInstanceUpdatePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "PUT", url, payload)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// Delete cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/DeleteInstance
func (ci *CloudComputingInstancesHandler) Delete(ctx context.Context, id string) error {
	url := ci.client.buildURL(cloudInstanceDeletePath, []interface{}{id}...)

	_, err := ci.client.buildAndExecRequest(ctx, "DELETE", url, nil)

	return err
}

// Reinstall cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ReinstallInstanceWithImage
func (ci *CloudComputingInstancesHandler) Reinstall(ctx context.Context, id string, input CloudComputingInstanceReinstallInput) (*CloudComputingInstance, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := ci.client.buildURL(cloudInstanceReinstallPath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, payload)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// Rescue cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/MoveInstanceToRescueState
func (ci *CloudComputingInstancesHandler) Rescue(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstanceRescuePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// Unrescue cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ExitFromRescueState
func (ci *CloudComputingInstancesHandler) Unrescue(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstanceUnrescuePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// Upgrade cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/UpgradeInstance
func (ci *CloudComputingInstancesHandler) Upgrade(ctx context.Context, id string, input CloudComputingInstanceUpgradeInput) (*CloudComputingInstance, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := ci.client.buildURL(cloudInstanceUpgradePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, payload)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// RevertUpgrade cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/RevertInstanceUpgrade
func (ci *CloudComputingInstancesHandler) RevertUpgrade(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstanceRevertUpgradePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// ApproveUpgrade cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ApproveInstanceUpgrade
func (ci *CloudComputingInstancesHandler) ApproveUpgrade(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstanceApproveUpgradePath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// PowerOn cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/SwitchPowerOn
func (ci *CloudComputingInstancesHandler) PowerOn(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstancePowerOnPath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// PowerOff cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/SwitchPowerOff
func (ci *CloudComputingInstancesHandler) PowerOff(ctx context.Context, id string) (*CloudComputingInstance, error) {
	url := ci.client.buildURL(cloudInstancePowerOffPath, []interface{}{id}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	var cloudInstance *CloudComputingInstance

	if err := json.Unmarshal(body, &cloudInstance); err != nil {
		return nil, err
	}

	return cloudInstance, nil
}

// PTRRecords builds a new CloudComputingInstancePTRRecordsCollection interface
func (ci *CloudComputingInstancesHandler) PTRRecords(id string) CloudComputingInstancePTRRecordsCollection {
	return NewCloudComputingInstancePTRRecordsCollection(ci.client, id)
}

// CreatePTRRecord creates ptr record for the cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/CreatePtrForInstance
func (ci *CloudComputingInstancesHandler) CreatePTRRecord(ctx context.Context, cloudInstanceID string, input PTRRecordCreateInput) (*PTRRecord, error) {
	url := ci.client.buildURL(cloudInstanceCreatePTRRecordPath, []interface{}{cloudInstanceID}...)

	body, err := ci.client.buildAndExecRequest(ctx, "POST", url, nil)

	if err != nil {
		return nil, err
	}

	ptrRecord := new(PTRRecord)

	if err := json.Unmarshal(body, &ptrRecord); err != nil {
		return nil, err
	}

	return ptrRecord, nil
}

// DeletePTRRecord deleted ptr record for the cloud instance
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/DetetePtrForInstance
func (ci *CloudComputingInstancesHandler) DeletePTRRecord(ctx context.Context, cloudInstanceID string, ptrRecordID string) error {
	url := ci.client.buildURL(cloudInstanceDeletePTRRecordPath, []interface{}{cloudInstanceID, ptrRecordID}...)

	_, err := ci.client.buildAndExecRequest(ctx, "DELETE", url, nil)

	return err
}
