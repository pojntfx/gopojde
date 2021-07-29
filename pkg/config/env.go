package config

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/alessio/shellescape"
	"github.com/joho/godotenv"
)

const (
	prefix = "POJDE_"
)

const (
	rootPasswordKey = "ROOT_PASSWORD"
	usernameKey     = "USERNAME"
	passwordKey     = "PASSWORD"
	emailKey        = "EMAIL"
	fullNameKey     = "FULL_NAME"
	sshKeyKey       = "SSH_KEY_URL"
	ipKey           = "IP"
	domainKey       = "DOMAIN"
	modulesKey      = "MODULES"
	servicesKey     = "SERVICES"
)

func getWithPrefix(key string) string {
	return prefix + key
}

func getSuffixForModuleName(moduleName string) (string, error) {
	parts := strings.Split(moduleName, ".")

	if len(parts) < 2 {
		return "", fmt.Errorf(`"." character is missing in module name`)
	}

	return parts[1], nil
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

func (c *EnvConfig) Unmarshal(envFileContents string) error {
	// Parse config file
	env, err := godotenv.Unmarshal(envFileContents)
	if err != nil {
		return err
	}

	// Basic configuration parameters
	ok := false
	c.RootPassword, ok = env[getWithPrefix(rootPasswordKey)]
	if !ok {
		return errors.New("could not get root password from config file")
	}

	c.UserName, ok = env[getWithPrefix(usernameKey)]
	if !ok {
		return errors.New("could not get username from config file")
	}

	c.UserPassword, ok = env[getWithPrefix(passwordKey)]
	if !ok {
		return errors.New("could not get password from config file")
	}

	c.UserEmail, ok = env[getWithPrefix(emailKey)]
	if !ok {
		return errors.New("could not get email from config file")
	}

	c.UserFullName, ok = env[getWithPrefix(fullNameKey)]
	if !ok {
		return errors.New("could not get full name from config file")
	}

	c.SSHKey, ok = env[getWithPrefix(sshKeyKey)]
	if !ok {
		return errors.New("could not get SSH key from config file")
	}

	// Addition IPs and domains
	additionalIPs, ok := env[getWithPrefix(ipKey)]
	if !ok {
		return errors.New("could not get additional IPs from config file")
	}
	c.AdditionalIPs = strings.Split(additionalIPs, " ")

	additionalDomains, ok := env[getWithPrefix(domainKey)]
	if !ok {
		return errors.New("could not get domains from config file")
	}
	c.AdditionalDomains = strings.Split(additionalDomains, " ")

	// Modules and services
	enabledModules, ok := env[getWithPrefix(modulesKey)]
	if !ok {
		return errors.New("could not get modules from config file")
	}
	for _, module := range strings.Split(enabledModules, " ") {
		enabledModule, err := getSuffixForModuleName(module)
		if err != nil {
			return err
		}

		c.EnabledModules = append(c.EnabledModules, enabledModule)
	}

	enabledServices, ok := env[getWithPrefix(modulesKey)]
	if !ok {
		return errors.New("could not get services from config file")
	}
	for _, service := range strings.Split(enabledServices, " ") {
		enabledService, err := getSuffixForModuleName(service)
		if err != nil {
			return err
		}

		c.EnabledServices = append(c.EnabledServices, enabledService)
	}

	return nil
}

func (c *EnvConfig) Marshal() string {
	env := map[string]string{}

	// Basic configuration parameters
	env[rootPasswordKey] = c.RootPassword
	env[usernameKey] = c.UserName
	env[passwordKey] = c.UserPassword
	env[emailKey] = c.UserEmail
	env[fullNameKey] = c.UserFullName
	env[sshKeyKey] = c.SSHKey

	// Addition IPs and domains
	env[ipKey] = strings.Join(c.AdditionalIPs, " ")
	env[domainKey] = strings.Join(c.AdditionalDomains, " ")

	// Modules and services
	env[modulesKey] = strings.Join(c.EnabledModules, " ")
	for _, moduleName := range c.EnabledModules {
		env[strings.ToUpper("MODULE_"+moduleName+"_ENABLED")] = "true"
	}
	env[servicesKey] = strings.Join(c.EnabledServices, " ")
	for _, serviceName := range c.EnabledServices {
		env[strings.ToUpper("MODULE_"+serviceName+"_ENABLED")] = "true"
	}

	// Marshal
	lines := make([]string, 0, len(env))
	for k, v := range env {
		lines = append(lines, fmt.Sprintf(`export %s="%s"`, k, shellescape.Quote(v)))
	}
	sort.Strings(lines)

	return strings.Join(lines, "\n")
}
