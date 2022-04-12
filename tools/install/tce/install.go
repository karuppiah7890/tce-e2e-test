package main

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/tce"
)

func main() {
	log.InitLogger("tce-install")
	tce.Install("0.11.0")
}
