import React from 'react'
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext'
import ProtectedRoute from './components/ProtectedRoute'
import Navigation from './components/Navigation'
import LoginPage from './pages/LoginPage'
import Dashboard from './pages/Dashboard'
import SystemInfoPage from './pages/SystemInfoPage'
import AuditResultsPage from './pages/AuditResultsPage'
import HardeningGuidesPage from './pages/HardeningGuidesPage'

export default function App() {
  return (
    <Router>
      <AuthProvider>
        <Routes>
          {/* Public Routes */}
          <Route path="/login" element={<LoginPage />} />

          {/* Protected Routes */}
          <Route
            path="/*"
            element={
              <ProtectedRoute>
                <div className="min-h-screen bg-gray-100">
                  <Navigation />
                  <Routes>
                    <Route path="/dashboard" element={<Dashboard />} />
                    <Route path="/system-info" element={<SystemInfoPage />} />
                    <Route path="/audit-results" element={<AuditResultsPage />} />
                    <Route path="/hardening-guides" element={<HardeningGuidesPage />} />
                    <Route path="/" element={<Navigate to="/dashboard" replace />} />
                    <Route path="*" element={<Navigate to="/dashboard" replace />} />
                  </Routes>
                </div>
              </ProtectedRoute>
            }
          />
        </Routes>
      </AuthProvider>
    </Router>
  )
}
