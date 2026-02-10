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
	LongToken string `json:"longToken"`
	MachineId string `json:"machineId"`
	Secret    string `json:"secret"`
	Token     string `json:"token"`
	Uid       int    `json:"uid"`
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

type FilesDto struct {
	Btim int    `json:"btim"`
	Dir  int    `json:"dir"`
	Mtim int    `json:"mtim"`
	Name string `json:"name"`
	Uid  int    `json:"uid"`
	V    int    `json:"v"`
}
type GetFileLsRetDto struct {
	FnRetBase
	Files []FilesDto `json:"files"`
}
