package docker

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"os/exec"
)

func CheckDockerInstallation() {

	log.Info("Checking docker CLI and Engine installation")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("docker Engine is not installed")
	}
	path, err := exec.LookPath("docker")
	if err != nil {
		log.Fatalf("Docker CLI is not installed")
	}
	log.Infof("Docker cli is available at path: %s", path)
	log.Infof("Docker version %s", cli.ClientVersion())
}

func StopRunningContainer(containerName string) {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Infof("Failed to find container with  name: %s", containerName)
	}
	log.Infof("Container stopped: %s", containerName)

}
func StopAllRunningContainer() {
	ctx := context.Background()
	cli, _ := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}
	log.Infof("Containers %s", containers)
	for _, container := range containers {
		if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Infof("Failed to find container with  name: %s", container.Names)
		}
		log.Infof("Container stopped: %s", container.Names)
	}
}
