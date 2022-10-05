package testutils

import (
	"path/filepath"
	"testing"

	"github.com/BLasan/APKCTL-Demo/integration/base"
	"github.com/stretchr/testify/assert"
)

func AddNewAPIWithSwagger(t *testing.T, swagerPath string) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	out, err := deployAPIWithSwagger(t, apiName, apiVersion, swagerPath)

	assert.Nil(t, err, "Error while deploying API")
	assert.Contains(t, out, "Successfully deployed")
}

func CreateNewAPIFromSwaggerWithDryRun(t *testing.T, swagerPath string) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	out, err := createAPIWithSwagger(t, apiName, apiVersion, swagerPath)

	assert.Nil(t, err, "Error while creating API from Swagger File")
	assert.Contains(t, out, "Successfully created")

	apiProjectDir := base.GetExportedPathFromOutput(out)

	httprouteconfig := filepath.Join(apiProjectDir, HttpRouteConfigFile)
	configmap := filepath.Join(apiProjectDir, ConfigMapFile)

	assert.True(t, base.IsFileAvailable(t, httprouteconfig), "HttpRouteConfig is not available")
	assert.True(t, base.IsFileAvailable(t, configmap), "ConfigMap is not available")
}

func AddNewAPIWithBackendServiceURL(t *testing.T) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	out, err := deployAPIWithBackendServiceURL(t, apiName, apiVersion, BackendServiceURL)

	assert.Nil(t, err, "Error while deploying API")
	assert.Contains(t, out, "Successfully deployed")
}

func CreateNewAPIFromBackendServiceURLWithDryRun(t *testing.T) {
	t.Helper()
	apiName := base.GenerateRandomName(15) + "API"
	apiVersion := APIVersion
	out, err := createAPIWithBackendServiceURL(t, apiName, apiVersion, BackendServiceURL)

	assert.Nil(t, err, "Error while creating AP from Backend Service URL")
	assert.Contains(t, out, "Successfully created")

	apiProjectDir := base.GetExportedPathFromOutput(out)

	httprouteconfig := filepath.Join(apiProjectDir, HttpRouteConfigFile)
	configmap := filepath.Join(apiProjectDir, ConfigMapFile)

	assert.True(t, base.IsFileAvailable(t, httprouteconfig), "HttpRouteConfig is not available")
	assert.True(t, base.IsFileAvailable(t, configmap), "ConfigMap is not available")
}

func ValidateInstallAPKComponents(t *testing.T) {
	t.Helper()

	out, err := installAPK(t)

	assert.Nil(t, err, "Error while installing APK components")
	assert.Contains(t, out, "Successfully installed!")
}

func installAPK(t *testing.T) (string, error) {
	output, err := base.Execute(t, "install", "platform", "-k", "--verbose")
	return output, err
}

// Creates API from swagger file
func deployAPIWithSwagger(t *testing.T, apiName, apiversion, swagger string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "-f", swagger, "-k", "--verbose")
	t.Cleanup(func() {
		removeAPI(t, apiName, apiversion)
	})
	return output, err
}

// Creates API from the Backend Service URL
func deployAPIWithBackendServiceURL(t *testing.T, apiName, apiversion, backendURL string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--service-url", backendURL, "-k", "--verbose")
	t.Cleanup(func() {
		removeAPI(t, apiName, apiversion)
	})
	return output, err
}

// Creates API from swagger file
func createAPIWithSwagger(t *testing.T, apiName, apiversion, swagger string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "-f", swagger, "--dry-run", "-k", "--verbose")
	t.Cleanup(func() {
		removeAPI(t, apiName, apiversion)
	})
	return output, err
}

// Creates API from the Backend Service URL
func createAPIWithBackendServiceURL(t *testing.T, apiName, apiversion, backendURL string) (string, error) {
	output, err := base.Execute(t, "create", "api", apiName, "--service-url", backendURL, "--dry-run", "-k", "--verbose")
	t.Cleanup(func() {
		removeAPI(t, apiName, apiversion)
	})
	return output, err
}

func removeAPI(t *testing.T, apiname, version string) {
	base.Execute(t, "delete", "api", apiname, "--version", version)
}
