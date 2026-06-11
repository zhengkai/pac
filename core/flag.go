package core

import (
	"flag"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zhengkai/pac/core/log"
)

const Version = "1.0.0"

var outputDir string
var pwd string
var configFile string
var configDir string

func Flag() bool {

	getPwd()

	cFile := flag.String(`config-file`, ``, `config file`)
	oDir := flag.String(`output-dir`, ``, `directory to save files`)
	showVersion := flag.Bool(`version`, false, `show version`)
	help := flag.Bool(`help`, false, `show help`)
	flag.Parse()

	log.J(`pac `, Version)
	if *showVersion {
		return false
	}

	if *help {
		flag.Usage()
		return false
	}

	flagConfig(cFile)
	log.J(`config file:`, configFile)

	flagOutput(oDir)
	log.J(`output dir:`, outputDir)
	return true
}

// 获取 pwd，失败则获取 $0 的目录
func getPwd() {
	var err error
	pwd, err = os.Getwd()
	if err != nil {
		log.W("get pwd fail:", err)
		execPath, err := os.Executable()
		if err != nil {
			pwd = ``
			return
		}
		pwd = path.Dir(execPath)
	}
}

func flagConfig(f *string) {
	if f == nil || *f == `` {
		configFile = `config.json`
		return
	}

	cf := *f

	if strings.HasPrefix(cf, `~/`) {
		home, err := os.UserHomeDir()
		if err != nil {
			return
		}
		cf = filepath.Join(home, cf[2:])
	}

	if path.IsAbs(cf) {
		configFile = cf
		return
	}

	configFile, _ = filepath.Abs(path.Join(pwd, cf))
}

func flagOutput(d *string) {

	if d == nil || *d == `` {
		outputDir = path.Join(pwd, `output`)
		return
	}
	s := *d

	if path.IsAbs(s) {
		outputDir = s
		return
	}

	outputDir, _ = filepath.Abs(path.Join(pwd, s))
}
