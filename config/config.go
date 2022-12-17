package config

var (
	Osprey_version = "v0.0.1"
	ConfigFilePath = "./config.yaml"
	Config         ConfigType
)

type ConfigType struct {
	SecretKey          string `yaml:"secretkey"`
	JSONRPCEndpointURL string `yaml:"JSONRPCEndpointURL"`
	PageSize           int    `yaml:"pagesize"`
}
