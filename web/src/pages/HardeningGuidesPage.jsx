import React, { useState, useEffect } from 'react'
import { BookOpen, Calendar, ArrowRight, Search } from 'lucide-react'

export default function HardeningGuidesPage() {
  const [blogs, setBlogs] = useState([])
  const [filteredBlogs, setFilteredBlogs] = useState([])
  const [selectedBlog, setSelectedBlog] = useState(null)
  const [searchTerm, setSearchTerm] = useState('')
  const [isLoading, setIsLoading] = useState(true)

  // Mock blog data - In production, this would come from an API endpoint
  const mockBlogs = [
    {
      id: 1,
      title: 'Linux Kernel Hardening Best Practices',
      category: 'Kernel',
      date: '2024-05-15',
      excerpt: 'Learn how to harden your Linux kernel with essential sysctl parameters and kernel module management.',
      content: `
# Linux Kernel Hardening Best Practices

## Overview
Kernel hardening is crucial for system security. This guide covers essential configurations.

## Key Parameters

### IP Forwarding
Disable IP forwarding if your system is not a router:
\`\`\`
net.ipv4.ip_forward = 0
\`\`\`

### SYN Cookies
Enable SYN cookies to protect against SYN flood attacks:
\`\`\`
net.ipv4.tcp_syncookies = 1
\`\`\`

### ASLR (Address Space Layout Randomization)
Enable ASLR to prevent memory-based exploits:
\`\`\`
kernel.randomize_va_space = 2
\`\`\`

### Restrict Kernel Modules
Disable unnecessary kernel modules like cramfs, freevxfs, hfs, hfsplus, jffs2, etc.

## Verification
Run \`solace audit\` to verify your kernel hardening configuration.
      `,
    },
    {
      id: 2,
      title: 'SSH Configuration Security',
      category: 'SSH',
      date: '2024-05-14',
      excerpt: 'Secure your SSH daemon with strong authentication and encryption settings.',
      content: `
# SSH Configuration Security

## Critical SSH Settings

### Disable Root Login
Never allow SSH access as root:
\`\`\`
PermitRootLogin no
\`\`\`

### Use Public Key Authentication
Disable password authentication and use SSH keys:
\`\`\`
PasswordAuthentication no
PubkeyAuthentication yes
\`\`\`

### Limit Authentication Attempts
Prevent brute-force attacks:
\`\`\`
MaxAuthTries 3
MaxSessions 3
\`\`\`

### Use Strong Ciphers and MACs
Configure modern, secure algorithms:
\`\`\`
Ciphers aes-256-gcm@openssh.com,aes-128-gcm@openssh.com,aes-256-ctr,aes-128-ctr
MACs hmac-sha2-512-etm@openssh.com,hmac-sha2-256-etm@openssh.com
KexAlgorithms curve25519-sha256,curve25519-sha256@libssh.org
\`\`\`

### X11 Forwarding
Disable X11 forwarding unless needed:
\`\`\`
X11Forwarding no
\`\`\`

## File Permissions
Ensure proper permissions on SSH config files:
\`\`\`bash
chmod 600 ~/.ssh/authorized_keys
chmod 700 ~/.ssh
chmod 600 /etc/ssh/sshd_config
\`\`\`
      `,
    },
    {
      id: 3,
      title: 'File System Hardening',
      category: 'Filesystem',
      date: '2024-05-13',
      excerpt: 'Harden your file system mount points with appropriate permissions and options.',
      content: `
# File System Hardening

## Mount Options for Security

### /tmp Mount Options
Make /tmp secure against malicious files:
\`\`\`
/tmp: nodev,noexec,nosuid
\`\`\`

### /var Mount Options
Prevent executable code in /var:
\`\`\`
/var: nodev,nosuid
\`\`\`

### /home Mount Options
Restrict script execution in user directories:
\`\`\`
/home: nodev,nosuid
\`\`\`

### /dev/shm Mount Options
Secure temporary file system:
\`\`\`
/dev/shm: nodev,noexec,nosuid
\`\`\`

## Disable Unnecessary Filesystems

### Bluetooth Filesystem
\`\`\`bash
echo "install bluetooth /bin/true" >> /etc/modprobe.d/bluetooth.conf
\`\`\`

### USB-Storage
\`\`\`bash
echo "install usb-storage /bin/true" >> /etc/modprobe.d/usb-storage.conf
\`\`\`

## fstab Configuration
Example secure /etc/fstab configuration:
\`\`\`
/dev/sda1 /         ext4  defaults,errors=remount-ro  0 1
/dev/sda2 /tmp      ext4  defaults,nodev,noexec,nosuid 0 2
/dev/sda3 /var      ext4  defaults,nodev,nosuid        0 2
/dev/sda4 /home     ext4  defaults,nodev,nosuid        0 2
\`\`\`
      `,
    },
    {
      id: 4,
      title: 'User Account Management',
      category: 'Users',
      date: '2024-05-12',
      excerpt: 'Implement strong user authentication and account security policies.',
      content: `
# User Account Management

## Audit User Accounts

### Check for Empty Passwords
\`\`\`bash
awk -F: '($2 == "" ) { print $1 " does not have a password!" }' /etc/shadow
\`\`\`

### Check for UID 0 Accounts
\`\`\`bash
awk -F: '($3 == 0) { print $1 }' /etc/passwd
\`\`\`

## Password Policy

### Set Password Expiration
\`\`\`bash
grep PASS_MAX_DAYS /etc/login.defs  # Should be ≤ 365
grep PASS_MIN_DAYS /etc/login.defs  # Should be ≥ 1
grep PASS_WARN_AGE /etc/login.defs  # Should be ≥ 7
\`\`\`

### Implement Password Complexity
Use PAM to enforce strong passwords:
\`\`\`bash
# Install libpam-pwquality
apt-get install libpam-pwquality
\`\`\`

## Sudo Access Control

### Audit Sudoers
\`\`\`bash
audit syslog for 'sudo' commands
\`\`\`

### Sudo Environment
\`\`\`bash
# Disable sudo with no password
# Disable PATH_MAX
\`\`\`

## Account Lockout Policy
Implement account lockout after failed attempts:
\`\`\`bash
account required pam_tally2.so onerr=fail audit silent deny=5 unlock_time=900
\`\`\`
      `,
    },
    {
      id: 5,
      title: 'Network Security Configuration',
      category: 'Network',
      date: '2024-05-11',
      excerpt: 'Configure network interfaces and firewall rules for maximum security.',
      content: `
# Network Security Configuration

## Network Protocol Hardening

### Disable IPv6 (if not used)
\`\`\`bash
net.ipv6.conf.all.disable_ipv6 = 1
net.ipv6.conf.default.disable_ipv6 = 1
\`\`\`

### Prevent IP Spoofing
\`\`\`bash
net.ipv4.conf.all.arp_ignore = 1
net.ipv4.conf.all.arp_announce = 2
\`\`\`

### Enable Reverse Path Filtering
Prevent spoofed packets:
\`\`\`bash
net.ipv4.conf.all.rp_filter = 1
\`\`\`

### Disable ICMP Redirects
\`\`\`bash
net.ipv4.conf.all.send_redirects = 0
net.ipv4.conf.all.accept_redirects = 0
\`\`\`

## Firewall Configuration

### UFW (Uncomplicated Firewall)
\`\`\`bash
# Enable UFW
ufw enable

# Set default policies
ufw default deny incoming
ufw default allow outgoing

# Allow SSH
ufw allow ssh

# Allow HTTP/HTTPS
ufw allow 80/tcp
ufw allow 443/tcp

# Enable verbose logging
ufw logging on
\`\`\`

### Fail2Ban Integration
Protect against brute-force attacks:
\`\`\`bash
apt-get install fail2ban
systemctl enable fail2ban
\`\`\`

## Network Monitoring
\`\`\`bash
# Monitor active connections
netstat -tlnp

# Check listening ports
ss -tlnp
\`\`\`
      `,
    },
  ]

  useEffect(() => {
    // Simulate loading
    setIsLoading(true)
    setTimeout(() => {
      setBlogs(mockBlogs)
      setFilteredBlogs(mockBlogs)
      setIsLoading(false)
    }, 500)
  }, [])

  useEffect(() => {
    const filtered = blogs.filter(
      blog =>
        blog.title.toLowerCase().includes(searchTerm.toLowerCase()) ||
        blog.excerpt.toLowerCase().includes(searchTerm.toLowerCase()) ||
        blog.category.toLowerCase().includes(searchTerm.toLowerCase())
    )
    setFilteredBlogs(filtered)
  }, [searchTerm, blogs])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (selectedBlog) {
    return (
      <div className="p-6 max-w-4xl mx-auto">
        <button
          onClick={() => setSelectedBlog(null)}
          className="text-blue-600 hover:text-blue-700 font-medium mb-6 flex items-center gap-1"
        >
          ← Back to Articles
        </button>

        <div className="card">
          <div className="mb-6 pb-6 border-b">
            <div className="flex items-center gap-3 mb-3">
              <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-semibold">
                {selectedBlog.category}
              </span>
              <span className="text-gray-600 text-sm flex items-center gap-1">
                <Calendar className="w-4 h-4" />
                {new Date(selectedBlog.date).toLocaleDateString()}
              </span>
            </div>
            <h1 className="text-4xl font-bold text-gray-900 mb-2">{selectedBlog.title}</h1>
          </div>

          <div className="prose prose-sm max-w-none">
            {selectedBlog.content.split('\n').map((line, idx) => {
              if (line.startsWith('# ')) {
                return <h1 key={idx} className="text-2xl font-bold mt-6 mb-4">{line.replace('# ', '')}</h1>
              }
              if (line.startsWith('## ')) {
                return <h2 key={idx} className="text-xl font-bold mt-4 mb-3">{line.replace('## ', '')}</h2>
              }
              if (line.startsWith('### ')) {
                return <h3 key={idx} className="text-lg font-bold mt-3 mb-2">{line.replace('### ', '')}</h3>
              }
              if (line.startsWith('```')) {
                return null
              }
              if (line.trim()) {
                return <p key={idx} className="text-gray-700 mb-2 leading-relaxed">{line}</p>
              }
              return <div key={idx} className="mb-2" />
            })}
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-6xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Hardening Guides & Advisory Content</h1>
        <p className="text-gray-600 mt-2">Best practices and security hardening recommendations</p>
      </div>

      {/* Search */}
      <div className="mb-8">
        <div className="relative">
          <Search className="absolute left-3 top-3 w-5 h-5 text-gray-400" />
          <input
            type="text"
            placeholder="Search guides..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
          />
        </div>
      </div>

      {/* Blog Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        {filteredBlogs.length > 0 ? (
          filteredBlogs.map(blog => (
            <div
              key={blog.id}
              onClick={() => setSelectedBlog(blog)}
              className="card cursor-pointer hover:shadow-lg transition group"
            >
              <div className="flex items-start justify-between mb-3">
                <span className="bg-blue-100 text-blue-800 px-3 py-1 rounded-full text-sm font-semibold">
                  {blog.category}
                </span>
              </div>

              <h3 className="text-lg font-bold text-gray-900 mb-2 group-hover:text-blue-600 transition">
                {blog.title}
              </h3>

              <p className="text-gray-600 text-sm mb-4">{blog.excerpt}</p>

              <div className="flex items-center justify-between">
                <span className="text-xs text-gray-500 flex items-center gap-1">
                  <Calendar className="w-4 h-4" />
                  {new Date(blog.date).toLocaleDateString()}
                </span>
                <ArrowRight className="w-4 h-4 text-blue-600 group-hover:translate-x-1 transition" />
              </div>
            </div>
          ))
        ) : (
          <div className="col-span-full py-12 text-center">
            <BookOpen className="w-12 h-12 text-gray-400 mx-auto mb-4" />
            <p className="text-gray-600">No guides found matching your search</p>
          </div>
        )}
      </div>

      {/* Categories Summary */}
      <div className="mt-12 card">
        <h2 className="text-lg font-bold text-gray-900 mb-4">Coverage Areas</h2>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
          {['Kernel', 'SSH', 'Filesystem', 'Users', 'Network', 'Services', 'Audit', 'IAM'].map(
            category => (
              <div
                key={category}
                className="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-lg text-center"
              >
                <p className="font-semibold text-gray-900">{category}</p>
              </div>
            )
          )}
        </div>
      </div>
    </div>
  )
}
