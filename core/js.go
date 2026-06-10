package core

import (
	"fmt"

	"github.com/zhengkai/pac/core/log"
)

const FileMode = 0664
const DirFileMode = 0775

var lanIP string

type JS struct {
	Config
	Buf
	isLan bool
}

func writeJS(cfg *Config) {

	log.J()

	for _, isLan := range []bool{true, false} {
		for name, row := range cfg.List {
			js := &JS{
				Config: *cfg,
				isLan:  isLan,
			}
			network := `local`
			if isLan {
				network = `lan`
			}
			js.build(row)

			file := fmt.Sprintf(`output/%s-%s.js`, network, name)
			writeFile(file, js.buf.Bytes())
		}
	}
}

func (js *JS) serverAddr(server string, withFallback bool) string {

	if server == `direct` {
		return `DIRECT`
	}
	if server == `reject` {
		return `PROXY 0.0.0.0:0`
	}

	s, ok := js.Server[server]
	if !ok && s.Socks5 == `` {
		if server != `direct` {
			log.W(`unknown server:`, server)
		}
		return `DIRECT`
	}
	so := s.Socks5
	if js.isLan {
		so = getLanIP() + so
	} else {
		so = `127.0.0.1` + so
	}
	re := `SOCKS5 ` + so

	/*
		if withFallback {
			for _, v := range s.Fallback {
				if v == server {
					log.W(`fallback self:`, server)
					continue
				}
				re += `; ` + js.serverAddr(v, false)
			}
		}
	*/

	return re
}

func (js *JS) server(server string) {

	s := js.serverAddr(server, true)

	js.wf("\treturn '%s';", s)
}

func (js *JS) build(list *List) {

	js.w(`function FindProxyForURL(_, raw) {`)
	js.w("\tconst host = '.' + raw;")

	for idx, r := range list.Rule {
		if idx > 0 {
			js.br()
		}
		js.rule(r)
	}

	js.br()
	js.server(list.Final)

	js.w(`}`)
}

func (js *JS) rule(r *Rule) {
	js.w("\tif (")
	js.domain(r.domain)
	if len(r.Extra) > 0 {
		js.ruleExtra(r.Extra, len(r.domain) > 0)
	}
	js.w("\t) {")
	js.buf.WriteString("\t")
	js.server(r.Server)
	js.w("\t}")
}

func (js *JS) domain(r []string) {
	for idx, d := range r {
		js.buf.WriteString("\t\t")
		if idx != 0 {
			js.buf.WriteString(`|| `)
		}
		js.wf(`host.endsWith('.%s')`, d)
	}
}

func (js *JS) ruleExtra(r []string, domain bool) {
	for idx, e := range r {
		js.buf.WriteString("\t\t")
		if idx != 0 || domain {
			js.buf.WriteString(`|| `)
		}
		js.w(e)
	}
}
