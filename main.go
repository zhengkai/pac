package main

import (
	"github.com/zhengkai/pac/core"
)

func main() {

	if !core.Flag() {
		return
	}

	core.Run()
}
