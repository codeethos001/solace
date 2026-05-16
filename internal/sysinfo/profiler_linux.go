package sysinfo

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LinuxProfiler struct{}

func NewLinuxProfiler() Profiler {
	return &LinuxProfiler{}
}

func (p *LinuxProfiler) CollectIdentity() (*SystemIdentity, error) {
	identity := &SystemIdentity{}

	hostname, err := os.Hostname()
	if err == nil {
		identity.Hostname = hostname
	}

	osReleaseInfo := parseOSRelease()
	identity.OS = osReleaseInfo["NAME"]
	identity.OSVersion = osReleaseInfo["VERSION_ID"]

	if output, err := exec.Command("uname", "-r").Output(); err == nil {
		identity.KernelVersion = strings.TrimSpace(string(output))
	}
	if output, err := exec.Command("uname", "-m").Output(); err == nil {
		identity.Architecture = strings.TrimSpace(string(output))
	}

	identity.Architecture = runtime.GOARCH

	if data, err := os.ReadFile("/etc/machine-id"); err == nil {
		identity.MachineID = strings.TrimSpace(string(data))
	}

	if data, err := os.ReadFile("/proc/sys/kernel/random/boot_id"); err == nil {
		identity.BootID = strings.TrimSpace(string(data))
	}

	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) > 0 {
			if uptime, err := strconv.ParseFloat(fields[0], 64); err == nil {
				identity.UptimeSeconds = uint64(uptime)
			}
		}
	}

	if currentUser, err := user.Current(); err == nil {
		identity.CurrentUser = currentUser.Username
	}

	identity.Timezone, _ = time.Now().Zone()

	identity.OS = "Linux"

	return identity, nil
}

func (p *LinuxProfiler) CollectCPU() (*CPUInfo, error) {
	cpuInfo := &CPUInfo{}

	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		coreCount := make(map[int]bool)
		physicalCoreCount := make(map[int]bool)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					cpuInfo.ModelName = strings.TrimSpace(parts[1])
				}
			} else if strings.HasPrefix(line, "processor") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					if id, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						coreCount[id] = true
					}
				}
			} else if strings.HasPrefix(line, "core id") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					if id, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
						physicalCoreCount[id] = true
					}
				}
			} else if strings.HasPrefix(line, "cpu MHz") {
				parts := strings.Split(line, ":")
				if len(parts) > 1 {
					if freq, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64); err == nil {
						cpuInfo.FrequencyMHz = freq
					}
				}
			}
		}

		cpuInfo.LogicalCores = len(coreCount)
		if cpuInfo.LogicalCores == 0 {
			cpuInfo.LogicalCores = runtime.NumCPU()
		}
		if len(physicalCoreCount) > 0 {
			cpuInfo.PhysicalCores = len(physicalCoreCount)
		} else {
			cpuInfo.PhysicalCores = cpuInfo.LogicalCores
		}
	}

	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		if len(fields) >= 3 {
			cpuInfo.LoadAverage = make([]float64, 3)
			for i := 0; i < 3; i++ {
				if val, err := strconv.ParseFloat(fields[i], 64); err == nil {
					cpuInfo.LoadAverage[i] = val
				}
			}
		}
	}

	cpuInfo.UsagePercent = getCPUUsagePercent()

	cpuInfo.PerCoreUsage = getPerCoreUsage()
	cpuInfo.StealTime = getStealTime()

	cpuInfo.Temperature = getCPUTemperature()

	return cpuInfo, nil
}

func (p *LinuxProfiler) CollectMemory() (*MemoryInfo, error) {
	memInfo := &MemoryInfo{}

	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		memData := make(map[string]uint64)

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				key := strings.TrimSuffix(parts[0], ":")
				if val, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
					memData[key] = val * 1024 // Convert KB to bytes
				}
			}
		}

		memInfo.TotalRAM = memData["MemTotal"]
		memInfo.FreeRAM = memData["MemFree"]
		memInfo.CachedRAM = memData["Cached"] + memData["SReclaimable"]
		memInfo.BufferedRAM = memData["Buffers"]
		memInfo.SwapTotal = memData["SwapTotal"]
		memInfo.SwapFree = memData["SwapFree"]

		memInfo.UsedRAM = memInfo.TotalRAM - memInfo.FreeRAM - memInfo.CachedRAM - memInfo.BufferedRAM
		memInfo.SwapUsed = memInfo.SwapTotal - memInfo.SwapFree

		if memInfo.TotalRAM > 0 {
			memInfo.UsagePercent = float64(memInfo.UsedRAM) / float64(memInfo.TotalRAM) * 100
		}
	}

	parsePressure("/proc/pressure/memory", memInfo)

	return memInfo, nil
}

func (p *LinuxProfiler) CollectDisk() ([]DiskInfo, error) {
	var disks []DiskInfo

	var mountData []byte
	var err error

	if mountData, err = os.ReadFile("/proc/mounts"); err != nil {
		if mountData, err = os.ReadFile("/etc/mtab"); err != nil {
			return disks, err
		}
	}

	scanner := bufio.NewScanner(bytes.NewReader(mountData))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}

		device := parts[0]
		mountPoint := parts[1]
		fsType := parts[2]
		mountOptions := parts[3]

		if strings.HasPrefix(device, "/dev/") || strings.HasPrefix(fsType, "tmpfs") ||
			strings.HasPrefix(fsType, "devtmpfs") || strings.HasPrefix(fsType, "sysfs") {
			if !strings.HasPrefix(device, "/dev/") {
				continue
			}
		}

		disk := DiskInfo{
			Device:       device,
			MountPoint:   mountPoint,
			FS:           fsType,
			MountOptions: strings.Split(mountOptions, ","),
		}

		for _, opt := range disk.MountOptions {
			switch opt {
			case "nodev":
				disk.Nodev = true
			case "nosuid":
				disk.Nosuid = true
			case "noexec":
				disk.Noexec = true
			case "ro":
				disk.ReadOnly = true
			}
		}

		if stat, err := os.Stat(mountPoint); err == nil && stat.IsDir() {
			if usage := getDiskUsage(mountPoint); usage != nil {
				disk.Total = usage.Total
				disk.Used = usage.Used
				disk.Free = usage.Free
			}
		}

		disks = append(disks, disk)
	}

	return disks, nil
}

func (p *LinuxProfiler) CollectNetwork() (*NetworkInfo, error) {
	netInfo := &NetworkInfo{
		Interfaces:  []NetInterface{},
		OpenPorts:   []PortInfo{},
		DNS:         []string{},
		Connections: []Connection{},
	}

	if data, err := os.ReadFile("/etc/resolv.conf"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "nameserver") {
				parts := strings.Fields(line)
				if len(parts) > 1 {
					netInfo.DNS = append(netInfo.DNS, parts[1])
				}
			}
		}
	}

	netInfo.Interfaces = getNetworkInterfaces()

	netInfo.OpenPorts = getOpenPorts()

	netInfo.Connections = getNetworkConnections()

	netInfo.TORIndicators = detectTOR(netInfo.Connections)

	return netInfo, nil
}

func (p *LinuxProfiler) CollectSecurity() (*SecurityInfo, error) {
	secInfo := &SecurityInfo{}

	if _, err := os.Stat("/etc/selinux/config"); err == nil {
		secInfo.SELinuxEnabled = isSELinuxEnabled()
	}

	if _, err := os.Stat("/etc/apparmor"); err == nil {
		secInfo.AppArmorEnabled = isAppArmorEnabled()
	}

	secInfo.UFWEnabled = isUFWEnabled()
	secInfo.FirewallEnabled = isFirewallEnabled()

	if data, err := os.ReadFile("/proc/sys/kernel/randomize_va_space"); err == nil {
		val := strings.TrimSpace(string(data))
		secInfo.ASLREnabled = val != "0"
	}

	if _, err := os.Stat("/var/run/auditd.pid"); err == nil {
		secInfo.AuditdRunning = true
	}

	return secInfo, nil
}

func (p *LinuxProfiler) CollectSSH() (*SSHInfo, error) {
	sshInfo := &SSHInfo{}

	if data, err := os.ReadFile("/etc/ssh/sshd_config"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		hardened := 0
		maxScore := 7

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if strings.HasPrefix(line, "#") || line == "" {
				continue
			}

			parts := strings.Fields(line)
			if len(parts) < 2 {
				continue
			}

			switch strings.ToLower(parts[0]) {
			case "permitrootlogin":
				sshInfo.PermitRootLogin = strings.ToLower(parts[1]) == "yes"
				if !sshInfo.PermitRootLogin {
					hardened++
				}
			case "passwordauthentication":
				sshInfo.PasswordAuth = strings.ToLower(parts[1]) == "yes"
				if !sshInfo.PasswordAuth {
					hardened++
				}
			case "maxauthtries":
				if val, err := strconv.Atoi(parts[1]); err == nil {
					sshInfo.MaxAuthTries = val
					if val <= 3 {
						hardened++
					}
				}
			case "allowusers":
				sshInfo.AllowUsers = parts[1:]
				if len(sshInfo.AllowUsers) > 0 {
					hardened++
				}
			case "ciphers":
				sshInfo.Ciphers = []string{strings.Join(parts[1:], ",")}
			case "macs":
				sshInfo.MACs = []string{strings.Join(parts[1:], ",")}
			case "kexalgorithms":
				sshInfo.KexAlgorithms = []string{strings.Join(parts[1:], ",")}
			}
		}

		sshInfo.HardeningScore = float64(hardened) / float64(maxScore) * 100
	} else {
		sshInfo.HardeningScore = 0
	}

	return sshInfo, nil
}

func (p *LinuxProfiler) CollectKernelSecurity() (*KernelSecurity, error) {
	kernelSec := &KernelSecurity{}

	kernelParams := map[string]*bool{
		"/proc/sys/net/ipv4/ip_forward":                      &kernelSec.IPForwarding,
		"/proc/sys/net/ipv4/tcp_syncookies":                  &kernelSec.TCPsyncookies,
		"/proc/sys/net/ipv4/conf/all/rp_filter":              &kernelSec.RPFilter,
		"/proc/sys/net/ipv4/conf/all/accept_redirects":       &kernelSec.AcceptRedirects,
		"/proc/sys/net/ipv4/conf/all/send_redirects":         &kernelSec.AcceptRedirects,
		"/proc/sys/net/ipv4/conf/all/accept_source_route":    &kernelSec.AcceptSourceRoute,
		"/proc/sys/net/ipv4/conf/all/reverse_path_filter":    &kernelSec.ReversePathFilter,
		"/proc/sys/net/ipv4/icmp_ignore_bogus_error_responses": &kernelSec.ICMPRedirects,
	}

	for path, target := range kernelParams {
		if data, err := os.ReadFile(path); err == nil {
			val := strings.TrimSpace(string(data))
			*target = val != "0"
		}
	}

	return kernelSec, nil
}

func (p *LinuxProfiler) CollectServices() ([]ServiceInfo, error) {
	var services []ServiceInfo

	if output, err := exec.Command("systemctl", "list-units", "--all", "--type=service", "--quiet").Output(); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(output))
		forbiddenServices := []string{"telnet", "ftp", "rsh", "snmp", "xinetd"}

		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Fields(line)
			if len(parts) < 3 {
				continue
			}

			service := ServiceInfo{
				Name:   strings.TrimSuffix(parts[0], ".service"),
				Status: parts[2],
			}

			if _, err := exec.Command("systemctl", "is-enabled", parts[0]).Output(); err == nil {
				service.Enabled = true
			}

			for _, forbidden := range forbiddenServices {
				if strings.Contains(service.Name, forbidden) {
					service.Forbidden = true
					break
				}
			}

			services = append(services, service)
		}
	}

	return services, nil
}

func (p *LinuxProfiler) CollectProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	procDir, err := os.Open("/proc")
	if err != nil {
		return processes, err
	}
	defer procDir.Close()

	entries, err := procDir.Readdirnames(-1)
	if err != nil {
		return processes, err
	}

	for _, entry := range entries {
		pid, err := strconv.Atoi(entry)
		if err != nil {
			continue
		}

		procInfo := getProcessInfo(pid)
		if procInfo != nil {
			processes = append(processes, *procInfo)
		}
	}

	return processes, nil
}

func (p *LinuxProfiler) CollectUsers() ([]UserInfo, error) {
	var users []UserInfo

	if data, err := os.ReadFile("/etc/passwd"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ":")
			if len(parts) < 7 {
				continue
			}

			uid, _ := strconv.Atoi(parts[2])
			gid, _ := strconv.Atoi(parts[3])

			userInfo := UserInfo{
				Username: parts[0],
				UID:      uid,
				GID:      gid,
				Home:     parts[5],
				Shell:    parts[6],
				IsRoot:   uid == 0,
			}

			if checkSudoAccess(parts[0]) {
				userInfo.Sudo = true
			}

			users = append(users, userInfo)
		}
	}

	return users, nil
}

func (p *LinuxProfiler) CollectAudit() (*AuditInfo, error) {
	auditInfo := &AuditInfo{}

	if _, err := exec.Command("systemctl", "is-active", "auditd").Output(); err == nil {
		auditInfo.AuditdEnabled = true
	}

	if output, err := exec.Command("auditctl", "-l").Output(); err == nil {
		auditInfo.AuditRulesLoaded = strings.Contains(string(output), "rule")
		auditInfo.AuditRuleCount = strings.Count(string(output), "rule")
	}

	if _, err := exec.Command("systemctl", "is-active", "systemd-journald").Output(); err == nil {
		auditInfo.JournaldActive = true
	}

	if _, err := exec.Command("systemctl", "is-active", "rsyslog").Output(); err == nil {
		auditInfo.RsyslogActive = true
	}

	return auditInfo, nil
}

func (p *LinuxProfiler) CollectPackages() ([]PackageInfo, error) {
	var packages []PackageInfo
	// This would require checking dpkg (Debian), rpm (RedHat), pacman (Arch), etc.
	// For now, returning empty slice
	return packages, nil
}

func (p *LinuxProfiler) CollectEnvironment() (*EnvironmentInfo, error) {
	envInfo := &EnvironmentInfo{}

	if isRunningInVM() {
		envInfo.IsVM = true
		envInfo.Hypervisor = detectHypervisor()
	}

	if isRunningInContainer() {
		envInfo.IsContainer = true
		envInfo.ContainerType = detectContainerType()
	}

	envInfo.CloudProvider = detectCloudProvider()

	return envInfo, nil
}

func (p *LinuxProfiler) CollectPersistence() ([]PersistenceMechanism, error) {
	var persistence []PersistenceMechanism
	// Check systemd services, cron jobs, startup scripts
	return persistence, nil
}

func (p *LinuxProfiler) CollectIntegrity() ([]IntegrityIndicator, error) {
	var integrity []IntegrityIndicator
	// Check executable hashes, SUID binaries, etc.
	return integrity, nil
}

func (p *LinuxProfiler) CollectHardware() (*HardwareInfo, error) {
	hwInfo := &HardwareInfo{}

	// Read from DMI/SMBIOS if available
	if data, err := os.ReadFile("/sys/class/dmi/id/sys_vendor"); err == nil {
		hwInfo.Vendor = strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile("/sys/class/dmi/id/product_name"); err == nil {
		hwInfo.Product = strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile("/sys/class/dmi/id/product_serial"); err == nil {
		hwInfo.Serial = strings.TrimSpace(string(data))
	}
	if data, err := os.ReadFile("/sys/class/dmi/id/bios_version"); err == nil {
		hwInfo.BIOSVersion = strings.TrimSpace(string(data))
	}

	if _, err := os.Stat("/dev/tpm0"); err == nil {
		hwInfo.TPMPresent = true
	}

	// secure boot
	if _, err := os.Stat("/sys/firmware/efi/efivars/SecureBoot-8be4df61-93ca-11d2-aa0d-00e098032b8c"); err == nil {
		hwInfo.SecureBoot = true
	}

	return hwInfo, nil
}


func parseOSRelease() map[string]string {
	osInfo := make(map[string]string)
	if data, err := os.ReadFile("/etc/os-release"); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, "=")
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.Trim(strings.TrimSpace(parts[1]), "\"")
				osInfo[key] = value
			}
		}
	}
	return osInfo
}

func getCPUUsagePercent() float64 {
	var stat1, stat2 string
	if data, err := os.ReadFile("/proc/stat"); err == nil {
		stat1 = string(data)
	}
	time.Sleep(100 * time.Millisecond)
	if data, err := os.ReadFile("/proc/stat"); err == nil {
		stat2 = string(data)
	}

	// Parse CPU times and calculate percentage
	// For simplicity, we'll return a placeholder
	_ = regexp.MustCompile(``)
	_ = stat1
	_ = stat2
	return 0
}

func getPerCoreUsage() []float64 {
	// This would require reading /proc/stat and parsing per-cpu entries
	return []float64{}
}

func getStealTime() float64 {
	// Read from /proc/stat for steal time
	return 0
}

func getCPUTemperature() float64 {
	entries, err := os.ReadDir("/sys/class/thermal")
	if err != nil {
		return 0
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "thermal_zone") {
			tempPath := filepath.Join("/sys/class/thermal", entry.Name(), "temp")
			if data, err := os.ReadFile(tempPath); err == nil {
				if temp, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64); err == nil {
					return temp / 1000 // Convert millidegrees to degrees
				}
			}
		}
	}
	return 0
}

func parsePressure(path string, memInfo *MemoryInfo) {
	if data, err := os.ReadFile(path); err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "some") {
				parts := strings.Fields(line)
				for _, part := range parts {
					if strings.HasPrefix(part, "avg10=") {
						if val, err := strconv.ParseFloat(strings.TrimPrefix(part, "avg10="), 64); err == nil {
							memInfo.PressureSome = val
						}
					}
				}
			} else if strings.HasPrefix(line, "full") {
				parts := strings.Fields(line)
				for _, part := range parts {
					if strings.HasPrefix(part, "avg10=") {
						if val, err := strconv.ParseFloat(strings.TrimPrefix(part, "avg10="), 64); err == nil {
							memInfo.PressureFull = val
						}
					}
				}
			}
		}
	}
}

type diskUsage struct {
	Total uint64
	Used  uint64
	Free  uint64
}

func getDiskUsage(path string) *diskUsage {
	// This would require syscall.Statfs on Linux
	// For now, returning nil
	return nil
}

func getNetworkInterfaces() []NetInterface {
	// Parse /proc/net/dev and /sys/class/net
	var interfaces []NetInterface
	return interfaces
}

func getOpenPorts() []PortInfo {
	// Parse /proc/net/tcp and /proc/net/udp
	var ports []PortInfo
	return ports
}

func getNetworkConnections() []Connection {
	// Parse /proc/net/tcp and /proc/net/udp
	var connections []Connection
	return connections
}

func detectTOR(connections []Connection) bool {
	// Check for TOR bridge IPs or ports
	return false
}

func isSELinuxEnabled() bool {
	// Check using getenforce command as its only for SELinux module
	if output, err := exec.Command("getenforce").Output(); err == nil {
		return !strings.Contains(string(output), "Disabled")
	}
	return false
}

func isAppArmorEnabled() bool {
	if _, err := os.Stat("/sys/module/apparmor"); err == nil {
		return true
	}
	return false
}

func isUFWEnabled() bool {
	if output, err := exec.Command("ufw", "status").Output(); err == nil {
		return strings.Contains(string(output), "active")
	}
	return false
}

func isFirewallEnabled() bool {
	// Check firewalld
	if output, err := exec.Command("firewall-cmd", "--state").Output(); err == nil {
		return strings.Contains(string(output), "running")
	}
	return false
}

func checkSudoAccess(username string) bool {
	// Check /etc/sudoers and /etc/sudoers.d
	output, err := exec.Command("sudo", "-l", "-U", username).Output()
	return err == nil && len(output) > 0
}

func getProcessInfo(pid int) *ProcessInfo {
	procInfo := &ProcessInfo{PID: pid}

	if data, err := os.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid)); err == nil {
		cmdline := string(data)
		procInfo.Cmdline = strings.ReplaceAll(cmdline, "\x00", " ")
	}

	if data, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid)); err == nil {
		parts := strings.Fields(string(data))
		if len(parts) > 3 {
			procInfo.Name = parts[1]
		}
		if len(parts) > 20 {
			if ppid, err := strconv.Atoi(parts[3]); err == nil {
				procInfo.PPID = ppid
			}
		}
	}

	return procInfo
}

func isRunningInVM() bool {
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		return strings.Contains(string(data), "hypervisor")
	}
	if _, err := os.Stat("/sys/hypervisor"); err == nil {
		return true
	}
	return false
}

func detectHypervisor() string {
	if data, err := os.ReadFile("/sys/class/dmi/id/sys_vendor"); err == nil {
		vendor := strings.TrimSpace(string(data))
		switch {
		case strings.Contains(vendor, "QEMU"):
			return "QEMU/KVM"
		case strings.Contains(vendor, "VMware"):
			return "VMware"
		case strings.Contains(vendor, "VirtualBox"):
			return "VirtualBox"
		case strings.Contains(vendor, "Xen"):
			return "Xen"
		case strings.Contains(vendor, "Hyper-V"):
			return "Hyper-V"
		}
	}
	return "Unknown"
}

func isRunningInContainer() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return true
	}
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "/docker") || strings.Contains(content, "/lxc") || strings.Contains(content, "/kube") {
			return true
		}
	}
	return false
}

func detectContainerType() string {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return "Docker"
	}
	if _, err := os.Stat("/run/.containerenv"); err == nil {
		return "Podman"
	}
	if data, err := os.ReadFile("/proc/1/cgroup"); err == nil {
		content := string(data)
		if strings.Contains(content, "/lxc") {
			return "LXC"
		}
		if strings.Contains(content, "/kube") {
			return "Kubernetes"
		}
	}
	return "Unknown"
}

func detectCloudProvider() string {
	if _, err := os.Stat("/sys/hypervisor/uuid"); err == nil {
		if data, err := os.ReadFile("/sys/hypervisor/uuid"); err == nil {
			uuid := string(data)
			if strings.HasPrefix(uuid, "ec2") {
				return "AWS"
			}
		}
	}
	return ""
}
