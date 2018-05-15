package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type response struct {
	Result interface{} `json:"result"`
	Err    string      `json:"err"`
}

type respHandler func(c *gin.Context) (interface{}, error)

func resp(handler respHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		result, err := handler(c)
		r := &response{}
		if err != nil {
			r.Err = err.Error()
		} else {
			r.Result = result
		}
		if r.Err == "Not authorized" {
			c.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			c.AbortWithStatus(401)
			c.JSON(http.StatusUnauthorized, r)
		} else {
			c.JSON(http.StatusOK, r)
		}
	}
}
