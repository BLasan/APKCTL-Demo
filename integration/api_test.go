package integration

import (
	"path/filepath"
	"testing"

	testutils "github.com/BLasan/APKCTL-Demo/integration/testUtils"
)

func TestInstallAPKComponents(t *testing.T) {
	t.Run("APKCTL Installation", func(t *testing.T) {
		testutils.ValidateInstallAPKComponents(t)
	})
}

func TestAPIDeploymentFromSwagger(t *testing.T) {
	t.Run("Deploy API with a Swagger File", func(t *testing.T) {
		swaggerPath := filepath.Join(testutils.SampleTestData, testutils.SampleTestSwaggerFile)
		testutils.AddNewAPIWithSwagger(t, swaggerPath)
	})
}

func TestAPIDeploymentFromServiceURL(t *testing.T) {
	t.Run("Deploy API with the backend service URL", func(t *testing.T) {
		testutils.AddNewAPIWithBackendServiceURL(t)
	})
}

func TestAPICreationtFromSwaggerWithDryRun(t *testing.T) {
	t.Run("Create API without deploying with a Swagger File", func(t *testing.T) {
		swaggerPath := filepath.Join(testutils.SampleTestData, testutils.SampleTestSwaggerFile)
		testutils.CreateNewAPIFromSwaggerWithDryRun(t, swaggerPath)
	})
}

func TestAPICreationtFromServiceURLWithDryRun(t *testing.T) {
	t.Run("Create API without deploying with the backend service URL", func(t *testing.T) {
		testutils.CreateNewAPIFromBackendServiceURLWithDryRun(t)
	})
}

func TestAPIConfigFiles(t *testing.T) {
	t.Run("Validate values getting overriden in Swagger definition", func(t *testing.T) {
		testutils.ValidateAPIConfigFiles(t)
	})
}

func TestAPKComponentsUninstallation(t *testing.T) {
	t.Run("Validate APK components are getting removed properly", func(t *testing.T) {
		testutils.ValidateUninstallAPKComponents(t)
	})
}
