package vclient

import (
	"context"
	"encoding/json"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (c *Client) GetDeviceOrganizationAddresses(ctx context.Context,
	deviceName string, organizationName string) (*DevObjectsAddressList, error) {

	client, _, err := vHttpClient(c.Config.ServerIP, c.Config.ServerPort, "")
	if err != nil {
		log.Printf("Unable to create http-client %v\n", err)
		return nil, err
	}

	httpUrl := "https://" + c.Config.ServerIP + ":" + strconv.Itoa(c.Config.ServerPort) + "/" +
		vmsDirectorDevicesURL + "/" +
		deviceName + "/" +
		vmsDirectorOrgServicesURL + "/" +
		organizationName + "/" +
		vmsDirectorObjectsAddressesURL + "/" + "address"
	tflog.Debug(ctx, "CLIENT-DATA GET URL: "+httpUrl)

	if data, err := c.vHttpHandleGetReq(ctx, client, httpUrl, nil); err != nil {
		log.Printf("HTTP GET failed for URL: %v, error: %v", httpUrl, err)
		return nil, err
	} else {
		tflog.Debug(ctx, "CLIENT-DATA GET SUCECSSFUL fot=r URL: "+httpUrl)
		addrListData := DevOjectsAddressListData{}
		if err := json.Unmarshal([]byte(data), &addrListData.AddrList); err != nil {
			log.Printf("Unmarshal failed for address data, error %v\n", err)
			return nil, err
		}
		//vDisplayResponse(data)
		return &addrListData.AddrList, nil
	}
}
