import { useState, useEffect } from 'react';
import type { ScanResult, SystemStatus } from '../types';
import { apiClient } from '../api/client';
import MetricsPanel from './MetricsPanel';
import ProcessDetailsPanel from './ProcessDetailsPanel';
import NetworkDetailsPanel from './NetworkDetailsPanel';
import FixWizardPanel from './FixWizardPanel';
import TrendHistoryPanel from './TrendHistoryPanel';
import './Dashboard.css';
import type { ScanResult, SystemStatus } from '../types';
import { apiClient } from '../api/client';
import MetricsPanel from './MetricsPanel';
import './Dashboard.css';

export default function Dashboard() {
  const [scanResult, setScanResult] = useState<ScanResult | null>(null);
  const [status, setStatus] = useState<SystemStatus | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadSystemStatus();
  }, []);

  const loadSystemStatus = async () => {
    try {
      setStatus(await apiClient.getStatus());
    } catch (err) {
      setError('Failed to load system status');
    }
  };

  const handleScan = async () => {
    setLoading(true);
    setError(null);
    
    try {
      const result = await apiClient.scan();
      setScanResult(result);
    } catch (err) {
      setError('Scan failed');
    } finally {
      setLoading(false);
    }
  };

  const getRiskColor = (level: string) => {
    switch (level) {
      case 'LOW': return '#10b981';
      case 'MEDIUM': return '#f59e0b';
      case 'HIGH': return '#ef4444';
      case 'CRITICAL': return '#7c3aed';
      default: return '#6b7280';
    }
  };

  const getRiskIcon = (level: string) => {
    switch (level) {
      case 'LOW': return '✅';
      case 'MEDIUM': return '⚠️';
      case 'HIGH': return '🔶';
      case 'CRITICAL': return '🛑';
      default: return '❓';
    }
  };

  return (
    <div className="dashboard">
      <header className="dashboard-header">
        <h1>🛡️ AI AgentGuard Dashboard</h1>
        {status && <span className="version">{status.version}</span>}
      </header>

		<MetricsPanel refreshInterval={30000} />

		<div className="dashboard-sections">
			<ProcessDetailsPanel />
			<NetworkDetailsPanel />
			<FixWizardPanel />
			<TrendHistoryPanel />
		</div>

      <div className="dashboard-controls">
        <button 
          onClick={handleScan} 
          disabled={loading}
          className="scan-button"
        >
          {loading ? '⏳ Scanning...' : '🔍 Run Security Scan'}
        </button>
      </div>

      {error && (
        <div className="error-message">
          ❌ {error}
        </div>
      )}

      {scanResult && (
        <div className="scan-results">
          <div className="results-header">
            <h2>Scan Results</h2>
            <span className="timestamp">
              {new Date(scanResult.timestamp).toLocaleString()}
            </span>
          </div>

          <div className="overall-risk">
            <h3>Overall Risk: {getRiskIcon(scanResult.overall)} {scanResult.overall}</h3>
          </div>

          <div className="permission-results">
            <h3>Permission Breakdown</h3>
            <div className="results-grid">
              <ResultCard title="Filesystem" risk={scanResult.results.filesystem} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="Shell" risk={scanResult.results.shell} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="Network" risk={scanResult.results.network} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="Secrets" risk={scanResult.results.secrets} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="File Content" risk={scanResult.results.filecontent} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="Dependencies (Go)" risk={scanResult.results.dependencies} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="NPM Dependencies" risk={scanResult.results.npm_deps || 'LOW'} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="Pip Dependencies" risk={scanResult.results.pip_deps || 'LOW'} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
              <ResultCard title="Cargo Dependencies" risk={scanResult.results.cargo_deps || 'LOW'} getRiskColor={getRiskColor} getRiskIcon={getRiskIcon} />
            </div>
          </div>

          {scanResult.details.length > 0 && (
            <div className="details-section">
              <h3>Detailed Findings</h3>
              <div className="details-list">
                {scanResult.details.map((detail, index) => (
                  <div key={index} className="detail-item" style={{ borderLeftColor: getRiskColor(detail.type) }}>
                    <span className="detail-risk">{getRiskIcon(detail.type)} {detail.type}</span>
                    <span className="detail-category">[{detail.category}]</span>
                    <span className="detail-description">{detail.description}</span>
                    {detail.path && <span className="detail-path">{detail.path}</span>}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

interface ResultCardProps {
  title: string;
  risk: string;
  getRiskColor: (level: string) => string;
  getRiskIcon: (level: string) => string;
}

const ResultCard: React.FC<ResultCardProps> = ({ title, risk, getRiskColor, getRiskIcon }) => (
  <div className="result-card" style={{ borderColor: getRiskColor(risk) }}>
    <h4>{title}</h4>
    <p className="risk-level" style={{ color: getRiskColor(risk) }}>
      {getRiskIcon(risk)} {risk}
    </p>
  </div>
);

