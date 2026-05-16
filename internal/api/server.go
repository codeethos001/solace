package api

// Need to rewrite this for proper RESTful design and to integrate with the engine and sysinfo packages. This is just a starting point.
// Add comments.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	// "strings"
	"solace/internal/engine"
)

type Server struct {
	engine       *engine.HardeningEngine
	port         string
	blogService  *BlogService
}

func NewServer(e *engine.HardeningEngine, port string) *Server {
	return &Server{
		engine:      e,
		port:        port,
		blogService: NewBlogService(),
	}
}

func (s *Server) Start() error {
	// CORS headers helper
	corsMiddleware := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next(w, r)
		}
	}

	// Public API routes
	http.HandleFunc("/api/status", corsMiddleware(s.handleStatus))
	http.HandleFunc("/api/audit", corsMiddleware(s.handleAudit))

	// Authentication routes
	http.HandleFunc("/api/auth/login", corsMiddleware(s.handleLogin))
	http.HandleFunc("/api/auth/logout", corsMiddleware(s.AuthMiddleware(s.handleLogout)))
	http.HandleFunc("/api/auth/user", corsMiddleware(s.AuthMiddleware(s.handleGetUser)))

	// System info routes (protected)
	http.HandleFunc("/api/sysinfo", corsMiddleware(s.AuthMiddleware(s.handleGetSystemInfo)))
	http.HandleFunc("/api/sysinfo/identity", corsMiddleware(s.AuthMiddleware(s.handleGetIdentity)))
	http.HandleFunc("/api/sysinfo/cpu", corsMiddleware(s.AuthMiddleware(s.handleGetCPU)))
	http.HandleFunc("/api/sysinfo/memory", corsMiddleware(s.AuthMiddleware(s.handleGetMemory)))
	http.HandleFunc("/api/sysinfo/disk", corsMiddleware(s.AuthMiddleware(s.handleGetDisk)))
	http.HandleFunc("/api/sysinfo/network", corsMiddleware(s.AuthMiddleware(s.handleGetNetwork)))
	http.HandleFunc("/api/sysinfo/security", corsMiddleware(s.AuthMiddleware(s.handleGetSecurity)))
	http.HandleFunc("/api/sysinfo/processes", corsMiddleware(s.AuthMiddleware(s.handleGetProcesses)))

	// Blog/Advisory routes (protected)
	http.HandleFunc("/api/blogs", corsMiddleware(s.AuthMiddleware(s.handleGetBlogs)))
	http.HandleFunc("/api/blogs/search", corsMiddleware(s.AuthMiddleware(s.handleSearchBlogs)))

	fmt.Printf("🌐 API Server running at http://localhost:%s\n", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}

func (s *Server) setJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// ======================
// Authentication Handlers
// ======================

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
		return
	}

	if !s.validateCredentials(loginReq.Username, loginReq.Password) {
		http.Error(w, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return
	}

	token := generateToken(loginReq.Username)
	response := LoginResponse{
		Token: token,
		User: User{
			ID:       1,
			Username: loginReq.Username,
			Role:     "admin",
		},
	}

	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	json.NewEncoder(w).Encode(map[string]string{"message": "logged out successfully"})
}

func (s *Server) handleGetUser(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	username := r.Header.Get("X-Username")
	user := User{
		ID:       1,
		Username: username,
		Role:     "admin",
	}

	json.NewEncoder(w).Encode(user)
}

// ======================
// System Info Handlers
// ======================

func (s *Server) handleGetSystemInfo(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Return comprehensive system info object
	// In production, this would call the sysinfo package
	sysInfo := SystemInfoResponse{
		Hostname: "solace-server",
		OS:       "Linux",
		Identity: map[string]interface{}{
			"hostname":        "solace-server",
			"os":              "Linux",
			"os_version":      "22.04 LTS",
			"kernel_version": "5.15.0-86-generic",
			"architecture":    "x86_64",
			"machine_id":      "a1b2c3d4e5f6g7h8",
			"boot_id":         "boot1234567890",
			"timezone":        "UTC",
		},
		CPU: map[string]interface{}{
			"model_name":    "Intel Core i7",
			"physical_cores": 4,
			"logical_cores": 8,
			"usage_percent": 15.5,
			"load_average":  "0.45 0.38 0.32",
			"frequency_mhz": 2400,
		},
		Memory: map[string]interface{}{
			"total_ram":     "16 GB",
			"used_ram":      "8.2 GB",
			"free_ram":      "7.8 GB",
			"cached_ram":    "2.5 GB",
			"swap_total":    "4 GB",
			"swap_used":     "0.5 GB",
			"usage_percent": 51.25,
		},
		Security: map[string]interface{}{
			"selinux_enabled":     true,
			"apparmor_enabled":    false,
			"firewall_enabled":    true,
			"secure_boot_enabled": true,
			"aslr_enabled":        true,
			"auditd_running":      true,
			"ufw_enabled":         true,
		},
	}

	json.NewEncoder(w).Encode(sysInfo)
}

func (s *Server) handleGetIdentity(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	identity := map[string]interface{}{
		"hostname":         "solace-server",
		"os":               "Linux",
		"os_version":       "22.04 LTS",
		"kernel_version":   "5.15.0-86-generic",
		"architecture":     "x86_64",
		"machine_id":       "a1b2c3d4e5f6g7h8",
		"boot_id":          "boot1234567890",
		"uptime_seconds":   1234567,
		"timezone":         "UTC",
		"current_user":     "admin",
	}
	json.NewEncoder(w).Encode(identity)
}

func (s *Server) handleGetCPU(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	cpu := map[string]interface{}{
		"model_name":    "Intel Core i7-10700K",
		"physical_cores": 8,
		"logical_cores": 16,
		"usage_percent": 15.5,
		"load_average":  "0.45 0.38 0.32",
		"frequency_mhz": 3800,
		"temperature":   "42.5°C",
		"count":         1,
	}
	json.NewEncoder(w).Encode(cpu)
}

func (s *Server) handleGetMemory(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	memory := map[string]interface{}{
		"total_ram":     "16777216 KB",
		"used_ram":      "8585932 KB",
		"free_ram":      "8191284 KB",
		"cached_ram":    "2621440 KB",
		"swap_total":    "4194304 KB",
		"swap_used":     "524288 KB",
		"usage_percent": 51.25,
	}
	json.NewEncoder(w).Encode(memory)
}

func (s *Server) handleGetDisk(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	disk := []map[string]interface{}{
		{
			"device":       "/dev/sda1",
			"mount_point":  "/",
			"fs":           "ext4",
			"total":        "100 GB",
			"used":         "45 GB",
			"free":         "55 GB",
			"mount_options": "rw,relatime",
			"read_only":    false,
		},
		{
			"device":       "/dev/sda2",
			"mount_point":  "/home",
			"fs":           "ext4",
			"total":        "200 GB",
			"used":         "120 GB",
			"free":         "80 GB",
			"mount_options": "rw,nodev,nosuid,relatime",
			"read_only":    false,
		},
	}
	json.NewEncoder(w).Encode(disk)
}

func (s *Server) handleGetNetwork(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	network := map[string]interface{}{
		"interfaces": []map[string]interface{}{
			{
				"name":     "eth0",
				"mac":      "00:0a:95:9d:68:16",
				"ipv4":     []string{"192.168.1.100"},
				"ipv6":     []string{"fe80::1"},
				"mtu":      1500,
				"up":       true,
				"wireless": false,
			},
		},
		"open_ports": []map[string]interface{}{
			{"port": 22, "protocol": "TCP", "service": "SSH"},
			{"port": 80, "protocol": "TCP", "service": "HTTP"},
			{"port": 443, "protocol": "TCP", "service": "HTTPS"},
		},
		"default_gateway": "192.168.1.1",
		"dns_servers":     []string{"8.8.8.8", "8.8.4.4"},
	}
	json.NewEncoder(w).Encode(network)
}

func (s *Server) handleGetSecurity(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	security := map[string]interface{}{
		"selinux_enabled":     true,
		"apparmor_enabled":    false,
		"firewall_enabled":    true,
		"secure_boot_enabled": true,
		"aslr_enabled":        true,
		"auditd_running":      true,
		"ufw_enabled":         true,
		"ssh_hardened":        true,
		"sudo_access":         true,
	}
	json.NewEncoder(w).Encode(security)
}

func (s *Server) handleGetProcesses(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)
	processes := []map[string]interface{}{
		{
			"pid":              1,
			"ppid":             0,
			"name":             "systemd",
			"cmdline":          "/lib/systemd/systemd",
			"user":             "root",
			"cpu_percent":      0.1,
			"memory_percent":   0.5,
			"start_time":       "2024-01-15 10:00:00",
		},
		{
			"pid":              1234,
			"ppid":             1,
			"name":             "sshd",
			"cmdline":          "/usr/sbin/sshd -D",
			"user":             "root",
			"cpu_percent":      0.0,
			"memory_percent":   0.2,
			"start_time":       "2024-01-15 10:01:00",
		},
	}
	json.NewEncoder(w).Encode(processes)
}

// ======================
// Blog/Advisory Handlers
// ======================

func (s *Server) handleGetBlogs(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Check for ID query parameter
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error":"invalid id"}`, http.StatusBadRequest)
			return
		}

		blog := s.blogService.GetBlogByID(id)
		if blog == nil {
			http.Error(w, `{"error":"blog not found"}`, http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(blog)
		return
	}

	// Check for category query parameter
	category := r.URL.Query().Get("category")
	var blogs []*Blog
	if category != "" {
		blogs = s.blogService.GetBlogsByCategory(category)
	} else {
		blogs = s.blogService.GetAllBlogs()
	}

	json.NewEncoder(w).Encode(blogs)
}

func (s *Server) handleSearchBlogs(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, `{"error":"missing search query"}`, http.StatusBadRequest)
		return
	}

	results := s.blogService.SearchBlogs(query)
	json.NewEncoder(w).Encode(results)
}

// ======================
// Audit/Hardening Handlers
// ======================

func (s *Server) handleAudit(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	results := s.engine.EvaluateRules()
	json.NewEncoder(w).Encode(results)
}

// Health check
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.setJSON(w)

	status := map[string]string{
		"status":  "online",
		"service": "Solace Security Toolkit",
		"version": "1.0.0",
	}
	json.NewEncoder(w).Encode(status)
}cd web
npm install
npm run dev