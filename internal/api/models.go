package api

import "time"

type Blog struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	Category string    `json:"category"`
	Content  string    `json:"content"`
	Date     time.Time `json:"date"`
	Excerpt  string    `json:"excerpt"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type AuthUser struct {
	ID       int
	Username string
	Role     string
}

type SystemInfoResponse struct {
	Hostname   string                 `json:"hostname"`
	OS         string                 `json:"os"`
	Identity   map[string]interface{} `json:"identity,omitempty"`
	CPU        map[string]interface{} `json:"cpu,omitempty"`
	Memory     map[string]interface{} `json:"memory,omitempty"`
	Disk       []map[string]interface{} `json:"disk,omitempty"`
	Network    map[string]interface{} `json:"network,omitempty"`
	Security   map[string]interface{} `json:"security,omitempty"`
	Processes  []map[string]interface{} `json:"processes,omitempty"`
	Services   []map[string]interface{} `json:"services,omitempty"`
}
