package sysinfo

import "time"

type SystemIdentity struct {
	Hostname      string
	OS            string
	OSVersion     string
	KernelVersion string
	Architecture  string
	MachineID     string
	BootID        string
	Domain        string
	Timezone      string
	CurrentUser   string
	UptimeSeconds uint64
}

type CPUInfo struct {
	ModelName       string
	PhysicalCores   int
	LogicalCores    int
	UsagePercent    float64
	LoadAverage     []float64
	FrequencyMHz    float64
	Temperature     float64
	PerCoreUsage    []float64
	StealTime       float64
	ThrottlingState string
}

type MemoryInfo struct {
	TotalRAM      uint64
	UsedRAM       uint64
	FreeRAM       uint64
	CachedRAM     uint64
	BufferedRAM   uint64
	SwapTotal     uint64
	SwapUsed      uint64
	SwapFree      uint64
	UsagePercent  float64
	PressureLight float64
	PressureSome  float64
	PressureFull  float64
}

type DiskInfo struct {
	Device       string
	MountPoint   string
	FS           string
	Total        uint64
	Used         uint64
	Free         uint64
	MountOptions []string
	ReadOnly     bool
	Nodev        bool
	Nosuid       bool
	Noexec       bool
}

type NetInterface struct {
	Name      string
	MAC       string
	IPv4      []string
	IPv6      []string
	MTU       int
	Up        bool
	Wireless  bool
	Promiscuous bool
	Tunneled  bool
	VPN       bool
	Docker    bool
}

type PortInfo struct {
	Port        int
	Protocol    string
	State       string
	Address     string
	Process     string
	PID         int
	Unexpected  bool
}

type Connection struct {
	Protocol  string
	LocalAddr string
	LocalPort int
	RemoteAddr string
	RemotePort int
	State     string
	PID       int
}

type NetworkInfo struct {
	Interfaces     []NetInterface
	OpenPorts      []PortInfo
	DefaultGateway string
	DNS            []string
	Connections    []Connection
	TORIndicators  bool
}

type ServiceInfo struct {
	Name      string
	Status    string
	Enabled   bool
	PID       int
	StartType string
	User      string
	Rogue     bool
	Forbidden bool
}

type ProcessInfo struct {
	PID            int
	PPID           int
	Name           string
	Cmdline        string
	User           string
	CPUPercent     float64
	MemoryPercent  float64
	StartTime      int64
	Connections    []Connection
	ExecHash       string
	Deleted        bool
	SuspiciousChain bool
	HighEntropy    bool
}

type SecurityInfo struct {
	SELinuxEnabled      bool
	AppArmorEnabled     bool
	FirewallEnabled     bool
	SecureBootEnabled   bool
	ASLREnabled         bool
	AuditdRunning       bool
	UFWEnabled          bool
	DefenderEnabled     bool
}

type SSHInfo struct {
	PermitRootLogin  bool
	PasswordAuth     bool
	MaxAuthTries     int
	AllowUsers       []string
	Ciphers          []string
	MACs             []string
	KexAlgorithms    []string
	HardeningScore   float64
}

type KernelSecurity struct {
	IPForwarding      bool
	TCPsyncookies     bool
	RPFilter          bool
	AcceptRedirects   bool
	AcceptSourceRoute bool
	ReversePathFilter bool
	ICMPRedirects     bool
	UnsealedModules   bool
}

type UserInfo struct {
	Username      string
	UID           int
	GID           int
	Shell         string
	Home          string
	LastLogin     string
	Locked        bool
	Sudo          bool
	IsRoot        bool
	DormantDays   int
	WeakShell     bool
	Unauthorized  bool
}

type AuditInfo struct {
	AuditdEnabled    bool
	AuditRulesLoaded bool
	JournaldActive   bool
	RsyslogActive    bool
	RemoteLogging    bool
	AuditRuleCount   int
}

type PackageInfo struct {
	Name      string
	Version   string
	Source    string
	Vulnerable bool
	CVEs      []string
}

type EnvironmentInfo struct {
	IsVM           bool
	Hypervisor     string
	IsContainer    bool
	ContainerType  string
	CloudProvider  string
	KubernetesNode bool
	LXC            bool
	Docker         bool
}

type PersistenceMechanism struct {
	Type     string
	Name     string
	Path     string
	Enabled  bool
	Rogue    bool
	Modified time.Time
}

type IntegrityIndicator struct {
	Check        string
	Expected     string
	Actual       string
	Hash         string
	Modified     bool
	Unsigned     bool
	Suspicious   bool
}

type HardwareInfo struct {
	Vendor      string
	Product     string
	Serial      string
	BIOSVersion string
	TPMPresent  bool
	TPMVersion  string
	SecureBoot  bool
}

type Snapshot struct {
	Timestamp     time.Time
	Identity      SystemIdentity
	CPU           CPUInfo
	Memory        MemoryInfo
	Disk          []DiskInfo
	Network       NetworkInfo
	Security      SecurityInfo
	SSH           SSHInfo
	KernelSec     KernelSecurity
	Services      []ServiceInfo
	Processes     []ProcessInfo
	Users         []UserInfo
	Audit         AuditInfo
	Packages      []PackageInfo
	Environment   EnvironmentInfo
	Persistence   []PersistenceMechanism
	Integrity     []IntegrityIndicator
	Hardware      HardwareInfo
}
