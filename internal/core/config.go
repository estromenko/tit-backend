package core

import (
	"os"
	"path"

	"github.com/tutorin-tech/tit-backend/internal/utils"
)

const (
	defaultPort           = 3000
	defaultPgPort         = 5432
	defaultJWTExpireHours = 24 * 3
)

type Config struct {
	Debug                        bool
	Port                         int
	PgHost                       string
	PgPort                       int
	PgName                       string
	PgUser                       string
	PgPassword                   string
	SecretKey                    string
	JWTExpireHours               int
	KubernetesUseInClusterConfig bool
	KubernetesConfigPath         string
	KubernetesDashboardNamespace string
	DashboardImage               string
	DashboardIngressDomain       string
}

func NewConfig() *Config {
	defaultKubernetesConfigPath := path.Join(os.Getenv("HOME"), ".kube", "config")

	return &Config{
		Debug:                        utils.GetEnvOrDefault("DEBUG", "false") == "true",
		Port:                         utils.GetEnvIntOrDefault("PORT", defaultPort),
		PgHost:                       utils.GetEnvOrDefault("PG_HOST", "localhost"),
		PgPort:                       utils.GetEnvIntOrDefault("PG_PORT", defaultPgPort),
		PgName:                       utils.GetEnvOrDefault("PG_NAME", "tutorintech"),
		PgUser:                       utils.GetEnvOrDefault("PG_USER", "postgres"),
		PgPassword:                   utils.GetEnvOrDefault("PG_PASSWORD", "secret"),
		SecretKey:                    utils.GetEnv("SECRET_KEY"),
		JWTExpireHours:               utils.GetEnvIntOrDefault("JWT_EXPIRE_HOURS", defaultJWTExpireHours),
		KubernetesConfigPath:         utils.GetEnvOrDefault("KUBERNETES_CONFIG_PATH", defaultKubernetesConfigPath),
		KubernetesUseInClusterConfig: utils.GetEnvOrDefault("KUBERNETES_USE_IN_CLUSTER_CONFIG", "false") == "true",
		KubernetesDashboardNamespace: utils.GetEnvOrDefault("KUBERNETES_DASHBOARD_NAMESPACE", "default"),
		DashboardImage:               utils.GetEnv("DASHBOARD_IMAGE"),
		DashboardIngressDomain:       utils.GetEnv("DASHBOARD_INGRESS_DOMAIN"),
	}
}
