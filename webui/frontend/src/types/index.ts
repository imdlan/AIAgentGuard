// API Types
export interface ScanResult {
  id: string;
  timestamp: string;
  duration: number;
  results: PermissionResult;
  overall: RiskLevel;
  details: RiskDetail[];
}

export interface PermissionResult {
  filesystem: RiskLevel;
  shell: RiskLevel;
  network: RiskLevel;
  secrets: RiskLevel;
  filecontent: RiskLevel;
  dependencies: RiskLevel;
  npm_deps: RiskLevel;
  pip_deps: RiskLevel;
  cargo_deps: RiskLevel;
}

export interface RiskDetail {
  type: RiskLevel;
  category: string;
  description: string;
  path?: string;
	remediation?: RemediationInfo[];
}

export interface RemediationInfo {
	command: string;
	description?: string;
	priority?: string;
}

export interface SuspiciousProcess {
	pid: number;
	name: string;
	command_line: string;
	user: string;
	risk_reason: string;
}

export interface PortDetail {
	port: number;
	protocol: string;
	service: string;
	risk_reason: string;
}

export interface ConnectionDetail {
	local_address: string;
	local_port: number;
	remote_address: string;
	remote_port: number;
	state: string;
	protocol: string;
}

export interface ProcessesResponse {
	processes: SuspiciousProcess[];
	total: number;
	high_risk: number;
}

export interface NetworkResponse {
	open_ports: PortDetail[];
	active_connections: ConnectionDetail[];
	total_ports: number;
	total_connections: number;
}

export interface FixRequest {
	dry_run: boolean;
	auto: boolean;
	category?: string;
}

export interface FixResponse {
	success: boolean;
	message: string;
	fixed?: string[];
	failed?: string[];
	skipped?: string[];
}

export interface TrendHistoryResponse {
	trend_data: TrendDataItem[];
	period: string;
}

export interface TrendDataItem {
	date: string;
	overall: string;
	filesystem?: string;
	shell?: string;
	network?: string;
	secrets?: string;
}
export type RiskLevel = 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';

export interface SystemStatus {
  version: string;
  status: string;
  uptime: string;
  scanners: Record<string, string>;
}

export interface TrendData {
  date: string;
  low: number;
  medium: number;
  high: number;
  critical: number;
}

export interface MetricsData {
  timestamp: string;
  scan_total: number;
  scan_rate: number;
  duration_avg: number;
}

export interface VulnerabilityMetrics {
  timestamp: string;
  vulnerabilities: {
    critical: number;
    high: number;
    medium: number;
    low: number;
  };
  by_severity: Record<string, number>;
  by_language: Record<string, number>;
}

export interface DurationMetrics {
  timestamp: string;
  duration: Record<string, number>;
  quantiles: {
    p50: number;
    p95: number;
    p99: number;
  };
}

export interface ChartData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    borderColor?: string;
    backgroundColor?: string;
  }[];
}
