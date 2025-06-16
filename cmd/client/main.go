package main

import (
	"github.com/gen-4/gorrent/config/client"
	"github.com/gen-4/gorrent/internal/client/ui"
)

func main() {
	config.Config()
	defer config.CloseConfig()
	ui.Run()
}
