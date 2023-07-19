package vclient

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

const vmsDirectorOrganizationsURL = "nextgen/organization"

/*
 * Organizations data received from director in json format.
 */
type VmsDirectorOrganization struct {
	Name              string `json:"name"`
	UUID              string `json:"uuid"`
	Parent            string `json:"parent"`
	SubscriptionPlan  string `json:"subscriptionPlan"`
	Id                int    `json:"id"`
	CpeDeploymentType string `json:"cpeDeploymentType"`
	Appliances        []struct {
		ApplianceUUID string   `json:"applianceuuid"`
		CustomParams  []string `json:"customParams"`
	} `json:"appliances"`
	VrfsGroups []struct {
		Id          int    `json:"id"`
		VrfId       int    `json:"vrfId"`
		Name        string `json:"name"`
		Description string `json:"description"`
		EnableVpn   bool   `json:"enable_vpn"`
	} `json:"vrfsGroups,omitempty"`
	WanNetworkGroups []struct {
		Id               int      `json:"id"`
		Name             string   `json:"name"`
		Description      string   `json:"description"`
		TransportDomains []string `json:"transport-domains"`
	}
	AnalyticsClusters       []string `json:"analyticsClusters,omitempty"`
	SharedControlPlane      bool     `json:"sharedControlPlane"`
	BlockInterRegionRouting bool     `json:"blockInterRegionRouting"`
}

func (c *Client) GetAllOrganizations(ctx context.Context) ([]VmsDirectorOrganization, error) {

	client, apiUrl, err := vHttpClient(c.Config.ServerIP, c.Config.ServerPort, vmsDirectorOrganizationsURL)
	if err != nil {
		log.Printf("Unable to create http-client %v\n", err)
		return nil, err
	}

	urlData := url.Values{}
	urlData.Set("limit", strconv.Itoa(10))
	urlData.Add("offset", strconv.Itoa(0))
	urlData.Add("uuidOnly", "false")

	if data, err := c.vHttpHandleGetReq(ctx, client, apiUrl, urlData); err != nil {
		return nil, err
	} else {
		vDisplayResponse(data)

		organizationsData := []VmsDirectorOrganization{}
		if err := json.Unmarshal([]byte(data), &organizationsData); err != nil {
			return nil, err
		}
		return organizationsData, nil
	}
}
