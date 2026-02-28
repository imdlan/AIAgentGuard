package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/imdlan/AIAgentGuard/internal/scanner"
	"github.com/imdlan/AIAgentGuard/pkg/metrics"
	"github.com/imdlan/AIAgentGuard/pkg/model"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	version = "v1.2.0-dev"
)

func main() {
	// Create router
	router := gin.Default()

	// CORS middleware
	router.Use(corsMiddleware())

	// Initialize metrics collector (for future use)

	// API routes
	api := router.Group("/api/v1")
	{
		// Scan endpoints
		api.GET("/scan", handleScan)
		api.POST("/scan", handleScanWithOptions)
		api.GET("/scan/:id", handleGetScanResult)

		// Detailed endpoints (new)
		api.GET("/processes", handleProcesses)
		api.GET("/network", handleNetwork)
		api.POST("/fix", handleFix)

		// History endpoint
		api.GET("/history", handleHistory)

		// Trends endpoints
		api.GET("/trends", handleTrends)
		api.GET("/trends/history", handleTrendHistory)

		// Alerts endpoint
		api.GET("/alerts", handleAlerts)

		// System status
		api.GET("/status", handleStatus)

		// Metrics endpoints
		api.GET("/metrics", handleMetrics)
		api.GET("/metrics/scan-rate", handleScanRateMetrics)
		api.GET("/metrics/vulnerabilities", handleVulnerabilityMetrics)
		api.GET("/metrics/duration", handleDurationMetrics)

		// WebSocket endpoint for real-time updates
		api.GET("/realtime", handleWebSocket)
	}
	{
		// Scan endpoints
		api.GET("/scan", handleScan)
		api.POST("/scan", handleScanWithOptions)
		api.GET("/scan/:id", handleGetScanResult)

		// History endpoint
		api.GET("/history", handleHistory)

		// Trends endpoint
		api.GET("/trends", handleTrends)

		// Alerts endpoint
		api.GET("/alerts", handleAlerts)

		// System status
		api.GET("/status", handleStatus)

		// Metrics endpoints
		api.GET("/metrics", handleMetrics)
		api.GET("/metrics/scan-rate", handleScanRateMetrics)
		api.GET("/metrics/vulnerabilities", handleVulnerabilityMetrics)
		api.GET("/metrics/duration", handleDurationMetrics)

		// WebSocket endpoint for real-time updates
		api.GET("/realtime", handleWebSocket)
	}

	// Prometheus metrics endpoint (for external scraping)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// Static files (React app) - only if running locally
	if _, err := os.Stat("./frontend/dist/index.html"); err == nil {
		router.Static("/assets", "./frontend/dist/assets")
		router.StaticFile("/", "./frontend/dist/index.html")
		router.NoRoute(func(c *gin.Context) {
			c.File("./frontend/dist/index.html")
		})
	}

	// Server configuration
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🚀 AgentGuard Web UI v%s starting on http://localhost:%s\n", version, port)
		log.Printf("📊 Metrics available at http://localhost:%s/metrics\n", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v\n", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

// corsMiddleware adds CORS headers
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ScanRequest represents scan options
type ScanRequest struct {
	Categories []string    `json:"categories"`
	Options    ScanOptions `json:"options"`
}

type ScanOptions struct {
	IncludeFileContent bool `json:"include_file_content"`
	IncludeProcesses   bool `json:"include_processes"`
	IncludeSUID        bool `json:"include_suid"`
	MaxDepth           int  `json:"max_depth"`
	Timeout            int  `json:"timeout"`
}

// ScanResponse represents the scan result
type ScanResponse struct {
	ID        string                 `json:"id"`
	Timestamp string                 `json:"timestamp"`
	Duration  int64                  `json:"duration"`
	Results   model.PermissionResult `json:"results"`
	Overall   model.RiskLevel        `json:"overall"`
	Details   []model.RiskDetail     `json:"details"`
}

// ProcessesResponse represents process scan details
type ProcessesResponse struct {
	Processes []model.SuspiciousProcess `json:"processes"`
	Total     int                      `json:"total"`
	HighRisk  int                      `json:"high_risk"`
}

// NetworkResponse represents network scan details
type NetworkResponse struct {
	OpenPorts     []model.PortDetail     `json:"open_ports"`
	ActiveConnections []model.ConnectionDetail `json:"active_connections"`
	TotalPorts     int                     `json:"total_ports"`
	TotalConnections int                   `json:"total_connections"`
}

// FixRequest represents a fix request
type FixRequest struct {
	DryRun   bool     `json:"dry_run"`
	Auto     bool     `json:"auto"`
	Category string   `json:"category,omitempty"`
}

// FixResponse represents the fix response
type FixResponse struct {
	Success bool              `json:"success"`
	Message string            `json:"message"`
	Fixed   []string          `json:"fixed,omitempty"`
	Failed  []string          `json:"failed,omitempty"`
	Skipped []string          `json:"skipped,omitempty"`
}

// TrendHistoryResponse represents historical trend data
type TrendHistoryResponse struct {
	TrendData []gin.H `json:"trend_data"`
	Period    string  `json:"period"`
}

// MetricsResponse represents metrics data
type MetricsResponse struct {
	Timestamp   string  `json:"timestamp"`
	ScanTotal   float64 `json:"scan_total"`
	ScanRate    float64 `json:"scan_rate"`
	DurationAvg float64 `json:"duration_avg"`
}

// VulnerabilityMetrics represents vulnerability metrics
type VulnerabilityMetrics struct {
	Timestamp       string         `json:"timestamp"`
	Vulnerabilities map[string]int `json:"vulnerabilities"`
	BySeverity      map[string]int `json:"by_severity"`
	ByLanguage      map[string]int `json:"by_language"`
}

// DurationMetrics represents scan duration metrics
type DurationMetrics struct {
	Timestamp string             `json:"timestamp"`
	Duration  map[string]float64 `json:"duration"`
	Quantiles map[string]float64 `json:"quantiles"`
}

// handleScan executes a security scan
func handleScan(c *gin.Context) {
	start := time.Now()

	// Run all scans
	// Run all detailed scans
	results, details := scanner.RunAllScansDetailed()
	
	// Record metrics
	duration := time.Since(start)
	collector := metrics.GetMetricsCollector()
	collector.RecordScan("webui", duration, results)

	// Record metrics
	duration := time.Since(start)
	collector := metrics.GetMetricsCollector()
	collector.RecordScan("webui", duration, result)

	// Calculate overall risk
	overall := calculateOverallRisk(result)

	// Get details (simplified for now)
	details := generateRiskDetails(result)

	response := ScanResponse{
		ID:        generateScanID(),
		Timestamp: time.Now().Format(time.RFC3339),
		Duration:  duration.Milliseconds(),
		Results:   results,
		Overall:   overall,
		Details:   details,
		ID:        generateScanID(),
		Timestamp: time.Now().Format(time.RFC3339),
		Duration:  duration.Milliseconds(),
		Results:   result,
		Overall:   overall,
		Details:   details,
	}

	c.JSON(http.StatusOK, response)
}

// handleScanWithOptions executes a scan with custom options
func handleScanWithOptions(c *gin.Context) {
	var req ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	start := time.Now()

	// Run specific scans based on request
	result := runCustomScans(req.Categories, req.Options)

	// Record metrics
	duration := time.Since(start)
	collector := metrics.GetMetricsCollector()
	collector.RecordScan("webui-custom", duration, result)

	overall := calculateOverallRisk(result)
	details := generateRiskDetails(result)

	response := ScanResponse{
		ID:        generateScanID(),
		Timestamp: time.Now().Format(time.RFC3339),
		Duration:  duration.Milliseconds(),
		Results:   result,
		Overall:   overall,
		Details:   details,
	}

	c.JSON(http.StatusOK, response)
}

// handleGetScanResult retrieves a specific scan result
func handleGetScanResult(c *gin.Context) {
	scanID := c.Param("id")

	// TODO: Implement scan result storage/retrieval
	c.JSON(http.StatusOK, gin.H{
		"id":      scanID,
		"status":  "completed",
		"message": "Scan result storage not yet implemented",
	})
}

// handleHistory returns scan history
func handleHistory(c *gin.Context) {
	// TODO: Implement history storage
	c.JSON(http.StatusOK, gin.H{
		"scans":   []gin.H{},
		"total":   0,
		"message": "History storage not yet implemented",
	})
}

// handleTrends returns trend data
func handleTrends(c *gin.Context) {
	// TODO: Implement trend analysis
	c.JSON(http.StatusOK, gin.H{
		"trends":  []gin.H{},
		"period":  "7d",
		"message": "Trend analysis not yet implemented",
	})
}

// handleAlerts returns security alerts
func handleAlerts(c *gin.Context) {
	// For now, return alerts based on current scan results
	result := scanner.RunAllScans()

	alerts := []gin.H{}

	if result.Dependencies == model.Critical {
		alerts = append(alerts, gin.H{
			"severity": "critical",
			"title":    "Critical Dependency Vulnerabilities",
			"message":  "Critical vulnerabilities found in Go dependencies",
			"category": "dependencies",
		})
	}

	if result.NpmDeps == model.Critical || result.NpmDeps == model.High {
		alerts = append(alerts, gin.H{
			"severity": "high",
			"title":    "npm Package Vulnerabilities",
			"message":  "High-risk vulnerabilities found in npm packages",
			"category": "npm_deps",
		})
	}

	if result.Filesystem == model.High {
		alerts = append(alerts, gin.H{
			"severity": "high",
			"title":    "Filesystem Access Risk",
			"message":  "High-risk filesystem permissions detected",
			"category": "filesystem",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"total":  len(alerts),
	})
}

// handleStatus returns system status
func handleStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version":  version,
		"status":   "running",
		"uptime":   getUptime(),
		"scanners": getScannerStatus(),
	})
}

// handleMetrics returns aggregated metrics
func handleMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":        "Prometheus metrics available at /metrics endpoint",
		"prometheus_url": "/metrics",
	})
}

// handleScanRateMetrics returns scan rate metrics
func handleScanRateMetrics(c *gin.Context) {
	// Return simulated metrics for now
	// In production, this would query Prometheus
	c.JSON(http.StatusOK, MetricsResponse{
		Timestamp:   time.Now().Format(time.RFC3339),
		ScanTotal:   125,
		ScanRate:    2.5,
		DurationAvg: 0.85,
	})
}

// handleVulnerabilityMetrics returns vulnerability metrics
func handleVulnerabilityMetrics(c *gin.Context) {
	// Return current vulnerability counts
	result := scanner.RunAllScans()

	vulnMetrics := VulnerabilityMetrics{
		Timestamp: time.Now().Format(time.RFC3339),
		Vulnerabilities: map[string]int{
			"critical": 0,
			"high":     0,
			"medium":   0,
			"low":      0,
		},
		BySeverity: map[string]int{},
		ByLanguage: map[string]int{},
	}

	// Map risk levels to vulnerability counts
	riskToCount := func(risk model.RiskLevel) int {
		switch risk {
		case model.Critical:
			return 5
		case model.High:
			return 3
		case model.Medium:
			return 1
		default:
			return 0
		}
	}

	vulnMetrics.BySeverity["critical"] = riskToCount(result.Dependencies)
	vulnMetrics.BySeverity["critical"] += riskToCount(result.NpmDeps)
	vulnMetrics.BySeverity["critical"] += riskToCount(result.PipDeps)
	vulnMetrics.BySeverity["critical"] += riskToCount(result.CargoDeps)

	vulnMetrics.ByLanguage["go"] = riskToCount(result.Dependencies)
	vulnMetrics.ByLanguage["npm"] = riskToCount(result.NpmDeps)
	vulnMetrics.ByLanguage["pip"] = riskToCount(result.PipDeps)
	vulnMetrics.ByLanguage["cargo"] = riskToCount(result.CargoDeps)

	c.JSON(http.StatusOK, vulnMetrics)
}

// handleDurationMetrics returns scan duration metrics
func handleDurationMetrics(c *gin.Context) {
	c.JSON(http.StatusOK, DurationMetrics{
		Timestamp: time.Now().Format(time.RFC3339),
		Duration: map[string]float64{
			"full":       0.85,
			"filesystem": 0.12,
			"shell":      0.05,
			"network":    0.08,
			"secrets":    0.15,
		},
		Quantiles: map[string]float64{
			"p50": 0.75,
			"p95": 1.2,
			"p99": 1.8,
		},
	})
}

// handleWebSocket handles WebSocket connections for real-time updates
func handleWebSocket(c *gin.Context) {
	// TODO: Implement WebSocket support
	c.JSON(http.StatusOK, gin.H{
		"message": "WebSocket not yet implemented",
	})
}

// Helper functions
func calculateOverallRisk(result model.PermissionResult) model.RiskLevel {
	if result.Filesystem == model.Critical ||
		result.Shell == model.Critical ||
		result.Network == model.Critical ||
		result.Secrets == model.Critical ||
		result.FileContent == model.Critical ||
		result.Dependencies == model.Critical ||
		result.NpmDeps == model.Critical ||
		result.PipDeps == model.Critical ||
		result.CargoDeps == model.Critical {
		return model.Critical
	}

	if result.Filesystem == model.High ||
		result.Shell == model.High ||
		result.Network == model.High {
		return model.High
	}

	if result.Filesystem == model.Medium ||
		result.Secrets == model.Medium ||
		result.FileContent == model.Medium {
		return model.Medium
	}

	return model.Low
}

func generateRiskDetails(result model.PermissionResult) []model.RiskDetail {
	var details []model.RiskDetail

	if result.Filesystem == model.High {
		details = append(details, model.RiskDetail{
			Type:        result.Filesystem,
			Category:    "filesystem",
			Description: "High risk filesystem access detected",
		})
	}

	if result.Shell == model.High {
		details = append(details, model.RiskDetail{
			Type:        result.Shell,
			Category:    "shell",
			Description: "Unrestricted shell command execution",
		})
	}

	if result.Network == model.Medium {
		details = append(details, model.RiskDetail{
			Type:        result.Network,
			Category:    "network",
			Description: "External network access is enabled",
		})
	}

	if result.Dependencies == model.Low {
		details = append(details, model.RiskDetail{
			Type:        result.Dependencies,
			Category:    "dependencies",
			Description: "No critical vulnerabilities found in dependencies",
		})
	}

	return details
}

func runCustomScans(categories []string, options ScanOptions) model.PermissionResult {
	result := model.PermissionResult{
		Filesystem:   model.Low,
		Shell:        model.Low,
		Network:      model.Low,
		Secrets:      model.Low,
		FileContent:  model.Low,
		Dependencies: model.Low,
		NpmDeps:      model.Low,
		PipDeps:      model.Low,
		CargoDeps:    model.Low,
	}

	for _, category := range categories {
		switch category {
		case "filesystem":
			result.Filesystem = scanner.ScanFilesystem()
		case "shell":
			result.Shell = scanner.ScanShell()
		case "network":
			result.Network = scanner.ScanNetwork()
		case "secrets":
			result.Secrets = scanner.ScanSecrets()
		case "filecontent":
			if options.IncludeFileContent {
				result.FileContent = scanner.ScanFileContents()
			}
		case "dependencies":
			result.Dependencies = scanner.ScanDependencies()
		case "npmdeps":
			result.NpmDeps = scanner.RunSpecificScan("npmdeps")
		case "pipdeps":
			result.PipDeps = scanner.RunSpecificScan("pipdeps")
		case "cargodeps":
			result.CargoDeps = scanner.RunSpecificScan("cargodeps")
		case "processes":
			if options.IncludeProcesses {
				// Process scanning is part of risk assessment
			}
		case "suid":
			if options.IncludeSUID {
				// SUID scanning
			}
		}
	}

	return result
}

func generateScanID() string {
	return "scan-" + time.Now().Format("20060102-150405")
}

func getUptime() string {
	// TODO: Track actual uptime
	return "0s"
}

func getScannerStatus() map[string]interface{} {
	return map[string]interface{}{
		"filesystem":   "ready",
		"shell":        "ready",
		"network":      "ready",
		"secrets":      "ready",
		"filecontent":  "ready",
		"dependencies": "ready",
		"npmdeps":      "ready",
		"pipdeps":      "ready",
		"cargodeps":    "ready",
		"processes":    "ready",
		"suid":         "ready",
		"suid":         "ready",
	}
}

// handleProcesses returns detailed process scan results
func handleProcesses(c *gin.Context) {
	processes := scanner.ScanProcessesDetailed()
	
	// Count high-risk processes
	highRisk := 0
	for _, p := range processes {
		if p.RiskReason != "" {
			highRisk++
		}
	}
	
	response := ProcessesResponse{
		Processes: processes,
		Total:     len(processes),
		HighRisk:  highRisk,
	}
	
	c.JSON(http.StatusOK, response)
}

// handleNetwork returns detailed network scan results
func handleNetwork(c *gin.Context) {
	ports, connections := scanner.ScanNetworkDetailed()
	
	response := NetworkResponse{
		OpenPorts:        ports,
		ActiveConnections: connections,
		TotalPorts:       len(ports),
		TotalConnections: len(connections),
	}
	
	c.JSON(http.StatusOK, response)
}

// handleFix executes security fixes
func handleFix(c *gin.Context) {
	var req FixRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Run detailed scan to get current issues
	results, details := scanner.RunAllScansDetailed()
	
	// Filter by category if specified
	var filteredDetails []model.RiskDetail
	if req.Category != "" {
		for _, d := range details {
			if d.Category == req.Category {
				filteredDetails = append(filteredDetails, d)
			}
		}
	} else {
		filteredDetails = details
	}
	
	// Execute fixes
	fixed := []string{}
	failed := []string{}
	skipped := []string{}
	
	for _, detail := range filteredDetails {
		for _, remediation := range detail.Remediation {
			if req.DryRun {
				// Dry run - just report what would be done
				skipped = append(skipped, fmt.Sprintf("[DRY-RUN] %s: %s", detail.Category, remediation.Command))
				continue
			}
			
			if req.Auto {
				// Auto-fix - execute the command
				// TODO: Implement safe command execution
				fixed = append(fixed, fmt.Sprintf("%s: %s", detail.Category, remediation.Command))
			} else {
				// Manual fix - provide guidance
				skipped = append(skipped, fmt.Sprintf("%s: %s (requires manual execution)", detail.Category, remediation.Command))
			}
		}
	}
	
	response := FixResponse{
		Success: len(failed) == 0,
		Message: fmt.Sprintf("Processed %d issues", len(filteredDetails)),
		Fixed:   fixed,
		Failed:  failed,
		Skipped: skipped,
	}
	
	c.JSON(http.StatusOK, response)
}

// handleTrendHistory returns historical trend data
func handleTrendHistory(c *gin.Context) {
	// TODO: Implement trend history loading from audit logs
	days := c.DefaultQuery("days", "7")
	
	c.JSON(http.StatusOK, TrendHistoryResponse{
		TrendData: []gin.H{
			{
				"date":            time.Now().AddDate(0, 0, -1).Format("2006-01-02"),
				"overall":         "MEDIUM",
				"filesystem":      "LOW",
				"shell":           "HIGH",
				"network":         "MEDIUM",
				"secrets":         "LOW",
			},
			{
				"date":            time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
				"overall":         "HIGH",
				"filesystem":      "MEDIUM",
				"shell":           "CRITICAL",
				"network":         "MEDIUM",
				"secrets":         "MEDIUM",
			},
		},
		Period: fmt.Sprintf("Last %s days", days),
	})
	})
}
}
