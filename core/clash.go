package core

type Clash struct {
	Config
	Buf
}

/*
proxies:
  - name: doll
    type: vmess
    server: 6.7.8.9
    port: 53005
    uuid: ffff
    alterId: 0
    cipher: auto
    tls: true
    servername: doll.9farm.com
    skip-cert-verify: false

proxy-groups:
  - name: Proxy
    type: select
    proxies:
      - doll

rules:
  - DOMAIN-SUFFIX,google.com,Proxy
  - MATCH,DIRECT
*/

func writeClash() {
}
