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

type MachineIdDto struct {
	MachineId string `json:"machineId" jsonschema:"设备ID"`
}

type GetMachineIdRetDto struct {
	FnRetBase
	Data MachineIdDto `json:"data"`
}

// IPv6Addr IPv6地址结构体
type IPv6Addr struct {
	Addr      string `json:"addr" jsonschema:"IPv6地址的字符串表示,格式为标准IPv6格式(压缩或完整形式),示例:2001:db8::1或fe80::1"`
	PrefixLen int    `json:"prefixLen" jsonschema:"IPv6地址的前缀长度(子网掩码位数),范围:0-128,常见值:64(本地链路/全局单播)或128(回环/点对点)"`
	Scope     string `json:"scope" jsonschema:"IPv6地址的作用域/范围,枚举值:global(全局可路由)/link(链路本地)/host(主机本地)/site(站点本地),常见值:global(公网)/link(fe80::/10开头)/host(::1)"`
}

// NetworkInterface 网络接口结构体
type NetworkInterface struct {
	Name        string     `json:"name" jsonschema:"操作系统识别的接口名称(系统唯一标识符),Linux示例:eth0/ens33/wlan0/docker0/br-xxx/vethxxx,特殊值:lo(回环接口)"`
	Index       int        `json:"index" jsonschema:"操作系统分配的内部接口索引号,系统级唯一标识,用于底层网络编程(如socket bind),类型:正整数,通常从1开始"`
	IfType      int        `json:"ifType" jsonschema:"接口的硬件/逻辑类型(IANA ARP协议编号),常见值:1(其他)/6(以太网)/23(PPP)/24(回环)/71(Wi-Fi),参考:IANA ifType定义(MIB-2)"`
	Enable      bool       `json:"enable" jsonschema:"该接口是否被管理员启用(软件层面启用状态),true:已启用(ifconfig up),false:已禁用(ifconfig down),注意:与Running不同,Enable是配置状态"`
	Running     bool       `json:"running" jsonschema:"接口是否处于实际运行状态(物理层就绪),true:接口已激活且物理连接正常(网线已插/Wi-Fi已连接),通常Enable为true且物理层正常时Running为true"`
	Onlink      bool       `json:"onlink" jsonschema:"该接口是否配置为直连(On-Link),在路由上下文中表示该接口连接的网络不需要网关即可到达,应用场景:静态路由配置/NDP/ARP解析"`
	State       int        `json:"state" jsonschema:"接口的详细操作状态(RFC 2863标准状态),枚举值:1(未知)/2(未启用)/3(下层关闭)/4(下层关闭但正在启动)/5(测试中)/6(休眠)/7(不存在)/8(第二层向下)"`
	Duplex      bool       `json:"duplex" jsonschema:"以太网接口的双工模式,true:全双工(Full Duplex),false:半双工(Half Duplex),仅对以太网类型接口有意义(IfType为6)"`
	Speed       int        `json:"speed" jsonschema:"接口的链路速率(带宽),单位:Mbps(兆比特每秒),常见值:10/100/1000/2500/10000/40000/100000,说明:0或-1表示未知或无法获取"`
	Mtu         int        `json:"mtu" jsonschema:"最大传输单元(Maximum Transmission Unit),单位:字节(bytes),常见值:1500(标准以太网)/9000(巨型帧Jumbo Frame)/65536(回环),说明:超过此值的数据包需要分片"`
	HwAddr      string     `json:"hwAddr" jsonschema:"硬件地址(MAC地址,Media Access Control),格式:48位十六进制6组双字符冒号或连字符分隔,示例:00:1a:2b:3c:4d:5e或00-1A-2B-3C-4D-5E,全0或空字符串表示虚拟接口"`
	IsOvsPort   bool       `json:"isOvsPort" jsonschema:"标识该接口是否为Open vSwitch虚拟端口,true:OVS管理的虚拟端口,false:物理网卡或普通虚拟网卡,应用场景:SDN(软件定义网络)/云平台网络虚拟化"`
	Wireless    bool       `json:"wireless" jsonschema:"标识是否为无线接口(Wi-Fi/802.11),true:无线网卡,false:有线网卡或其他类型,通常IfType为71时此值为true"`
	Ipv4Dhcp    bool       `json:"ipv4Dhcp" jsonschema:"IPv4地址是否通过DHCP协议自动获取,true:DHCP自动分配,false:静态配置(Static/Manual),与Ipv4Mode字段相关"`
	Ipv4Mode    string     `json:"ipv4Mode" jsonschema:"IPv4地址的配置模式,枚举值:static(静态配置)/dhcp(自动获取)/pppoe(拨号)/disabled(禁用),决定Ipv4相关字段的配置来源"`
	Ipv4Broad   string     `json:"ipv4Broad" jsonschema:"IPv4广播地址(子网广播),格式:标准IPv4点分十进制如192.168.1.255,通常通过IP地址和子网掩码计算得出(如192.168.1.0/24的广播地址)"`
	Ipv4Mask    string     `json:"ipv4Mask" jsonschema:"IPv4子网掩码(点分十进制表示),格式:如255.255.255.0或255.255.0.0,对应CIDR:/24或/16等"`
	Ipv4Gateway string     `json:"ipv4Gateway" jsonschema:"IPv4默认网关地址,格式:标准IPv4地址如192.168.1.1,说明:发往非本地子网的数据包转发目标,通常为本网段第一个可用IP"`
	Ipv4Dns     string     `json:"ipv4Dns" jsonschema:"IPv4 DNS服务器地址,格式:单个IP地址或多个IP逗号分隔如8.8.8.8,114.114.114.114,说明:域名解析服务器,可为空表示未配置"`
	Ipv4        []string   `json:"ipv4" jsonschema:"接口绑定的所有IPv4地址列表(包含别名/secondary地址),格式:CIDR表示法列表如[192.168.1.10/24,10.0.0.5/8],主地址通常放在第一位,包含子网前缀长度"`
	Ipv4Addr    string     `json:"ipv4Addr" jsonschema:"主IPv4地址(首选地址,通常不带CIDR后缀),格式:点分十进制如192.168.1.10,与Ipv4字段不同这是去除了掩码的纯地址,通常用于展示或API调用"`
	Ipv6Mode    string     `json:"ipv6Mode" jsonschema:"IPv6地址的配置模式,枚举值:static(静态)/auto(SLAAC自动配置)/dhcpv6(DHCPv6)/disabled(禁用),决定Ipv6地址的获取方式"`
	Ipv6Gateway string     `json:"ipv6Gateway" jsonschema:"IPv6默认网关地址,格式:标准IPv6地址如fe80::1或2001:db8::1,说明:通常可能是链路本地地址(fe80::/10)或全局地址"`
	Ipv6Dns     string     `json:"ipv6Dns" jsonschema:"IPv6 DNS服务器地址,格式:单个IPv6或多个逗号分隔如2001:4860:4860::8888,2408:8888::8"`
	Ipv6        []IPv6Addr `json:"ipv6" jsonschema:"接口的完整IPv6地址配置列表(包含地址/前缀/作用域),说明:一个接口可能有多个IPv6地址(链路本地+全局地址+临时地址),结构:包含Addr/PrefixLen/Scope的结构体数组"`
}

// NetworkInfo 网络信息结构体
type NetworkInfo struct {
	Ifs []NetworkInterface `json:"ifs" jsonschema:"系统上所有网络接口的列表,包含:物理网卡/虚拟网卡/回环接口/Docker网桥等,排序:通常按Index或名称排序"`
}

// HardwareInfoData 硬件信息数据结构体
type HardwareInfoData struct {
	Net NetworkInfo `json:"net" jsonschema:"设备的网络子系统完整信息,作为设备信息上报或配置的顶层容器,应用场景:设备信息上报/网络配置备份/远程诊断/自动化配置"`
}

// HardwareInfoResponse 硬件信息响应结构体
type HardwareInfoResponse struct {
	FnRetBase
	Data HardwareInfoData `json:"data"`
	Rev  string           `json:"rev"`
	Req  string           `json:"req"`
}
