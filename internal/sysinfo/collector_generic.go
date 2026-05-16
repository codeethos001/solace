package sysinfo

type GenericProfiler struct{}

func NewGenericProfiler() Profiler {
	return &GenericProfiler{}
}

func (p *GenericProfiler) CollectIdentity() (*SystemIdentity, error) { return &SystemIdentity{}, nil }
func (p *GenericProfiler) CollectCPU() (*CPUInfo, error)              { return &CPUInfo{}, nil }
func (p *GenericProfiler) CollectMemory() (*MemoryInfo, error)        { return &MemoryInfo{}, nil }
func (p *GenericProfiler) CollectDisk() ([]DiskInfo, error)           { return []DiskInfo{}, nil }
func (p *GenericProfiler) CollectNetwork() (*NetworkInfo, error)      { return &NetworkInfo{}, nil }
func (p *GenericProfiler) CollectSecurity() (*SecurityInfo, error)    { return &SecurityInfo{}, nil }
func (p *GenericProfiler) CollectSSH() (*SSHInfo, error)              { return &SSHInfo{}, nil }
func (p *GenericProfiler) CollectKernelSecurity() (*KernelSecurity, error) {
	return &KernelSecurity{}, nil
}
func (p *GenericProfiler) CollectServices() ([]ServiceInfo, error)     { return []ServiceInfo{}, nil }
func (p *GenericProfiler) CollectProcesses() ([]ProcessInfo, error)    { return []ProcessInfo{}, nil }
func (p *GenericProfiler) CollectUsers() ([]UserInfo, error)           { return []UserInfo{}, nil }
func (p *GenericProfiler) CollectAudit() (*AuditInfo, error)           { return &AuditInfo{}, nil }
func (p *GenericProfiler) CollectPackages() ([]PackageInfo, error)     { return []PackageInfo{}, nil }
func (p *GenericProfiler) CollectEnvironment() (*EnvironmentInfo, error) {
	return &EnvironmentInfo{}, nil
}
func (p *GenericProfiler) CollectPersistence() ([]PersistenceMechanism, error) {
	return []PersistenceMechanism{}, nil
}
func (p *GenericProfiler) CollectIntegrity() ([]IntegrityIndicator, error) {
	return []IntegrityIndicator{}, nil
}
func (p *GenericProfiler) CollectHardware() (*HardwareInfo, error) { return &HardwareInfo{}, nil }
