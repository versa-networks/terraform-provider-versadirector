package vclient

import (
	"context"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
)

const vmsDirectorAppliancesURL = "vnms/appliance/appliance"

/*
 * Appliances data received from director.
 */
type VmsDirectorAppliances struct {
	TotalCount int `json:"totalCount"`
	Appliances []struct {
		Name              string `json:"name"`
		UUID              string `json:"uuid"`
		ApplianceLocation struct {
			ApplianceName string `json:"applianceName"`
			ApplianceUUID string `json:"applianceUuid"`
			LocationId    string `json:"locationId"`
			Latitude      string `json:"latitude"`
			Longitude     string `string:"longitude"`
			Type          string `json:"type"`
		} `json:"applianceLocation,omitempty"`
		LastUpdatedTime         string `json:"last-updated-time"`
		PingStatus              string `json:"ping-status"`
		SyncStatus              string `json:"sync-status"`
		CreatedAt               string `json:"createdAt"`
		YangCompatibilityStatus string `json:"yang-compatibility-status"`
		ServicesStatus          string `json:"services-status"`
		OverallStatus           string `json:"overall-status"`
		ControllStatus          string `json:"controll-status"`
		PathStatus              string `json:"path-status"`
		InterChassisHaStatus    struct {
			HaConfigured bool `json:"ha-configured"`
		} `json:"inter-chassis-ha-status"`
		TemplateStatus  string   `json:"templateStatus"`
		OwnerOrgUuid    string   `json:"ownerOrgUuid"`
		OwnerOrg        string   `json:"ownerOrg"`
		Type            string   `json:"type"`
		Deployment      string   `json:"deployment"`
		CmsOrg          string   `json:"cmsOrg"`
		Orgs            []string `json:"orgs"`
		SngCount        int      `json:"sngCount"`
		SoftwareVersion string   `json:"softwareVersion"`
		Connector       string   `json:"connector"`
		ConnectorType   string   `json:"connectorType"`
		BranchId        string   `json:"branchId"`
		Services        []string `json:"services"`
		IpAddress       string   `json:"ipAddress"`
		Location        string   `json:"location"`
		StartTime       string   `json:"startTime"`
		Hardware        struct {
			Name                         string `json:"name"`
			Model                        string `json:"model"`
			CpuCores                     int    `json:"cpuCores"`
			Memory                       string `json:"memory"`
			FreeMemory                   string `json:"freeMemory"`
			DiskSize                     string `json:"diskSize"`
			FreeDisk                     string `json:"freeDisk"`
			Lpm                          bool   `json:"lpm"`
			Fanless                      bool   `json:"fanless"`
			IntelQuickAssistAcceleration bool   `json:"intelQuickAssistAcceleration"`
			FirmwareVersion              string `json:"firmwareVersion"`
			Manufacturer                 string `json:"manufacturer"`
			SerialNo                     string `json:"serialNo"`
			HardWareSerialNo             string `json:"hardWareSerialNo"`
			CpuModel                     string `json:"cpuModel"`
			CpuCount                     int    `json:"cpuCount"`
			CpuLoad                      int    `json:"cpuLoad"`
			InterfaceCount               int    `json:"interfaceCount"`
			PackageName                  string `json:"packageName"`
			Sku                          string `json:"sku"`
			Ssd                          bool   `json:"ssd"`
		}
		SPack struct {
			Name         string `json:"name"`
			SpackVersion string `json:"spackVersion"`
			ApiVersion   string `json:"apiVersion"`
			Flavor       string `json:"flavor"`
			ReleaseDate  string `json:"releaseDate"`
			UpdateType   string `json:"updateType"`
		} `json:"SPack"`
		OssPack struct {
			Name           string `json:"name"`
			OsspackVersion string `json:"osspackVersion"`
			UpdateType     string `jsoon:"updateType"`
		}
		AppIdDetails struct {
			AppIdInstalledEngineVersion string `json:"appIdInstalledEngineVersion"`
			AppIdInstalledBundleVersion string `json:"appIdInstalledBundleVersion"`
			AppIdAvailableBundleVersion string `json:"appIdAvailableBundleVersion"`
		} `json:"appIdDetails"`
		AlarmSummary struct {
			TableId     string   `json:"tableId"`
			TableName   string   `json:"tableName"`
			MonitorType string   `json:"monitorType"`
			ColumnNames []string `json:"columnNames"`
			Rows        []struct {
				FirstColumnValue string `json:"firstColumnValue"`
				columnValues     []int
			}
		} `json:"alarmSummary"`
		CpeHealth struct {
			ColumnNames []string `json:"columnNames"`
			Rows        []struct {
				FirstColumnValue string `json:"firstColumnValue"`
				ColumnValues     []int  `json:"columnValues"`
			} `json:"rows"`
		} `json:"cpeHealth"`
		Controllers           []string `json:"controllers"`
		RefreshCycleCount     int      `json:"refreshCycleCount"`
		SubType               string   `json:"subType"`
		BranchMaintenanceMode bool     `json:"branch-maintenance-mode"`
		ApplianceCapabilities struct {
			Capabilities []string `json:"capabilities"`
		} `json:"applianceCapabilities"`
		LockDetails struct {
			User     string `json:"user"`
			LockType string `json:"lockType"`
		} `json:"lockDetails"`
		BranchInMaintenanceMode bool `json:"branchInMaintenanceMode"`
		Unreachable             bool `json:"unreachable"`
	} `json:"appliances"`
}

func (c *Client) GetAllAppliances(ctx context.Context) (*VmsDirectorAppliances, error) {

	client, apiUrl, err := vHttpClient(c.Config.ServerIP, c.Config.ServerPort, vmsDirectorAppliancesURL)
	if err != nil {
		log.Printf("Unable to create http-client %v\n", err)
		return nil, err
	}

	urlData := url.Values{}
	urlData.Set("limit", strconv.Itoa(10))
	urlData.Add("offset", strconv.Itoa(0))

	if data, err := c.vHttpHandleGetReq(ctx, client, apiUrl, urlData); err != nil {
		return nil, err
	} else {
		vDisplayResponse(data)

		applianceData := VmsDirectorAppliances{}
		if err := json.Unmarshal([]byte(data), &applianceData); err != nil {
			return nil, err
		}
		return &applianceData, nil
	}
}
