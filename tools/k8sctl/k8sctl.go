package main

import (
	"flag"
	"path/filepath"

	"github.com/karuppiah7890/tce-e2e-test/testutils/kubeclient"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"k8s.io/client-go/util/homedir"
)

func main() {
	log.InitLogger("k8sctl")

	// TODO: We can also use the $KUBECONFIG env var for kubeconfig.
	// But yeah flags take higher precedence! We gotta check what we support!

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	kclient, err := kubeclient.GetKubeClient(*kubeconfig, "")
	if err != nil {
		log.Fatalf("error getting kubernetes client: %v", err)
	}

	pods, err := kclient.GetAllPodsFromAllNamespaces()
	if err != nil {
		log.Fatalf("error getting all pods from all namespaces: %v", err)
	}

	for index, pod := range pods.Items {
		log.Infof("Pod %d: %s", index, pod.Name)
	}
}
