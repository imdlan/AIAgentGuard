// API Client for communicating with backend
import type { ScanResult, SystemStatus, MetricsData, VulnerabilityMetrics, DurationMetrics, ProcessesResponse, NetworkResponse, FixRequest, FixResponse, TrendHistoryResponse } from '../types';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl;
  }

  async scan(): Promise<ScanResult> {
    const response = await fetch(`${this.baseUrl}/api/v1/scan`);
    if (!response.ok) {
      throw new Error('Scan failed');
    }
    return await response.json();
  }

  async scanWithOptions(categories: string[], options: any): Promise<ScanResult> {
    const response = await fetch(`${this.baseUrl}/api/v1/scan`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ categories, options }),
    });
    if (!response.ok) {
      throw new Error('Scan failed');
    }
    return await response.json();
  }

  async getScanResult(id: string): Promise<ScanResult> {
    const response = await fetch(`${this.baseUrl}/api/v1/scan/${id}`);
    if (!response.ok) {
      throw new Error('Failed to get scan result');
    }
    return await response.json();
  }

  async getHistory(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/v1/history`);
    if (!response.ok) {
      throw new Error('Failed to get history');
    }
    return await response.json();
  }

  async getTrends(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/v1/trends`);
    if (!response.ok) {
      throw new Error('Failed to get trends');
    }
    return await response.json();
  }

  async getAlerts(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/v1/alerts`);
    if (!response.ok) {
      throw new Error('Failed to get alerts');
    }
    return await response.json();
  }

  async getStatus(): Promise<SystemStatus> {
    const response = await fetch(`${this.baseUrl}/api/v1/status`);
    if (!response.ok) {
      throw new Error('Failed to get status');
    }
    return await response.json();
  }

  // Metrics endpoints
  async getMetrics(): Promise<any> {
    const response = await fetch(`${this.baseUrl}/api/v1/metrics`);
    if (!response.ok) {
      throw new Error('Failed to get metrics');
    }
    return await response.json();
  }

  async getScanRateMetrics(): Promise<MetricsData> {
    const response = await fetch(`${this.baseUrl}/api/v1/metrics/scan-rate`);
    if (!response.ok) {
      throw new Error('Failed to get scan rate metrics');
    }
    return await response.json();
  }

  async getVulnerabilityMetrics(): Promise<VulnerabilityMetrics> {
    const response = await fetch(`${this.baseUrl}/api/v1/metrics/vulnerabilities`);
    if (!response.ok) {
      throw new Error('Failed to get vulnerability metrics');
    }
    return await response.json();
  }

  async getDurationMetrics(): Promise<DurationMetrics> {
    const response = await fetch(`${this.baseUrl}/api/v1/metrics/duration`);
    if (!response.ok) {
      throw new Error('Failed to get duration metrics');
    }
    return await response.json();
  }


	// New endpoints for detailed security information
	async getProcesses(): Promise<ProcessesResponse> {
		const response = await fetch(`${this.baseUrl}/api/v1/processes`);
		if (!response.ok) {
			throw new Error('Failed to get process details');
		}
		return await response.json();
	}

	async getNetwork(): Promise<NetworkResponse> {
		const response = await fetch(`${this.baseUrl}/api/v1/network`);
		if (!response.ok) {
			throw new Error('Failed to get network details');
		}
		return await response.json();
	}

	async fix(request: FixRequest): Promise<FixResponse> {
		const response = await fetch(`${this.baseUrl}/api/v1/fix`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(request),
		});
		if (!response.ok) {
			throw new Error('Failed to execute fix');
		}
		return await response.json();
	}

	async getTrendHistory(days?: number): Promise<TrendHistoryResponse> {
		const url = days 
			? `${this.baseUrl}/api/v1/trends/history?days=${days}`
			: `${this.baseUrl}/api/v1/trends/history`;
		const response = await fetch(url);
		if (!response.ok) {
			throw new Error('Failed to get trend history');
		}
		return await response.json();
	}
}

export const apiClient = new ApiClient();
