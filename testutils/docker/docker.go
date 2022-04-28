package docker

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

// TODO: Let's provide a Docker Client struct with `NewDockerClient` or `GetDockerClient`
// and struct methods on it so that we don't have to create a new client every time using
// client.NewClientWithOpts()
var ctx context.Context = context.Background()

func GetDockerClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		// TODO: Handle errors by returning them? Should we log them here too or let caller decide about the logging?
		log.Fatalf("error creating docker client: %v", err)
	}
	return cli
}

func CheckDockerInstallation() {

	log.Info("Checking docker CLI and Docker Engine installation")
	cli := GetDockerClient()
	path, err := exec.LookPath("docker")
	if err != nil {
		// TODO: Handle errors by returning them? Should we log them here too or let caller decide about the logging?
		log.Fatalf("docker CLI is not installed")
	}
	serverVersionInfo, err := cli.ServerVersion(context.TODO())
	if err != nil {
		// TODO: Should the below be info? And should we always show it or only when there's an error in
		// checking Docker Engine version?
		log.Warn("Ensure Docker Engine is installed and accessible")
		// TODO: Handle errors by returning them? Should we log them here too or let caller decide about the logging?
		log.Fatalf("Error checking Docker Engine version: %v .", err)
	}
	log.Infof("docker CLI is available at path: %s", path)
	log.Infof("E2E test Docker client's API version: %s", cli.ClientVersion())
	log.Infof("Docker Engine's API version: %s", serverVersionInfo.APIVersion)
	log.Infof("Docker Engine's version: %s", serverVersionInfo.Version)
}

// TODO: rename this to RemoveRunningContainer ? Or StopAndRemoveRunningContainer? or ForceRemoveRunningContainer?
// containerName - name of the container / container ID (full ID or unique partial ID)
func StopRunningContainer(containerName string) error {
	cli := GetDockerClient()

	// TODO: Handle errors by returning them? Should we log them here too or let caller decide about the logging?
	if err := cli.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{Force: true}); err != nil {
		log.Infof("Failed to find container with  name: %s", containerName)
		return fmt.Errorf("failed to find container with  name: %s", containerName)
	}
	log.Infof("Container removed: %s", containerName)
	return nil
}

// TODO: rename this to RemoveAllRunningContainers ? Or StopAndRemoveAllRunningContainers? or ForceRemoveAllRunningContainers?
func ForceRemoveAllRunningContainers() {
	cli := GetDockerClient()
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		log.Infof("Failed to List Containers: %s", err)
	}
	log.Infof("Containers %s", containers)
	for _, container := range containers {
		// TODO: Handle errors by returning them? Should we log them here too or let caller decide about the logging?
		if err := cli.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{Force: true}); err != nil {
			log.Infof("Failed to find container with  name: %s", container.Names)
		}
		log.Infof("Container removed: %s", container.Names)
	}
}
