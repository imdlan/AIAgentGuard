import { useState, useEffect } from 'react';
import type { ProcessesResponse } from '../types';
import { apiClient } from '../api/client';
import './ProcessDetailsPanel.css';

export default function ProcessDetailsPanel() {
  const [data, setData] = useState<ProcessesResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadProcesses = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.getProcesses();
      setData(result);
    } catch (err) {
      setError('Failed to load process details');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadProcesses();
  }, []);

  return (
    <div className="process-details-panel">
      <div className="panel-header">
        <h3>🔍 Process Security Analysis</h3>
        <button onClick={loadProcesses} disabled={loading} className="refresh-button">
          {loading ? '⏳' : '🔄'}
        </button>
      </div>

      {error && <div className="error-message">❌ {error}</div>}

      {data && (
        <div className="process-details-content">
          <div className="summary-stats">
            <div className="stat">
              <span className="stat-label">Total Processes:</span>
              <span className="stat-value">{data.total}</span>
            </div>
            <div className="stat">
              <span className="stat-label">High Risk:</span>
              <span className="stat-value risk-high">{data.high_risk}</span>
            </div>
          </div>

          {data.processes.length > 0 ? (
            <div className="processes-table">
              <table>
                <thead>
                  <tr>
                    <th>PID</th>
                    <th>Name</th>
                    <th>User</th>
                    <th>Command</th>
                    <th>Risk</th>
                  </tr>
                </thead>
                <tbody>
                  {data.processes.map((proc, index) => (
                    <tr key={index} className={proc.risk_reason ? 'high-risk' : ''}>
                      <td>{proc.pid}</td>
                      <td>{proc.name}</td>
                      <td>{proc.user}</td>
                      <td className="command-line">{proc.command_line}</td>
                      <td>
                        {proc.risk_reason ? (
                          <span className="risk-badge high">⚠️ {proc.risk_reason}</span>
                        ) : (
                          <span className="risk-badge low">✅ Safe</span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="no-data">✅ No suspicious processes detected</div>
          )}
        </div>
      )}
    </div>
  );
}
