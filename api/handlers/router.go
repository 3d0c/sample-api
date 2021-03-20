package handlers

import (
	m "github.com/3d0c/sample-api/api/middleware"

	"github.com/julienschmidt/httprouter"
)

// SetupRouter sets up endpoints
func SetupRouter() *httprouter.Router {
	r := httprouter.New()

	// Create new uesr
	// {'name': 'example', 'password': 'password'}
	r.POST("/users", m.Chain(users().create))

	// Login user
	// {'name': 'example', 'password': 'password'}
	r.POST("/users/login", m.Chain(users().login))

	// Add flight (Protected method)
	r.POST("/flights", m.Chain(m.Auth, flights().create))

	// Delete flight (Protected method)
	r.DELETE("/flights/:id", m.Chain(m.Auth, flights().remove))

	// Update flight (Protected method)
	r.PUT("/flights/:id", m.Chain(m.Auth, flights().update))

	// Search for flights
	r.GET("/flights", m.Chain(m.Auth, flights().search))

	return r
}
