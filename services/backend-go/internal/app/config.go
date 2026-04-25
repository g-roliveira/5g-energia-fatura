package app

import "os"

type Config struct {
	Host               string
	Port               string
	APIKey             string
	BackofficePGURL    string
	ExtractorBaseURL   string
	NeoenergiaBaseURL  string
	ArtifactsDirectory string
	DatabaseURL        string
	IntegrationPGURL   string
	EncryptionKey      string
	BootstrapPythonBin string
	BootstrapScript    string
}

func LoadConfigFromEnv() Config {
	return Config{
		Host:               envOrDefault("BACKEND_HOST", "127.0.0.1"),
		Port:               envOrDefault("BACKEND_PORT", "8080"),
		APIKey:             envOrDefault("BACKEND_API_KEY", ""),
		BackofficePGURL:    envOrDefault("BACKOFFICE_PG_URL", ""),
		ExtractorBaseURL:   envOrDefault("EXTRACTOR_BASE_URL", "http://127.0.0.1:8090"),
		NeoenergiaBaseURL:  envOrDefault("NEOENERGIA_API_BASE_URL", "https://apineprd.neoenergia.com"),
		ArtifactsDirectory: envOrDefault("ARTIFACTS_DIR", "./artifacts"),
		DatabaseURL:        envOrDefault("BACKEND_DATABASE_URL", "file:data/backend-go.db"),
		IntegrationPGURL: envOrDefault(
			"BACKEND_INTEGRATION_PG_URL",
			envOrDefault("INTEGRATION_PG_URL", ""),
		),
		EncryptionKey:      envOrDefault("ENCRYPTION_KEY", ""),
		BootstrapPythonBin: envOrDefault("BOOTSTRAP_PYTHON_BIN", "./.venv/bin/python"),
		BootstrapScript:    envOrDefault("BOOTSTRAP_SCRIPT_PATH", "scripts/bootstrap_neoenergia_token.py"),
	}
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
