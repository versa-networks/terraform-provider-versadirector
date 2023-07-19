package vclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// https://10.40.73.242:9182/api/config/devices/device/Branch-1/config/orgs/org-services/ACME/objects/addresses
const (
	vmsDirectorDevicesURL          = "api/config/devices/device"
	vmsDirectorOrgServicesURL      = "config/orgs/org-services"
	vmsDirectorObjectsAddressesURL = "objects/addresses"
)

type DevObjectAddress struct {
	Name string `json:"name"`
	FQDN string `json:"fqdn,optional"`
}

type DevObjectsAddressList struct {
	DeviceName       string             `json:"-"`
	OrganizationName string             `json:"-"`
	Count            int                `json:"-"`
	Addresses        []DevObjectAddress `json:"address"`
}

type DevOjectsAddressListData struct {
	AddrList DevObjectsAddressList `json:"addresses"`
}

func (c *Client) CreateDevOrgServiceObjAddresses(ctx context.Context,
	addrListData DevOjectsAddressListData) error {

	addrList := addrListData.AddrList

	if addrList.Count <= 0 || len(addrList.Addresses) <= 0 {
		tflog.Trace(ctx, "Device Orgs Address creation failed as addresses count is 0")
		return errors.New("Device Orgs Address creation failed as addresses count is 0")
	}

	tflog.Trace(ctx, "Device-Name "+addrList.DeviceName+" OrgName "+addrList.OrganizationName)
	for key, val := range addrList.Addresses {
		tflog.Trace(ctx, "Address["+strconv.Itoa(key)+"]: Name "+val.Name+" FQDN "+val.FQDN)
	}

	client, _, _ := vHttpClient(c.Config.ServerIP, c.Config.ServerPort, "")
	httpUrl := "https://" + c.Config.ServerIP + ":" + strconv.Itoa(c.Config.ServerPort) + "/" +
		vmsDirectorDevicesURL + "/" +
		addrList.DeviceName + "/" +
		vmsDirectorOrgServicesURL + "/" +
		addrList.OrganizationName + "/" +
		vmsDirectorObjectsAddressesURL

	jsonData, err := json.Marshal(addrList)
	if err != nil {
		tflog.Error(ctx, "POST Addresses request failed, json marshal error: "+err.Error())
		return err
	}

	if _, err := c.vHttpHandlePostReq(ctx, client, httpUrl, jsonData, nil); err != nil {
		tflog.Error(ctx, "POST Addresses request failed for URL: "+httpUrl+" Error: "+err.Error())
		return err
	}

	return err
}

func (c *Client) UpdateDevOrgServiceObjAddresses(ctx context.Context,
	addrListData DevOjectsAddressListData) error {

	addrList := addrListData.AddrList

	if addrList.Count <= 0 || len(addrList.Addresses) <= 0 {
		tflog.Trace(ctx, "PUT addresses request failed as address count is 0")
		return errors.New("Device Orgs Address update failed as addresses count is 0")
	}

	tflog.Trace(ctx, "Device-Name "+addrList.DeviceName+" OrgName "+addrList.OrganizationName)
	for key, val := range addrList.Addresses {
		tflog.Trace(ctx, "Address["+strconv.Itoa(key)+"]: Name "+val.Name+" FQDN "+val.FQDN)
	}

	client, _, _ := vHttpClient(c.Config.ServerIP, c.Config.ServerPort, "")
	httpUrl := "https://" + c.Config.ServerIP + ":" + strconv.Itoa(c.Config.ServerPort) + "/" +
		vmsDirectorDevicesURL + "/" +
		addrList.DeviceName + "/" +
		vmsDirectorOrgServicesURL + "/" +
		addrList.OrganizationName + "/" +
		vmsDirectorObjectsAddressesURL + "/" +
		"address" + "/"

	// Modify expects individual objects, send one after another
	for _, val := range addrList.Addresses {
		curHttpUrl := httpUrl + val.Name
		var curAddrData DevObjectsAddressList
		curAddrData.Addresses = append(curAddrData.Addresses, val)
		if jsonData, err := json.Marshal(curAddrData); err != nil {
			tflog.Error(ctx, "PUT Addresses request failed for "+val.Name+" Error: "+err.Error())
			return err
		} else {
			if _, err := c.vHttpHandlePutReq(ctx, client, curHttpUrl, jsonData, nil); err != nil {
				tflog.Error(ctx, "PUT Addresses request failed, error: "+err.Error())
				return err
			}
		}
	}
	return nil
}

func (c *Client) DeleteDevOrgServiceObjAddresses(ctx context.Context,
	addrListData DevOjectsAddressListData) error {

	addrList := addrListData.AddrList

	if addrList.Count <= 0 || len(addrList.Addresses) <= 0 {
		tflog.Trace(ctx, "Device Orgs Address deletion failed as addresses count is 0")
		return errors.New("Device Orgs Address deletion failed as addresses count is 0")
	}

	tflog.Trace(ctx, "Device-Name "+addrList.DeviceName+" OrgName "+addrList.OrganizationName)
	for key, val := range addrList.Addresses {
		tflog.Trace(ctx, "Address["+strconv.Itoa(key)+"]: Name "+val.Name+" FQDN "+val.FQDN)
	}

	client, _, _ := vHttpClient(c.Config.ServerIP, c.Config.ServerPort, "")
	httpUrl := "https://" + c.Config.ServerIP + ":" + strconv.Itoa(c.Config.ServerPort) + "/" +
		vmsDirectorDevicesURL + "/" +
		addrList.DeviceName + "/" +
		vmsDirectorOrgServicesURL + "/" +
		addrList.OrganizationName + "/" +
		vmsDirectorObjectsAddressesURL + "/" +
		"address" + "/"

	// Delete expects individual objects, send one after another
	for _, val := range addrList.Addresses {
		curHttpUrl := httpUrl + val.Name
		var curAddrData DevObjectsAddressList
		curAddrData.Addresses = append(curAddrData.Addresses, val)
		if jsonData, err := json.Marshal(curAddrData); err != nil {
			tflog.Error(ctx, "Address Delete failed for "+val.Name+" Error: "+err.Error())
			return err
		} else {
			if _, err := c.vHttpHandleDeleteReq(ctx, client, curHttpUrl, jsonData, nil); err != nil {
				fmt.Printf("DELETE Addresses request failed, error %v\n", err)
				tflog.Error(ctx, "DELETE Addresses request failed, error: "+err.Error())
				return err
			}
		}
	}
	return nil
}
