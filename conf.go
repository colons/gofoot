// internal stuff for dealing with configurations

package main

import "github.com/thoj/go-ircevent"

type UserConfig struct {
	global map[string]string
	network map[string]map[string]string
	channel map[string]map[string]map[string]string
}

type config struct {
	userConfig UserConfig
	ourNetwork string
}

func GetConfig(network string) config {
	return config{GetUserConfig(), network}
}

// get config for key applying globally
func (c config) Global(key string) string {
	if value, ok := c.userConfig.global[key]; ok {
		return value
	} else {
		return GlobalConfig.Global(key)
	}
}

// get config for key applying to the attached network
func (c config) Network(key string) string {
	if value, ok := c.userConfig.network[c.ourNetwork][key]; ok {
		return value
	} else {
		return c.Global(key)
	}
}

// Get config for key applying to the target for whom event was issued
func (c config) Event(event *irc.Event, key string) string {
	if value, ok := c.userConfig.channel[c.ourNetwork][getTarget(event)][key]; ok {
		return value
	} else {
		return c.Network(key)
	}
}
