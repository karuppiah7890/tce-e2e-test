package vsphere

import (
	"context"
	"fmt"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"net/url"
	"os"
)

//TODO: Use VSPHERE_URL style Env Vars for now used to make use of same env vars what govc uses
const (
	envURL      = "GOVC_URL"
	envUserName = "GOVC_USERNAME"
	envPassword = "GOVC_PASSWORD"
)

var ctx context.Context = context.Background()

/*
TODO: Idea for OVA upload
	* create Library
		`govc library.info test`
	* upload OVA
		`govc library.import test photon-3-kube-v1.22.8+vmware.1-tkg.1-d69148b2a4aa7ef6d5380cc365cac8cd.ova`
	* deploy Vm from Library
		To confirm on command
		`govc library.deploy -folder=/SDDC-Datacenter/vm/VMs-tce-test/ /tce-test/photon-3-kube-v1.22.8+vmware.1-tkg.1-d69148b2a4aa7ef6d5380cc365cac8cd`
	* make Vm as template
		`govc vm.markastemplate /SDDC-Datacenter/vm/VMs-tce-test/photon-3-kube-v1.22.8`
	* delete Library
		`govc library.rm test`
*/

// NewClient creates a vim25.Client
func GetGovmomiClient() *vim25.Client {
	//TODO: To make use of common creds function or struct to avoid redundant vars
	envUserName := os.Getenv(envUserName)
	envPassword := os.Getenv(envPassword)
	envURL := os.Getenv(envURL)
	u := &url.URL{
		Scheme: "https",
		Host:   envURL,
		Path:   "/sdk",
	}
	u.User = url.UserPassword(envUserName, envPassword)
	client, err := govmomi.NewClient(ctx, u, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Login to vsphere failed, %v", err)
		os.Exit(1)
	}
	return client.Client
}

// Rest Client this is being used by library module
func GetRestClient(client *vim25.Client) *rest.Client {
	//TODO: To make use of common creds function or struct to avoid redundant vars
	envUserName := os.Getenv(envUserName)
	envPassword := os.Getenv(envPassword)
	envURL := os.Getenv(envURL)
	u := &url.URL{
		Scheme: "https",
		Host:   envURL,
		Path:   "/sdk",
	}
	u.User = url.UserPassword(envUserName, envPassword)
	rc := rest.NewClient(client)
	if err := rc.Login(ctx, u.User); err != nil {
		fmt.Fprintf(os.Stderr, "rc.Login failed, %v", err)
		os.Exit(1)
	}
	return rc
}

func ListVmsTemplates(client *vim25.Client) []string {
	m := view.NewManager(client)
	v, err := m.CreateContainerView(ctx, client.ServiceContent.RootFolder, []string{"VirtualMachine"}, true)
	if err != nil {
		log.Errorf("%s", err)
	}
	defer v.Destroy(ctx)

	var vms []mo.VirtualMachine
	err = v.Retrieve(ctx, []string{"VirtualMachine"}, []string{"summary"}, &vms)
	if err != nil {
		log.Errorf("%s", err)
	}

	vmTemplates := []string{}
	for _, vm := range vms {
		if vm.Summary.Config.Template == true {
			//log.Infof("%s: %s", vm.Summary.Config.Name, vm.Summary.Config.GuestFullName)
			vmTemplates = append(vmTemplates, vm.Summary.Config.Name)
		}
	}
	return vmTemplates

}

func GetLibraryItem(rc *rest.Client) (*library.Item, error) {
	const (
		libraryName     = ""
		libraryItemName = ""
		libraryItemType = "ovf"
	)

	//  rc   library.Manager
	m := library.NewManager(rc)
	//m.CreateLibrary(ctx,library.Library{Name: "libraryName"})
	libraries, err := m.FindLibrary(ctx, library.Find{Name: libraryName})
	if err != nil {
		fmt.Printf("Find library by name %s failed, %v", libraryName, err)
		return nil, err
	}

	if len(libraries) == 0 {
		fmt.Printf("Library %s was not found", libraryName)
		return nil, fmt.Errorf("library %s was not found", libraryName)
	}

	if len(libraries) > 1 {
		fmt.Printf("There are multiple libraries with the name %s", libraryName)
		return nil, fmt.Errorf("there are multiple libraries with the name %s", libraryName)
	}

	//  ovf   ovf
	items, err := m.FindLibraryItems(ctx, library.FindItem{Name: libraryItemName,
		Type: libraryItemType, LibraryID: libraries[0]})

	if err != nil {
		fmt.Printf("Find library item by name %s failed", libraryItemName)
		return nil, fmt.Errorf("find library item by name %s failed", libraryItemName)
	}

	if len(items) == 0 {
		fmt.Printf("Library item %s was not found", libraryItemName)
		return nil, fmt.Errorf("library item %s was not found", libraryItemName)
	}

	if len(items) > 1 {
		fmt.Printf("There are multiple library items with the name %s", libraryItemName)
		return nil, fmt.Errorf("there are multiple library items with the name %s", libraryItemName)
	}

	item, err := m.GetLibraryItem(ctx, items[0])
	if err != nil {
		fmt.Printf("Get library item by %s failed, %v", items[0], err)
		return nil, err
	}

	return item, nil
}
