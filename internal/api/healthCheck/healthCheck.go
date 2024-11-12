package healthCheck

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Alive(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}
