package models

// Host represents the data structure for a DHCP host entry.
// Corresponds to the Host Pydantic model in the Python application.
type Host struct {
	Name                    string `json:"name" binding:"required"`
	HardwareEthernet        string `json:"hardware_ethernet" binding:"required"`
	OptionRouters           string `json:"option_routers" binding:"required"`
	OptionSubnetMask        string `json:"option_subnet_mask" binding:"required"`
	FixedAddress            string `json:"fixed_address" binding:"required"`
	OptionDomainNameServers string `json:"option_domain_name_servers" binding:"required"`
}

// HostUpdate represents the data structure for updating a DHCP host entry.
// Fields are optional, as in the HostUpdate Pydantic model.
// We use pointers to strings to distinguish between an empty string and a field not provided.
type HostUpdate struct {
	Name                    *string `json:"name,omitempty"`
	HardwareEthernet        *string `json:"hardware_ethernet,omitempty"`
	OptionRouters           *string `json:"option_routers,omitempty"`
	OptionSubnetMask        *string `json:"option_subnet_mask,omitempty"`
	FixedAddress            *string `json:"fixed_address,omitempty"`
	OptionDomainNameServers *string `json:"option_domain_name_servers,omitempty"`
}

// InterfaceOperation represents the data structure for adding or deleting an interface.
// Corresponds to the InterfaceOperation Pydantic model in the Python application.
type InterfaceOperation struct {
	Type      string `json:"type" binding:"required,oneof=v4 v6"` // "v4" or "v6"
	Interface string `json:"interface" binding:"required"`
}
