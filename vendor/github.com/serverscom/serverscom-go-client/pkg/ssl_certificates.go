package serverscom

import (
	"context"
	"encoding/json"
)

const (
	sslCreatificatedCreatePath = "/ssl_certificates/custom"
	sslCertificatePath         = "/ssl_certificates/custom/%s"
)

// SSLCertificatesService is an interface to interfacing with the SSL Certificate endpoints
// API documentation: https://developers.servers.com/api-documentation/v1/#tag/SSL-Certificate
type SSLCertificatesService interface {
	// Primary collection
	Collection() SSLCertificatesCollection

	// Generic operations
	CreateCustom(ctx context.Context, input SSLCertificateCreateCustomInput) (*SSLCertificateCustom, error)
	GetCustom(ctx context.Context, id string) (*SSLCertificateCustom, error)
}

// SSLCertificatesHandler handles operations around ssl certificates
type SSLCertificatesHandler struct {
	client *Client
}

// Collection builds a new SSLCertificatesCollection interface
func (c *SSLCertificatesHandler) Collection() SSLCertificatesCollection {
	return NewSSLCertificatesCollection(c.client)
}

// CreateCustom creates a custom ssl certificate
func (c *SSLCertificatesHandler) CreateCustom(ctx context.Context, input SSLCertificateCreateCustomInput) (*SSLCertificateCustom, error) {
	payload, err := json.Marshal(input)

	if err != nil {
		return nil, err
	}

	url := c.client.buildURL(sslCreatificatedCreatePath)

	body, err := c.client.buildAndExecRequest(ctx, "POST", url, payload)

	if err != nil {
		return nil, err
	}

	SSLCertificateCustom := new(SSLCertificateCustom)

	if err := json.Unmarshal(body, &SSLCertificateCustom); err != nil {
		return nil, err
	}

	return SSLCertificateCustom, nil
}

// GetCustom returns a custom ssl certificate
func (c *SSLCertificatesHandler) GetCustom(ctx context.Context, id string) (*SSLCertificateCustom, error) {
	url := c.client.buildURL(sslCertificatePath, id)

	body, err := c.client.buildAndExecRequest(ctx, "GET", url, nil)

	if err != nil {
		return nil, err
	}

	SSLCertificate := new(SSLCertificateCustom)

	if err := json.Unmarshal(body, &SSLCertificate); err != nil {
		return nil, err
	}

	return SSLCertificate, nil
}
