package sysinfo

import (
	"runtime"
	"sync"
	"time"
)

type Profiler interface {
	CollectIdentity() (*SystemIdentity, error)
	CollectCPU() (*CPUInfo, error)
	CollectMemory() (*MemoryInfo, error)
	CollectDisk() ([]DiskInfo, error)
	CollectNetwork() (*NetworkInfo, error)
	CollectSecurity() (*SecurityInfo, error)
	CollectSSH() (*SSHInfo, error)
	CollectKernelSecurity() (*KernelSecurity, error)
	CollectServices() ([]ServiceInfo, error)
	CollectProcesses() ([]ProcessInfo, error)
	CollectUsers() ([]UserInfo, error)
	CollectAudit() (*AuditInfo, error)
	CollectPackages() ([]PackageInfo, error)
	CollectEnvironment() (*EnvironmentInfo, error)
	CollectPersistence() ([]PersistenceMechanism, error)
	CollectIntegrity() ([]IntegrityIndicator, error)
	CollectHardware() (*HardwareInfo, error)
}

type Collector struct {
	profiler Profiler
	mu       sync.RWMutex
}

func NewCollector() *Collector {
	var profiler Profiler

	switch runtime.GOOS {
	case "linux":
		profiler = NewLinuxProfiler()
	case "windows":
		profiler = NewWindowsProfiler()
	case "darwin":
		profiler = NewDarwinProfiler()
	default:
		// Fallback to a basic profiler
		profiler = NewGenericProfiler()
	}

	return &Collector{
		profiler: profiler,
	}
}

func (c *Collector) CollectIdentity() (*SystemIdentity, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectIdentity()
}

func (c *Collector) CollectCPU() (*CPUInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectCPU()
}

func (c *Collector) CollectMemory() (*MemoryInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectMemory()
}

func (c *Collector) CollectDisk() ([]DiskInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectDisk()
}

func (c *Collector) CollectNetwork() (*NetworkInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectNetwork()
}

func (c *Collector) CollectSecurity() (*SecurityInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectSecurity()
}

func (c *Collector) CollectSSH() (*SSHInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectSSH()
}

func (c *Collector) CollectKernelSecurity() (*KernelSecurity, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectKernelSecurity()
}

func (c *Collector) CollectServices() ([]ServiceInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectServices()
}

func (c *Collector) CollectProcesses() ([]ProcessInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectProcesses()
}

func (c *Collector) CollectUsers() ([]UserInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectUsers()
}

func (c *Collector) CollectAudit() (*AuditInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectAudit()
}

func (c *Collector) CollectPackages() ([]PackageInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectPackages()
}

func (c *Collector) CollectEnvironment() (*EnvironmentInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectEnvironment()
}

func (c *Collector) CollectPersistence() ([]PersistenceMechanism, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectPersistence()
}

func (c *Collector) CollectIntegrity() ([]IntegrityIndicator, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectIntegrity()
}

func (c *Collector) CollectHardware() (*HardwareInfo, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.profiler.CollectHardware()
}

func (c *Collector) CollectSnapshot() (*Snapshot, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	snapshot := &Snapshot{
		Timestamp: time.Now(),
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	identity, err := c.profiler.CollectIdentity()
	if err != nil {
		errs = append(errs, err)
	} else {
		snapshot.Identity = *identity
	}

	collectors := []func() error{
		func() error {
			cpu, err := c.profiler.CollectCPU()
			if err == nil {
				snapshot.CPU = *cpu
			}
			return err
		},
		func() error {
			mem, err := c.profiler.CollectMemory()
			if err == nil {
				snapshot.Memory = *mem
			}
			return err
		},
		func() error {
			disk, err := c.profiler.CollectDisk()
			if err == nil {
				snapshot.Disk = disk
			}
			return err
		},
		func() error {
			net, err := c.profiler.CollectNetwork()
			if err == nil {
				snapshot.Network = *net
			}
			return err
		},
		func() error {
			sec, err := c.profiler.CollectSecurity()
			if err == nil {
				snapshot.Security = *sec
			}
			return err
		},
		func() error {
			ssh, err := c.profiler.CollectSSH()
			if err == nil {
				snapshot.SSH = *ssh
			}
			return err
		},
		func() error {
			kernelsec, err := c.profiler.CollectKernelSecurity()
			if err == nil {
				snapshot.KernelSec = *kernelsec
			}
			return err
		},
		func() error {
			svcs, err := c.profiler.CollectServices()
			if err == nil {
				snapshot.Services = svcs
			}
			return err
		},
		func() error {
			procs, err := c.profiler.CollectProcesses()
			if err == nil {
				snapshot.Processes = procs
			}
			return err
		},
		func() error {
			users, err := c.profiler.CollectUsers()
			if err == nil {
				snapshot.Users = users
			}
			return err
		},
		func() error {
			audit, err := c.profiler.CollectAudit()
			if err == nil {
				snapshot.Audit = *audit
			}
			return err
		},
		func() error {
			pkgs, err := c.profiler.CollectPackages()
			if err == nil {
				snapshot.Packages = pkgs
			}
			return err
		},
		func() error {
			env, err := c.profiler.CollectEnvironment()
			if err == nil {
				snapshot.Environment = *env
			}
			return err
		},
		func() error {
			persist, err := c.profiler.CollectPersistence()
			if err == nil {
				snapshot.Persistence = persist
			}
			return err
		},
		func() error {
			integrity, err := c.profiler.CollectIntegrity()
			if err == nil {
				snapshot.Integrity = integrity
			}
			return err
		},
		func() error {
			hw, err := c.profiler.CollectHardware()
			if err == nil {
				snapshot.Hardware = *hw
			}
			return err
		},
	}

	for _, collector := range collectors {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				mu.Lock()
				errs = append(errs, err)
				mu.Unlock()
			}
		}(collector)
	}

	wg.Wait()

	if len(errs) > 0 {
		return snapshot, errs[0]
	}

	return snapshot, nil
}
