package api

// This is temporary and will be replaced with a proper CMS or database in the future.

import (
	"strings"
	"time"
)

type BlogService struct {
	blogs map[int]*Blog
	nextID int
}

func NewBlogService() *BlogService {
	bs := &BlogService{
		blogs:  make(map[int]*Blog),
		nextID: 1,
	}
	bs.initializeBlogs()
	return bs
}

// Crappy blogs made using content generation.

func (bs *BlogService) initializeBlogs() {
	blogs := []Blog{
		{
			Title:    "Linux Kernel Hardening Best Practices",
			Category: "Kernel",
			Date:     time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Learn how to harden your Linux kernel with essential sysctl parameters and kernel module management.",
			Content: `# Linux Kernel Hardening Best Practices

## Overview
Kernel hardening is crucial for system security. Focus on disabling unnecessary modules and configuring secure sysctl parameters.

## Key Parameters

### IP Forwarding
Disable IP forwarding if your system is not a router:
net.ipv4.ip_forward = 0

### SYN Cookies
Enable SYN cookies to protect against SYN flood attacks:
net.ipv4.tcp_syncookies = 1

### ASLR (Address Space Layout Randomization)
Enable ASLR to prevent memory-based exploits:
kernel.randomize_va_space = 2

### Disable Unnecessary Kernel Modules
Disable: cramfs, freevxfs, hfs, hfsplus, jffs2, vfat, udf, usb-storage`,
		},
		{
			Title:    "SSH Configuration Security",
			Category: "SSH",
			Date:     time.Date(2024, 5, 14, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Secure your SSH daemon with strong authentication and encryption settings.",
			Content: `# SSH Configuration Security

## Critical SSH Settings

### Disable Root Login
Never allow SSH access as root:
PermitRootLogin no

### Use Public Key Authentication
Disable password authentication:
PasswordAuthentication no
PubkeyAuthentication yes

### Limit Authentication Attempts
MaxAuthTries 3
MaxSessions 3

### Use Strong Ciphers
Configure modern algorithms for ciphers, MACs, and KEX algorithms.

### X11 Forwarding
Disable X11 forwarding:
X11Forwarding no`,
		},
		{
			Title:    "File System Hardening",
			Category: "Filesystem",
			Date:     time.Date(2024, 5, 13, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Harden your file system mount points with appropriate permissions and options.",
			Content: `# File System Hardening

## Mount Options for Security

### /tmp Mount Options
Make /tmp secure:
/tmp: nodev,noexec,nosuid

### /var Mount Options
Prevent executable code in /var:
/var: nodev,nosuid

### /home Mount Options
Restrict script execution in user directories:
/home: nodev,nosuid

### /dev/shm Mount Options
Secure temporary file system:
/dev/shm: nodev,noexec,nosuid`,
		},
		{
			Title:    "User Account Management",
			Category: "Users",
			Date:     time.Date(2024, 5, 12, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Implement strong user authentication and account security policies.",
			Content: `# User Account Management

## Audit User Accounts

### Check for Empty Passwords
Identify accounts without passwords.

### Check for UID 0 Accounts
Ensure only root has UID 0.

## Password Policy

### Set Password Expiration
PASS_MAX_DAYS ≤ 365
PASS_MIN_DAYS ≥ 1
PASS_WARN_AGE ≥ 7

### Implement Password Complexity
Use PAM (libpam-pwquality) for complex passwords.`,
		},
		{
			Title:    "Network Security Configuration",
			Category: "Network",
			Date:     time.Date(2024, 5, 11, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Configure network interfaces and firewall rules for maximum security.",
			Content: `# Network Security Configuration

## Network Protocol Hardening

### IP Spoofing Prevention
net.ipv4.conf.all.rp_filter = 1

### ICMP Redirects
Disable redirect handling:
net.ipv4.conf.all.accept_redirects = 0

## Firewall Configuration

### UFW (Uncomplicated Firewall)
- Enable UFW
- Set default deny incoming
- Set default allow outgoing
- Allow only necessary services

### Fail2Ban
Protect against brute-force attacks.`,
		},
		{
			Title:    "Service Hardening",
			Category: "Services",
			Date:     time.Date(2024, 5, 10, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Remove unnecessary services and harden critical ones.",
			Content: `# Service Hardening

## Disable Unnecessary Services

### Telnet, FTP, Rsh
These services transmit credentials in plain text. Use SSH instead.

### X11
Disable X11 server if not needed:
systemctl disable x11-common

### DNS/DHCP
Only enable if providing these services.

### Webmin, Plesk
Disable if not in use - they increase attack surface.

## Service Running Checks
Regularly audit running services:
systemctl list-unit-files --type=service`,
		},
		{
			Title:    "Audit and Logging Configuration",
			Category: "Audit",
			Date:     time.Date(2024, 5, 9, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Implement comprehensive logging and audit trails for security monitoring.",
			Content: `# Audit and Logging Configuration

## Auditd Configuration

### Enable Auditd
Ensure auditd service is enabled and running.

### Key Audit Rules
Monitor critical system calls:
- execve (process execution)
- open (file access)
- unlink (file deletion)
- sysctl (system configuration changes)

### Log File Monitoring
Enable rsyslog for centralized logging.

## Remote Logging
Forward logs to secure remote server for protection against tampering.`,
		},
		{
			Title:    "SELinux and AppArmor",
			Category: "IAM",
			Date:     time.Date(2024, 5, 8, 0, 0, 0, 0, time.UTC),
			Excerpt:  "Implement mandatory access controls for defense-in-depth security.",
			Content: `# SELinux and AppArmor

## SELinux (Security-Enhanced Linux)

### Check SELinux Status
getenforce

### Common Modes
- Enforcing: Policies are enforced
- Permissive: Violations logged but not enforced
- Disabled: SELinux is off

### Policy Types
- targeted: Protects selected services
- strict: Complete protection

## AppArmor

### Check AppArmor Status
apparmor_status

### Profile Modes
- enforce: Policy is enforced
- complain: Violations logged, not enforced

### Important Profiles
Enable profiles for: sudo, apache, nginx, mysql`,
		},
	}

	for _, blog := range blogs {
		blog.ID = bs.nextID
		bs.blogs[bs.nextID] = &blog
		bs.nextID++
	}
}

func (bs *BlogService) GetAllBlogs() []*Blog {
	var result []*Blog
	for _, blog := range bs.blogs {
		result = append(result, blog)
	}
	return result
}

func (bs *BlogService) GetBlogByID(id int) *Blog {
	return bs.blogs[id]
}

func (bs *BlogService) SearchBlogs(query string) []*Blog {
	query = strings.ToLower(query)
	var result []*Blog

	for _, blog := range bs.blogs {
		if strings.Contains(strings.ToLower(blog.Title), query) ||
			strings.Contains(strings.ToLower(blog.Content), query) ||
			strings.Contains(strings.ToLower(blog.Category), query) ||
			strings.Contains(strings.ToLower(blog.Excerpt), query) {
			result = append(result, blog)
		}
	}

	return result
}

func (bs *BlogService) GetBlogsByCategory(category string) []*Blog {
	var result []*Blog
	for _, blog := range bs.blogs {
		if blog.Category == category {
			result = append(result, blog)
		}
	}
	return result
}

func (bs *BlogService) AddBlog(blog *Blog) int {
	blog.ID = bs.nextID
	bs.blogs[bs.nextID] = blog
	bs.nextID++
	return blog.ID
}
