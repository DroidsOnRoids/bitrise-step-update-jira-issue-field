package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"

	"bytes"
	"encoding/json"
	"github.com/bitrise-io/go-utils/log"
	"io/ioutil"
	"net/http"
	"strings"
)

// -----------------------
// --- Models
// -----------------------

// ConfigsModel ...
type ConfigsModel struct {
	JiraUsername    string
	JiraPassword    string
	JiraInstanceURL string
	IssueIDOrKey    string
	FieldKey        string
	FieldValue      string
}

func main() {
	configs := createConfigsModelFromEnvs()
	configs.dump()
	if err := configs.validate(); err != nil {
		log.Errorf("Issue with input: %s", err)
		os.Exit(1)
	}

	payload := map[string]interface{}{
		"fields": map[string]interface{}{
			configs.FieldKey: configs.FieldValue,
		},
	}

	if err := sendRequest(configs, payload); err != nil {
		log.Errorf("Could not update issue, error: %s", err)
		os.Exit(2)
	}
}

func sendRequest(configs ConfigsModel, payload map[string]interface{}) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	requestURL := fmt.Sprintf("%s/rest/api/2/issue/%s", configs.JiraInstanceURL, configs.IssueIDOrKey)
	request, err := http.NewRequest("PUT", requestURL, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.SetBasicAuth(configs.JiraUsername, configs.JiraPassword)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	fmt.Println()
	log.Infof("Performing request")

	client := http.Client{}
	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			log.Warnf("Failed to close response body, error: %s", err)
		}
	}()

	if response.StatusCode != http.StatusNoContent {
		log.Warnf("JIRA API response status: %s", response.Status)
		contents, readErr := ioutil.ReadAll(response.Body)
		if readErr != nil {
			return errors.New("could not read JIRA API response")
		}
		log.Warnf("JIRA API response: %s", contents)

		if response.Header.Get("X-Seraph-LoginReason") == "AUTHENTICATION_DENIED" {
			log.Warnf("CAPTCHA triggered")
		}
		return errors.New("JIRA API request failed")
	}

	log.Infof("Issue %s updated successfully", configs.IssueIDOrKey)
	return nil
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		JiraUsername:    os.Getenv("jira_username"),
		JiraPassword:    os.Getenv("jira_password"),
		JiraInstanceURL: os.Getenv("jira_instance_url"),
		IssueIDOrKey:    os.Getenv("issue_id_or_key"),
		FieldKey:        os.Getenv("field_key"),
		FieldValue:      os.Getenv("field_value"),
	}
}

func (configs ConfigsModel) dump() {
	fmt.Println()
	log.Infof("Configs:")
	log.Printf(" - JiraUsername: %s", configs.JiraUsername)
	log.Printf(" - JiraPassword (hidden): %s", strings.Repeat("*", len(configs.JiraPassword)))
	log.Printf(" - JiraInstanceURL: %s", configs.JiraInstanceURL)
	log.Printf(" - IssueIdOrKey: %s", configs.IssueIDOrKey)
	log.Printf(" - FieldKey: %s", configs.FieldKey)
	log.Printf(" - FieldValue: %s", configs.FieldValue)
}

func (configs ConfigsModel) validate() error {
	if configs.JiraUsername == "" {
		return errors.New("no Jira Username specified")
	}
	if configs.JiraPassword == "" {
		return errors.New("no Jira Password specified")
	}
	_, err := url.ParseRequestURI(configs.JiraInstanceURL)
	if err != nil {
		return fmt.Errorf("invalid Jira instance URL, error %s", err)
	}
	if configs.IssueIDOrKey == "" {
		return errors.New("no Jira issue ID nor key specified")
	}
	if configs.FieldKey == "" {
		return errors.New("no field key specified")
	}
	return nil
}
