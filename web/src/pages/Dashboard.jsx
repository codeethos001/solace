import React, { useState, useEffect } from 'react'
import { apiClient } from '../api/client'
import { Activity, AlertTriangle, CheckCircle, TrendingUp, Zap } from 'lucide-react'

const StatCard = ({ icon: Icon, title, value, unit = '', color = 'blue' }) => {
  const colorClass = {
    blue: 'bg-blue-50 text-blue-600',
    green: 'bg-green-50 text-green-600',
    red: 'bg-red-50 text-red-600',
    yellow: 'bg-yellow-50 text-yellow-600',
  }[color]

  return (
    <div className="card">
      <div className="flex items-center justify-between">
        <div>
          <p className="text-gray-600 text-sm font-medium">{title}</p>
          <p className="text-3xl font-bold text-gray-900 mt-2">
            {value}
            <span className="text-lg text-gray-500 ml-1">{unit}</span>
          </p>
        </div>
        <div className={`p-3 rounded-full ${colorClass}`}>
          <Icon className="w-6 h-6" />
        </div>
      </div>
    </div>
  )
}

export default function Dashboard() {
  const [stats, setStats] = useState(null)
  const [auditStats, setAuditStats] = useState(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    fetchDashboardData()
  }, [])

  const fetchDashboardData = async () => {
    try {
      // Fetch system info and audit results
      const [sysInfoRes, auditRes] = await Promise.all([
        apiClient.getSystemInfo().catch(() => null),
        apiClient.getAuditResults(),
      ])

      if (sysInfoRes?.data) {
        const info = sysInfoRes.data
        setStats({
          hostname: info.hostname || 'N/A',
          os: info.os || 'Linux',
          uptime: info.uptime_seconds ? Math.floor(info.uptime_seconds / 86400) : 0,
          users: info.users_count || 0,
        })
      }

      // Calculate audit stats from results
      const results = auditRes.data || []
      const passed = results.filter(r => r.status === 'pass').length
      const failed = results.filter(r => r.status === 'fail').length
      const warning = results.filter(r => r.status === 'warning').length

      setAuditStats({
        total: results.length,
        passed,
        failed,
        warning,
        passPercentage: results.length > 0 ? Math.round((passed / results.length) * 100) : 0,
      })
    } catch (error) {
      console.error('Error fetching dashboard data:', error)
    } finally {
      setIsLoading(false)
    }
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="text-gray-600 mt-2">System overview and security status</p>
      </div>

      {/* System Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {stats && (
          <>
            <StatCard
              icon={Activity}
              title="Hostname"
              value={stats.hostname}
              color="blue"
            />
            <StatCard
              icon={Zap}
              title="Uptime"
              value={stats.uptime}
              unit="days"
              color="green"
            />
            <StatCard
              icon={Activity}
              title="Active Users"
              value={stats.users}
              color="blue"
            />
            <StatCard
              icon={Activity}
              title="OS"
              value={stats.os}
              color="blue"
            />
          </>
        )}
      </div>

      {/* Audit Status Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <div className="card">
          <h2 className="text-lg font-bold text-gray-900 mb-4">Security Posture</h2>
          <div className="space-y-4">
            <div>
              <div className="flex justify-between items-center mb-2">
                <span className="text-gray-700 font-medium">Compliance Score</span>
                <span className="text-2xl font-bold text-green-600">
                  {auditStats?.passPercentage || 0}%
                </span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div
                  className="bg-green-600 h-2 rounded-full transition-all"
                  style={{ width: `${auditStats?.passPercentage || 0}%` }}
                ></div>
              </div>
            </div>
            <div className="text-sm text-gray-600 mt-4">
              <p>
                <span className="font-semibold">Total Rules:</span> {auditStats?.total}
              </p>
            </div>
          </div>
        </div>

        <div className="card">
          <h2 className="text-lg font-bold text-gray-900 mb-4">Audit Results</h2>
          <div className="space-y-3">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <CheckCircle className="w-5 h-5 text-green-600" />
                <span className="text-gray-700">Passed</span>
              </div>
              <span className="font-bold text-green-600">{auditStats?.passed}</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <AlertTriangle className="w-5 h-5 text-yellow-600" />
                <span className="text-gray-700">Warning</span>
              </div>
              <span className="font-bold text-yellow-600">{auditStats?.warning}</span>
            </div>
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <AlertTriangle className="w-5 h-5 text-red-600" />
                <span className="text-gray-700">Failed</span>
              </div>
              <span className="font-bold text-red-600">{auditStats?.failed}</span>
            </div>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="card">
        <h2 className="text-lg font-bold text-gray-900 mb-4">Quick Actions</h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <a
            href="/system-info"
            className="p-4 border-2 border-blue-200 rounded-lg hover:border-blue-600 hover:bg-blue-50 transition"
          >
            <h3 className="font-semibold text-gray-900 mb-1">View System Info</h3>
            <p className="text-sm text-gray-600">CPU, memory, disk, network</p>
          </a>
          <a
            href="/audit-results"
            className="p-4 border-2 border-orange-200 rounded-lg hover:border-orange-600 hover:bg-orange-50 transition"
          >
            <h3 className="font-semibold text-gray-900 mb-1">Audit Results</h3>
            <p className="text-sm text-gray-600">All security checks</p>
          </a>
          <a
            href="/hardening-guides"
            className="p-4 border-2 border-green-200 rounded-lg hover:border-green-600 hover:bg-green-50 transition"
          >
            <h3 className="font-semibold text-gray-900 mb-1">Hardening Guides</h3>
            <p className="text-sm text-gray-600">Best practices & advisories</p>
          </a>
        </div>
      </div>
    </div>
  )
}
