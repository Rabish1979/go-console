package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type Module struct {
	Name string `json:"Name"`
	Path string `json:"Path"`
}

type Process struct {
	ID         string   `json:"ID"`
	Path       string   `json:"Path"`
	ParentID   string   `json:"ParentID"`
	ModuleList []Module `json:"ModuleList"`
}

var processes []Process

var API_KEY = "^%&%*&*&^BDBKBDDJVJHDVJHDVDV"

type Response struct {
	Status  int
	Message []string
	Error   []string
}

func ValidateAPIKey(c *gin.Context) bool {
	apiKey := c.Request.Header.Get("Authorization")

	if apiKey == API_KEY {
		return true
	}

	return false
}

func SendResponse(c *gin.Context, response Response) {
	if len(response.Message) > 0 {
		c.JSON(response.Status, map[string]interface{}{"message": strings.Join(response.Message, "; ")})
	} else if len(response.Error) > 0 {
		c.JSON(response.Status, map[string]interface{}{"error": strings.Join(response.Error, "; ")})
	}
}

// getProcesses responds with the list of all processes as JSON.
func getProcesses(c *gin.Context) {
	if !ValidateAPIKey(c) {
		SendResponse(c, Response{Status: http.StatusUnauthorized, Error: []string{"Invalid API key"}})
		return
	}

	c.IndentedJSON(http.StatusOK, processes)
}

func postProcesses(c *gin.Context) {
	if !ValidateAPIKey(c) {
		SendResponse(c, Response{Status: http.StatusForbidden, Error: []string{"unauthorized access"}})
		return
	}

	if err := c.BindJSON(&processes); err != nil {
		SendResponse(c, Response{Status: http.StatusBadRequest, Error: []string{"One or more params are wrong"}})
		return
	}

	c.IndentedJSON(http.StatusCreated, processes)
}

func main() {
	router := gin.Default()

	router.GET("/processes", getProcesses)
	router.POST("/processes", postProcesses)
	router.Run("0.0.0.0:8080")
}
