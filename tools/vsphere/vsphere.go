package main

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/vsphere"
	"path/filepath"
)

// This is just for Checking and using govmomi SDK  and getting VM names
func main() {
	log.InitLogger("vSphere-OVA-Upload-Testing")
	tce := "0120"
	//Setting Vsphere Clients
	dir := "/tmp/"
	fileName := ""
	client := vsphere.GetGovmomiClient()
	rs := vsphere.GetRestClient(client)
	requiredOvaFile := vsphere.GetOvaFileNameFromTanzuFramework()
	log.Info(requiredOvaFile)
	vmTemplates := vsphere.ListVmsTemplates(client)
	filesAvailableToDownload := vsphere.RetriveVersion(tce)
	for _, file := range requiredOvaFile {
		for _, template := range vmTemplates {
			if file == template {
				log.Infof("File Already present on VC as %s", template)
			} else {
				for _, download := range filesAvailableToDownload {
					if file == download {
						vsphere.RetrieveAndDownload(tce, dir, file)
						fileName = file
					}
				}

			}
		}

	}
	//vsphere.RetrieveAndDownload(tce, "photon-3-kube-v1.22.8+vmware.1-tkg.1-d69148b2a4aa7ef6d5380cc365cac8cd.ova")
	log.Info("Retrived info")
	//vsphere.ListVMs()

	//for _, y := range vmTemplates {
	//	log.Info(y)
	//}
	vsphere.CreateLibrary("test", rs, client)
	lib := vsphere.GetLibrary("test", rs)
	err := vsphere.ImportOVAFromLibrary(rs, client, lib, filepath.Join(dir+fileName))
	if err != nil {
		log.Errorf("%s", err)
	}
	vm, err := vsphere.DeployVmFromLibrary(rs, client, lib)
	vm.Name()
	vsphere.MarkAsTemplate(client, "testing")
	for _, vm := range vsphere.ListVms(client) {
		//err:=vsphere.DeleteVM(client,vm)
		log.Infof("%s Vm will be deleted", vm)
		//if err != nil {
		//	log.Errorf("Unable to delete Vm")
		//}
	}
	vsphere.DeleteLibrary(lib, rs)
	item, err := vsphere.GetLibraryItem(rs)
	if err != nil {
		log.Errorf("something went wrong")
	}
	log.Info(item)

}
