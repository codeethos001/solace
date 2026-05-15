package osops

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"os/exec"
	"strings"
	"regexp"
	"syscall"
)

type linuxEngine struct {
	osName string
}

func NewEngine() (OSEngine, error) {
	return &linuxEngine{
		osName: "Linux",
	}, nil
}

func (l *linuxEngine) GetOSName() string {
	return l.osName
}

func (l *linuxEngine) CheckPrivileges() (bool, error) {
	currentUser, err := user.Current()
	if err != nil {
		return false, err
	}
	isRoot := currentUser.Uid == "0"
	return isRoot, nil
}

func (l *linuxEngine) GetLogPaths() []string {
	return []string{"/var/log/auth.log", "/var/log/syslog"}
}

func (l *linuxEngine) CheckKernelModuleLoaded(moduleName string) (bool, error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		return false, fmt.Errorf("failed to open /proc/modules: %v", err)
	}
	defer file.Close()

	searchName := strings.ReplaceAll(moduleName, "-", "_")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, searchName+" ") {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}
	return false, nil 
}

func (l *linuxEngine) CheckMountPoint(targetPath string) (bool, []string, error) {
	file, err := os.Open("/proc/mounts")
	if err != nil {
		return false, nil, fmt.Errorf("failed to open /proc/mounts: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) >= 4 {
			mountPoint := fields[1]
			if mountPoint == targetPath {
				options := strings.Split(fields[3], ",")
				return true, options, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return false, nil, err
	}
	return false, nil, nil
}

func (l *linuxEngine) CheckServiceStatus(name string) (string, error) {
	// Using systemctl
	cmd := exec.Command("systemctl", "is-active", name)
	out, _ := cmd.Output()
	status := strings.TrimSpace(string(out))
	if status == "active" {
		return "running", nil
	}
	// check if failed or unknown
	if strings.Contains(status, "unknown") || strings.Contains(status, "not-found") {
		return "not_found", nil
	}
	return "stopped", nil
}

func (l *linuxEngine) GetSysctlValue(key string) (string, error) {
	path := "/proc/sys/" + strings.ReplaceAll(key, ".", "/")
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("sysctl key %s not found", key)
	}
	return strings.TrimSpace(string(content)), nil
}

func (l *linuxEngine) GetFilePermissions(path string) (string, string, string, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return "", "", "", err
	}
	
	sysStat := stat.Sys().(*syscall.Stat_t)
	mode := fmt.Sprintf("%04o", stat.Mode().Perm())
	
	owner := fmt.Sprint(sysStat.Uid)
	if u, err := user.LookupId(owner); err == nil {
		owner = u.Username
	}
	
	group := fmt.Sprint(sysStat.Gid)
	if g, err := user.LookupGroupId(group); err == nil {
		group = g.Name
	}
	
	return mode, owner, group, nil
}

func (l *linuxEngine) CheckFileRegex(path string, regexPattern string) (bool, string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return false, "", err
	}
	
	re, err := regexp.Compile("(?m)" + regexPattern)
	if err != nil {
		return false, "", fmt.Errorf("invalid regex pattern: %v", err)
	}
	
	match := re.FindString(string(content))
	if match != "" {
		return true, strings.TrimSpace(match), nil
	}
	return false, "", nil
}

func (l *linuxEngine) RunCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

// dummy implementation for Linux to satisfy interface
func (l *linuxEngine) GetSeceditValue(key string) (string, error) {
	return "", fmt.Errorf("secedit not supported on linux")
}

func (l *linuxEngine) GetRegistryValue(path string, key string) (string, error) {
	return "", fmt.Errorf("registry not supported on linux")
}
