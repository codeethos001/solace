package sysinfo

type WindowsProfiler struct{}

func NewWindowsProfiler() Profiler {
	return &WindowsProfiler{}
}

func (p *WindowsProfiler) CollectIdentity() (*SystemIdentity, error) { return &SystemIdentity{}, nil }
func (p *WindowsProfiler) CollectCPU() (*CPUInfo, error)              { return &CPUInfo{}, nil }
func (p *WindowsProfiler) CollectMemory() (*MemoryInfo, error)        { return &MemoryInfo{}, nil }
func (p *WindowsProfiler) CollectDisk() ([]DiskInfo, error)           { return []DiskInfo{}, nil }
func (p *WindowsProfiler) CollectNetwork() (*NetworkInfo, error)      { return &NetworkInfo{}, nil }
func (p *WindowsProfiler) CollectSecurity() (*SecurityInfo, error)    { return &SecurityInfo{}, nil }
func (p *WindowsProfiler) CollectSSH() (*SSHInfo, error)              { return &SSHInfo{}, nil }
func (p *WindowsProfiler) CollectKernelSecurity() (*KernelSecurity, error) {
	return &KernelSecurity{}, nil
}
func (p *WindowsProfiler) CollectServices() ([]ServiceInfo, error)     { return []ServiceInfo{}, nil }
func (p *WindowsProfiler) CollectProcesses() ([]ProcessInfo, error)    { return []ProcessInfo{}, nil }
func (p *WindowsProfiler) CollectUsers() ([]UserInfo, error)           { return []UserInfo{}, nil }
func (p *WindowsProfiler) CollectAudit() (*AuditInfo, error)           { return &AuditInfo{}, nil }
func (p *WindowsProfiler) CollectPackages() ([]PackageInfo, error)     { return []PackageInfo{}, nil }
func (p *WindowsProfiler) CollectEnvironment() (*EnvironmentInfo, error) {
	return &EnvironmentInfo{}, nil
}
func (p *WindowsProfiler) CollectPersistence() ([]PersistenceMechanism, error) {
	return []PersistenceMechanism{}, nil
}
func (p *WindowsProfiler) CollectIntegrity() ([]IntegrityIndicator, error) {
	return []IntegrityIndicator{}, nil
}
func (p *WindowsProfiler) CollectHardware() (*HardwareInfo, error) { return &HardwareInfo{}, nil }
