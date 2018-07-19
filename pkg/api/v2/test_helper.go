package v2

import (
	"net/http"
	"testing"

	"github.com/alejandroEsc/golang-maas-client/pkg/api/client"
	"github.com/stretchr/testify/assert"
)

// createTestServerController creates a ControllerInterface backed on to a test server
// that has sufficient knowledge of versions and users to be able to create a
// valid ControllerInterface.
func createTestServerController(t *testing.T) (*client.SimpleTestServer, *Controller) {
	server := client.NewSimpleServer()
	server.AddGetResponse("/api/2.0/users/?op=whoami", http.StatusOK, `"captain awesome"`)
	server.AddGetResponse("/api/2.0/version/", http.StatusOK, versionResponse)
	server.Start()

	controller, err := NewController(ControllerArgs{
		BaseURL: server.URL,
		APIKey:  "fake:as:key",
	})
	assert.Nil(t, err)
	return server, controller
}
