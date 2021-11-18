package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/lshsuper/go-pkg/src/oauth/jwt"
	"net/http"
	"time"
)


//Authorization
func Authorization(ctx *gin.Context){
    //授权
	expire:= time.Now().Add(time.Second*1000).Unix()
	token,_:=jwt.GetToken(jwt.TokenRequest{
		SigningKey: "ABCDEF",
        Expire: int(expire),
	})

	ctx.JSON(http.StatusOK,gin.H{
		"token":token,
	})
	return

}
