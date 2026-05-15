package osops

type OSEngine interface {
	GetOSName() string
	CheckPrivileges() (bool, error)
	GetLogPaths() []string

	// Common
	CheckServiceStatus(serviceName string) (string, error)

	// linux specific
	CheckKernelModuleLoaded(moduleName string) (bool, error)
	CheckMountPoint(path string) (isSeparate bool, options []string, err error)
	GetSysctlValue(key string) (string, error)
	GetFilePermissions(path string) (mode string, owner string, group string, err error)
	CheckFileRegex(path string, regexPattern string) (bool, string, error)
	RunCommand(name string, args ...string) (string, error)

	// windows specific
	GetSeceditValue(key string) (string, error)
	GetRegistryValue(path string, key string) (string, error)

}