package main

import (
	server "github.com/snowdrop/generator/pkg/server"
)

var (
	// VERSION is set during build
	VERSION string

	// GITCOMMIT is hash of the commit that wil be displayed when running ./odo version
	// this will be overwritten when running  build like this: go build -ldflags="-X github.com/redhat-developer/odo/cmd.GITCOMMIT=$(GITCOMMIT)"
	// HEAD is default indicating that this was not set during build
	GITCOMMIT = "HEAD"
)

func main() {
	server.Run(VERSION, GITCOMMIT)
}
