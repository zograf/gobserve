package serviceregistry

type ServiceInfo struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

//func (sd *ServiceInfo) ToJson() string {
//	b, _ := json.Marshal(sd)
//	return string(b)
//}
//
//func FromJson(s string) (*ServiceInfo, error) {
//	sd := &ServiceInfo{}
//	err := json.Unmarshal([]byte(s), sd)
//	return sd, err
//}
