package serverscom

import (
	"context"
	"encoding/json"
)

const (
	sshKetCreatePath = "/ssh_keys"
	sshKeyPath       = "/ssh_keys/%s"
)

// SSHKeysService is an interface to interfacing with the SSH Key endpoints
// API documentation: https://developers.servers.com/api-documentation/v1/#tag/SSH-Key
type SSHKeysService interface {
	// Primary collection
	Collection() SSHKeysCollection

	// Generic operations
	Get(ctx context.Context, fingerprint string) (*SSHKey, error)
	Create(ctx context.Context, input SSHKeyCreateInput) (*SSHKey, error)
	Update(ctx context.Context, fingerprint string, input SSHKeyUpdateInput) (*SSHKey, error)
	Delete(ctx context.Context, fingerprint string) error
}

// SSHKeysHandler handles operations around ssh keys
type SSHKeysHandler struct {
	client *Client
}

// Collection builds a new SSHKeysCollection interface
func (s *SSHKeysHandler) Collection() SSHKeysCollection {
	return NewSSHKeysCollection(s.client)
}

// Get ssh key
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/ShowSshKey
func (s *SSHKeysHandler) Get(ctx context.Context, fingerprint string) (*SSHKey, error) {
	url := s.client.buildURL(sshKeyPath, []interface{}{fingerprint}...)

	body, err := s.client.buildAndExecRequest(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	SSHKey := new(SSHKey)

	if err := json.Unmarshal(body, &SSHKey); err != nil {
		return nil, err
	}

	return SSHKey, nil
}

// Create ssh key
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/AddNewSshKey
func (s *SSHKeysHandler) Create(ctx context.Context, input SSHKeyCreateInput) (*SSHKey, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := s.client.buildURL(sshKetCreatePath)

	body, err := s.client.buildAndExecRequest(ctx, "POST", url, payload)

	if err != nil {
		return nil, err
	}

	var SSHKey *SSHKey

	if err := json.Unmarshal(body, &SSHKey); err != nil {
		return nil, err
	}

	return SSHKey, nil
}

// Update ssh key
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/UpdateTheNameOfSshKey
func (s *SSHKeysHandler) Update(ctx context.Context, fingerprint string, input SSHKeyUpdateInput) (*SSHKey, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := s.client.buildURL(sshKeyPath, []interface{}{fingerprint}...)

	body, err := s.client.buildAndExecRequest(ctx, "PUT", url, payload)

	if err != nil {
		return nil, err
	}

	var SSHKey *SSHKey

	if err := json.Unmarshal(body, &SSHKey); err != nil {
		return nil, err
	}

	return SSHKey, nil
}

// Delete ssh key
// Endpoint: https://developers.servers.com/api-documentation/v1/#operation/DeleteSshKey
func (s *SSHKeysHandler) Delete(ctx context.Context, fingerprint string) error {
	url := s.client.buildURL(sshKeyPath, []interface{}{fingerprint}...)

	_, err := s.client.buildAndExecRequest(ctx, "DELETE", url, nil)

	return err
}
