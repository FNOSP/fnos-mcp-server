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

// NetworkNetListData 硬件信息数据结构体
type NetworkNetListData struct {
	Net NetworkInfo `json:"net" jsonschema:"设备的网络子系统完整信息,作为设备信息上报或配置的顶层容器,应用场景:设备信息上报/网络配置备份/远程诊断/自动化配置"`
}

// NetworkNetListResponse 硬件信息响应结构体
type NetworkNetListResponse struct {
	FnRetBase
	Data HardwareInfoData `json:"data"`
	Rev  string           `json:"rev"`
	Req  string           `json:"req"`
}

// HardwareInfoResponse 硬件信息查询响应结构体
type HardwareInfoResponse struct {
	FnRetBase
	Data HardwareInfoData `json:"data" jsonschema:"硬件信息数据体,包含CPU/内存/虚拟化/BIOS/磁盘等详细信息"`
	Rev  string           `json:"rev" jsonschema:"接口版本号,示例:0.1/1.0,用于API兼容性管理"`
	Req  string           `json:"req" jsonschema:"请求接口名称,示例:appcgi.sysinfo.getHardwareInfo,标识具体的API端点"`
}

// HardwareInfoData 硬件信息数据体
type HardwareInfoData struct {
	CPU     CPUInfo     `json:"cpu" jsonschema:"中央处理器(CPU)信息,包含型号/核心数/线程数等"`
	Mem     MemoryInfo  `json:"mem" jsonschema:"内存(Memory)信息,包含容量/频率/类型/制造商"`
	VM      VMInfo      `json:"vm" jsonschema:"虚拟化特性支持状态,包含VT-x/IOMMU/SR-IOV支持情况"`
	BIOS    BIOSInfo    `json:"bios" jsonschema:"固件/BIOS信息,包含厂商/版本/主板/系统DMI信息"`
	SysDisk SysDiskInfo `json:"sysdisk" jsonschema:"系统磁盘信息,包含容量/型号/协议/序列号"`
}

// CPUInfo 中央处理器信息
type CPUInfo struct {
	Name   string `json:"name" jsonschema:"CPU型号名称,包含厂商/系列/频率信息,示例:Intel(R) Core(TM) i7-10750H CPU @ 2.60GHz/AMD Ryzen 9 5900X"`
	Num    int    `json:"num" jsonschema:"物理CPU插槽数量(封装个数),示例:1(单路)/2(双路),普通消费级通常为1"`
	Core   int    `json:"core" jsonschema:"每个物理CPU的核心数(Core),示例:4(四核)/8(八核)"`
	Thread int    `json:"thread" jsonschema:"每个核心的逻辑线程数(Thread),考虑超线程技术,示例:4(无超线程)/8(开启超线程)"`
}

// MemoryInfo 内存信息
type MemoryInfo struct {
	Num       int    `json:"num" jsonschema:"物理内存条数量(DIMM插槽占用数),示例:1(单条)/2(双通道)"`
	Total     int64  `json:"total" jsonschema:"总物理内存容量,单位:字节(bytes),换算:4294967296表示4GB,8589934592表示8GB"`
	Frequency int    `json:"frequency" jsonschema:"内存运行频率,单位:MHz,0表示未知或未识别,示例:2666/3200/4800"`
	Type      string `json:"type" jsonschema:"内存技术类型,示例:DDR3/DDR4/DDR5/RAM/DRAM/LPDDR,虚拟环境可能显示为QEMU"`
	Vendor    string `json:"vendor" jsonschema:"内存制造商/供应商,示例:Samsung/Micron/Hynix/QEMU/Unknown"`
}

// VMInfo 虚拟化技术支持状态
type VMInfo struct {
	Available int `json:"available" jsonschema:"硬件虚拟化技术支持状态,0:不支持/1:支持,指CPU虚拟化扩展(VT-x/AMD-V),是运行KVM/VMware/Hyper-V的前提"`
	Iommu     int `json:"iommu" jsonschema:"IOMMU(输入输出内存管理单元)支持状态,0:不支持/1:支持,用于PCI设备直通(GPU直通/NVMe直通等SR-IOV场景)"`
	Sriov     int `json:"sriov" jsonschema:"SR-IOV(单根I/O虚拟化)支持状态,0:不支持/1:支持,允许单个物理PCI设备虚拟成多个轻量级VF(虚拟功能)供虚拟机使用"`
}

// BIOSInfo 固件和系统DMI信息
type BIOSInfo struct {
	Vendor    string        `json:"vendor" jsonschema:"BIOS/UEFI固件厂商,示例:SeaBIOS/AMI/Insyde/Phoenix/ American Megatrends"`
	Version   string        `json:"version" jsonschema:"BIOS固件版本号,示例:1.16.3-debian-1.16.3-2/F2.20"`
	Baseboard BaseboardInfo `json:"baseboard" jsonschema:"主板(Baseboard)信息,包含制造商/型号/序列号等DMI信息"`
	System    SystemInfo    `json:"system" jsonschema:"系统(System)信息,包含整机厂商/产品型号/UUID等DMI信息,用于资产管理和识别"`
}

// BaseboardInfo 主板信息
type BaseboardInfo struct {
	Vendor       string `json:"vendor,omitempty" jsonschema:"主板制造商,示例:QEMU/ASUS/Gigabyte/MSI/Dell/HP"`
	Name         string `json:"name,omitempty" jsonschema:"主板型号名称,示例:Z390-A PRO/Prime B550-Plus"`
	Version      string `json:"version,omitempty" jsonschema:"主板版本/修订号"`
	SerialNumber string `json:"serialNumber,omitempty" jsonschema:"主板序列号,用于资产追踪和保修验证"`
	AssetTag     string `json:"assetTag,omitempty" jsonschema:"主板资产标签,企业资产管理(Asset Management)标识"`
}

// SystemInfo 系统DMI信息(整机信息)
type SystemInfo struct {
	Vendor       string `json:"vendor" jsonschema:"系统厂商/品牌,示例:QEMU/Dell/HP/Lenovo/ASUSTeK/VMware"`
	ProductName  string `json:"productName" jsonschema:"产品名称/机型型号,示例:Standard PC (Q35 + ICH9, 2009)/OptiPlex 7080/XPS 15 9500"`
	Version      string `json:"version" jsonschema:"系统版本/修订,示例:pc-q35-10.0/Not Defined"`
	Sku          string `json:"sku" jsonschema:"SKU(库存单位)编号,商业产品标识符,用于区分配置版本"`
	SerialNumber string `json:"serialNumber" jsonschema:"整机序列号/服务标签(Service Tag),唯一标识符,保修和支持的关键依据"`
	Uuid         string `json:"uuid" jsonschema:"全局唯一标识符(UUID/GUID),格式:8-4-4-4-12十六进制字符,示例:1f6396b0-cb27-4e33-8088-948940cea6fb,用于操作系统授权和资产识别"`
	Family       string `json:"family" jsonschema:"产品家族系列,示例:ThinkPad X1 Carbon/Precision Workstation/undefined"`
}

// SysDiskInfo 系统磁盘信息
type SysDiskInfo struct {
	Size         int    `json:"size" jsonschema:"磁盘总容量,单位:字节(bytes),0表示容量未知或获取失败,示例:256060514304(256GB SSD)/1000204886016(1TB HDD)"`
	Model        string `json:"model" jsonschema:"磁盘型号/部件号,示例:Samsung SSD 970 EVO Plus 500GB/WDC WD10EZEX-00BN5A0/Unknown"`
	Protocol     string `json:"protocol" jsonschema:"磁盘连接协议/总线类型,枚举值:Unknown/SATA(SATA III)/NVMe(PCIe NVMe)/SAS/IDE"`
	SerialNumber string `json:"serialNumber" jsonschema:"磁盘序列号,唯一标识符,用于资产管理和保修查询,示例:S5H9NS0N123456/Unknown"`
}
