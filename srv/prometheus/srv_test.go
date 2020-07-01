package prometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Service

	srv.readiness = false
	assert.Equal(t, "prometheus srv is't ready yet", srv.Check().Error())
}
