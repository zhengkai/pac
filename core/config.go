package core

import (
	"encoding/json"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/zhengkai/pac/core/log"
)

type Rule struct {
	Server string   `json:"server"`
	File   []string `json:"file"`
	Extra  []string `json:"extra"`
	domain []string
}

type List struct {
	Rule  []*Rule `json:"rule"`
	Final string  `json:"final"`
}

type Server struct {
	Socks5 string `json:"socks5"`
	// VMess  string `json:"vmess"`
	Fallback []string `json:"fallback"`
}

type Config struct {
	Server map[string]*Server `json:"server"`
	List   map[string]*List   `json:"list"`
}

func LoadCfg(file string) (*Config, error) {

	ab, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(ab, cfg)
	if err != nil {
		return nil, err
	}

	dir := path.Dir(file)
	err = parseCfg(cfg, dir)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func parseCfg(cfg *Config, dir string) error {
	for name, v := range cfg.List {
		for _, r := range v.Rule {
			for _, f := range r.File {
				re, err := loadDomain(path.Join(dir, f))
				if err != nil {
					log.W(err)
					continue
				}
				r.domain = append(r.domain, re...)
				log.J(name, r.Server, f, len(re))
			}
			r.domain = unique(r.domain)
			slices.Sort(r.domain)
		}
	}
	return nil
}

func loadDomain(file string) ([]string, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	content := strings.Split(string(data), "\n")
	result := make([]string, 0, len(content))
	for _, line := range content {
		line = strings.TrimSpace(line)
		if line == `` {
			continue
		}
		if strings.HasPrefix(line, `#`) {
			continue
		}
		line = strings.ToLower(line)
		result = append(result, line)
	}
	return result, nil
}
