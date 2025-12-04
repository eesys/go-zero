package config

import "github.com/eesys/go-zero/core/logx"

// Config defines a service configure for goctl update
type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}
