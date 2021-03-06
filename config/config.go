package config

import "github.com/astaxie/beego"
import "github.com/humpback/common/models"
import "github.com/humpback/gounits/network"

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var config *models.Config

// Init - Load config info
func Init() {
	envEndpoint := os.Getenv("DOCKER_ENDPOINT")
	if envEndpoint == "" {
		envEndpoint = beego.AppConfig.DefaultString("DOCKER_ENDPOINT", "unix:///var/run/docker.sock")
	}

	envAPIVersion := os.Getenv("DOCKER_API_VERSION")
	if envAPIVersion == "" {
		envAPIVersion = beego.AppConfig.DefaultString("DOCKER_API_VERSION", "v1.20")
	}

	envRegistryAddr := os.Getenv("DOCKER_REGISTRY_ADDRESS")
	if envRegistryAddr == "" {
		envRegistryAddr = beego.AppConfig.DefaultString("DOCKER_REGISTRY_ADDRESS", "docker.neg")
	}

	envNodeHTTPAddr := os.Getenv("DOCKER_NODE_HTTPADDR")
	if envNodeHTTPAddr == "" {
		envNodeHTTPAddr = beego.AppConfig.DefaultString("DOCKER_NODE_HTTPADDR", "0.0.0.0:8500")
	}

	envContainerPortsRange := os.Getenv("DOCKER_CONTAINER_PORTS_RANGE")
	if envContainerPortsRange == "" {
		envContainerPortsRange = beego.AppConfig.DefaultString("DOCKER_CONTAINER_PORTS_RANGE", "0-0")
	}

	var envEnableBuildImg bool
	if tempEnableBuildImg := os.Getenv("ENABLE_BUILD_IMAGE"); tempEnableBuildImg != "" {
		if tempEnableBuildImg == "1" || tempEnableBuildImg == "true" {
			envEnableBuildImg = true
		}
	} else {
		envEnableBuildImg = beego.AppConfig.DefaultBool("ENABLE_BUILD_IMAGE", false)
	}

	envComposePath := os.Getenv("DOCKER_COMPOSE_PATH")
	if envComposePath == "" {
		envComposePath = beego.AppConfig.DefaultString("DOCKER_COMPOSE_PATH", "./compose_files")
	}

	var envComposePackageMaxSize int64
	packageMaxSize := os.Getenv("DOCKER_COMPOSE_PACKAGE_MAXSIZE")
	if packageMaxSize == "" {
		envComposePackageMaxSize = beego.AppConfig.DefaultInt64("DOCKER_COMPOSE_PACKAGE_MAXSIZE", 67108864)
	} else {
		value, err := strconv.ParseInt(packageMaxSize, 10, 64)
		if err != nil {
			envComposePackageMaxSize = 67108864
		} else {
			envComposePackageMaxSize = value
		}
	}

	envClusterEnabled := false
	enabled := os.Getenv("DOCKER_CLUSTER_ENABLED")
	if enabled == "" {
		envClusterEnabled = beego.AppConfig.DefaultBool("DOCKER_CLUSTER_ENABLED", false)
	} else {
		var err error
		if envClusterEnabled, err = strconv.ParseBool(enabled); err != nil {
			envClusterEnabled = false
		}
	}

	envClusterURIs := os.Getenv("DOCKER_CLUSTER_URIS")
	if envClusterURIs == "" {
		envClusterURIs = beego.AppConfig.DefaultString("DOCKER_CLUSTER_URIS", "zk://127.0.0.1:2181")
	}

	envClusterName := os.Getenv("DOCKER_CLUSTER_NAME")
	if envClusterName == "" {
		envClusterName = beego.AppConfig.DefaultString("DOCKER_CLUSTER_NAME", "humpback/center")
	}

	envClusterHeartBeat := os.Getenv("DOCKER_CLUSTER_HEARTBEAT")
	if envClusterHeartBeat == "" {
		envClusterHeartBeat = beego.AppConfig.DefaultString("DOCKER_CLUSTER_HEARTBEAT", "10s")
	}

	envClusterTTL := os.Getenv("DOCKER_CLUSTER_TTL")
	if envClusterTTL == "" {
		envClusterTTL = beego.AppConfig.DefaultString("DOCKER_CLUSTER_TTL", "35s")
	}

	var logLevel int
	if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
		logLevel, _ = strconv.Atoi(envLogLevel)
	} else {
		logLevel = beego.AppConfig.DefaultInt("LOG_LEVEL", 3)
	}

	config = &models.Config{
		DockerEndPoint:              envEndpoint,
		DockerAPIVersion:            envAPIVersion,
		DockerRegistryAddress:       envRegistryAddr,
		EnableBuildImage:            envEnableBuildImg,
		DockerComposePath:           envComposePath,
		DockerComposePackageMaxSize: envComposePackageMaxSize,
		DockerNodeHTTPAddr:          envNodeHTTPAddr,
		DockerContainerPortsRange:   envContainerPortsRange,
		DockerClusterEnabled:        envClusterEnabled,
		DockerClusterURIs:           envClusterURIs,
		DockerClusterName:           envClusterName,
		DockerClusterHeartBeat:      envClusterHeartBeat,
		DockerClusterTTL:            envClusterTTL,
		LogLevel:                    logLevel,
	}
}

// GetConfig - return config struct
func GetConfig() models.Config {
	return *config
}

// SetVersion - set app version
func SetVersion(version string) {
	config.AppVersion = version
}

// GetNodeHTTPAddrIPPort - return local agent httpaddr info
func GetNodeHTTPAddrIPPort() (string, int, error) {

	httpAddr := strings.TrimSpace(config.DockerNodeHTTPAddr)
	if strings.Index(httpAddr, ":") < 0 {
		httpAddr = httpAddr + ":"
	}

	pAddrStr := strings.SplitN(httpAddr, ":", 2)
	if pAddrStr[1] == "" {
		httpAddr = httpAddr + "8500"
	}

	strHost, strPort, err := net.SplitHostPort(httpAddr)
	if err != nil {
		return "", 0, err
	}

	if strHost == "" || strHost == "0.0.0.0" {
		strHost = network.GetDefaultIP()
	}

	ip := net.ParseIP(strHost)
	if ip == nil {
		return "", 0, fmt.Errorf("httpAddr ip invalid")
	}

	port, err := strconv.Atoi(strPort)
	if err != nil || port > 65535 || port <= 0 {
		return "", 0, fmt.Errorf("httpAddr port invalid")
	}
	return ip.String(), port, nil
}
