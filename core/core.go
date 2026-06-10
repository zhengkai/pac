// Package core ...
package core

import (
	"github.com/zhengkai/pac/core/log"
)

func Run() {

	lanIP = getLanIP()

	cfgFile := `data/config.json`

	cfg, err := LoadCfg(cfgFile)
	if err != nil {
		log.W(`load config fail:`, cfgFile, err)
		return
	}

	writeJS(cfg)
	writeClash(cfg)
}
