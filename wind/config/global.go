package config

import (
	"fmt"
	env2 "mond/wind/env"
)

func GetAppid() string {
	return env2.GetAppId()
}

func GetHostName() string {
	return env2.GetHostName()
}

func GetEnv() string {
	return env2.GetEnv()
}

func SetNamespaceId() {
	base := "base"
	if env2.GetEnv() != "" {
		base += "_" + env2.GetEnv()
	}
	nsId := fmt.Sprintf("%s", _viper.GetStringMap(base)["namespaceid"])
	env2.SetNamespaceId(nsId)
}

func SetServerIp() {
	base := "base"
	if env2.GetEnv() != "" {
		base += "_" + env2.GetEnv()
	}
	serverIp := fmt.Sprintf("%s", _viper.GetStringMap(base)["serverip"])
	env2.SetServerIp(serverIp)
}
