package srvhttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestService_StatusCheckReadiness(t *testing.T) {
	var srv Service

	srv.readiness = false
	assert.Equal(t, "http srv is't ready yet", srv.Check().Error())
}
