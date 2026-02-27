package model

// RiskLevel represents the security risk level
type RiskLevel string

const (
	Low      RiskLevel = "LOW"
	Medium   RiskLevel = "MEDIUM"
	High     RiskLevel = "HIGH"
	Critical RiskLevel = "CRITICAL"
)

// PermissionResult represents the scanning result for each permission type
type PermissionResult struct {
	Filesystem  RiskLevel `json:"filesystem"`
	Shell       RiskLevel `json:"shell"`
	Network     RiskLevel `json:"network"`
	Secrets     RiskLevel `json:"secrets"`
	FileContent RiskLevel `json:"filecontent"`
	Dependencies RiskLevel `json:"dependencies"`
}


// ScanReport represents the complete security scan report
type ScanReport struct {
	ToolName string           `json:"tool_name"`
	Results  PermissionResult `json:"results"`
	Overall  RiskLevel        `json:"overall"`
	Details  []RiskDetail     `json:"details,omitempty"`
}

// RiskDetail provides additional information about detected risks
type RiskDetail struct {
	Type        RiskLevel `json:"type"`
	Category    string    `json:"category"`
	Description string    `json:"description"`
	Path        string    `json:"path,omitempty"`
}

// Plugin represents a plugin or tool that can be scanned
type Plugin struct {
	Name        string   `json:"name"`
	Path        string   `json:"path"`
	Exec        string   `json:"exec,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// PluginScanResult represents the result of a plugin security scan
type PluginScanResult struct {
	PluginName string    `json:"plugin_name"`
	Path       string    `json:"path"`
	Risk       RiskLevel `json:"risk"`
	Reason     string    `json:"reason"`
	Detected   []string  `json:"detected,omitempty"`
}

// PolicyConfig represents the security policy configuration
type PolicyConfig struct {
	Version    int         `yaml:"version"`
	Filesystem FSRule      `yaml:"filesystem"`
	Shell      ShellRule   `yaml:"shell"`
	Network    NetworkRule `yaml:"network"`
	Secrets    SecretsRule `yaml:"secrets"`
	Sandbox    SandboxRule `yaml:"sandbox"`
}

// FSRule defines filesystem access rules
type FSRule struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

// ShellRule defines shell command execution rules
type ShellRule struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

// NetworkRule defines network access rules
type NetworkRule struct {
	Allow []string `yaml:"allow"`
	Deny  []string `yaml:"deny"`
}

// SecretsRule defines secrets protection rules
type SecretsRule struct {
	BlockEnv []string `yaml:"block_env"`
}

// SandboxRule defines sandbox execution rules
type SandboxRule struct {
	DisableNetwork bool `yaml:"disable_network"`
	ReadonlyRoot   bool `yaml:"readonly_root"`
}
