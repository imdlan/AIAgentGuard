package metrics

import (
	"net/http"
	"sync"
	"time"

	"github.com/imdlan/AIAgentGuard/pkg/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Singleton instance
	instance *MetricsCollector
	once     sync.Once
)

// MetricsCollector collects and exports Prometheus metrics
type MetricsCollector struct {
	registry *prometheus.Registry

	// Scan metrics
	scanTotal       *prometheus.CounterVec
	scanDuration    *prometheus.HistogramVec
	scanResults     *prometheus.GaugeVec
	vulnerabilities *prometheus.CounterVec

	// Component metrics
	componentScans    *prometheus.CounterVec
	componentFailures *prometheus.CounterVec

	// System metrics
	memoryUsage  prometheus.Gauge
	uptime       prometheus.Counter
	lastScanTime prometheus.Gauge
}

// GetMetricsCollector returns the singleton MetricsCollector instance
func GetMetricsCollector() *MetricsCollector {
	once.Do(func() {
		instance = NewMetricsCollector()
	})
	return instance
}

// NewMetricsCollector creates a new MetricsCollector
func NewMetricsCollector() *MetricsCollector {
	registry := prometheus.NewRegistry()

	m := &MetricsCollector{
		registry: registry,

		scanTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "agent_guard_scans_total",
				Help: "Total number of security scans performed",
			},
			[]string{"scan_type"}, // full, specific category
		),

		scanDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "agent_guard_scan_duration_seconds",
				Help:    "Duration of security scans in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"scan_type"},
		),

		scanResults: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "agent_guard_scan_result",
				Help: "Result of the most recent scan by category",
			},
			[]string{"category"}, // filesystem, shell, network, secrets, dependencies, npm_deps, pip_deps, cargo_deps
		),

		vulnerabilities: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "agent_guard_vulnerabilities_total",
				Help: "Total number of vulnerabilities discovered",
			},
			[]string{"severity", "language"}, // severity: critical, high, medium, low; language: go, npm, pip, cargo
		),

		componentScans: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "agent_guard_component_scans_total",
				Help: "Total number of component scans",
			},
			[]string{"component"}, // scanner, analyzer, reporter
		),

		componentFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "agent_guard_component_failures_total",
				Help: "Total number of component failures",
			},
			[]string{"component", "error_type"},
		),

		memoryUsage: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "agent_guard_memory_usage_bytes",
				Help: "Current memory usage in bytes",
			},
		),

		uptime: prometheus.NewCounter(
			prometheus.CounterOpts{
				Name: "agent_guard_uptime_seconds",
				Help: "Uptime of the agent-guard service in seconds",
			},
		),

		lastScanTime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "agent_guard_last_scan_timestamp",
				Help: "Timestamp of the last completed scan",
			},
		),
	}

	// Register metrics
	registry.MustRegister(
		m.scanTotal,
		m.scanDuration,
		m.scanResults,
		m.vulnerabilities,
		m.componentScans,
		m.componentFailures,
		m.memoryUsage,
		m.uptime,
		m.lastScanTime,
	)

	// Start uptime counter
	go m.trackUptime()

	return m
}

// RecordScan records metrics for a completed scan
func (m *MetricsCollector) RecordScan(scanType string, duration time.Duration, results model.PermissionResult) {
	// Increment scan counter
	m.scanTotal.WithLabelValues(scanType).Inc()

	// Record duration
	m.scanDuration.WithLabelValues(scanType).Observe(duration.Seconds())

	// Record results
	m.recordScanResults(results)

	// Update last scan time
	m.lastScanTime.SetToCurrentTime()
}

// recordScanResults records individual scan category results
func (m *MetricsCollector) recordScanResults(results model.PermissionResult) {
	// Convert risk level to numeric value (0=LOW, 1=MEDIUM, 2=HIGH, 3=CRITICAL)
	riskToValue := func(risk model.RiskLevel) float64 {
		switch risk {
		case model.Low:
			return 0
		case model.Medium:
			return 1
		case model.High:
			return 2
		case model.Critical:
			return 3
		default:
			return 0
		}
	}

	m.scanResults.WithLabelValues("filesystem").Set(riskToValue(results.Filesystem))
	m.scanResults.WithLabelValues("shell").Set(riskToValue(results.Shell))
	m.scanResults.WithLabelValues("network").Set(riskToValue(results.Network))
	m.scanResults.WithLabelValues("secrets").Set(riskToValue(results.Secrets))
	m.scanResults.WithLabelValues("filecontent").Set(riskToValue(results.FileContent))
	m.scanResults.WithLabelValues("dependencies").Set(riskToValue(results.Dependencies))
	m.scanResults.WithLabelValues("npm_deps").Set(riskToValue(results.NpmDeps))
	m.scanResults.WithLabelValues("pip_deps").Set(riskToValue(results.PipDeps))
	m.scanResults.WithLabelValues("cargo_deps").Set(riskToValue(results.CargoDeps))
}

// RecordVulnerabilities records discovered vulnerabilities
func (m *MetricsCollector) RecordVulnerabilities(severity string, language string, count int) {
	for i := 0; i < count; i++ {
		m.vulnerabilities.WithLabelValues(severity, language).Inc()
	}
}

// RecordComponentScan records a component scan
func (m *MetricsCollector) RecordComponentScan(component string) {
	m.componentScans.WithLabelValues(component).Inc()
}

// RecordComponentFailure records a component failure
func (m *MetricsCollector) RecordComponentFailure(component string, errorType string) {
	m.componentFailures.WithLabelValues(component, errorType).Inc()
}

// trackUptime tracks the uptime of the service
func (m *MetricsCollector) trackUptime() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	startTime := time.Now()
	for range ticker.C {
		elapsed := time.Since(startTime).Seconds()
		m.uptime.Add(elapsed)
	}
}

// GetHandler returns the HTTP handler for serving metrics
func (m *MetricsCollector) GetHandler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

// GetHandlerWithRegistry returns a custom metrics handler
func GetHandlerWithRegistry(registry *prometheus.Registry) http.Handler {
	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
}

// StartMetricsServer starts an HTTP server to serve metrics
func (m *MetricsCollector) StartMetricsServer(addr string) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", m.GetHandler())

	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return server.ListenAndServe()
}

// StartMetricsServerAsync starts the metrics server in a goroutine
func (m *MetricsCollector) StartMetricsServerAsync(addr string) {
	go func() {
		if err := m.StartMetricsServer(addr); err != nil {
			// Log error but don't crash
			panic(err)
		}
	}()
}
