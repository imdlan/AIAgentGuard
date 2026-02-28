import React, { useState, useEffect } from 'react';
import { apiClient } from '../api/client';
import type { MetricsData, VulnerabilityMetrics, DurationMetrics } from '../types';
import './MetricsPanel.css';

interface MetricsPanelProps {
  refreshInterval?: number;
}

export const MetricsPanel: React.FC<MetricsPanelProps> = ({ refreshInterval = 30000 }) => {
  const [scanMetrics, setScanMetrics] = useState<MetricsData | null>(null);
  const [vulnMetrics, setVulnMetrics] = useState<VulnerabilityMetrics | null>(null);
  const [durationMetrics, setDurationMetrics] = useState<DurationMetrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchMetrics = async () => {
      try {
        setLoading(true);
        const [scanRate, vulns, duration] = await Promise.all([
          apiClient.getScanRateMetrics(),
          apiClient.getVulnerabilityMetrics(),
          apiClient.getDurationMetrics(),
        ]);
        
        setScanMetrics(scanRate);
        setVulnMetrics(vulns);
        setDurationMetrics(duration);
        setError(null);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch metrics');
      } finally {
        setLoading(false);
      }
    };

    fetchMetrics();
    const interval = setInterval(fetchMetrics, refreshInterval);
    return () => clearInterval(interval);
  }, [refreshInterval]);

  if (loading) {
    return <div className="metrics-panel loading">Loading metrics...</div>;
  }

  if (error) {
    return <div className="metrics-panel error">Error: {error}</div>;
  }

  return (
    <div className="metrics-panel">
      <div className="metrics-header">
        <h2>📊 Security Metrics</h2>
        <span className="last-update">Last updated: {new Date().toLocaleTimeString()}</span>
      </div>

      {/* Scan Metrics */}
      <div className="metrics-section">
        <h3>Scan Statistics</h3>
        <div className="metrics-cards">
          <MetricCard
            title="Total Scans"
            value={scanMetrics?.scan_total || 0}
            icon="🔍"
            trend={scanMetrics?.scan_rate || 0}
            trendLabel="scans/sec"
          />
          <MetricCard
            title="Avg Duration"
            value={`${scanMetrics?.duration_avg || 0}s`}
            icon="⏱️"
          />
          <MetricCard
            title="Scan Rate"
            value={`${scanMetrics?.scan_rate?.toFixed(2) || 0}`}
            icon="📈"
            trendLabel="scans/sec"
          />
        </div>
      </div>

      {/* Vulnerability Metrics */}
      <div className="metrics-section">
        <h3>Vulnerability Overview</h3>
        <div className="vulnerability-grid">
          <VulnerabilityCard
            severity="critical"
            count={vulnMetrics?.by_severity?.critical || 0}
            total={Object.values(vulnMetrics?.by_severity || {}).reduce((a, b) => a + b, 0)}
          />
          <VulnerabilityCard
            severity="high"
            count={vulnMetrics?.by_severity?.high || 0}
            total={Object.values(vulnMetrics?.by_severity || {}).reduce((a, b) => a + b, 0)}
          />
          <VulnerabilityCard
            severity="medium"
            count={vulnMetrics?.by_severity?.medium || 0}
            total={Object.values(vulnMetrics?.by_severity || {}).reduce((a, b) => a + b, 0)}
          />
          <VulnerabilityCard
            severity="low"
            count={vulnMetrics?.by_severity?.low || 0}
            total={Object.values(vulnMetrics?.by_severity || {}).reduce((a, b) => a + b, 0)}
          />
        </div>
      </div>

      {/* Language Breakdown */}
      {vulnMetrics && Object.keys(vulnMetrics.by_language || {}).length > 0 && (
        <div className="metrics-section">
          <h3>Vulnerabilities by Language</h3>
          <div className="language-bars">
            {Object.entries(vulnMetrics.by_language || {}).map(([lang, count]) => (
              <LanguageBar key={lang} language={lang} count={count} />
            ))}
          </div>
        </div>
      )}

      {/* Duration Metrics */}
      {durationMetrics && (
        <div className="metrics-section">
          <h3>Scan Duration</h3>
          <div className="duration-stats">
            <div className="duration-stat">
              <span className="label">Full Scan</span>
              <span className="value">{durationMetrics.duration?.full?.toFixed(2)}s</span>
            </div>
            <div className="duration-stat">
              <span className="label">p95</span>
              <span className="value">{durationMetrics.quantiles?.p95?.toFixed(2)}s</span>
            </div>
            <div className="duration-stat">
              <span className="label">p99</span>
              <span className="value">{durationMetrics.quantiles?.p99?.toFixed(2)}s</span>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

interface MetricCardProps {
  title: string;
  value: string | number;
  icon: string;
  trend?: number;
  trendLabel?: string;
}

const MetricCard: React.FC<MetricCardProps> = ({ title, value, icon, trend, trendLabel }) => (
  <div className="metric-card">
    <div className="metric-icon">{icon}</div>
    <div className="metric-content">
      <div className="metric-title">{title}</div>
      <div className="metric-value">{value}</div>
      {trend !== undefined && (
        <div className="metric-trend">
          <span className="trend-value">{trend > 0 ? '↑' : '↓'} {Math.abs(trend)}</span>
          {trendLabel && <span className="trend-label">{trendLabel}</span>}
        </div>
      )}
    </div>
  </div>
);

interface VulnerabilityCardProps {
  severity: 'critical' | 'high' | 'medium' | 'low';
  count: number;
  total: number;
}

const VulnerabilityCard: React.FC<VulnerabilityCardProps> = ({ severity, count, total }) => {
  const percentage = total > 0 ? (count / total) * 100 : 0;
  
  return (
    <div className={`vuln-card vuln-${severity}`}>
      <div className="vuln-count">{count}</div>
      <div className="vuln-label">{severity.toUpperCase()}</div>
      <div className="vuln-bar">
        <div className="vuln-fill" style={{ width: `${percentage}%` }} />
      </div>
      <div className="vuln-percentage">{percentage.toFixed(1)}%</div>
    </div>
  );
};

interface LanguageBarProps {
  language: string;
  count: number;
}

const LanguageBar: React.FC<LanguageBarProps> = ({ language, count }) => {
  const maxCount = 10; // For scaling
  const percentage = Math.min((count / maxCount) * 100, 100);
  
  return (
    <div className="language-bar">
      <span className="lang-name">{language}</span>
      <div className="bar-container">
        <div className="bar-fill" style={{ width: `${percentage}%` }} />
      </div>
      <span className="lang-count">{count}</span>
    </div>
  );
};

export default MetricsPanel;
