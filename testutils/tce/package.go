package tce

import (
	"fmt"
	"os"

	"github.com/karuppiah7890/tce-e2e-test/testutils/clirunner"
	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
)

type Package struct {
	Name         string
	Version      string
	ManualCreate bool
}

func PackageE2Etest(packageDetails Package, workloadClusterKubeContext string) error {
	err := kubeclient.UseKubeConfigContext(workloadClusterKubeContext)
	if err != nil {
		return fmt.Errorf("error occurred while using the workload cluster context. error: %v", err)
	}
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"package",
			"repository",
			"add",
			"tce-repo",
			"--url",
			"projects.registry.vmware.com/tce/main:0.12.0",
			"--namespace",
			"tanzu-package-repo-global",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while adding package repo. exit code: %v. error: %v", exitCode, err)
	}

	//Prerequisites(packageDetails)
	if packageDetails.ManualCreate {
		err := InstallPackage(packageDetails)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	err = os.Chdir("community-edition/addons/packages/" + packageDetails.Name + "/" + packageDetails.Version + "/test")
	if err != nil {
		return fmt.Errorf("error while changing directory to community-edition: %v", err)
	}
	exitCode, err = clirunner.Run(clirunner.Cmd{
		Name: "make",
		Args: []string{
			"e2e-test",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})

	if packageDetails.ManualCreate {
		err := DeletePackage(packageDetails)
		if err != nil {
			return fmt.Errorf("%v", err)
		}

	}

	if err != nil {
		return fmt.Errorf("error occurred while E2E test for %v. Exit code: %v. Error: %v", packageDetails.Name, exitCode, err)
	}

	return nil
}

func InstallPackage(packageDetails Package) error {
	wd, _ := os.Getwd()
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"package",
			"install",
			packageDetails.Name,
			"--package-name",
			packageDetails.Name + ".community.tanzu.vmware.com",
			"--version", packageDetails.Version,
			"--values-file", wd + "/testutils/tce/testdata/" + packageDetails.Name + "_values.yaml",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while installing %v package. exit code: %v. error: %v", packageDetails.Name, exitCode, err)
	}
	return nil
}

func DeletePackage(packageDetails Package) error {
	exitCode, err := clirunner.Run(clirunner.Cmd{
		Name: "tanzu",
		Args: []string{
			"package",
			"installed",
			"delete",
			packageDetails.Name,
			"-y",
		},
		Stdout: log.InfoWriter,
		// TODO: Should we log standard errors as errors in the log? Because tanzu prints other information also
		// to standard error, which are kind of like information, apart from actual errors, so showing
		// everything as error is misleading. Gotta think what to do about this. The main problem is
		// console has only standard output and standard error, and tanzu is using standard output only for
		// giving output for things like --dry-run when it needs to print yaml content, but everything else
		// is printed to standard error
		Stderr: log.ErrorWriter,
	})
	if err != nil {
		return fmt.Errorf("error occurred while deleting %v package. exit code: %v. error: %v", packageDetails.Name, exitCode, err)
	}
	return nil
}

/*func TestPrerequisites(t *testing.T) {
	packageDetails := Package{}
	packageDetails.Name = "velero"   //os.Getenv("PACKAGE_NAME")
	packageDetails.Version = "1.8.0" //os.Getenv("PACKAGE_VERSION")
	packageDetails.ManualCreate = true
	wd, _ := os.Getwd()
	configFilePath := wd + "/testdata/" + packageDetails.Name + "_values.yaml"
	configTempFilePath := wd + "/testdata/" + packageDetails.Name + "_temp_values.yaml"
	if packageDetails.Name == "velero" {
		err := testutils.Copy(configFilePath, configTempFilePath)
		//t.Errorf("%v", err)
		file, err := ioutil.ReadFile(configFilePath)
		if err != nil {
			log.Fatal(err)
		}

		data := make(map[interface{}]interface{})
		error := yaml.Unmarshal([]byte(file), &data)
		if error != nil {
			fmt.Println("testingggg ", configTempFilePath, error)
			log.Fatal(err)
		}

		for key, value := range data {

			fmt.Println(key, " //// ", value)
		}

		//dbURL := os.ExpandEnv(data)
		//fmt.Println(data)
	}

	//return nil
}*/
