package handlers

import (
	"github.com/gin-gonic/gin"
	"log"
)

func getAllCountries(c *gin.Context) {
	log.Print("getAllCountries")
}

func postCountry(c *gin.Context) {
	log.Print("postCountry")
}
