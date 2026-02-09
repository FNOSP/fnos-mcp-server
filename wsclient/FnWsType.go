package wsclient

type FnRetBase struct {
	ReqId  string `json:"reqid"`
	Result string `json:"result"`
}

type GetRSAPubRetDto struct {
	FnRetBase
	Pub string `json:"pub"`
	Si  string `json:"si"`
}

type LoginRetDto struct {
	FnRetBase
	Admin     bool   `json:"admin"`
	BackId    string `json:"backId"`
	LongToken bool   `json:"longToken"`
	MachineId bool   `json:"machineId"`
	Secret    bool   `json:"secret"`
	Token     bool   `json:"token"`
	Uid       bool   `json:"uid"`
}

type GetHostNameRetDto struct {
	FnRetBase
	Data GetHostNameDataRetDto `json:"data"`
}

type GetHostNameDataRetDto struct {
	HasUsers    bool   `json:"hasUsers"`
	HostName    string `json:"hostName"`
	TrimVersion string `json:"trimVersion"`
}
