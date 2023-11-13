package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// SecretItem represents a configuration key and its value.
type SecretItem struct {
	Name  string `json:"Name,omitempty"`
	Value string `json:"Value,omitempty"`
}

// SecretGroup represents a configuration group with cluster ID, system secret indicator, and secrets.
type SecretGroup struct {
	ClusterID    string       `json:"ClusterId"`
	SystemSecret bool         `json:"SystemSecret"`
	Secrets      []SecretItem `json:"Secrets"`
}

// SecretReference represents a configuration reference with name, value, version, and ID.
type SecretReference struct {
	Name    string `json:"Name,omitempty"`
	Value   string `json:"Value,omitempty"`
	Version string `json:"Version,omitempty"`
	Id      string `json:"Id,omitempty"`
}

// SecretManagerResponseData represents the response data from the secret manager.
type SecretManagerResponseData struct {
	Data                []SecretReference `json:"Data"`
	SecretStoreProvider string            `json:"SecretStoreProvider"`
	SecretStoreName     string            `json:"SecretStoreName"`
}

// ClusterData represents the data structure for cluster information.
type ClusterData struct {
	ID                   string `json:"id"`
	CreatedAt            string `json:"created_at"`
	OrganizationID       int    `json:"organization_id"`
	OrganizationUUID     string `json:"organization_uuid"`
	EnvName              string `json:"env_name"`
	Region               string `json:"region"`
	ChoreoEnv            string `json:"choreo_env"`
	ClusterID            string `json:"cluster_id"`
	DockerCredentialUUID string `json:"docker_credential_uuid"`
	ExternalApimEnvName  string `json:"external_apim_env_name"`
	InternalApimEnvName  string `json:"internal_apim_env_name"`
	SandboxApimEnvName   string `json:"sandbox_apim_env_name"`
	Critical             bool   `json:"critical"`
	PdpWebAppDnsPrefix   string `json:"pdp_web_app_dns_prefix"`
}

// CloudManagerResponseData represents the response data from the cloud manager.
type CloudManagerResponseData struct {
	Data ClusterData `json:"data"`
}

func main() {
	envList := []string{
		"49eae34b-fb17-479f-a657-0bba998a4e79",
		"433a2e50-5d4b-4efb-a994-680173bc1079",
		"0af31d8b-9f84-4e20-b42d-fb4a5927940d",
		"ad5a37e4-9c92-4952-9694-c3269fdd8cdc",
		"842e5f7b-eed1-423b-98c2-e074c008c949",
	}

	fmt.Println("Before calling cloud manager client")
	_, err := getClusterID(envList[1])
	if err != nil {
		fmt.Println("Error calling cloud manager client:", err)
	}
	fmt.Println("Before calling secret manager client")
	_, errSecret := retrieveConfigValuesFromKV()
	if errSecret != nil {
		fmt.Println("Error calling secret manager client:", err)
	}
	_, errSecret2 := retrieveConfigValuesFromKV()
	if errSecret2 != nil {
		fmt.Println("Error calling secret manager client:", err)
	}
	fmt.Println("Before calling cloud manager client 2")
	_, err2 := getClusterID(envList[1])
	if err2 != nil {
		fmt.Println("Error calling cloud manager client:", err)
	}
	fmt.Println("###### After calling cloud manager client ######")
}

func retrieveConfigValuesFromKV() (string, error) {
	secretManagerEndpoint := fmt.Sprintf("http://%s:5006", "localhost")
	secrets := []SecretItem{
		{Name: "01ee7ebe-c770-1ab8-ade5-584ac5adffb9"},
		{Name: "01ee7ebe-c770-1ab8-b16a-4188708e8251"},
	}
	secretGroup := SecretGroup{
		ClusterID:    "7eca5163-6a37-ee11-b8f0-000d3adac5f0",
		SystemSecret: false,
		Secrets:      secrets,
	}
	secretGroupJSON, err := json.Marshal(secretGroup)
	if err != nil {
		return "", err
	}

	resp, err := http.Post(secretManagerEndpoint+"/api/v1/secrets/get", "application/json", bytes.NewBuffer(secretGroupJSON))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Secret manager request failed with status code: %d", resp.StatusCode)
	}

	return string(body), nil
}

func getClusterID(envTemplateUUID string) (string, error) {
	cloudManagerEndpoint := fmt.Sprintf("http://%s:5009", "localhost")
	url := fmt.Sprintf("%s/api/v1/env-templates/%s", cloudManagerEndpoint, envTemplateUUID)
	// Define headers
	headers := map[string]string{
		"x-organization-id": "0",
		"x-project-id":      "global",
	}
	// Create a new HTTP request with headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// Set headers on the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Cloud manager request failed with status code: %d", resp.StatusCode)
	}

	return string(body), nil
}
