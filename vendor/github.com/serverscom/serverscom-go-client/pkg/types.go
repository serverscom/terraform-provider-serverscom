package serverscom

import (
	"time"
)

// Location represents location
type Location struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// SSLCertificate represents ssl certificate
type SSLCertificate struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Sha1Fingerprint string     `json:"sha1_fingerprint"`
	Expires         *time.Time `json:"expires_at"`
	Created         time.Time  `json:"created_at"`
	Updated         time.Time  `json:"updated_at"`
}

// SSLCertificateCustom represents custom ssl certificate
type SSLCertificateCustom struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Sha1Fingerprint string     `json:"sha1_fingerprint"`
	Expires         *time.Time `json:"expires_at"`
	Created         time.Time  `json:"created_at"`
	Updated         time.Time  `json:"updated_at"`
}

// SSLCertificateCreateCustomInput represents custom ssl certificate create input
type SSLCertificateCreateCustomInput struct {
	Name       string `json:"name"`
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	ChainKey   string `json:"chain_key"`
}

// Host represents host
type Host struct {
	ID                 string     `json:"id"`
	Title              string     `json:"title"`
	LocationID         int64      `json:"location_id"`
	LocationCode       string     `json:"location_code"`
	Status             string     `json:"status"`
	Configuration      string     `json:"configuration"`
	PrivateIPv4Address *string    `json:"private_ipv4_address"`
	PublicIPv4Address  *string    `json:"public_ipv4_address"`
	ScheduledRelease   *time.Time `json:"scheduled_release_at"`
	Created            time.Time  `json:"created_at"`
	Updated            time.Time  `json:"updated_at"`
}

// ConfigurationDetails represents host configuration details
type ConfigurationDetails struct {
	RAMSize                 int     `json:"ram_size"`
	ServerModelID           *int64  `json:"server_model_id"`
	ServerModelName         *string `json:"server_model_name"`
	PublicUplinkID          *int64  `json:"public_uplink_id"`
	PublicUplinkName        *string `json:"public_uplink_name"`
	PrivateUplinkID         *int64  `json:"private_uplink_id"`
	PrivateUplinkName       *string `json:"private_uplink_name"`
	BandwidthID             *int64  `json:"bandwidth_id"`
	BandwidthName           *string `json:"bandwidth_name"`
	OperatingSystemID       *int64  `json:"operating_system_id"`
	OperatingSystemFullName *string `json:"operating_system_full_name"`
}

// DedicatedServer represents dedicated server
type DedicatedServer struct {
	ID                   string               `json:"id"`
	Title                string               `json:"title"`
	LocationID           int64                `json:"location_id"`
	LocationCode         string               `json:"location_code"`
	Status               string               `json:"status"`
	Configuration        string               `json:"configuration"`
	PrivateIPv4Address   *string              `json:"private_ipv4_address"`
	PublicIPv4Address    *string              `json:"public_ipv4_address"`
	ScheduledRelease     *time.Time           `json:"scheduled_release_at"`
	ConfigurationDetails ConfigurationDetails `json:"configuration_details"`
	Created              time.Time            `json:"created_at"`
	Updated              time.Time            `json:"updated_at"`
}

// DedicatedServerLayoutPartitionInput represents partition for DedicatedServerLayoutInput
type DedicatedServerLayoutPartitionInput struct {
	Target string  `json:"target"`
	Size   int     `json:"size"`
	Fs     *string `json:"fs,omitempty"`
	Fill   bool    `json:"fill,omitempty"`
}

// DedicatedServerLayoutInput represents layout for DedicatedServerDrivesInput
type DedicatedServerLayoutInput struct {
	SlotPositions []int                                 `json:"slot_positions"`
	Raid          *int                                  `json:"raid,omitempty"`
	Partitions    []DedicatedServerLayoutPartitionInput `json:"partitions"`
}

// DedicatedServerSlotInput represents slot for DedicatedServerDrivesInput
type DedicatedServerSlotInput struct {
	Position     int    `json:"position"`
	DriveModelID *int64 `json:"drive_model_id,omitempty"`
}

// DedicatedServerDrivesInput represents drives for DedicatedServerCreateInput
type DedicatedServerDrivesInput struct {
	Slots  []DedicatedServerSlotInput   `json:"slots"`
	Layout []DedicatedServerLayoutInput `json:"layout"`
}

// DedicatedServerPublicUplinkInput represents public uplink for DedicatedServerUplinkModelsInput
type DedicatedServerPublicUplinkInput struct {
	ID               int64 `json:"id"`
	BandwidthModelID int64 `json:"bandwidth_model_id"`
}

// DedicatedServerPrivateUplinkInput represents private uplink for DedicatedServerUplinkModelsInput
type DedicatedServerPrivateUplinkInput struct {
	ID int64 `json:"id"`
}

// DedicatedServerUplinkModelsInput represents uplinks for DedicatedServerCreateInput
type DedicatedServerUplinkModelsInput struct {
	Public  *DedicatedServerPublicUplinkInput `json:"public,omitempty"`
	Private DedicatedServerPrivateUplinkInput `json:"private"`
}

// DedicatedServerHostInput represents hosts for DedicatedServerCreateInput
type DedicatedServerHostInput struct {
	Hostname string `json:"hostname"`
}

// DedicatedServerCreateInput represents dedicated server create input, example:
//
//  driveModelID := int64(1)
//  osUbuntuServerID := int64(1)
//  rootFilesystem := "ext4"
//  raidLevel := 0
//
//  input := DedicatedServerCreateInput{
//    ServerModelID: int64(1),
//    LocationID: int64(1),
//    RAMSize: 32,
//    UplinkModels: DedicatedServerUplinkModelInput{
//      PublicUplink &DedicatedServerPublicUplinkInput{ID: int64(1), BandwidthModelID: int64(1)},
//      PrivateUplink: DedicatedServerPrivateUplinkInput{ID: int64(2)},
//    },
//    Drives: DedicatedServerDrivesInput{
//      Slots: []DedicatedServerSlotInput{
//        DedicatedServerSlotInput{Position: 0, DriveModelID: &driveModelID},
//        DedicatedServerSlotInput{Position: 1, DriveModelID: &driveModelID},
//      },
//      Layout: []DedicatedServerLayoutInput{
//        DedicatedServerLayoutInput{
//          SlotPositions: []int{0, 1},
//          Riad:          &raidLevel,
//          Partitions:    []DedicatedServerLayoutPartitionInput{
//            DedicatedServerLayoutPartitionInput{Target: "swap", Size: 4096, Fill: false},
//            DedicatedServerLayoutPartitionInput{Target: "/", Fs: &rootFilesystem, Size: 100000, Fill: true},
//          },
//        },
//      },
//    },
//    IPv6: true,
//    OperatingSystemID: &osUbuntuServerID,
//    SSHKeyFingerprints: []string{
//      "48:81:0c:43:99:12:71:5e:ba:fd:e7:2f:20:d7:95:e8"
//    },
//    Hosts: []DedicatedServerHostInput{
//      Hostname: "example-host",
//    },
//  }
type DedicatedServerCreateInput struct {
	ServerModelID      int64                            `json:"server_model_id"`
	LocationID         int64                            `json:"location_id"`
	RAMSize            int                              `json:"ram_size"`
	UplinkModels       DedicatedServerUplinkModelsInput `json:"uplink_models"`
	Drives             DedicatedServerDrivesInput       `json:"drives"`
	Features           []string                         `json:"features,omitempty"`
	IPv6               bool                             `json:"ipv6"`
	Hosts              []DedicatedServerHostInput       `json:"hosts"`
	OperatingSystemID  *int64                           `json:"operating_system_id"`
	SSHKeyFingerprints []string                         `json:"ssh_key_fingerprints,omitempty"`
}

// ServerModelOption represents server model option
type ServerModelOption struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	RAM  int    `json:"ram"`
}

// RAMOption represents ram option
type RAMOption struct {
	RAM  int    `json:"ram"`
	Type string `json:"type"`
}

// OperatingSystemOption represents operating system option
type OperatingSystemOption struct {
	ID          int64    `json:"id"`
	FullName    string   `json:"full_name"`
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Arch        string   `json:"arch"`
	Filesystems []string `json:"filesystems"`
}

// UplinkOption represents uplink option
type UplinkOption struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Speed      int    `json:"speed"`
	Redundancy bool   `json:"redundancy"`
}

// BandwidthOption represents bandwidth option
type BandwidthOption struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Commit *int64 `json:"commit,omitempty"`
}

// DriveModel represents drive model
type DriveModel struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Capacity   int    `json:"capacity"`
	Interface  string `json:"interface"`
	FormFactor string `json:"form_factor"`
	MediaType  string `json:"media_type"`
}

// SSHKey represents ssh key
type SSHKey struct {
	Name        string    `json:"name"`
	Fingerprint string    `json:"fingerprint"`
	Created     time.Time `json:"created_at"`
	Updated     time.Time `json:"updated_at"`
}

// SSHKeyCreateInput represents ssh key create input
type SSHKeyCreateInput struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

// SSHKeyUpdateInput represents ssh key update input
type SSHKeyUpdateInput struct {
	Name string `json:"name"`
}

// CloudComputingInstance represents cloud instance
type CloudComputingInstance struct {
	Name               string    `json:"name"`
	ID                 string    `json:"id"`
	OpenstackUUID      string    `json:"openstack_uuid"`
	Status             string    `json:"status"`
	FlavorID           string    `json:"flavor_id"`
	FlavorName         string    `json:"flavor_name"`
	ImageID            string    `json:"image_id"`
	ImageName          *string   `json:"image_name"`
	PublicIPv4Address  *string   `json:"public_ipv4_address"`
	PrivateIPv4Address *string   `json:"private_ipv4_address"`
	PublicIPv6Address  *string   `json:"public_ipv6_address"`
	GPNEnabled         bool      `json:"gpn_enabled"`
	IPv6Enabled        bool      `json:"ipv6_enabled"`
	Created            time.Time `json:"created_at"`
	Updated            time.Time `json:"updated_at"`
}

// CloudComputingInstanceCreateInput represents cloud instance create input
type CloudComputingInstanceCreateInput struct {
	Name              string  `json:"name"`
	RegionID          int64   `json:"region_id"`
	FlavorID          string  `json:"flavor_id"`
	ImageID           string  `json:"image_id"`
	GPNEnabled        *bool   `json:"gpn_enabled,omitempty"`
	IPv6Enabled       *bool   `json:"ipv6_enabled,omitempty"`
	SSHKeyFingerprint *string `json:"ssh_key_fingerprint,omitempty"`
	BackupCopies      *int    `json:"backup_copies,omitempty"`
}

// CloudComputingInstanceUpdateInput represents cloud instance update input
type CloudComputingInstanceUpdateInput struct {
	Name         *string `json:"name,omitempty"`
	BackupCopies *int    `json:"backup_copies,omitempty"`
	GPNEnabled   *bool   `json:"gpn_enabled,omitempty"`
	IPv6Enabled  *bool   `json:"ipv6_enabled,omitempty"`
}

// CloudComputingInstanceReinstallInput represents cloud instance reinstall input
type CloudComputingInstanceReinstallInput struct {
	ImageID string `json:"image_id"`
}

// CloudComputingInstanceUpgradeInput represents cloud instance upgrade input
type CloudComputingInstanceUpgradeInput struct {
	FlavorID string `json:"flavor_id"`
}

// CloudComputingRegion represents cloud computing region
type CloudComputingRegion struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// CloudComputingImage represents cloud computing image
type CloudComputingImage struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CloudComputingFlavor represents cloud computing flavor
type CloudComputingFlavor struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// L2Segment represents l2 segment
type L2Segment struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Type              string    `json:"type"`
	Status            string    `json:"status"`
	LocationGroupID   int64     `json:"location_group_id"`
	LocationGroupCode string    `json:"location_group_code"`
	Created           time.Time `json:"created_at"`
	Updated           time.Time `json:"updated_at"`
}

// L2SegmentMemberInput represents l2 segment member input for L2SegmentCreateInput and L2SegmentUpdateInput
type L2SegmentMemberInput struct {
	ID   string `json:"id"`
	Mode string `json:"mode"`
}

// L2SegmentCreateInput represents l2 segment create input
type L2SegmentCreateInput struct {
	Name            *string                `json:"name,omitempty"`
	Type            string                 `json:"type"`
	LocationGroupID int64                  `json:"location_group_id"`
	Members         []L2SegmentMemberInput `json:"members"`
}

// L2SegmentUpdateInput represents l2 segment update input
type L2SegmentUpdateInput struct {
	Name    *string                `json:"name,omitempty"`
	Members []L2SegmentMemberInput `json:"members,omitempty"`
}

// L2Member respresents l2 segment member
type L2Member struct {
	ID      string    `json:"id"`
	Title   string    `json:"title"`
	Mode    string    `json:"mode"`
	Vlan    *int      `json:"vlan"`
	Status  string    `json:"status"`
	Created time.Time `json:"created_at"`
	Updated time.Time `json:"updated_at"`
}

// L2SegmentCreateNetworksInput represents input to create networks for L2SegmentChangeNetworksInput
type L2SegmentCreateNetworksInput struct {
	Mask               int    `json:"mask"`
	DistributionMethod string `json:"distribution_method"`
}

// L2SegmentChangeNetworksInput represents input to change networks
type L2SegmentChangeNetworksInput struct {
	Create []L2SegmentCreateNetworksInput `json:"create,omitempty"`
	Delete []string                       `json:"delete,omitempty"`
}

// Network represents network
type Network struct {
	ID                 string    `json:"id"`
	Title              *string   `json:"title,omitempty"`
	Status             string    `json:"status"`
	Cidr               *string   `json:"cidr,omitempty"`
	Family             string    `json:"family"`
	InterfaceType      string    `json:"interface_type"`
	DistributionMethod string    `json:"distribution_method"`
	Additional         bool      `json:"additional"`
	Created            time.Time `json:"created_at"`
	Updated            time.Time `json:"updated_at"`

	// DEPRECATED: should be replaced by Statu
	State string `json:"state"`
}

// L2LocationGroup represents l2 location groups
type L2LocationGroup struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Code        string  `json:"code"`
	GroupType   string  `json:"group_type"`
	LocationIDs []int64 `json:"location_ids"`
}

// HostPowerFeed represents feed status
type HostPowerFeed struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

// HostConnection represents host connection
type HostConnection struct {
	Port       string  `json:"port"`
	Type       string  `json:"type"`
	MACAddress *string `json:"macaddr"`
}

// PTRRecord represents ptr record
type PTRRecord struct {
	ID       string `json:"id"`
	IP       string `json:"ip"`
	Domain   string `json:"domain"`
	Priority int    `json:"priority"`
	TTL      int    `json:"ttl"`
}

// PTRRecordCreateInput represents ptr record create input
type PTRRecordCreateInput struct {
	IP       string `json:"ip"`
	Domain   string `json:"domain"`
	Priority *int   `json:"priority"`
	TTL      *int   `json:"ttl"`
}

// OperatingSystemReinstallPartitionInput represents partition for os reinstallation layout input
type OperatingSystemReinstallPartitionInput struct {
	Target string  `json:"target"`
	Size   int     `json:"size"`
	Fs     *string `json:"fs,omitempty"`
	Fill   bool    `json:"fill,omitempty"`
}

// OperatingSystemReinstallLayoutInput represents layout for os reinstallation drives input
type OperatingSystemReinstallLayoutInput struct {
	SlotPositions []int                                    `json:"slot_positions"`
	Raid          *int                                     `json:"raid,omitempty"`
	Ignore        *bool                                    `json:"ignore,omitempty"`
	Partitions    []OperatingSystemReinstallPartitionInput `json:"partitions,omitempty"`
}

// OperatingSystemReinstallDrivesInput represents drives for os reinstallation input
type OperatingSystemReinstallDrivesInput struct {
	Layout []OperatingSystemReinstallLayoutInput `json:"layout,omitempty"`
}

// OperatingSystemReinstallInput represents os reinstallation input
type OperatingSystemReinstallInput struct {
	Hostname           string                              `json:"hostname"`
	Drives             OperatingSystemReinstallDrivesInput `json:"drives"`
	OperatingSystemID  *int64                              `json:"operating_system_id,omitempty"`
	SSHKeyFingerprints []string                            `json:"ssh_key_fingerprints,omitempty"`
}

// HostDriveSlot represents host drive slot
type HostDriveSlot struct {
	Position   int         `json:"position"`
	Interface  string      `json:"interface"`
	FormFactor string      `json:"form_factor"`
	DriveModel *DriveModel `json:"drive_model"`
}
