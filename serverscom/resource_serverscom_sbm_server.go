package serverscom

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	scgo "github.com/serverscom/serverscom-go-client/pkg"
)

var (
	serverscomSBMDefaultCreateTimeout = 5 * time.Minute
	serverscomSBMDefaultDeleteTimeout = 1 * time.Minute
)

func resourceServerscomSBM() *schema.Resource {
	return &schema.Resource{
		Read:   resourceServerscomSBMRead,
		Update: resourceServerscomSBMUpdate,
		Delete: resourceServerscomSBMDelete,
		Create: resourceServerscomSBMCreate,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(serverscomSBMDefaultCreateTimeout),
			Delete: schema.DefaultTimeout(serverscomSBMDefaultDeleteTimeout),
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"location": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareStrings,
				ValidateFunc:     validation.NoZeroValues,
			},
			"flavor": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: compareStrings,
				ValidateFunc:     validation.NoZeroValues,
			},
			"operating_system": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareStrings,
				ValidateFunc:     validation.NoZeroValues,
			},
			"ssh_key_fingerprints": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"user_data": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				StateFunc:    HashStringStateFunc(),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return new != "" && old == d.Get("user_data")
				},
			},
			"private_ipv4_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"public_ipv4_network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceServerscomSBMRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)

	ctx := context.TODO()

	sbm, err := client.Hosts.GetSBMServer(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom SBM server (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving SBM server: %s", err)
		}
	}

	d.Set("hostname", sbm.Title)
	d.Set("status", sbm.Status)
	d.Set("operating_system", sbm.ConfigurationDetails.OperatingSystemFullName)
	d.Set("location", sbm.LocationCode)
	d.Set("private_ipv4_address", sbm.PrivateIPv4Address)
	d.Set("public_ipv4_address", sbm.PublicIPv4Address)

	if sbm.Status != "active" {
		return nil
	}

	return nil
}

func resourceServerscomSBMUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceServerscomSBMRead(d, meta)
}

func resourceServerscomSBMDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*scgo.Client)
	ctx := context.TODO()

	sbm, err := client.Hosts.GetSBMServer(ctx, d.Id())
	if err != nil {
		switch err.(type) {
		case *scgo.NotFoundError:
			log.Printf("[WARN] Serverscom SBM server (%s) not found", d.Id())
			d.SetId("")
			return nil
		default:
			return fmt.Errorf("Error retrieving SBM server: %s", err.Error())
		}
	}

	if sbm.Status == "pending" || sbm.Status == "init" {
		_, err = waitForSBMAttribute(ctx, d, "active", []string{"init", "pending"}, "status", meta, schema.TimeoutDelete)
		if err != nil {
			return fmt.Errorf("Error waiting for SBM server (%s) to become ready: %s", d.Id(), err)
		}
	}

	if _, err := client.Hosts.ReleaseSBMServer(ctx, d.Id()); err != nil {
		return err
	}

	return nil
}

func resourceServerscomSBMCreate(d *schema.ResourceData, meta interface{}) error {
	var (
		publicIpv4NetworkId  *string
		privateIpv4NetworkId *string
	)

	input := &SBMServerCreateInput{}

	if id, ok := d.GetOk("public_ipv4_network_id"); ok {
		publicIpv4NetworkIdValue := id.(string)
		publicIpv4NetworkId = &publicIpv4NetworkIdValue
	}

	if id, ok := d.GetOk("private_ipv4_network_id"); ok {
		privateIpv4NetworkIdValue := id.(string)
		privateIpv4NetworkId = &privateIpv4NetworkIdValue
	}

	hostname := d.Get("hostname").(string)
	input.Hosts = []scgo.SBMServerHostInput{
		{
			Hostname:             hostname,
			PublicIPv4NetworkID:  publicIpv4NetworkId,
			PrivateIPv4NetworkID: privateIpv4NetworkId,
		},
	}

	location, err := getLocation(d.Get("location").(string))
	if err != nil {
		return err
	}

	input.LocationID = location.ID

	flavor, err := getSBMFlavor(location.ID, d.Get("flavor").(string))
	if err != nil {
		return err
	}

	input.FlavorModelID = flavor.ID

	if operatingSystemName, ok := d.GetOk("operating_system"); ok {
		operatingSystem, err := getSBMOperatingSystem(location.ID, flavor.ID, operatingSystemName.(string))
		if err != nil {
			return err
		}

		input.OperatingSystemID = &operatingSystem.ID
	}

	if val, ok := d.GetOk("ssh_key_fingerprints"); ok {
		input.SSHKeyFingerprints = expandedStringList(val.([]interface{}))
	}

	if userData, ok := d.GetOk("user_data"); ok {
		userDataValue := userData.(string)
		input.UserData = &userDataValue
	}

	ctx := context.TODO()

	resultChan, err := serverCollector.AddRequest(ctx, "sbm", input)
	if err != nil {
		return err
	}

	// waiting for result from collector
	result := <-resultChan
	if result.Error != nil {
		return result.Error
	}

	if result.Servers.Count() == 0 {
		return fmt.Errorf("Invalid SBM servers count returned by api")
	}

	// find corresponding server by title matching hostname
	id := result.Servers.GetIdByHostname(hostname)
	if id == "" {
		return fmt.Errorf("Can't find the server with title '%s' in api response", hostname)
	}

	d.SetId(id)

	_, err = waitForSBMAttribute(ctx, d, "active", []string{"init", "pending"}, "status", meta, schema.TimeoutCreate)
	if err != nil {
		return fmt.Errorf("Error waiting for SBM server (%s) to become ready: %s", d.Id(), err)
	}

	return nil
}

func waitForSBMAttribute(ctx context.Context, d *schema.ResourceData, target string, pending []string, attribute string, meta interface{}, timeoutKey string) (interface{}, error) {
	log.Printf(
		"[INFO] Waiting for SBM server (%s) to have %s of %s",
		d.Id(), attribute, target,
	)

	stateConf := &retry.StateChangeConf{
		Pending:    pending,
		Target:     []string{target},
		Refresh:    newSBMStateRefreshFunc(d, attribute, meta),
		Timeout:    d.Timeout(timeoutKey),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForStateContext(ctx)
}

func newSBMStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := resourceServerscomSBMRead(d, meta)
		if err != nil {
			return nil, "", err
		}

		// See if we can access our attribute
		if attr, ok := d.GetOk(attribute); ok {
			switch attr.(type) {
			case bool:
				return d, strconv.FormatBool(attr.(bool)), nil
			default:
				return d, attr.(string), nil
			}
		}

		return nil, "", nil
	}
}

func getSBMOperatingSystem(locationID int64, sbmFlavorModelID int64, name string) (*scgo.OperatingSystemOption, error) {
	operatingSystems, err := cache.SBMOperatingSystems(locationID, sbmFlavorModelID)
	if err != nil {
		return nil, err
	}

	for _, os := range operatingSystems {
		fullName := fmt.Sprintf("%s %s %s", os.Name, os.Version, os.Arch)

		if normalizeString(fullName) == normalizeString(name) {
			return &os, nil
		}
	}

	return nil, fmt.Errorf("Can't find operating system by: %s", name)
}

func getSBMFlavor(regionID int64, name string) (*scgo.SBMFlavor, error) {
	flavors, err := cache.SBMFlavors(regionID)
	if err != nil {
		return nil, err
	}

	for _, flavor := range flavors {
		if normalizeString(flavor.Name) == normalizeString(name) {
			return &flavor, nil
		}
	}

	return nil, fmt.Errorf("Can't find SBM flavor by: %s", name)
}

// SBMServerCreateInput implements ServerCreateInput interface
type SBMServerCreateInput struct {
	scgo.SBMServerCreateInput
}

// GetHosts returns hosts from server input
func (s *SBMServerCreateInput) GetHosts() []interface{} {
	hosts := make([]interface{}, len(s.Hosts))
	for i, h := range s.Hosts {
		hosts[i] = h
	}
	return hosts
}

// SetHosts sets hosts for server create input
func (s *SBMServerCreateInput) SetHosts(hosts []interface{}) {
	if hosts == nil {
		s.Hosts = nil
		return
	}
	sbmHosts := make([]scgo.SBMServerHostInput, len(hosts))
	for i, h := range hosts {
		sbmHosts[i] = h.(scgo.SBMServerHostInput)
	}
	s.Hosts = sbmHosts
}

// SBMServerResponse implements ServersResponse interface
type SBMServerResponse struct {
	servers []scgo.SBMServer
}

// GetIdByHostname returns server id if hostname match
func (r *SBMServerResponse) GetIdByHostname(h string) string {
	for _, s := range r.servers {
		if s.Title == h {
			return s.ID
		}
	}
	return ""
}

// Count returns amount of servers in response
func (r *SBMServerResponse) Count() int {
	return len(r.servers)
}
