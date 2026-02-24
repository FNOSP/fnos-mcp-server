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
	Btim int    `json:"btim" jsonschema:"文件创建时间"`
	Dir  int    `json:"dir" jsonschema:"是否为目录，0表示普通文件，1表示目录"`
	Mtim int    `json:"mtim" jsonschema:"文件最后修改时间"`
	Name string `json:"name" jsonschema:"文件或目录的名称"`
	Uid  int    `json:"uid" jsonschema:"文件所属用户的ID"`
	V    int    `json:"v" jsonschema:"存储空间标识"`
}
type GetFileLsRetDto struct {
	FnRetBase
	Files []FilesDto `json:"files"`
}
