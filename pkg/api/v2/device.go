package v2

type Device struct {
	ResourceURI  string              `json:"resource_uri,omitempty"`
	SystemID     string              `json:"system_id,omitempty"`
	Hostname     string              `json:"Hostname,omitempty"`
	FQDN         string              `json:"FQDN,omitempty"`
	Parent       string              `json:"Parent,omitempty"`
	Owner        string              `json:"Owner,omitempty"`
	IPAddresses  []string            `json:"ip_addresses,omitempty"`
	InterfaceSet []*NetworkInterface `json:"interface_set,omitempty"`
	Zone         *Zone               `json:"Zone,omitempty"`
}
