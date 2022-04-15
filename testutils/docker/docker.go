package docker

import (
	"context"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func CheckDockerInstallation() {

	log.Info("Checking docker CLI and Docker Engine installation")
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalf("error creating docker client: %v", err)
	}
	path, err := exec.LookPath("docker")
	if err != nil {
		log.Fatalf("docker CLI is not installed")
	}
	serverVersionInfo, err := cli.ServerVersion(context.TODO())
	if err != nil {
		// TODO: Should the below be info? And should we always show it or only when there's an error in
		// checking Docker Engine version?
		log.Warn("Ensure Docker Engine is installed and accessible")
		log.Fatalf("Error checking Docker Engine version: %v .", err)
	}
	log.Infof("docker CLI is available at path: %s", path)
	log.Infof("E2E test Docker client's API version: %s", cli.ClientVersion())
	log.Infof("Docker Engine's API version: %s", serverVersionInfo.APIVersion)
	log.Infof("Docker Engine's version: %s", serverVersionInfo.Version)
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
