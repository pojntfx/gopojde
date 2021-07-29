package config

import (
	"errors"
	"strings"

	"github.com/joho/godotenv"
)

const (
	prefix = "POJDE_"
)

func getWithPrefix(key string) string {
	return prefix + key
}

type EnvConfig struct {
	RootPassword string
	UserName     string
	UserPassword string

	UserEmail    string
	UserFullName string
	SSHKey       string

	AdditionalIPs     []string
	AdditionalDomains []string

	EnabledModules  []string
	EnabledServices []string
}

func NewConfig() *EnvConfig {
	return &EnvConfig{
		AdditionalIPs:     []string{},
		AdditionalDomains: []string{},

		EnabledModules:  []string{},
		EnabledServices: []string{},
	}
}

func (c *EnvConfig) Read(envFileContents string) error {
	env, err := godotenv.Unmarshal(envFileContents)
	if err != nil {
		return err
	}

	ok := false
	c.RootPassword, ok = env[getWithPrefix("ROOT_PASSWORD")]
	if !ok {
		return errors.New("could not get root password from config file")
	}

	c.UserName, ok = env[getWithPrefix("USERNAME")]
	if !ok {
		return errors.New("could not get username from config file")
	}

	c.UserPassword, ok = env[getWithPrefix("PASSWORD")]
	if !ok {
		return errors.New("could not get password from config file")
	}

	c.UserEmail, ok = env[getWithPrefix("EMAIL")]
	if !ok {
		return errors.New("could not get email from config file")
	}

	c.UserFullName, ok = env[getWithPrefix("FULL_NAME")]
	if !ok {
		return errors.New("could not get full name from config file")
	}

	c.SSHKey, ok = env[getWithPrefix("SSH_KEY_URL")]
	if !ok {
		return errors.New("could not get SSH key from config file")
	}

	additionalIPs, ok := env[getWithPrefix("IP")]
	if !ok {
		return errors.New("could not get additional IPs from config file")
	}
	c.AdditionalIPs = strings.Split(additionalIPs, " ")

	additionalDomains, ok := env[getWithPrefix("DOMAIN")]
	if !ok {
		return errors.New("could not get domains from config file")
	}
	c.AdditionalDomains = strings.Split(additionalDomains, " ")

	enabledModules, ok := env[getWithPrefix("MODULES")]
	if !ok {
		return errors.New("could not get modules from config file")
	}
	c.EnabledModules = strings.Split(enabledModules, " ")

	enabledServices, ok := env[getWithPrefix("SERVICES")]
	if !ok {
		return errors.New("could not get services from config file")
	}
	c.EnabledServices = strings.Split(enabledServices, " ")

	return nil
}
