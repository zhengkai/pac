package core

import (
	"flag"
	"fmt"

	"github.com/zhengkai/pac/core/log"
)

const Version = "1.0.0"

func Flag() bool {
	outputDir := flag.String(`output-dir`, ``, `directory to save files`)
	showVersion := flag.Bool(`version`, false, `show version`)

	if *showVersion {
		fmt.Println(Version)
		return false
	}

	if outputDir != nil {
		log.F("output-dir: %s\n", *outputDir)
	}
	return true
}
