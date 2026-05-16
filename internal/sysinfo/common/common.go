package common

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ParseProcValue(data string, key string) (string, error) {
	lines := strings.Split(data, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, key) {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}
	return "", fmt.Errorf("key %s not found", key)
}

func ParseProcUintValue(data string, key string) (uint64, error) {
	val, err := ParseProcValue(data, key)
	if err != nil {
		return 0, err
	}
	// remove unit suffix if present (e.g., "12345 kB" -> "12345")
	parts := strings.Fields(val)
	if len(parts) > 0 {
		return strconv.ParseUint(parts[0], 10, 64)
	}
	return strconv.ParseUint(val, 10, 64)
}

func ParseMountOptions(opts string) map[string]bool {
	result := make(map[string]bool)
	for _, opt := range strings.Split(opts, ",") {
		opt = strings.TrimSpace(opt)
		if opt != "" {
			result[opt] = true
		}
	}
	return result
}

func HasMountOption(opts string, option string) bool {
	optMap := ParseMountOptions(opts)
	return optMap[option]
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func NormalizePercentage(val float64) float64 {
	if val < 0 {
		return 0
	}
	if val > 100 {
		return 100
	}
	return val
}

func BytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

func BytesToMB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024)
}

func BytesToKB(bytes uint64) float64 {
	return float64(bytes) / 1024
}

func GBToBytes(gb float64) uint64 {
	return uint64(gb * 1024 * 1024 * 1024)
}

func MBToBytes(mb float64) uint64 {
	return uint64(mb * 1024 * 1024)
}

func KBToBytes(kb float64) uint64 {
	return uint64(kb * 1024)
}
