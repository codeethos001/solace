import React, { useState, useEffect } from 'react'
import { apiClient } from '../api/client'
import { Cpu, HardDrive, Network, Lock, Users, Zap } from 'lucide-react'

const InfoSection = ({ icon: Icon, title, data }) => (
  <div className="card mb-6">
    <div className="flex items-center gap-2 mb-4 border-b pb-4">
      <Icon className="w-5 h-5 text-blue-600" />
      <h2 className="text-lg font-bold text-gray-900">{title}</h2>
    </div>
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {Object.entries(data).map(([key, value]) => (
        <div key={key} className="bg-gray-50 p-3 rounded">
          <p className="text-xs font-medium text-gray-600 uppercase tracking-wide">
            {key.replace(/_/g, ' ')}
          </p>
          <p className="text-base font-semibold text-gray-900 mt-1">
            {typeof value === 'boolean' ? (value ? 'Yes' : 'No') : String(value)}
          </p>
        </div>
      ))}
    </div>
  </div>
)

export default function SystemInfoPage() {
  const [systemInfo, setSystemInfo] = useState(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetchSystemInfo()
  }, [])

  const fetchSystemInfo = async () => {
    try {
      const response = await apiClient.getSystemInfo()
      setSystemInfo(response.data)
    } catch (err) {
      setError('Failed to load system information: ' + err.message)
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

  if (error) {
    return (
      <div className="p-6 max-w-7xl mx-auto">
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-red-800">{error}</p>
        </div>
      </div>
    )
  }

  if (!systemInfo) {
    return (
      <div className="p-6 max-w-7xl mx-auto">
        <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
          <p className="text-yellow-800">No system information available</p>
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-7xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">System Information</h1>
        <p className="text-gray-600 mt-2">Complete hardware and OS details</p>
      </div>

      {/* System Identity */}
      {systemInfo.identity && (
        <InfoSection icon={Zap} title="System Identity" data={systemInfo.identity} />
      )}

      {/* CPU Info */}
      {systemInfo.cpu && (
        <InfoSection icon={Cpu} title="CPU Information" data={systemInfo.cpu} />
      )}

      {/* Memory Info */}
      {systemInfo.memory && (
        <InfoSection icon={HardDrive} title="Memory Information" data={systemInfo.memory} />
      )}

      {/* Disk Info */}
      {systemInfo.disk && Array.isArray(systemInfo.disk) ? (
        <div className="card mb-6">
          <div className="flex items-center gap-2 mb-4 border-b pb-4">
            <HardDrive className="w-5 h-5 text-blue-600" />
            <h2 className="text-lg font-bold text-gray-900">Disk Information</h2>
          </div>
          <div className="overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b">
                  <th className="text-left py-2 px-2 font-semibold">Device</th>
                  <th className="text-left py-2 px-2 font-semibold">Mount Point</th>
                  <th className="text-left py-2 px-2 font-semibold">Total</th>
                  <th className="text-left py-2 px-2 font-semibold">Used</th>
                  <th className="text-left py-2 px-2 font-semibold">Free</th>
                </tr>
              </thead>
              <tbody>
                {systemInfo.disk.map((disk, idx) => (
                  <tr key={idx} className="border-b hover:bg-gray-50">
                    <td className="py-2 px-2">{disk.device}</td>
                    <td className="py-2 px-2">{disk.mount_point}</td>
                    <td className="py-2 px-2">{disk.total}</td>
                    <td className="py-2 px-2">{disk.used}</td>
                    <td className="py-2 px-2">{disk.free}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      ) : null}

      {/* Network Info */}
      {systemInfo.network && (
        <InfoSection icon={Network} title="Network Information" data={systemInfo.network} />
      )}

      {/* Security Info */}
      {systemInfo.security && (
        <InfoSection icon={Lock} title="Security Configuration" data={systemInfo.security} />
      )}

      {/* Refresh Button */}
      <div className="mt-8 flex justify-center">
        <button
          onClick={fetchSystemInfo}
          className="btn-primary"
        >
          Refresh Information
        </button>
      </div>
    </div>
  )
}
