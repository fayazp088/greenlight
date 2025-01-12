package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")
				app.serverErrorResponse(c, fmt.Errorf("%v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}
