package tools


import (
	"github.com/gin-gonic/gin"
	"fmt"
)

func AbortWithJSONError(c *gin.Context, error_code int, err_str string, message gin.H){
  fmt.Printf("%s" , err_str)
  c.AbortWithStatusJSON(error_code, message)
}
