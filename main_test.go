package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateConfigMissingJiraUsername(t *testing.T) {
	configs := createDummyConfigs()
	configs.JiraUsername = ""
	require.Error(t, configs.validate())
}

func TestValidateConfigMissingJiraPassword(t *testing.T) {
	configs := createDummyConfigs()
	configs.JiraPassword = ""
	require.Error(t, configs.validate())
}

func TestValidateConfigMissingJiraInstanceURL(t *testing.T) {
	configs := createDummyConfigs()
	configs.JiraInstanceURL = ""
	require.Error(t, configs.validate())
}

func TestValidateConfigEmptyIssueIDOrKeyList(t *testing.T) {
	configs := createDummyConfigs()
	configs.IssueIDOrKeyList = []string{""}
	require.Error(t, configs.validate())
}

func TestValidateConfigInvalidIssueIDOrKeyList(t *testing.T) {
	configs := createDummyConfigs()
	configs.IssueIDOrKeyList = []string{"TEST-1", ""}
	require.Error(t, configs.validate())
}

func TestValidateConfigMissingFieldKey(t *testing.T) {
	configs := createDummyConfigs()
	configs.FieldKey = ""
	require.Error(t, configs.validate())
}

func createDummyConfigs() ConfigsModel {
	return ConfigsModel{
		JiraUsername:     "login",
		JiraPassword:     "password",
		JiraInstanceURL:  "http://jira.invalid",
		IssueIDOrKeyList: []string{"TEST-1"},
		FieldKey:         "test",
		FieldValue:       "this is a test",
	}
}
