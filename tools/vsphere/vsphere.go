package main

import (
	"github.com/karuppiah7890/tce-e2e-test/testutils/log"
	"github.com/karuppiah7890/tce-e2e-test/testutils/vsphere"
)

// This is just for Checking and using govmomi SDK  and getting VM names
func main() {
	log.InitLogger("vSphere-OVA-Upload-Testing")

	ovaFiles := vsphere.GetOvaFileNameFromTanzuFramework()
	log.Info(ovaFiles)
	vsphere.RetriveAndDownlaod("0120", "photon-3-kube-v1.22.8+vmware.1-tkg.1-d69148b2a4aa7ef6d5380cc365cac8cd.ova")
	log.Info("Retrived info")
	//vsphere.ListVMs()
	client := vsphere.GetGovmomiClient()
	rs := vsphere.GetRestClient(client)
	vmTemplates := vsphere.ListVmsTemplates(client)

	for _, y := range vmTemplates {
		log.Info(y)
	}
	vsphere.CreateLibrary("test", rs, client)
	lib := vsphere.GetLibrary("test", rs)
	vsphere.ImportOVAFromLibrary(rs, client, lib, ovaFiles[0])
	vsphere.DeployVmFromLibrary(rs, client, lib)
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
	vsphere.RetriveVersion("0110")

}
