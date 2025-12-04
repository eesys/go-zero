package main

import (
	"github.com/eesys/go-zero/core/load"
	"github.com/eesys/go-zero/core/logx"
	"github.com/eesys/go-zero/tools/goctl/cmd"
)

func main() {
	logx.Disable()
	load.Disable()
	cmd.Execute()
}
