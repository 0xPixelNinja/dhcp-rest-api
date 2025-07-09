package models

// Host represents a DHCP host entry
type Host struct {
	Name                    string `json:"name" binding:"required"`
	HardwareEthernet        string `json:"hardware_ethernet" binding:"required"`
	OptionRouters           string `json:"option_routers" binding:"required"`
	OptionSubnetMask        string `json:"option_subnet_mask" binding:"required"`
	FixedAddress            string `json:"fixed_address" binding:"required"`
	OptionDomainNameServers string `json:"option_domain_name_servers" binding:"required"`
}

// HostUpdate contains optional fields for updating a host
// Using pointers to distinguish between empty string and not provided
type HostUpdate struct {
	Name                    *string `json:"name,omitempty"`
	HardwareEthernet        *string `json:"hardware_ethernet,omitempty"`
	OptionRouters           *string `json:"option_routers,omitempty"`
	OptionSubnetMask        *string `json:"option_subnet_mask,omitempty"`
	FixedAddress            *string `json:"fixed_address,omitempty"`
	OptionDomainNameServers *string `json:"option_domain_name_servers,omitempty"`
}

// InterfaceOperation is used for adding or deleting network interfaces
type InterfaceOperation struct {
	Type      string `json:"type" binding:"required,oneof=v4 v6"`
	Interface string `json:"interface" binding:"required"`
}
