package serviceregistry

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}
