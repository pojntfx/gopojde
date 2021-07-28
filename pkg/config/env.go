package config

import (
	"errors"

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
	EnabledServices []string // <- DON'T FORGET THESE!
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

	// TODO: Add rest of parameters

	return nil
}
