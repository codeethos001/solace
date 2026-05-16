# Solace Web Dashboard - Complete Setup and Documentation

A comprehensive security dashboard for the Solace Security Toolkit with system profiler integration, hardening rules evaluation, and security advisory content.

## 📋 Table of Contents

- [Features](#current-features)
- [Architecture](#architecture)
- [Installation & Setup](#installation--setup)
- [API Endpoints](#api-endpoints)
- [Authentication](#authentication)
- [Development](#development)
- [Deployment](#deployment)

## Current Features

### Dashboard
- **System Overview**: Real-time system information (hostname, uptime, users)
- **Security Posture Score**: Compliance percentage based on audit results
- **Quick Stats**: Pass/Fail/Warning summary of all security checks

### System Information
- CPU details (cores, usage, frequency, temperature)
- Memory statistics (RAM, swap, usage percentage)
- Disk information (mount points, usage, permissions)
- Network interfaces (IPs, MACs, connection status)
- Security configuration status

### Audit Results
- Complete rule evaluation results with filtering
- Status indicators (Pass, Fail, Warning, Skipped)
- Severity levels (Critical, High, Medium, Low)
- Category-based organization
- Reference links to security documentation
- Rule-by-rule details with expected vs current values

### Hardening Guides & Advisories
- 8+ comprehensive security hardening guides
- Coverage areas: Kernel, SSH, Filesystem, Users, Network, Services, Audit, IAM
- Full markdown formatted content
- Search functionality across all guides
- Category-based browsing

### Security Features
- Admin-only access with login authentication
- Simple token-based authentication (JWT-ready)
- CORS-enabled for cross-origin requests

## Architecture

### Frontend Stack
- **React 18**: UI framework
- **Vite**: Build tool and dev server
- **React Router**: Client-side routing
- **Axios**: HTTP client
- **Tailwind CSS**: Styling
- **Lucide React**: Icons

### Backend Stack
- **Go HTTP Server**: RESTful API
- **CORS Middleware**: Cross-origin support
- **Authentication**: Token-based (Bearer tokens)
- **Blog Service**: In-memory blog management (temporary, we may use more robust one later)

### Endpoints Overview
```
Public Routes:
  POST   /api/auth/login              - Login (no auth required)
  GET    /api/status                  - Health check (public)

Protected Routes (Require Bearer Token):
  POST   /api/auth/logout             - Logout
  GET    /api/auth/user               - Get current user info
  
System Info:
  GET    /api/sysinfo                 - Complete system info
  GET    /api/sysinfo/identity        - System identity
  GET    /api/sysinfo/cpu             - CPU info
  GET    /api/sysinfo/memory          - Memory info
  GET    /api/sysinfo/disk            - Disk info
  GET    /api/sysinfo/network         - Network info
  GET    /api/sysinfo/security        - Security info
  GET    /api/sysinfo/processes       - Process list
  
Audit:
  GET    /api/audit                   - All audit results
  GET    /api/audit?category=<cat>    - Filtered by category
  
Blogs:
  GET    /api/blogs                   - All blogs/advisories
  GET    /api/blogs?id=<id>           - Get specific blog
  GET    /api/blogs?category=<cat>    - Get blogs by category
  GET    /api/blogs/search?q=<query>  - Search blogs
```

## Installation & Setup

### Prerequisites
- Node.js 16+ and npm/yarn
- Go 1.19+ (for backend)
- Linux system (tested on Ubuntu 22.04)

### Frontend Setup

1. Install dependencies:
```bash
cd web
npm install
```

2. Start development server:
```bash
npm run dev
```
The frontend runs on `http://localhost:3000`

3. Build for production:
```bash
npm run build
```
This creates optimized files in `web/dist`

### Backend Setup

1. Ensure the API server is running on port 8080:
```bash
# In the main project root
make build-linux
./bin/solace-linux --api --port 8080
```

2. The API will be available at `http://localhost:8080`

3. The frontend dev server proxies API calls to the backend

### Connecting Frontend to Backend

The frontend is configured to proxy `/api` requests to `http://localhost:8080`:

In `vite.config.js`:
```javascript
server: {
  proxy: {
    '/api': {
      target: 'http://localhost:8080',
      changeOrigin: true,
    }
  }
}
```

## Authentication

### Default Credentials
- **Username**: `admin`
- **Password**: `password`

### How It Works

1. User submits login form with username/password
2. Frontend calls `POST /api/auth/login`
3. Backend validates credentials and returns Bearer token
4. Token is stored in localStorage
5. All subsequent API calls include `Authorization: Bearer <token>` header
6. Server validates token before processing requests

### Token Format
Current implementation uses Base64 encoding (production should use JWT):
```
Token = Base64("username:authenticated")
```

## Development

### Adding New API Endpoints

1. **Backend (Go)**:
   - Add handler method to `Server` struct
   - Register route in `Start()` method
   - Handle authentication if needed

2. **Frontend (React)**:
   - Add method to `apiClient` in `src/api/client.js`
   - Create page component in `src/pages/`
   - Add route in `App.jsx`
   - Add navigation link in `Navigation.jsx`

### Development Tips

- Hot module reloading (HMR) enabled in Vite
- API calls show in browser network tab/console
- Error messages displayed in UI
- Use browser DevTools for React debugging
- Check terminal output for API logs

## Frontend Build Output

Build creates optimized production files:
```bash
npm run build
# Output: web/dist/
```

Files include:
- `index.html` - Minified HTML entry
- `assets/index-*.js` - Optimized JavaScript (split)
- `assets/index-*.css` - Optimized CSS
- All minified and with cache-busting hashes

## Deployment

### Development Deployment
```bash
# Terminal 1: Backend
cd /home/ciph3r/Programs/Go/solace
make build-linux
./bin/solace-linux --api --port 8080

# Terminal 2: Frontend
cd web
npm install
npm run dev
```

Access at: `http://localhost:3000`

### Production Deployment

1. **Build Frontend**:
```bash
cd web
npm run build
```

2. **Set up backend** (if not already running):
```bash
make build-linux
./bin/solace-linux --api --port 8080
```

3. **Serve static files** (option 1 - via Go):
   - Copy `web/dist` to a known location
   - Serve via Go HTTP server

4. **Or serve separately** (option 2 - via nginx/apache):
```nginx
server {
    listen 80;
    server_name your-domain.com;
    root /path/to/web/dist;
    
    try_files $uri $uri/ /index.html;
    
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

## Integration Points

### With /internal/sysinfo
The system info endpoints currently return mock data. To integrate with the actual sysinfo package:

1. Import sysinfo package in `internal/api/server.go`
2. Replace mock data with actual calls to sysinfo collectors
3. Load system data on each `/api/sysinfo*` request

### With /internal/engine
Already integrated - audit results come directly from engine.EvaluateRules()

### With /rules directory
Rules are loaded by the engine and evaluated, results displayed in UI

## Configuration

### Environment Variables (Frontend)
Current config uses proxy to localhost:8080. For different backends:

Edit `vite.config.js`:
```javascript
proxy: {
  '/api': {
    target: 'http://your-api-server:port',
    changeOrigin: true,
  }
}
```

### Backend Port
Change API port by modifying main.go or command-line args

## TODO

1. Complete integration with `/internal/sysinfo` package
2. Implement database for persistent blog/config storage
3. Add user management and role-based access
4. Implement JWT authentication
5. Add audit logging for all API calls
6. Create admin settings page
7. Add report generation features
8. Implement real-time monitoring with WebSockets
9. Add data export functionality (PDF, CSV)
10. Create system health trend charts

### Security
Not our priority
- Default credentials should be changed in production
- Consider implementing JWT tokens instead of Base64
- Add rate limiting to API
- Use HTTPS in production
- Implement proper user management and roles
- Audit log API access
- Validate all user inputs
- Use secure password hashing