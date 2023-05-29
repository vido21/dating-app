package swipes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoutes(t *testing.T) {
	println("There should be 1 route defined")
	authController := SwipesController{}
	routes := authController.Routes()
	assert.Equal(t, len(routes), 1)
}
