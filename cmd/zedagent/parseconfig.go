package main

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/zededa/api/zconfig"
	"github.com/zededa/go-provision/types"
	"io/ioutil"
	"log"
	"net"
	"os"
	"reflect"
	"strings"
)

func parseConfig(config *zconfig.EdgeDevConfig) {

	log.Println("Applying new config")
	parseBaseOsConfig(config)
	parseAppInstanceConfig(config)
}

func parseBaseOsConfig(config *zconfig.EdgeDevConfig) {

	log.Println("Applying Base Os config")

	cfgOsList := config.GetBase()

	if len(cfgOsList) != 0 {

		baseOsList := make([]types.BaseOsConfig, len(cfgOsList))

		for idx, cfgOs := range cfgOsList {

			baseOs := new(types.BaseOsConfig)

			baseOs.UUIDandVersion.UUID, _ = uuid.FromString(cfgOs.Uuidandversion.Uuid)
			baseOs.UUIDandVersion.Version = cfgOs.Uuidandversion.Version

			baseOs.Activate = cfgOs.GetActivate()
			baseOs.BaseOsVersion = cfgOs.GetBaseOSVersion()

			cfgOsDetails := cfgOs.GetBaseOSDetails()
			cfgOsParamList := cfgOsDetails.GetBaseOSParams()

			for jdx, cfgOsDetail := range cfgOsParamList {
				param := new(types.OsVerParams)
				param.OSVerKey = cfgOsDetail.GetOSVerKey()
				param.OSVerValue = cfgOsDetail.GetOSVerValue()
				baseOs.OsParams[jdx] = *param
			}

			var imageCount int
			for _, drive := range cfgOs.Drives {
				if drive.Image != nil {
					imageId := drive.Image.DsId

					for _, dsEntry := range config.Datastores {
						if dsEntry.Id == imageId {
							imageCount++
							break
						}
					}
				}
			}

			if imageCount != 0 {
				baseOs.StorageConfigList = make([]types.StorageConfig, imageCount)
				parseStorageConfigList(config, baseOsObj, baseOs.StorageConfigList,
					cfgOs.Drives)
			}

			baseOsList[idx] = *baseOs

			getCerts(baseOs.StorageConfigList)

			// Dump the config content
			bytes, err := json.Marshal(baseOs)

			if err == nil {
				log.Printf("New/updated BaseOs %d: %s\n", idx, bytes)
			}
		}

		if validateBaseOsConfig(baseOsList) == true {
			createBaseOsConfig(baseOsList)
		}
	}
}

func parseAppInstanceConfig(config *zconfig.EdgeDevConfig) {

	var appInstance = types.AppInstanceConfig{}

	log.Println("Applying App Instance config")

	Apps := config.GetApps()

	for _, cfgApp := range Apps {

		log.Printf("New/updated app instance %v\n", cfgApp)

		appInstance.UUIDandVersion.UUID, _ = uuid.FromString(cfgApp.Uuidandversion.Uuid)
		appInstance.UUIDandVersion.Version = cfgApp.Uuidandversion.Version
		appInstance.DisplayName = cfgApp.Displayname
		appInstance.Activate = cfgApp.Activate

		appInstance.FixedResources.Kernel = cfgApp.Fixedresources.Kernel
		appInstance.FixedResources.BootLoader = cfgApp.Fixedresources.Bootloader
		appInstance.FixedResources.Ramdisk = cfgApp.Fixedresources.Ramdisk
		appInstance.FixedResources.MaxMem = int(cfgApp.Fixedresources.Maxmem)
		appInstance.FixedResources.Memory = int(cfgApp.Fixedresources.Memory)
		appInstance.FixedResources.RootDev = cfgApp.Fixedresources.Rootdev
		appInstance.FixedResources.VCpus = int(cfgApp.Fixedresources.Vcpus)

		var imageCount int
		for _, drive := range cfgApp.Drives {
			if drive.Image != nil {
				imageId := drive.Image.DsId

				for _, dsEntry := range config.Datastores {
					if dsEntry.Id == imageId {
						imageCount++
						break
					}
				}
			}
		}

		if imageCount != 0 {
			appInstance.StorageConfigList = make([]types.StorageConfig, imageCount)
			parseStorageConfigList(config, appImgObj, appInstance.StorageConfigList,
				cfgApp.Drives)
		}

		// fill the overlay/underlay config
		parseNetworkConfig(&appInstance, cfgApp, config.Networks)

		// get the certs for image sha verification
		getCerts(appInstance.StorageConfigList)

		if validateAppInstanceConfig(appInstance) == true {

			// write to zedmanager config directory
			appFilename := cfgApp.Uuidandversion.Uuid

			writeAppInstanceConfig(appInstance, appFilename)
		}
	}
}

func parseStorageConfigList(config *zconfig.EdgeDevConfig, objType string,
	storageList []types.StorageConfig,
	drives []*zconfig.Drive) {

	var idx int = 0

	for _, drive := range drives {

		found := false

		image := new(types.StorageConfig)
		for _, ds := range config.Datastores {

			if drive.Image != nil &&
				drive.Image.DsId == ds.Id {

				found = true
				image.DownloadURL = ds.Fqdn + "/" + ds.Dpath + "/" + drive.Image.Name
				image.TransportMethod = ds.DType.String()
				image.ApiKey = ds.ApiKey
				image.Password = ds.Password
				image.Dpath = ds.Dpath
				break
			}
		}

		if found == false {
			continue
		}

		image.Format = strings.ToLower(drive.Image.Iformat.String())
		image.MaxSize = uint(drive.Maxsize)
		image.ReadOnly = drive.Readonly
		image.Preserve = drive.Preserve
		image.Target = strings.ToLower(drive.Target.String())
		image.Devtype = strings.ToLower(drive.Drvtype.String())
		image.ImageSignature = drive.Image.Siginfo.Signature
		image.ImageSha256 = drive.Image.Sha256
		image.ObjType = objType

		// copy the certificates
		if drive.Image.Siginfo.Signercerturl != "" {
			image.SignatureKey = drive.Image.Siginfo.Signercerturl
			image.NeedVerification = true
		}

		// XXX:FIXME certificates can be many
		// this list, currently contains the certUrls
		// should be the sha/uuid of cert filenames
		// as proper DataStore Entries

		if drive.Image.Siginfo.Intercertsurl != "" {
			image.CertificateChain = make([]string, 1)
			image.CertificateChain[0] = drive.Image.Siginfo.Intercertsurl
		}

		storageList[idx] = *image
		idx++
	}
}

func parseNetworkConfig(appInstance *types.AppInstanceConfig,
	cfgApp *zconfig.AppInstanceConfig,
	cfgNetworks []*zconfig.NetworkConfig) {

	var ulMaxIdx int = 0
	var olMaxIdx int = 0

	// count the interfaces and allocate
	for _, intfEnt := range cfgApp.Interfaces {
		for _, netEnt := range cfgNetworks {

			if intfEnt.NetworkId == netEnt.Id {

				switch strings.ToLower(netEnt.Type.String()) {
				// underlay interface
				case "v4", "v6":
					{
						ulMaxIdx++
						break
					}
				// overlay interface
				case "lisp":
					{
						olMaxIdx++
						break
					}
				}
			}
		}
	}

	if ulMaxIdx != 0 {
		appInstance.UnderlayNetworkList = make([]types.UnderlayNetworkConfig, ulMaxIdx)
		parseUnderlayNetworkConfig(appInstance, cfgApp, cfgNetworks)
	}

	if olMaxIdx != 0 {
		appInstance.OverlayNetworkList = make([]types.EIDOverlayConfig, olMaxIdx)
		parseOverlayNetworkConfig(appInstance, cfgApp, cfgNetworks)
	}
}

func parseUnderlayNetworkConfig(appInstance *types.AppInstanceConfig,
	cfgApp *zconfig.AppInstanceConfig,
	cfgNetworks []*zconfig.NetworkConfig) {

	var ulIdx int = 0

	for _, intfEnt := range cfgApp.Interfaces {
		for _, netEnt := range cfgNetworks {

			if intfEnt.NetworkId == netEnt.Id &&
				(strings.ToLower(netEnt.Type.String()) == "v4" ||
					strings.ToLower(netEnt.Type.String()) == "v6") {

				nv4 := netEnt.GetNv4() //XXX not required now...
				if nv4 != nil {
					booValNv4 := nv4.Dhcp
					log.Println("booValNv4: ", booValNv4)
				}
				nv6 := netEnt.GetNv6() //XXX not required now...
				if nv6 != nil {
					booValNv6 := nv6.Dhcp
					log.Println("booValNv6: ", booValNv6)
				}

				ulCfg := new(types.UnderlayNetworkConfig)
				ulCfg.ACLs = make([]types.ACE, len(intfEnt.Acls))

				for aclIdx, acl := range intfEnt.Acls {

					aclCfg := new(types.ACE)
					aclCfg.Matches = make([]types.ACEMatch, len(acl.Matches))
					aclCfg.Actions = make([]types.ACEAction, len(acl.Actions))

					for matchIdx, match := range acl.Matches {
						matchCfg := new(types.ACEMatch)
						matchCfg.Type = match.Type
						matchCfg.Value = match.Value
						aclCfg.Matches[matchIdx] = *matchCfg
					}

					for actionIdx, action := range acl.Actions {
						actionCfg := new(types.ACEAction)
						actionCfg.Limit = action.Limit
						actionCfg.LimitRate = int(action.Limitrate)
						actionCfg.LimitUnit = action.Limitunit
						actionCfg.LimitBurst = int(action.Limitburst)
						// XXX:FIXME actionCfg.Drop = <TBD>
						aclCfg.Actions[actionIdx] = *actionCfg
					}
					ulCfg.ACLs[aclIdx] = *aclCfg
				}
				appInstance.UnderlayNetworkList[ulIdx] = *ulCfg
				ulIdx++
			}
		}
	}
}

func parseOverlayNetworkConfig(appInstance *types.AppInstanceConfig,
	cfgApp *zconfig.AppInstanceConfig,
	cfgNetworks []*zconfig.NetworkConfig) {
	var olIdx int = 0

	for _, intfEnt := range cfgApp.Interfaces {
		for _, netEnt := range cfgNetworks {

			if intfEnt.NetworkId == netEnt.Id &&
				strings.ToLower(netEnt.Type.String()) == "lisp" {

				olCfg := new(types.EIDOverlayConfig)
				olCfg.ACLs = make([]types.ACE, len(intfEnt.Acls))

				for aclIdx, acl := range intfEnt.Acls {

					aclCfg := new(types.ACE)
					aclCfg.Matches = make([]types.ACEMatch, len(acl.Matches))
					aclCfg.Actions = make([]types.ACEAction, len(acl.Actions))

					for matchIdx, match := range acl.Matches {
						matchCfg := new(types.ACEMatch)
						matchCfg.Type = match.Type
						matchCfg.Value = match.Value
						aclCfg.Matches[matchIdx] = *matchCfg
					}

					for actionIdx, action := range acl.Actions {
						actionCfg := new(types.ACEAction)
						actionCfg.Limit = action.Limit
						actionCfg.LimitRate = int(action.Limitrate)
						actionCfg.LimitUnit = action.Limitunit
						actionCfg.LimitBurst = int(action.Limitburst)
						aclCfg.Actions[actionIdx] = *actionCfg
					}
					olCfg.ACLs[aclIdx] = *aclCfg
				}

				olCfg.EIDConfigDetails.EID = net.ParseIP(intfEnt.Addr)
				olCfg.EIDConfigDetails.LispSignature = intfEnt.Lispsignature
				olCfg.EIDConfigDetails.PemCert = intfEnt.Pemcert
				olCfg.EIDConfigDetails.PemPrivateKey = intfEnt.Pemprivatekey

				nlisp := netEnt.GetNlisp()

				if nlisp != nil {

					if nlisp.Eidalloc != nil {

						olCfg.EIDConfigDetails.IID = nlisp.Iid
						olCfg.EIDConfigDetails.EIDAllocation.Allocate = nlisp.Eidalloc.Allocate
						olCfg.EIDConfigDetails.EIDAllocation.ExportPrivate = nlisp.Eidalloc.Exportprivate
						olCfg.EIDConfigDetails.EIDAllocation.AllocationPrefix = nlisp.Eidalloc.Allocationprefix
						olCfg.EIDConfigDetails.EIDAllocation.AllocationPrefixLen = int(nlisp.Eidalloc.Allocationprefixlen)
					}

					if len(nlisp.Nmtoeid) != 0 {

						olCfg.NameToEidList = make([]types.NameToEid, len(nlisp.Nmtoeid))

						for nameIdx, nametoeid := range nlisp.Nmtoeid {

							nameCfg := new(types.NameToEid)
							nameCfg.HostName = nametoeid.Hostname
							nameCfg.EIDs = make([]net.IP, len(nametoeid.Eids))

							for eIdx, eid := range nametoeid.Eids {
								nameCfg.EIDs[eIdx] = net.ParseIP(eid)
							}

							olCfg.NameToEidList[nameIdx] = *nameCfg
						}
					}
				} else {
					log.Printf("No Nlisp in for %v\n", netEnt.Id)
				}

				appInstance.OverlayNetworkList[olIdx] = *olCfg
				olIdx++
			}
		}
	}
}

func writeAppInstanceConfig(appInstance types.AppInstanceConfig,
	appFilename string) {

	log.Printf("Writing app instance UUID %s\n", appFilename)
	bytes, err := json.Marshal(appInstance)
	if err != nil {
		log.Fatal(err, "json Marshal AppInstanceConfig")
	}
	configFilename := zedmanagerConfigDirname + "/" + appFilename + ".json"
	err = ioutil.WriteFile(configFilename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func writeBaseOsConfig(baseOsConfig types.BaseOsConfig, baseOsFilename string) {

	bytes, err := json.Marshal(baseOsConfig)

	if err != nil {
		log.Fatal(err, "json Marshal BaseOsConfig")
	}

	log.Printf("Writing baseOs config UUID %s, %s\n", baseOsFilename, bytes)

	configFilename := zedagentBaseOsConfigDirname + "/" + baseOsFilename + ".json"

	err = ioutil.WriteFile(configFilename, bytes, 0644)

	if err != nil {
		log.Fatal(err)
	}
}

func writeBaseOsStatus(baseOsStatus *types.BaseOsStatus,
	statusFilename string) {

	log.Printf("Writing baseOs status UUID %s\n", statusFilename)

	bytes, err := json.Marshal(baseOsStatus)
	if err != nil {
		log.Fatal(err, "json Marshal BaseOsStatus")
	}

	err = ioutil.WriteFile(statusFilename, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func getCerts(drives []types.StorageConfig) {

	for _, image := range drives {

		writeCertConfig(image, image.SignatureKey)

		for _, certUrl := range image.CertificateChain {
			writeCertConfig(image, certUrl)
		}
	}
}

func writeCertConfig(image types.StorageConfig, certUrl string) {

	if certUrl == "" {
		return
	}

	var safename = types.UrlToSafename(certUrl, "")

	var configFilename = fmt.Sprintf("%s/%s.json", certConfigDirname, safename)

	// XXX:FIXME dpath/key/pwd from image storage
	// should be coming from Drive
	// also the sha for the cert should be set
	// for now shortcutting it to only downloader
	var config = &types.DownloaderConfig{
		Safename:         safename,
		DownloadURL:      certUrl,
		MaxSize:          image.MaxSize,
		TransportMethod:  image.TransportMethod,
		Dpath:            "zededa-cert-repo",
		ApiKey:           image.ApiKey,
		Password:         image.Password,
		ImageSha256:      "",
		ObjType:          certObj,
		NeedVerification: false,
		DownloadObjDir:   certsPendingDirname,
		FinalObjDir:      certsDirname,
		RefCount:         1,
	}

	validateAndWriteDownloaderConfig(*config, configFilename)
}

func validateAndWriteVerifierConfig(config types.VerifyImageConfig,
	configFilename string) {
	changed := false

	if _, err := os.Stat(configFilename); err != nil {
		changed = true
	} else {

		var dConfig types.DownloaderConfig

		log.Printf("%s file exists\n", configFilename)

		cb, err := ioutil.ReadFile(configFilename)
		if err == nil {

			if err := json.Unmarshal(cb, &dConfig); err != nil {
				changed = true
			}

			if !reflect.DeepEqual(dConfig, config) {
				fmt.Printf("configChanged\n", configFilename)
				changed = true
			}
		} else {
			changed = true
		}
	}

	if changed == true {

		bytes, err := json.Marshal(config)

		if err != nil {
			log.Fatal(err, "json Marshal downloaderConfig")
		}

		err = ioutil.WriteFile(configFilename, bytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func validateAndWriteDownloaderConfig(config types.DownloaderConfig,
	configFilename string) {
	changed := false

	if _, err := os.Stat(configFilename); err != nil {
		changed = true
	} else {

		var dConfig types.DownloaderConfig

		log.Printf("%s file exists\n", configFilename)

		cb, err := ioutil.ReadFile(configFilename)
		if err == nil {

			if err := json.Unmarshal(cb, &dConfig); err != nil {
				changed = true
			}

			if !reflect.DeepEqual(dConfig, config) {
				fmt.Printf("configChanged\n", configFilename)
				changed = true
			}
		} else {
			changed = true
		}
	}

	if changed == true {

		bytes, err := json.Marshal(config)

		if err != nil {
			log.Fatal(err, "json Marshal downloaderConfig")
		}

		err = ioutil.WriteFile(configFilename, bytes, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func validateBaseOsConfig(baseOsList []types.BaseOsConfig) bool {

	var osCount, activateCount int
	if len(baseOsList) > 2 {
		return false
	}

	//count base os instance activate count
	for _, baseOsInstance := range baseOsList {

		osCount++
		if baseOsInstance.Activate == true {
			activateCount++
		}
	}

	if osCount == 0 {
		return true
	}

	if osCount == 2 {
		if activateCount != 1 {
			return false
		}
	}

	return true
}

func createBaseOsConfig(baseOsList []types.BaseOsConfig) {

	for _, baseOsInstance := range baseOsList {

		baseOsFilename := baseOsInstance.UUIDandVersion.UUID.String()
		writeBaseOsConfig(baseOsInstance, baseOsFilename)
	}
}

func validateAppInstanceConfig(appInstance types.AppInstanceConfig) bool {
	return true
}
