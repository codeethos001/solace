import React, { useState } from 'react'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'
import {
  Menu,
  X,
  Home,
  Shield,
  AlertCircle,
  BookOpen,
  LogOut,
  Settings,
} from 'lucide-react'

export default function Navigation() {
  const [isOpen, setIsOpen] = useState(false)
  const { logout, user } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()

  const handleLogout = async () => {
    await logout()
    navigate('/login')
  }

  const isActive = (path) => {
    return location.pathname === path
      ? 'bg-blue-700 text-white'
      : 'text-gray-300 hover:bg-blue-700/50'
  }

  const navItems = [
    { path: '/dashboard', label: 'Dashboard', icon: Home },
    { path: '/system-info', label: 'System Info', icon: Shield },
    { path: '/audit-results', label: 'Audit Results', icon: AlertCircle },
    { path: '/hardening-guides', label: 'Hardening Guides', icon: BookOpen },
  ]

  return (
    <nav className="bg-blue-600 text-white shadow-lg">
      <div className="max-w-7xl mx-auto px-4">
        <div className="flex justify-between items-center h-16">
          {/* Logo */}
          <Link to="/dashboard" className="flex items-center gap-2 font-bold text-xl">
            <Shield className="w-6 h-6" />
            Solace
          </Link>

          {/* Desktop Navigation */}
          <div className="hidden md:flex items-center gap-1">
            {navItems.map(({ path, label, icon: Icon }) => (
              <Link
                key={path}
                to={path}
                className={`flex items-center gap-2 px-3 py-2 rounded transition ${isActive(path)}`}
              >
                <Icon className="w-4 h-4" />
                {label}
              </Link>
            ))}
          </div>

          {/* User Info & Logout */}
          <div className="hidden md:flex items-center gap-4">
            <span className="text-sm">{user?.username}</span>
            <button
              onClick={handleLogout}
              className="flex items-center gap-2 px-3 py-2 rounded hover:bg-blue-700 transition"
            >
              <LogOut className="w-4 h-4" />
              Logout
            </button>
          </div>

          {/* Mobile Menu Button */}
          <button
            onClick={() => setIsOpen(!isOpen)}
            className="md:hidden p-2 hover:bg-blue-700 rounded"
          >
            {isOpen ? <X className="w-6 h-6" /> : <Menu className="w-6 h-6" />}
          </button>
        </div>

        {/* Mobile Navigation */}
        {isOpen && (
          <div className="md:hidden bg-blue-700 py-2">
            {navItems.map(({ path, label, icon: Icon }) => (
              <Link
                key={path}
                to={path}
                className="flex items-center gap-2 px-4 py-2 hover:bg-blue-800 transition text-gray-100"
                onClick={() => setIsOpen(false)}
              >
                <Icon className="w-4 h-4" />
                {label}
              </Link>
            ))}
            <div className="border-t border-blue-500 mt-2 pt-2 px-4 py-2">
              <p className="text-sm text-gray-200 mb-2">{user?.username}</p>
              <button
                onClick={handleLogout}
                className="flex items-center gap-2 text-gray-100 hover:text-white transition"
              >
                <LogOut className="w-4 h-4" />
                Logout
              </button>
            </div>
          </div>
        )}
      </div>
    </nav>
  )
}
