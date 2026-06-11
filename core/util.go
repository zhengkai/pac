package core

import (
	"bytes"
	"crypto/sha256"
	"net"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/zhengkai/pac/core/log"
	"golang.org/x/sys/unix"
)

const xattrHashKey = `user.sha256hash`

func unique(ss []string) []string {
	m := make(map[string]struct{}, len(ss))

	for _, s := range ss {
		m[s] = struct{}{}
	}

	result := make([]string, 0, len(m))
	for s := range m {
		result = append(result, s)
	}

	return result
}

func getLanIP() string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return ``
	}

	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			ip := ipNet.IP.To4()
			if ip == nil {
				continue
			}

			s := ip.String()
			if strings.HasPrefix(s, `192.168.`) {
				return s
			}
		}
	}

	return ``
}

func writeFile(file string, ab []byte) error {

	file = filepath.Join(outputDir, file)

	hash := sha256.Sum256(ab)

	buf := make([]byte, sha256.Size)
	_, err := unix.Getxattr(file, xattrHashKey, buf)
	if err == nil && bytes.Compare(buf, hash[:]) == 0 {
		// hash 一致，无需写入
		log.J(`skip write`, file, len(ab), "\n")
		return nil
	}
	log.J(`write file`, file, len(ab), "\n")

	fh, err := os.CreateTemp(path.Dir(file), `.pac-*.tmp`)
	if err != nil {
		return err
	}
	tmpName := fh.Name()
	fh.Chmod(FileMode)

	fd := int(fh.Fd())
	unix.Fsetxattr(fd, xattrHashKey, hash[:], 0)

	fh.Write(ab)
	fh.Close()

	err = os.Rename(tmpName, file)
	if err != nil {
		os.Remove(tmpName)
	}
	return err
}
