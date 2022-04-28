package vsphere

import (
	"context"
	"fmt"
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/govc/importx"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/library"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vapi/vcenter"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"net/url"
	"os"
)

//TODO: Use VSPHERE_URL style Env Vars for now used to make use of same env vars what govc uses
const (
	envURL          = "GOVC_URL"
	envUserName     = "GOVC_USERNAME"
	envPassword     = "GOVC_PASSWORD"
	envDataStore    = "GOVC_DATASTORE"
	envFolder       = "GOVC_FOLDER"
	envNetwork      = "GOVC_NETWORK"
	envResourcePool = "GOVC_RESOURCE_POOL"
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

func GetLibraryItem(rc *rest.Client) ([]string, error) {
	const (
		libraryName = "test"
	)
	m := library.NewManager(rc)
	libraries, err := m.FindLibrary(ctx, library.Find{Name: libraryName})
	if err != nil {
		fmt.Printf("Find library by name %s failed, %v", libraryName, err)
	}

	return libraries, nil
}

func CreateLibrary(libraryName string, rc *rest.Client, client *vim25.Client) {
	m := library.NewManager(rc)
	envDataStore := os.Getenv(envDataStore)
	ds, err := find.NewFinder(client).Datastore(ctx, envDataStore)
	if err != nil {
		log.Info(err)
	}
	res, err := m.CreateLibrary(ctx, library.Library{
		Name: libraryName,
		Type: "LOCAL",
		Storage: []library.StorageBackings{{
			DatastoreID: ds.Reference().Value,
			Type:        "DATASTORE",
		}},
		Publication: &library.Publication{
			AuthenticationMethod: "NONE",
			Published:            &[]bool{true}[0],
		},
	})
	if err != nil {
		log.Fatalf("Unable to create Library %s", err)
	}
	l, err := m.GetLibraryByID(ctx, res)
	if err != nil {
		log.Fatalf("Unable to create Library %s", err)
	}

	log.Infof("Library created Name : %s and ID :", l.Name, res)

}

// To Get Library as Library Struct
func GetLibrary(libraryName string, rc *rest.Client) *library.Library {
	m := library.NewManager(rc)
	res, err := m.GetLibraryByName(ctx, libraryName)
	if err != nil {
		log.Errorf("Unable to create Library")
	}
	log.Infof("Library created %s ", res.Name)
	return res
}

// To Upload OVA to Library
func ImportLibrary(rc *rest.Client, client *vim25.Client, item *library.Library, file string) error {

	//base := filepath.Base(file)
	//ext := filepath.Ext(base)
	//mf := strings.Replace(base, ext, ".mf", 1)
	kind := library.ItemTypeOVF
	opener := importx.Opener{Client: client}
	archive := &importx.ArchiveFlag{Archive: &importx.FileArchive{Path: file, Opener: opener}}
	archive.Archive = &importx.TapeArchive{Path: file, Opener: opener}

	m := library.NewManager(rc)
	lib, err := m.CreateLibraryItem(ctx, library.Item{Name: item.Name, ID: item.ID, Type: kind})
	if err != nil {
		log.Errorf("%s", err)
	}
	session, err := m.CreateLibraryItemUpdateSession(ctx, library.Session{
		ID: lib,
	})
	return m.CompleteLibraryItemUpdateSession(ctx, session)
}

// To Deploy Vm from Library will be used similar to
// command `govc library.deploy -folder=/SDDC-Datacenter/vm/VMs-tce-test/ /tce-test/photon-3-kube-v1.22.8+vmware.1-tkg.1-d69148b2a4aa7ef6d5380cc365cac8cd`

func DeployVmFromLibrary(rc *rest.Client, client *vim25.Client, item *library.Library) (*object.VirtualMachine, error) {
	//TODO to use common Env Vars
	envResourcePool := os.Getenv(envResourcePool)
	envFolder := os.Getenv(envFolder)
	envNetwork := os.Getenv(envNetwork)
	envDataStore := os.Getenv(envDataStore)

	finder := find.NewFinder(client)
	resourcePools, err := finder.ResourcePoolList(ctx, envResourcePool)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list resource pool at vc  %v", err)
		os.Exit(1)
	}
	//hosts, err := finder.HostSystemList(ctx, "*")
	datastores, err := finder.DatastoreList(ctx, envDataStore)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list datastore at vc , %v", err)
		os.Exit(1)
	}

	networks, err := finder.NetworkList(ctx, envNetwork)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list network at vc  %v", err)
		os.Exit(1)
	}

	folders, err := finder.FolderList(ctx, envFolder)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to list folder at vc  %v", err)
		os.Exit(1)
	}

	m := vcenter.NewManager(rc)
	fr := vcenter.FilterRequest{Target: vcenter.Target{
		ResourcePoolID: resourcePools[0].Reference().Value,
		FolderID:       folders[0].Reference().Value,
	},
	}
	r, err := m.FilterLibraryItem(ctx, item.ID, fr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FilterLibraryItem error, %v\n", err)
		os.Exit(1)
	}
	networkKey := r.Networks[0]
	storageKey := r.StorageGroups[0]

	deploy := vcenter.Deploy{
		DeploymentSpec: vcenter.DeploymentSpec{
			Name:               "test-ova-vm",
			DefaultDatastoreID: datastores[0].Reference().Value,
			AcceptAllEULA:      true,
			NetworkMappings: []vcenter.NetworkMapping{{
				Key:   networkKey,
				Value: networks[0].Reference().Value,
			}},
			StorageMappings: []vcenter.StorageMapping{{
				Key: storageKey,
				Value: vcenter.StorageGroupMapping{
					Type:         "DATASTORE",
					DatastoreID:  datastores[0].Reference().Value,
					Provisioning: "thin",
				},
			}},
			StorageProvisioning: "thin",
		},
		Target: vcenter.Target{
			ResourcePoolID: resourcePools[0].Reference().Value,
			FolderID:       folders[0].Reference().Value,
		},
	}

	ref, err := vcenter.NewManager(rc).DeployLibraryItem(ctx, item.ID, deploy)
	if err != nil {
		fmt.Printf("Deploy vm from library failed, %v", err)
		return nil, err
	}

	f := find.NewFinder(client)
	obj, err := f.ObjectReference(ctx, *ref)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Find vm failed, %v\n", err)
		os.Exit(1)
	}
	vm := obj.(*object.VirtualMachine)
	return vm, nil

}

func MarkAsTemplate(client *vim25.Client, vmName string) error {
	finder := find.NewFinder(client)
	vms, err := finder.VirtualMachine(context.TODO(), vmName)
	if err != nil {
		if _, ok := err.(*find.NotFoundError); ok {
			log.Errorf("Unable To find VM")
			return err
		}
	}
	errs := vms.MarkAsTemplate(ctx)
	if errs != nil {
		log.Errorf("Unable To Make Vm As template")
		return errs
	}
	return nil
}

// To Delete Libary after create OVF templates
func DeleteLibrary(libraryName *library.Library, rc *rest.Client) {
	m := library.NewManager(rc)
	err := m.DeleteLibrary(ctx, libraryName)
	if err != nil {
		log.Errorf("Unable to Delete Library")
	}
	log.Infof("Library Deleted %s Successfully ", libraryName.Name)
}
