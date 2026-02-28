package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/imdlan/AIAgentGuard/internal/audit"
	"github.com/imdlan/AIAgentGuard/internal/scanner"
	"github.com/imdlan/AIAgentGuard/internal/risk"
	"github.com/imdlan/AIAgentGuard/pkg/model"
	"github.com/spf13/cobra"
)

var (
	trendDays     int
	trendJSON     bool
	trendCategory string
)

var trendCmd = &cobra.Command{
	Use:   "trend",
	Short: "Show security risk trends over time",
	Long: `Analyze historical scan data to show risk trends.

Compares current scan results with historical data from audit logs
to show improvements or deteriorations in security posture.

Supports:
  - Time range selection (last N days)
  - Category-specific trends
  - JSON output for automation
  - Visual trend indicators`,
	Example: `  agent-guard trend
  agent-guard trend --days 7
  agent-guard trend --category filesystem --json`,
	RunE: runTrend,
}

func init() {
	rootCmd.AddCommand(trendCmd)

	trendCmd.Flags().IntVarP(&trendDays, "days", "d", 7, "Number of days to analyze")
	trendCmd.Flags().BoolVarP(&trendJSON, "json", "j", false, "Output trends in JSON format")
	trendCmd.Flags().StringVarP(&trendCategory, "category", "c", "", "Filter by category (filesystem, shell, network)")
}

type TrendData struct {
	Date     string               `json:"date"`
	Overall  model.RiskLevel      `json:"overall"`
	Category map[string]RiskValue `json:"category,omitempty"`
}

type RiskValue struct {
	Level model.RiskLevel `json:"level"`
	Score int             `json:"score"`
}

type TrendAnalysis struct {
	Period   string    `json:"period"`
	Current  TrendData `json:"current"`
	Previous TrendData `json:"previous"`
	Trend    string    `json:"trend"` // "improving", "worsening", "stable"
	Changes  []string  `json:"changes"`
	Summary  string    `json:"summary"`
}

func runTrend(cmd *cobra.Command, args []string) error {
	fmt.Println("📈 Security Risk Trend Analysis")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Get current scan results
	fmt.Println("📊 Analyzing current security state...")
	results, _ := scanner.RunAllScansDetailed()
	
	// Generate scan report (includes Overall and Details)
	report := risk.Analyze(results)
	
	currentData := TrendData{
		Date:    time.Now().Format("2006-01-02"),
		Overall: report.Overall,
		Category: make(map[string]RiskValue),
	}
	
	// Map risk levels to scores
	scoreMap := map[model.RiskLevel]int{
		model.Low:      25,
		model.Medium:   50,
		model.High:     75,
		model.Critical: 100,
	}
	
	// Extract category scores from details
	for _, detail := range report.Details {
		currentData.Category[detail.Category] = RiskValue{
			Level: detail.Type,
			Score: scoreMap[detail.Type],
		}
	}


	// Load historical data
	auditLogPath := audit.GetLogFilePath()
	if auditLogPath == "" {
		// Try default path
		homeDir, _ := os.UserHomeDir()
		auditLogPath = filepath.Join(homeDir, ".agent-guard", "audit.log")
	}

	historicalData, err := loadHistoricalScans(auditLogPath, trendDays)
	if err != nil {
		fmt.Printf("⚠️  Could not load historical data: %v\n", err)
		fmt.Println("   Showing only current state...")
		displayTrend(currentData, TrendData{}, "N/A")
		return nil
	}

	// Find previous scan for comparison
	previousData := findComparableScan(historicalData, currentData)

	// Analyze trends
	analysis := analyzeTrends(currentData, previousData)

	// Display results
	if trendJSON {
		return displayTrendJSON(analysis)
	}

	return displayTrendAnalysis(analysis)
}

// loadHistoricalScans loads scan results from audit log
func loadHistoricalScans(logPath string, days int) ([]TrendData, error) {
	var trends []TrendData

	file, err := os.Open(logPath)
	if err != nil {
		return trends, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cutoffDate := time.Now().AddDate(0, 0, -days)

	for {
		var event audit.AuditEvent
		if err := decoder.Decode(&event); err != nil {
			break
		}

		// Only consider scan completion events
		// Only consider scan completion events
		if event.EventType == "scan_completed" {
			// Timestamp is already time.Time, no need to parse
			eventDate := event.Timestamp
			
			if eventDate.Before(cutoffDate) {
				continue
			}

			// Extract scan results from event details
			if scanResults, ok := event.Details["results"].(map[string]interface{}); ok {
				data := TrendData{
					Date:    eventDate.Format("2006-01-02"),
					Overall: getRiskLevelFromMap(scanResults),
				}

				trends = append(trends, data)
			}
		}
	}

	// Sort by date (newest first)
	sort.Slice(trends, func(i, j int) bool {
		return trends[i].Date > trends[j].Date
	})

	return trends, nil
}

// findComparableScan finds the most recent comparable scan
func findComparableScan(historical []TrendData, current TrendData) TrendData {
	for _, data := range historical {
		if data.Date != current.Date {
			return data
		}
	}
	return TrendData{}
}

// analyzeTrends analyzes trends between current and previous scans
func analyzeTrends(current, previous TrendData) TrendAnalysis {
	if previous.Date == "" {
		return TrendAnalysis{
			Period:   "No historical data",
			Current:  current,
			Previous: previous,
			Trend:    "unknown",
			Changes:  []string{"First scan - no baseline for comparison"},
			Summary:  "Run another scan in the future to establish trends",
		}
	}

	scoreMap := map[model.RiskLevel]int{
		model.Low:      25,
		model.Medium:   50,
		model.High:     75,
		model.Critical: 100,
	}

	currentScore := scoreMap[current.Overall]
	previousScore := scoreMap[previous.Overall]

	// Determine trend direction
	trend := "stable"
	changes := []string{}

	if currentScore > previousScore {
		trend = "worsening"
		changes = append(changes, fmt.Sprintf("Risk increased from %s to %s", previous.Overall, current.Overall))
	} else if currentScore < previousScore {
		trend = "improving"
		changes = append(changes, fmt.Sprintf("Risk decreased from %s to %s", previous.Overall, current.Overall))
	} else {
		changes = append(changes, "Risk level unchanged")
	}

	// Compare categories
	for cat := range current.Category {
		if prevCat, ok := previous.Category[cat]; ok {
			if current.Category[cat].Score > prevCat.Score {
				changes = append(changes, fmt.Sprintf("%s: worsened from %s to %s", cat, prevCat.Level, current.Category[cat].Level))
			} else if current.Category[cat].Score < prevCat.Score {
				changes = append(changes, fmt.Sprintf("%s: improved from %s to %s", cat, prevCat.Level, current.Category[cat].Level))
			}
		}
	}

	// Generate summary
	summary := generateTrendSummary(trend, changes)

	return TrendAnalysis{
		Period:   fmt.Sprintf("%s to %s", previous.Date, current.Date),
		Current:  current,
		Previous: previous,
		Trend:    trend,
		Changes:  changes,
		Summary:  summary,
	}
}

// generateTrendSummary creates a human-readable summary
func generateTrendSummary(trend string, changes []string) string {
	switch trend {
	case "improving":
		return "✅ Security posture is improving! Keep up the good work."
	case "worsening":
		return "⚠️  Security posture is deteriorating. Review recent changes."
	case "stable":
		return "➡️  Security posture is stable. Continue monitoring."
	default:
		return "❓ Insufficient data to determine trend."
	}
}

// displayTrendAnalysis displays trend analysis in console format
func displayTrendAnalysis(analysis TrendAnalysis) error {
	fmt.Printf("\n📅 Analysis Period: %s\n", analysis.Period)
	fmt.Printf("📊 Trend Direction: %s\n\n", getTrendSymbol(analysis.Trend))

	// Show comparison
	fmt.Println("Risk Level Comparison:")
	fmt.Printf("  Previous: %s\n", analysis.Previous.Overall)
	fmt.Printf("  Current:  %s\n", analysis.Current.Overall)

	// Show changes
	if len(analysis.Changes) > 0 {
		fmt.Println("\n📝 Changes Detected:")
		for _, change := range analysis.Changes {
			fmt.Printf("  • %s\n", change)
		}
	}

	// Show summary
	fmt.Printf("\n%s\n\n", analysis.Summary)

	// Show category breakdown
	if len(analysis.Current.Category) > 0 {
		fmt.Println("📂 Category Breakdown:")
		for cat, data := range analysis.Current.Category {
			fmt.Printf("  %s: %s (score: %d)\n", cat, data.Level, data.Score)
		}
	}

	return nil
}

// displayTrendJSON outputs trend analysis in JSON format
func displayTrendJSON(analysis TrendAnalysis) error {
	data, err := json.MarshalIndent(analysis, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}

// getTrendSymbol returns a symbol for the trend direction
func getTrendSymbol(trend string) string {
	switch trend {
	case "improving":
		return "📈 Improving ↗️"
	case "worsening":
		return "📉 Worsening ↘️"
	case "stable":
		return "➡️  Stable →"
	default:
		return "❓ Unknown ?"
	}
}

// Helper function to get risk level from map
func getRiskLevelFromMap(m map[string]interface{}) model.RiskLevel {
	if v, ok := m["overall"].(string); ok {
		return model.RiskLevel(v)
	}
	return model.Low
}

// displayTrend displays simple trend data
func displayTrend(current, previous TrendData, trend string) error {
	fmt.Printf("Current Scan (%s):\n", current.Date)
	fmt.Printf("  Overall Risk: %s\n\n", current.Overall)

	if previous.Date != "" {
		fmt.Printf("Previous Scan (%s):\n", previous.Date)
		fmt.Printf("  Overall Risk: %s\n\n", previous.Overall)
		fmt.Printf("Trend: %s\n", trend)
	}

	return nil
}
