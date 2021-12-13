package main

import (
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	blacklist  []string `ini:"blacklist"`
	nameserver string   `ini:"nameserver"`
}

func filter(ss []string, test func(string) bool) (ret []string) {
	for _, s := range ss {
		if test(s) {
			ret = append(ret, s)
		}
	}
	return
}

func isNotEmpty(str string) bool {
	return len(str) > 0
}

func readConfig(name string) *Config {
	cfg, err := ini.Load(name)
	exitOnError(err, "Fail to read file: %v")

	config := new(Config)
	config.blacklist = filter(cfg.Section("").Key("blacklist").Strings("\n"), isNotEmpty)
	config.nameserver = cfg.Section("").Key("nameserver").String()
	return config
}

func (config Config) isBlacklisted(hostname string) bool {
	for _, blacklistedDomain := range config.blacklist {
		if hostname == blacklistedDomain {
			return true
		} else if strings.HasSuffix(hostname, "."+blacklistedDomain) {
			return true
		}
	}
	return false
}
