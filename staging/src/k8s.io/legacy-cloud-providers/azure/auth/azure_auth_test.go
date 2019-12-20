/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package auth

import (
	"testing"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/assert"
)

func TestGetServicePrincipalTokenFromMSIWithUserAssignedID(t *testing.T) {
	config := &AzureAuthConfig{
		UseManagedIdentityExtension: true,
		UserAssignedIdentityID:      "UserAssignedIdentityID",
	}
	env := &azure.PublicCloud

	token, err := GetServicePrincipalToken(config, env)
	assert.NoError(t, err)

	msiEndpoint, err := adal.GetMSIVMEndpoint()
	assert.NoError(t, err)

	spt, err := adal.NewServicePrincipalTokenFromMSIWithUserAssignedID(msiEndpoint,
		env.ServiceManagementEndpoint, config.UserAssignedIdentityID)
	assert.NoError(t, err)

	assert.Equal(t, token, spt)
}

func TestGetServicePrincipalTokenFromMSI(t *testing.T) {
	config := &AzureAuthConfig{
		UseManagedIdentityExtension: true,
	}
	env := &azure.PublicCloud

	token, err := GetServicePrincipalToken(config, env)
	assert.NoError(t, err)

	msiEndpoint, err := adal.GetMSIVMEndpoint()
	assert.NoError(t, err)

	spt, err := adal.NewServicePrincipalTokenFromMSI(msiEndpoint, env.ServiceManagementEndpoint)
	assert.NoError(t, err)

	assert.Equal(t, token, spt)
}

func TestGetServicePrincipalToken(t *testing.T) {
	config := &AzureAuthConfig{
		TenantID:        "TenantID",
		AADClientID:     "AADClientID",
		AADClientSecret: "AADClientSecret",
	}
	env := &azure.PublicCloud

	token, err := GetServicePrincipalToken(config, env)
	assert.NoError(t, err)

	oauthConfig, err := adal.NewOAuthConfig(env.ActiveDirectoryEndpoint, config.TenantID)
	assert.NoError(t, err)

	spt, err := adal.NewServicePrincipalToken(*oauthConfig, config.AADClientID, config.AADClientSecret, env.ServiceManagementEndpoint)

	assert.Equal(t, token, spt)
}

func TestParseAzureEngironment(t *testing.T) {
	cases := []struct {
		cloudName               string
		resourceManagerEndpoint string
		identitySystem          string
		expected                *azure.Environment
	}{
		{
			cloudName:               "",
			resourceManagerEndpoint: "",
			identitySystem:          "",
			expected:                &azure.PublicCloud,
		},
		{
			cloudName:               "AZURECHINACLOUD",
			resourceManagerEndpoint: "",
			identitySystem:          "",
			expected:                &azure.ChinaCloud,
		},
	}

	for _, c := range cases {
		env, err := ParseAzureEnvironment(c.cloudName, c.resourceManagerEndpoint, c.identitySystem)
		assert.NoError(t, err)
		assert.Equal(t, env, c.expected)
	}
}
