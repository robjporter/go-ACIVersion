package app

import (
	"github.com/robjporter/go-functions/logrus"
	"github.com/robjporter/go-functions/viper"
)

type Application struct {
	ConfigFile string
	Debug      bool
	Config     *viper.Viper
	Logger     *logrus.Logger
	ACI        []ACISystemInfo
	Key        []byte
	Version    string
}

type ACISystemInfo struct {
	ip                  string
	username            string
	password            string
	cookie              string
	name                string
	version             string
	status              bool
	suggestedVersion    string
	deferredVersion     bool
	releaseVersion      string
	releaseTrainVersion string
}
