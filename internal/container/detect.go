package container

import (
	"os"
	"os/exec"
	"strings"
)

type ContainerRuntime string

const (
	RuntimeNone      ContainerRuntime = "none"
	RuntimeDocker    ContainerRuntime = "docker"
	RuntimeKubernetes ContainerRuntime = "kubernetes"
	RuntimePodman    ContainerRuntime = "podman"
	RuntimeLXC       ContainerRuntime = "lxc"
	RuntimeWasm      ContainerRuntime = "wasm"
)

type ContainerInfo struct {
	Runtime      ContainerRuntime `json:"runtime"`
	IsContainer  bool             `json:"is_container"`
	Namespace    string           `json:"namespace,omitempty"`
	PodName      string           `json:"pod_name,omitempty"`
	PodID        string           `json:"pod_id,omitempty"`
	ContainerID  string           `json:"container_id,omitempty"`
	NodeName     string           `json:"node_name,omitempty"`
}

func DetectContainer() ContainerInfo {
	info := ContainerInfo{
		Runtime:     RuntimeNone,
		IsContainer: false,
	}

	if isDocker() {
		info.Runtime = RuntimeDocker
		info.IsContainer = true
		info.ContainerID = getContainerID()
		return info
	}

	if isKubernetes() {
		info.Runtime = RuntimeKubernetes
		info.IsContainer = true
		info.Namespace = getKubernetesNamespace()
		info.PodName = getKubernetesPodName()
		info.PodID = getKubernetesPodID()
		info.ContainerID = getContainerID()
		info.NodeName = getKubernetesNodeName()
		return info
	}

	if isPodman() {
		info.Runtime = RuntimePodman
		info.IsContainer = true
		info.ContainerID = getContainerID()
		return info
	}

	if isLXC() {
		info.Runtime = RuntimeLXC
		info.IsContainer = true
		return info
	}

	if isWasm() {
		info.Runtime = RuntimeWasm
		info.IsContainer = true
		return info
	}

	return info
}

func isDocker() bool {
	if _, err := os.Stat("/.dockerenv"); err == nil {
		return true
	}

	if _, err := os.Stat("/.dockerinit"); err == nil {
		return true
	}

	cgroupPath := "/proc/1/cgroup"
	if data, err := os.ReadFile(cgroupPath); err == nil {
		content := string(data)
		if strings.Contains(content, "docker") || strings.Contains(content, "/docker/") {
			return true
		}
	}

	return false
}

func isKubernetes() bool {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {
		return true
	}

	cgroupPath := "/proc/1/cgroup"
	if data, err := os.ReadFile(cgroupPath); err == nil {
		content := string(data)
		if strings.Contains(content, "kubepods") || strings.Contains(content, "kubernetes") {
			return true
		}
	}

	return false
}

func isPodman() bool {
	cgroupPath := "/proc/1/cgroup"
	if data, err := os.ReadFile(cgroupPath); err == nil {
		content := string(data)
		if strings.Contains(content, "podman") || strings.Contains(content, "libpod") {
			return true
		}
	}

	if _, err := os.Stat("/run/podman/podman.sock"); err == nil {
		return true
	}

	return false
}

func isLXC() bool {
	cgroupPath := "/proc/1/cgroup"
	if data, err := os.ReadFile(cgroupPath); err == nil {
		content := string(data)
		if strings.Contains(content, "lxc") || strings.Contains(content, "lxcfs") {
			return true
		}
	}

	if _, err := os.Stat("/proc/self/cgroup"); err == nil {
		if data, err := os.ReadFile("/proc/self/cgroup"); err == nil {
			if strings.Contains(string(data), "lxc") {
				return true
			}
		}
	}

	return false
}

func isWasm() bool {
	return os.Getenv("WASMER_BACKEND") != "" || os.Getenv("WASI_SDK") != ""
}

func getContainerID() string {
	cgroupPath := "/proc/self/cgroup"
	if data, err := os.ReadFile(cgroupPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.Contains(line, "docker") || strings.Contains(line, "podman") {
				parts := strings.Split(line, "/")
				if len(parts) > 0 {
					id := parts[len(parts)-1]
					if len(id) >= 12 {
						return id[:12]
					}
					return id
				}
			}
		}
	}

	if _, err := os.Stat("/etc/hostname"); err == nil {
		if data, err := os.ReadFile("/etc/hostname"); err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return "unknown"
}

func getKubernetesNamespace() string {
	if ns := os.Getenv("KUBERNETES_NAMESPACE"); ns != "" {
		return ns
	}

	if data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace"); err == nil {
		return strings.TrimSpace(string(data))
	}

	return "default"
}

func getKubernetesPodName() string {
	if pod := os.Getenv("POD_NAME"); pod != "" {
		return pod
	}

	if hostname, err := os.Hostname(); err == nil {
		return hostname
	}

	return "unknown"
}

func getKubernetesPodID() string {
	cgroupPath := "/proc/self/cgroup"
	if data, err := os.ReadFile(cgroupPath); err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")
		for _, line := range lines {
			if strings.Contains(line, "kubepods") {
				parts := strings.Split(line, "/")
				for i, part := range parts {
					if strings.Contains(part, "pod") {
						if i+1 < len(parts) {
							return parts[i+1][:12]
						}
					}
				}
			}
		}
	}

	return "unknown"
}

func getKubernetesNodeName() string {
	if node := os.Getenv("NODE_NAME"); node != "" {
		return node
	}

	if node := os.Getenv("KUBE_NODE_NAME"); node != "" {
		return node
	}

	return "unknown"
}

func GetContainerRuntimeInfo() string {
	info := DetectContainer()
	if !info.IsContainer {
		return "Not running in a container"
	}

	result := "Running in container: " + string(info.Runtime)
	if info.Namespace != "" {
		result += "\nNamespace: " + info.Namespace
	}
	if info.PodName != "" {
		result += "\nPod: " + info.PodName
	}
	if info.ContainerID != "" {
		result += "\nContainer ID: " + info.ContainerID
	}
	return result
}

func IsRunningInContainer() bool {
	return DetectContainer().IsContainer
}

func HasDockerSocket() bool {
	if _, err := os.Stat("/var/run/docker.sock"); err == nil {
		return true
	}
	if _, err := os.Stat("/run/docker.sock"); err == nil {
		return true
	}
	return false
}

func HasContainerdSocket() bool {
	if _, err := os.Stat("/run/containerd/containerd.sock"); err == nil {
		return true
	}
	if _, err := os.Stat("/var/run/containerd/containerd.sock"); err == nil {
		return true
	}
	return false
}

func GetDockerVersion() string {
	if _, err := exec.LookPath("docker"); err != nil {
		return "not installed"
	}

	cmd := exec.Command("docker", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "error"
	}

	return strings.TrimSpace(string(output))
}

func GetPodmanVersion() string {
	if _, err := exec.LookPath("podman"); err != nil {
		return "not installed"
	}

	cmd := exec.Command("podman", "--version")
	output, err := cmd.Output()
	if err != nil {
		return "error"
	}

	return strings.TrimSpace(string(output))
}
