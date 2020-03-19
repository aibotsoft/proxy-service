package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// API constructs an http.Handler with all application routes defined.
func API() http.Handler {
	r := gin.Default()
	r.GET("/ping", pingHandler)

	r.GET("/countries", getAllCountries)
	r.POST("/countries", postCountry)
	return r
}
