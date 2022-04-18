package clirunner

import (
	"os/exec"

	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

func Run(command Cmd) (int, error) {
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Stdout = command.Stdout
	cmd.Stderr = command.Stderr
	cmd.Env = command.Env

	// TODO: Maybe set cmd.Env explicitly to a narrow set of env vars to just inject the secrets
	// that we want to inject and nothing else. But system level env vars maybe needed for the CLI.
	// Think about how to inject the env vars. Use single struct as function argument?
	// Use a struct to store the data including env vars and use that as data/context in
	// it's methods where method runs the command by injecting the env vars
	// Or something like
	// Tanzu({ env: []string{"key=value", "key2=value2"}, command: "management-cluster version" })
	// But the above is not exactly readable, hmm

	log.Infof("Running the command `%v`", cmd.String())

	err := cmd.Run()
	if err != nil {
		// TODO: Handle the error by returning it?
		log.Infof("Error occurred while running the command `%v`: %v", cmd.String(), err)
		return cmd.ProcessState.ExitCode(), err
	}

	return cmd.ProcessState.ExitCode(), nil
}
