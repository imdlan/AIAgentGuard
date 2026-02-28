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
	Filesystem   RiskLevel `json:"filesystem"`
	Shell        RiskLevel `json:"shell"`
	Network      RiskLevel `json:"network"`
	Secrets      RiskLevel `json:"secrets"`
	FileContent  RiskLevel `json:"filecontent"`
	Dependencies RiskLevel `json:"dependencies"`
	NpmDeps      RiskLevel `json:"npm_deps"`
	PipDeps      RiskLevel `json:"pip_deps"`
	CargoDeps    RiskLevel `json:"cargo_deps"`
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
	Type        RiskLevel        `json:"type"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Path        string           `json:"path,omitempty"`
	Details     RiskSpecificInfo `json:"details,omitempty"`
	Remediation RemediationInfo  `json:"remediation,omitempty"`
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

// RiskSpecificInfo contains detailed information about detected risks
type RiskSpecificInfo struct {
	// Shell risk details
	ShellAvailable  string   `json:"shell_available,omitempty"`
	HasSudoAccess   bool     `json:"has_sudo_access,omitempty"`
	SudoSource      string   `json:"sudo_source,omitempty"`
	SudoRules       []string `json:"sudo_rules,omitempty"`

	// Filesystem risk details
	AffectedPaths   []PathDetail `json:"affected_paths,omitempty"`

	// Network risk details
	OpenPorts       []PortDetail `json:"open_ports,omitempty"`
	ActiveConns     []ConnectionDetail `json:"active_connections,omitempty"`

	// Process risk details
	SuspiciousProcs []ProcessInfo `json:"suspicious_processes,omitempty"`

	// Secrets risk details
	ExposedSecrets  []SecretInfo `json:"exposed_secrets,omitempty"`
}

// PathDetail contains detailed information about a filesystem path
type PathDetail struct {
	Path        string `json:"path"`
	Permission  string `json:"permission"`
	IsWritable  bool   `json:"is_writable"`
	Owner       string `json:"owner"`
	RiskReason  string `json:"risk_reason"`
}

// ProcessInfo contains information about a running process
type ProcessInfo struct {
	PID         int      `json:"pid"`
	Name        string   `json:"name"`
	CommandLine string   `json:"command_line"`
	User        string   `json:"user"`
	RiskReason  string   `json:"risk_reason"`
	StartTime   string   `json:"start_time"`
}

// PortDetail contains information about an open network port
type PortDetail struct {
	Port        int      `json:"port"`
	Protocol    string   `json:"protocol"`
	Service     string   `json:"service"`
	Process     string   `json:"process"`
	RiskReason  string   `json:"risk_reason"`
}

// ConnectionDetail contains information about an active network connection
type ConnectionDetail struct {
	LocalAddr   string `json:"local_addr"`
	RemoteAddr  string `json:"remote_addr"`
	Protocol    string `json:"protocol"`
	State       string `json:"state"`
	Process     string `json:"process"`
	RiskReason  string `json:"risk_reason"`
}

// SecretInfo contains information about an exposed secret
type SecretInfo struct {
	Type        string `json:"type"`
	Location    string `json:"location"`
	Value       string `json:"value,omitempty"` // Only show last 4 chars
	RiskReason  string `json:"risk_reason"`
}

// RemediationInfo provides actionable suggestions for fixing detected risks
type RemediationInfo struct {
	Summary     string            `json:"summary"`
	Steps       []RemediationStep `json:"steps"`
	Commands    []string          `json:"commands"`
	Priority    string            `json:"priority"`
	RiskAfter   RiskLevel         `json:"risk_after,omitempty"`
}

// RemediationStep represents a single step in the remediation process
type RemediationStep struct {
	Step         int    `json:"step"`
	Action       string `json:"action"`
	Command      string `json:"command,omitempty"`
	Explanation  string `json:"explanation"`
}
