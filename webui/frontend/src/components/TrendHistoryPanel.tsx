import { useState, useEffect } from 'react';
import type { TrendHistoryResponse } from '../types';
import { apiClient } from '../api/client';
import './TrendHistoryPanel.css';

export default function TrendHistoryPanel() {
  const [data, setData] = useState<TrendHistoryResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [days, setDays] = useState(7);

  const loadTrends = async () => {
    setLoading(true);
    setError(null);
    try {
      const result = await apiClient.getTrendHistory(days);
      setData(result);
    } catch (err) {
      setError('Failed to load trend history');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTrends();
  }, [days]);

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'LOW': return '#10b981';
      case 'MEDIUM': return '#f59e0b';
      case 'HIGH': return '#ef4444';
      case 'CRITICAL': return '#7c3aed';
      default: return '#6b7280';
    }
  };

  const getTrendIcon = (current: string, previous: string) => {
    const scoreMap: Record<string, number> = {
      'LOW': 25,
      'MEDIUM': 50,
      'HIGH': 75,
      'CRITICAL': 100,
    };
    
    const currentScore = scoreMap[current] || 0;
    const previousScore = scoreMap[previous] || 0;
    
    if (currentScore < previousScore) return '📈';
    if (currentScore > previousScore) return '📉';
    return '➡️';
  };

  return (
    <div className="trend-history-panel">
      <div className="panel-header">
        <h3>📈 Security Trend History</h3>
        <div className="controls">
          <select
            value={days}
            onChange={(e) => setDays(Number(e.target.value))}
            className="days-selector"
          >
            <option value={7}>Last 7 days</option>
            <option value={14}>Last 14 days</option>
            <option value={30}>Last 30 days</option>
          </select>
          <button onClick={loadTrends} disabled={loading} className="refresh-button">
            {loading ? '⏳' : '🔄'}
          </button>
        </div>
      </div>

      {error && <div className="error-message">❌ {error}</div>}

      {data && data.trend_data.length > 0 && (
        <div className="trend-content">
          <div className="period-info">
            <span className="period-label">Period:</span>
            <span className="period-value">{data.period}</span>
          </div>

          <div className="trend-table">
            <table>
              <thead>
                <tr>
                  <th>Date</th>
                  <th>Overall</th>
                  <th>Filesystem</th>
                  <th>Shell</th>
                  <th>Network</th>
                  <th>Secrets</th>
                  <th>Trend</th>
                </tr>
              </thead>
              <tbody>
                {data.trend_data.map((item, index) => {
                  const prevItem = index < data.trend_data.length - 1 ? data.trend_data[index + 1] : null;
                  return (
                    <tr key={index}>
                      <td>{item.date}</td>
                      <td>
                        <span
                          className="risk-badge"
                          style={{ backgroundColor: getRiskColor(item.overall), color: 'white' }}
                        >
                          {item.overall}
                        </span>
                      </td>
                      <td>
                        {item.filesystem && (
                          <span
                            className="risk-badge small"
                            style={{ backgroundColor: getRiskColor(item.filesystem), color: 'white' }}
                          >
                            {item.filesystem}
                          </span>
                        )}
                      </td>
                      <td>
                        {item.shell && (
                          <span
                            className="risk-badge small"
                            style={{ backgroundColor: getRiskColor(item.shell), color: 'white' }}
                          >
                            {item.shell}
                          </span>
                        )}
                      </td>
                      <td>
                        {item.network && (
                          <span
                            className="risk-badge small"
                            style={{ backgroundColor: getRiskColor(item.network), color: 'white' }}
                          >
                            {item.network}
                          </span>
                        )}
                      </td>
                      <td>
                        {item.secrets && (
                          <span
                            className="risk-badge small"
                            style={{ backgroundColor: getRiskColor(item.secrets), color: 'white' }}
                          >
                            {item.secrets}
                          </span>
                        )}
                      </td>
                      <td>
                        {prevItem ? (
                          <span className="trend-icon">
                            {getTrendIcon(item.overall, prevItem.overall)}
                          </span>
                        ) : (
                          <span className="trend-icon">—</span>
                        )}
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}

      {data && data.trend_data.length === 0 && (
        <div className="no-data">No trend data available</div>
      )}
    </div>
  );
}
