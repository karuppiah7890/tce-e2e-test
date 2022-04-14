package kubescheme

import (
	kubeRuntime "k8s.io/apimachinery/pkg/runtime"
	// TODO: Rename imports in a better manner, like, use capzABC and capiDEF etc for naming
	// CAPZ and CAPI stuff

	infrav1alpha3 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha3"
	infrav1alpha4 "sigs.k8s.io/cluster-api-provider-azure/api/v1alpha4"
	capzv1beta1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	infrav1alpha3exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1alpha3"
	infrav1alpha4exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1alpha4"
	infrav1beta1exp "sigs.k8s.io/cluster-api-provider-azure/exp/api/v1beta1"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	clusterv1alpha3 "sigs.k8s.io/cluster-api/api/v1alpha3"
	clusterv1alpha4 "sigs.k8s.io/cluster-api/api/v1alpha4"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

	capiBootstrapKubeadmv1beta "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1beta1"
	capiControlplaneKubeadmv1beta "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1beta1"

	expv1alpha3 "sigs.k8s.io/cluster-api/exp/api/v1alpha3"
	expv1alpha4 "sigs.k8s.io/cluster-api/exp/api/v1alpha4"
	expv1 "sigs.k8s.io/cluster-api/exp/api/v1beta1"

	addonsv1alpha3 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha3"
	addonsv1alpha4 "sigs.k8s.io/cluster-api/exp/addons/api/v1alpha4"
	addonsv1 "sigs.k8s.io/cluster-api/exp/addons/api/v1beta1"
)

var scheme = kubeRuntime.NewScheme()

func init() {
	registerCustomResources()
}

func GetScheme() *kubeRuntime.Scheme {
	return scheme
}

func registerCustomResources() {
	err := clientgoscheme.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = apiextensionsv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = clusterv1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = clusterv1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = clusterv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = expv1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = expv1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = expv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = addonsv1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = addonsv1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = addonsv1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = infrav1alpha3.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1alpha4.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = capzv1beta1.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1alpha3exp.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1alpha4exp.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
	err = infrav1beta1exp.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = capiControlplaneKubeadmv1beta.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}

	err = capiBootstrapKubeadmv1beta.AddToScheme(scheme)
	if err != nil {
		panic(err)
	}
}
