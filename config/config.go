package config

type Config struct {
	ServiceName    string   `json:"service_name"`
	DependServices []string `json:"depend_services"`
	ConfigFile     string   `json:"config_file"`
	ConfigRegistry string   `json:"config_registry"`
	Host           string   `json:"host"`
	Port           int      `json:"port"`
	EnableHTTP     bool     `json:"enable_http"`
	Database       string   `json:"database"`
	Logger         string   `json:"logger"`
}

var conf = Config{
	Host:           "0.0.0.0",
	Port:           8080,
	EnableHTTP:     true,
	ServiceName:    "my_test",
	DependServices: []string{"my_test"},
}

func GetConfig() Config {
	return conf
}
