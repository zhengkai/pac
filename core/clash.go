package core

import (
	"fmt"
	"slices"
	"strings"

	"github.com/zhengkai/pac/core/log"
)

type Clash struct {
	Config
	Buf
	server    []string
	serverMap map[string]*Server
}

func writeClash(cfg *Config) {

	for name, list := range cfg.List {
		c := &Clash{
			Config: *cfg,
		}
		c.run(name, list)
		file := fmt.Sprintf(`clash-%s.yaml`, name)
		writeFile(file, c.buf.Bytes())
	}
}

func (c *Clash) run(name string, list *List) {
	c.scanServer(list)
	c.buildServer()
	c.buildGroup()
	c.buildRule(list)
}

func (c *Clash) scanServer(list *List) {
	m := make(map[string]bool)
	c.serverMap = make(map[string]*Server)
	for _, r := range list.Rule {
		if !r.CanClash() {
			continue
		}
		m[r.Server] = true
		s, ok := c.Server[r.Server]
		if ok {
			for _, v := range s.Fallback {
				m[v] = true
			}
		}
	}
	m[list.Final] = true
	li := []string{}
	for k := range m {
		if k == `direct` || k == `reject` {
			continue
		}
		s, ok := c.Server[k]
		if !ok {
			continue
		}
		c.serverMap[k] = s
		li = append(li, k)
	}
	slices.Sort(li)
	c.server = li
}

func (c *Clash) buildServer() {
	c.w(`proxies:`)
	for _, name := range c.server {
		c.br()
		s := c.serverMap[name]
		if s.VMess == nil {
			li := strings.SplitN(s.Socks5, `:`, 2)
			if len(li) < 2 {
				log.WF(`server %s socks5 fail`, s.Socks5)
				continue
			}
			host := li[0]
			if host == `` {
				host = lanIP
			}
			c.wf(`  - name: %s`, name)
			c.w(`    type: socks5`)
			c.wf(`    server: %s`, host)
			c.wf(`    port: %s`, li[1])
		} else {
			vm := s.VMess
			c.wf(`  - name: %s`, name)
			c.w(`    type: vmess`)
			c.wf(`    server: %s`, vm.Host)
			c.wf(`    port: %d`, vm.Port)
			c.wf(`    uuid: %s`, vm.UUID)
			c.w(`    alterId: 0`)
			c.w(`    cipher: auto`)
			c.w(`    tls: true`)
			c.wf(`    servername: %s`, vm.Domain)
			c.w(`    skip-cert-verify: false`)
		}
	}
	c.br()
}

func (c *Clash) groupName(s string) string {

	if s == `direct` {
		return `DIRECT`
	}
	if s == `reject` {
		return `REJECT`
	}
	s = strings.ToUpper(s[:1]) + s[1:]
	return fmt.Sprintf(`G%s`, s)
}

func (c *Clash) buildGroup() {
	c.w(`proxy-groups:`)
	for _, name := range c.server {
		c.br()
		c.wf(`  - name: %s`, c.groupName(name))
		c.w(`    type: fallback`)
		c.w(`    proxies:`)
		c.wf(`      - %s`, name)
		server := c.serverMap[name]
		for _, s := range server.Fallback {
			c.wf(`      - %s`, s)
		}
	}
	c.br()
}

func (c *Clash) buildRule(list *List) {
	c.w(`rules:`)

	for _, v := range list.Rule {
		c.br()
		for _, domain := range v.domain {
			c.wf(`  - DOMAIN-SUFFIX,%s,%s`, domain, c.groupName(v.Server))

		}
	}
	c.br()
	c.wf(`  - MATCH,%s`, list.Final)
}
