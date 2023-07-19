package vclient

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

/* path needs to be appended along with server-ip and port to get token */
const (
	vOauthConfigFile      = "../../vOauth2Config.json"
	vOauthTokenFile       = "../../vOauth2Token.json"
	vOauthServerTokenPath = "auth/token"
)

/*
 * OUTH2 server response received for token request. This response includes
 * token used for subsequent http transactions.
 */
type vOauthServerToken struct {
	AccessToken  string `json:"access_token"`
	IssuedAt     string `json:"issued_at,omitempty"`
	ExpiresIn    string `json:"expires_in"`
	TokenType    string `json:"token_type,omitempty"`
	RefreshToken string `json:"refresh_token"`
	User         struct {
		Name            string   `json:"name"`
		IsExternalUser  bool     `json:"is_external_user,omitempty"`
		EnableTwoFactor bool     `json:"enable_two_factor,omitempty"`
		IdleTimeOut     int      `json:"idle_time_out,omitempty"`
		Roles           []string `json:"roles,omitempty"`
		PrimaryRole     string   `json:"primaryrole"`
	} `json:"user"`
}

/*
 * Configuration parameters needed to get auth-token from oauth server.
 * These parameters are read from either file or environment varibales.
 */
type vOauthConfig struct {
	ServerIP     string `json:"ServerIP,omitempty"`
	ServerPort   int    `json:"ServerPort,omitempty"`
	UserName     string `json:"UserName,omitempty"`
	Password     string `json:"Password,omitempty"`
	GrantType    string `json:"GrantType,omitempty"`
	ClientID     string `json:"ClientID,omitempty"`
	ClientSecret string `json:"ClientSecret,omitempty"`
	Scopes       []string
}

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Config     vOauthConfig
	Token      vOauthServerToken
}

/*
 * Utility function to display token information received from server.
 */
func vOauthTokenDisplay(tokenData vOauthServerToken) {
	log.Printf("----- OAUTH Token Data -----\n")
	log.Printf("Access Token         : %v\n", tokenData.AccessToken)
	log.Printf("Issued At            : %v\n", tokenData.IssuedAt)
	log.Printf("Expires In           : %v\n", tokenData.ExpiresIn)
	log.Printf("Token Type           : %v\n", tokenData.TokenType)
	log.Printf("Refresh Token        : %v\n", tokenData.RefreshToken)
	log.Printf("  User               : %v\n", tokenData.User.Name)
	log.Printf("  External User      : %v\n", tokenData.User.IsExternalUser)
	log.Printf("  Two Factor Enabled : %v\n", tokenData.User.EnableTwoFactor)
	log.Printf("  Idle Timeout       : %v\n", tokenData.User.IdleTimeOut)
	log.Printf("  Roles              : %v\n", tokenData.User.Roles)
	log.Printf("  Primary Role       : %v\n", tokenData.User.PrimaryRole)
	log.Printf("----------------------------\n")
}

/*
 * API to get OAUTH2 token from server. It sends POST request to server with
 * configuration data read from json file and/or from environment variables.
 * The token received from server is copied to golbal data structure to be
 * used for subsequent api calls with server.
 */
func vOauthGetToken(client *Client) error {

	config := &client.Config

	oauthParams := map[string]string{
		"client_id":     config.ClientID,
		"client_secret": config.ClientSecret,
		"grant_type":    config.GrantType,
		"username":      config.UserName,
		"password":      config.Password,
	}

	log.Printf("Get OAUTH tone for versadirector for %v\n", oauthParams)
	if requestBody, err := json.Marshal(oauthParams); err != nil {
		log.Printf("Unable to marshal oauth-parameters %v\n", err)
	} else {
		oauthServerUrl := "https://" +
			config.ServerIP + ":" +
			strconv.Itoa(config.ServerPort) + "/" +
			vOauthServerTokenPath

		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		tlsTransport := &http.Transport{TLSClientConfig: tlsConfig}
		httpTransport := &http.Client{Transport: tlsTransport}
		resp, err := httpTransport.Post(oauthServerUrl,
			"application/json", bytes.NewBuffer(requestBody))
		if err != nil {
			log.Printf("Unable to send POST request to get token %v\n", err)
			return err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Unable to read response from OAUTH server for token %v\n", err)
			return err
		}

		var tokenData vOauthServerToken
		if err := json.Unmarshal([]byte(body), &tokenData); err != nil {
			log.Printf("Unable to unmarshal token response %v\n", err)
			return err
		}
		client.Token = tokenData

		vOauthTokenDisplay(tokenData)

		if err := ioutil.WriteFile(vOauthTokenFile, body, 0777); err != nil {
			log.Printf("Failed to write token data to file\n")
		}
	}

	return nil
}

/*
 * Read token data from file. This file is created  after getting token from
 * director and same will be used until it is expired.
 */
func vOauthReadToken(fileName string, client *Client) error {

	var tokenData vOauthServerToken
	if len(fileName) > 0 {
		if fd, err := os.Open(fileName); err != nil {
			log.Printf("Unable to open config file %v, error %v\n", fileName, err)
			return err
		} else {
			defer fd.Close()
			if fileData, err := ioutil.ReadAll(fd); err != nil {
				log.Printf("Unable to read data from file %v fd %v error %v\n",
					fileName, fd, err)
				return err
			} else {
				if err := json.Unmarshal([]byte(fileData), &tokenData); err != nil {
					log.Printf("JSON unmarshal failed for file %v error %v\n",
						fileName, err)
					return err
				}
			}
		}
	}

	/* check validity of token */
	if tokenData.ExpiresIn != "-1" {
		log.Printf("Validate expiration date %v, get new token if expired\n",
			tokenData.ExpiresIn)
		return errors.New("Token expired, Get new one")
	}

	client.Token = tokenData

	log.Printf("Read AUTH2 token from file %v\n", fileName)
	vOauthTokenDisplay(tokenData)

	return nil
}

// NewClient -
func NewClient(host, username, password, port, oauthClientId,
	oauthClientSecret, oauthGrantType *string) (*Client, error) {

	c := Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	// If username or password not provided, return empty client
	if username == nil || password == nil {
		return &c, nil
	}

	log.Printf("Create new client for host %v:%v user %v client-id %v client-secret %v\n",
		*host, *port, *username, *oauthClientId, *oauthClientSecret)
	config := &c.Config
	config.ServerIP = *host
	config.ServerPort, _ = strconv.Atoi(*port)
	config.UserName = *username
	config.Password = *password
	config.ClientID = *oauthClientId
	config.ClientSecret = *oauthClientSecret
	config.GrantType = *oauthGrantType

	if err := vOauthReadToken(vOauthTokenFile, &c); err != nil {
		/* get auth token from server */
		if err := vOauthGetToken(&c); err != nil {
			log.Fatal("Authentication token isn't available, aborting")
		}
	}

	return &c, nil
}

/*
 * Utility function to convert response received from versadirector to
 * printable pretty string.
 */
func vDisplayResponse(body []byte) {
	var jsonData bytes.Buffer
	error := json.Indent(&jsonData, body, "", "\t")
	if error != nil {
		return
	}
	fmt.Printf("%v\n", jsonData.String())
}

/*
 * Common utility function to create http client and url string
 * needed form http requests.
 */
func vHttpClient(host string, port int, urlPath string) (*http.Client,
	string, error) {

	/* disable certificate validation */
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	var httpUrl string
	if len(urlPath) > 0 {
		httpUrl = "https://" + host + ":" + strconv.Itoa(port) + "/" + urlPath
	}

	return client, httpUrl, nil
}

func (c *Client) vHttpHandleGetReq(ctx context.Context,
	client *http.Client,
	apiUrl string,
	urlData url.Values) ([]byte, error) {

	/* form http GET request */
	httpReq, _ := url.ParseRequestURI(apiUrl)
	if len(urlData) > 0 {
		httpReq.RawQuery = urlData.Encode()
	}
	urlStr := fmt.Sprintf("%v", httpReq)

	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		fmt.Printf("Error creating http request %v\n", err)
		return nil, err
	}
	req.Header.Add("Accept", `application/json`)
	req.Header.Set("Content-Type", "application/json")

	bearer := "Bearer " + c.Token.AccessToken
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Unable to send http GET request %v\n", err)
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		tflog.Debug(ctx, "http response returned error "+resp.Status)
		return nil, errors.New("http response error")
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 0x10000))
	if err != nil {
		fmt.Printf("Failed to read http response %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	//fmt.Printf("GET REQ DATA %v\n", string(body))

	return body, nil
}

func (c *Client) vHttpHandlePostReq(ctx context.Context,
	client *http.Client,
	apiUrl string,
	request []byte,
	urlData url.Values) ([]byte, error) {

	tflog.Debug(ctx, "POST Request: "+apiUrl)
	tflog.Debug(ctx, "POST Request body: "+string(request))
	/* form http POST request */
	httpReq, _ := url.ParseRequestURI(apiUrl)
	if len(urlData) > 0 {
		httpReq.RawQuery = urlData.Encode()
	}
	urlStr := fmt.Sprintf("%v", httpReq)

	req, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewBuffer(request))
	if err != nil {
		tflog.Error(ctx, "Error in creating http request for POST")
		return nil, err
	}
	req.Header.Add("Accept", `application/json`)
	req.Header.Set("Content-Type", "application/json")

	bearer := "Bearer " + c.Token.AccessToken
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		tflog.Debug(ctx, "Error in sending POST request: "+err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		tflog.Debug(ctx, "Error response for POST request: "+resp.Status)
		return nil, errors.New("http response error")
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 0x10000))
	if err != nil {
		tflog.Debug(ctx, "Error reading POST response: "+err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	tflog.Debug(ctx, "POST request handled successfully")
	return body, nil
}

func (c *Client) vHttpHandlePutReq(ctx context.Context,
	client *http.Client,
	apiUrl string,
	request []byte,
	urlData url.Values) ([]byte, error) {

	/* form http PUT request */
	httpReq, _ := url.ParseRequestURI(apiUrl)
	if len(urlData) > 0 {
		httpReq.RawQuery = urlData.Encode()
	}
	urlStr := fmt.Sprintf("%v", httpReq)

	req, err := http.NewRequest(http.MethodPut, urlStr, bytes.NewBuffer(request))
	if err != nil {
		tflog.Error(ctx, "Error in creating http request for PUT")
		return nil, err
	}
	req.Header.Add("Accept", `application/json`)
	req.Header.Set("Content-Type", "application/json")

	bearer := "Bearer " + c.Token.AccessToken
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		tflog.Debug(ctx, "Error in sending PUT request: "+err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {
		tflog.Debug(ctx, "Error response for PUT request: "+resp.Status)
		return nil, errors.New("http response error")
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 0x10000))
	if err != nil {
		tflog.Debug(ctx, "Error reading PUT response: "+err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	tflog.Debug(ctx, "PUT request handled successfully")
	return body, nil
}

func (c *Client) vHttpHandleDeleteReq(ctx context.Context,
	client *http.Client,
	apiUrl string,
	request []byte,
	urlData url.Values) ([]byte, error) {

	/* form http POST request */
	httpReq, _ := url.ParseRequestURI(apiUrl)
	if len(urlData) > 0 {
		httpReq.RawQuery = urlData.Encode()
	}
	urlStr := fmt.Sprintf("%v", httpReq)

	req, err := http.NewRequest(http.MethodDelete, urlStr, bytes.NewBuffer(request))
	if err != nil {
		tflog.Error(ctx, "Error in creating http request for DELETE")
		return nil, err
	}
	req.Header.Add("Accept", `application/json`)
	req.Header.Set("Content-Type", "application/json")

	bearer := "Bearer " + c.Token.AccessToken
	req.Header.Add("Authorization", bearer)

	resp, err := client.Do(req)
	if err != nil {
		tflog.Debug(ctx, "Error in sending DELETE request: "+err.Error())
		return nil, err
	}
	if resp.StatusCode != http.StatusOK &&
		resp.StatusCode != http.StatusNoContent {
		tflog.Debug(ctx, "Error response for DELETE request: "+resp.Status)
		return nil, errors.New("http response error")
	}

	body, err := ioutil.ReadAll(io.LimitReader(resp.Body, 0x10000))
	if err != nil {
		tflog.Debug(ctx, "Error reading DELETE response: "+err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	tflog.Debug(ctx, "DELETE request handled successfully")
	return body, nil
}
