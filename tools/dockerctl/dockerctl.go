package main

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils/docker"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Maybe we should start writing tests instead of writing demo tools etc
// to do our testing?

func main() {
	log.InitLogger("dockerctl")
	docker.CheckDockerInstallation()
	docker.ForceRemoveAllRunningContainers()
	docker.ForceRemoveRunningContainer("kind")
}
