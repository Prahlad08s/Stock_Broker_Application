package setuptest

import "github.com/gin-gonic/gin"

func SetContext(ctx *gin.Context) {
	ctx.Request.Header.Set("Content-Type", "application/json")

}
