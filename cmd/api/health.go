package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ENV int

const (
	Dev ENV = iota
	Stage
	Prod
)

func (e ENV) String() string {
	return [...]string{
		"development", "staging", "production",
	}[e]
}

type Health struct {
	Status      int
	Environment string
	Version     string
}

func (app *application) Health(c *gin.Context) {

	health := Health{
		Status:      http.StatusOK,
		Environment: Dev.String(),
		Version:     version,
	}

	time.Sleep(4 * time.Second)

	app.writeJSON(c, http.StatusOK, envelope{"system_info": health}, nil)
}
