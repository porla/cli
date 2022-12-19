package config

import (
	"osprey/i18n"
)

var (
	Osprey_version = "v0.0.3"
	ConfigFilePath = "./config.yaml"
	Config         ConfigType
)

type ConfigType struct {
	SecretKey          string `yaml:"secretkey"`
	JSONRPCEndpointURL string `yaml:"JSONRPCEndpointURL"`
	PageSize           int    `yaml:"pagesize"`
	I18nLanguage       string `yaml:"i18nlanguage"`
}

var Currenti18n i18n.I18n
