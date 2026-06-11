// Package core ...
package core

import (
	"github.com/zhengkai/pac/core/log"
)

func Run() {

	lanIP = getLanIP()

	cfg, err := LoadCfg(configFile)
	if err != nil {
		log.W(`load config fail:`, configFile, err)
		return
	}

	writeJS(cfg)
	writeClash(cfg)
}
