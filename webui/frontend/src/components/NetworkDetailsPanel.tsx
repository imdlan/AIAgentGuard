import { useState, useEffect } from 'react';
import type { NetworkResponse } from '../types';
import { apiClient } from '../api/client';
import './NetworkDetailsPanel.css';

export default function NetworkDetailsPanel() {
  const [data, setData] = useState<NetworkResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadNetwork = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.getNetwork();
      setData(result);
    } catch (err) {
      setError('Failed to load network details');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadNetwork();
  }, []);

  return (
    <div className="network-details-panel">
      <div className="panel-header">
        <h3>🌐 Network Security Analysis</h3>
        <button onClick={loadNetwork} disabled={loading} className="refresh-button">
          {loading ? '⏳' : '🔄'}
        </button>
      </div>

      {error && <div className="error-message">❌ {error}</div>}

      {data && (
        <div className="network-details-content">
          <div className="summary-stats">
            <div className="stat">
              <span className="stat-label">Open Ports:</span>
              <span className="stat-value">{data.total_ports}</span>
            </div>
            <div className="stat">
              <span className="stat-label">Active Connections:</span>
              <span className="stat-value">{data.total_connections}</span>
            </div>
          </div>

          {data.open_ports.length > 0 && (
            <div className="section">
              <h4>🔓 Open Ports</h4>
              <div className="ports-table">
                <table>
                  <thead>
                    <tr>
                      <th>Port</th>
                      <th>Protocol</th>
                      <th>Service</th>
                      <th>Risk</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.open_ports.map((port, index) => (
                      <tr key={index}>
                        <td>{port.port}</td>
                        <td>{port.protocol}</td>
                        <td>{port.service}</td>
                        <td>
                          {port.risk_reason ? (
                            <span className="risk-badge high">⚠️ {port.risk_reason}</span>
                          ) : (
                            <span className="risk-badge low">✅ Safe</span>
                          )}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {data.active_connections.length > 0 && (
            <div className="section">
              <h4>🔗 Active Connections</h4>
              <div className="connections-table">
                <table>
                  <thead>
                    <tr>
                      <th>Local Address</th>
                      <th>Remote Address</th>
                      <th>State</th>
                      <th>Protocol</th>
                    </tr>
                  </thead>
                  <tbody>
                    {data.active_connections.map((conn, index) => (
                      <tr key={index}>
                        <td>{conn.local_address}:{conn.local_port}</td>
                        <td>{conn.remote_address}:{conn.remote_port}</td>
                        <td>
                          <span className={`state-badge ${conn.state.toLowerCase()}`}>
                            {conn.state}
                          </span>
                        </td>
                        <td>{conn.protocol}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            </div>
          )}

          {data.open_ports.length === 0 && data.active_connections.length === 0 && (
            <div className="no-data">✅ No suspicious network activity detected</div>
          )}
        </div>
      )}
    </div>
  );
}
