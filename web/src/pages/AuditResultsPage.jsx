import React, { useState, useEffect } from 'react'
import { apiClient } from '../api/client'
import { CheckCircle, AlertTriangle, XCircle, Filter, ExternalLink } from 'lucide-react'

const StatusBadge = ({ status }) => {
  const config = {
    pass: { color: 'badge-pass', icon: CheckCircle, label: 'Passed' },
    fail: { color: 'badge-fail', icon: XCircle, label: 'Failed' },
    warning: { color: 'badge-warning', icon: AlertTriangle, label: 'Warning' },
    skipped: { color: 'badge-skipped', icon: AlertTriangle, label: 'Skipped' },
  }

  const cfg = config[status] || config.skipped
  const Icon = cfg.icon

  return (
    <div className={`${cfg.color} inline-flex items-center gap-1`}>
      <Icon className="w-3 h-3" />
      {cfg.label}
    </div>
  )
}

const SeverityBadge = ({ severity }) => {
  const severityColors = {
    critical: 'bg-red-700 text-white',
    high: 'bg-red-600 text-white',
    medium: 'bg-orange-500 text-white',
    low: 'bg-yellow-500 text-white',
  }

  return (
    <span className={`px-2 py-1 rounded text-xs font-semibold ${severityColors[severity] || 'bg-gray-500 text-white'}`}>
      {severity?.toUpperCase()}
    </span>
  )
}

export default function AuditResultsPage() {
  const [results, setResults] = useState([])
  const [filteredResults, setFilteredResults] = useState([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)
  const [filters, setFilters] = useState({
    status: 'all',
    severity: 'all',
    category: 'all',
    searchTerm: '',
  })

  const [categories, setCategories] = useState([])

  useEffect(() => {
    fetchAuditResults()
  }, [])

  useEffect(() => {
    applyFilters()
  }, [results, filters])

  const fetchAuditResults = async () => {
    try {
      const response = await apiClient.getAuditResults()
      const auditResults = response.data || []
      setResults(auditResults)

      // Extract unique categories
      const uniqueCategories = [...new Set(auditResults.map(r => r.category).filter(Boolean))]
      setCategories(uniqueCategories)
    } catch (err) {
      setError('Failed to load audit results: ' + err.message)
    } finally {
      setIsLoading(false)
    }
  }

  const applyFilters = () => {
    let filtered = results

    if (filters.status !== 'all') {
      filtered = filtered.filter(r => r.status === filters.status)
    }

    if (filters.severity !== 'all') {
      filtered = filtered.filter(r => r.severity === filters.severity)
    }

    if (filters.category !== 'all') {
      filtered = filtered.filter(r => r.category === filters.category)
    }

    if (filters.searchTerm) {
      const term = filters.searchTerm.toLowerCase()
      filtered = filtered.filter(
        r =>
          r.rule_id?.toLowerCase().includes(term) ||
          r.title?.toLowerCase().includes(term) ||
          r.message?.toLowerCase().includes(term)
      )
    }

    setFilteredResults(filtered)
  }

  const handleFilterChange = (e) => {
    const { name, value } = e.target
    setFilters(prev => ({ ...prev, [name]: value }))
  }

  if (isLoading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-6 max-w-7xl mx-auto">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">{error}</p>
        </div>
      </div>
    )
  }

  const stats = {
    total: results.length,
    passed: results.filter(r => r.status === 'pass').length,
    failed: results.filter(r => r.status === 'fail').length,
    warning: results.filter(r => r.status === 'warning').length,
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Audit Results</h1>
        <p className="text-gray-600 mt-2">Security rule evaluation results</p>
      </div>

      {/* Summary Stats */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <div className="card">
          <p className="text-gray-600 text-sm">Total Rules</p>
          <p className="text-2xl font-bold text-gray-900">{stats.total}</p>
        </div>
        <div className="card">
          <p className="text-gray-600 text-sm">Passed</p>
          <p className="text-2xl font-bold text-green-600">{stats.passed}</p>
        </div>
        <div className="card">
          <p className="text-gray-600 text-sm">Failed</p>
          <p className="text-2xl font-bold text-red-600">{stats.failed}</p>
        </div>
        <div className="card">
          <p className="text-gray-600 text-sm">Warnings</p>
          <p className="text-2xl font-bold text-yellow-600">{stats.warning}</p>
        </div>
      </div>

      {/* Filters */}
      <div className="card mb-6">
        <div className="flex items-center gap-2 mb-4 pb-4 border-b">
          <Filter className="w-5 h-5 text-blue-600" />
          <h2 className="font-bold text-gray-900">Filters</h2>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              name="status"
              value={filters.status}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
            >
              <option value="all">All Status</option>
              <option value="pass">Passed</option>
              <option value="fail">Failed</option>
              <option value="warning">Warning</option>
              <option value="skipped">Skipped</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Severity</label>
            <select
              name="severity"
              value={filters.severity}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
            >
              <option value="all">All Severities</option>
              <option value="critical">Critical</option>
              <option value="high">High</option>
              <option value="medium">Medium</option>
              <option value="low">Low</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Category</label>
            <select
              name="category"
              value={filters.category}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
            >
              <option value="all">All Categories</option>
              {categories.map(cat => (
                <option key={cat} value={cat}>
                  {cat}
                </option>
              ))}
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">Search</label>
            <input
              type="text"
              name="searchTerm"
              value={filters.searchTerm}
              onChange={handleFilterChange}
              placeholder="Search rules..."
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 outline-none"
            />
          </div>
        </div>
      </div>

      {/* Results Table */}
      <div className="card overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="border-b bg-gray-50">
              <th className="text-left py-3 px-4 font-semibold text-gray-700">ID</th>
              <th className="text-left py-3 px-4 font-semibold text-gray-700">Title</th>
              <th className="text-left py-3 px-4 font-semibold text-gray-700">Status</th>
              <th className="text-left py-3 px-4 font-semibold text-gray-700">Severity</th>
              <th className="text-left py-3 px-4 font-semibold text-gray-700">Expected</th>
              <th className="text-left py-3 px-4 font-semibold text-gray-700">Current</th>
              <th className="text-left py-3 px-4 font-semibold text-gray-700">Action</th>
            </tr>
          </thead>
          <tbody>
            {filteredResults.length > 0 ? (
              filteredResults.map((result, idx) => (
                <tr key={idx} className="border-b hover:bg-gray-50">
                  <td className="py-3 px-4 font-mono text-xs text-gray-600">{result.rule_id}</td>
                  <td className="py-3 px-4">
                    <div className="max-w-xs">
                      <p className="font-medium text-gray-900">{result.title}</p>
                      {result.message && (
                        <p className="text-xs text-gray-600 mt-1">{result.message}</p>
                      )}
                    </div>
                  </td>
                  <td className="py-3 px-4">
                    <StatusBadge status={result.status} />
                  </td>
                  <td className="py-3 px-4">
                    <SeverityBadge severity={result.severity} />
                  </td>
                  <td className="py-3 px-4 font-mono text-xs">{result.expected_value}</td>
                  <td className="py-3 px-4 font-mono text-xs">{result.current_value || '-'}</td>
                  <td className="py-3 px-4">
                    {result.reference_link && (
                      <a
                        href={result.reference_link}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="text-blue-600 hover:underline inline-flex items-center gap-1"
                      >
                        <ExternalLink className="w-3 h-3" />
                        Learn
                      </a>
                    )}
                  </td>
                </tr>
              ))
            ) : (
              <tr>
                <td colSpan="7" className="py-8 px-4 text-center text-gray-600">
                  No results found
                </td>
              </tr>
            )}
          </tbody>
        </table>
      </div>

      {/* Refresh Button */}
      <div className="mt-8 flex justify-center">
        <button onClick={fetchAuditResults} className="btn-primary">
          Refresh Results
        </button>
      </div>
    </div>
  )
}
