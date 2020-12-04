package heeglibs

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type CookieInfo struct {
	Jyauth string `json:"jyauth" form:"jyauth"`
	Token  string `json:"__RequestVerificationToken" form:"__RequestVerificationToken"`
}

// 解析cookie数据
// 将头部的cookie或url中的cookie数据解析出来
// 有限解析head中的数据
//
func parseCookie(c *gin.Context) (auth CookieInfo, err error) {
	var cookie CookieInfo
	cookie.Jyauth, err = c.Cookie("jyauth")
	if nil != err {
		goto query
	}

	cookie.Token, err = c.Cookie("__RequestVerificationToken")
	if nil != err {
		goto query
	}

	goto token

query:
	err = c.ShouldBindQuery(&cookie)
	if nil != err {
		c.String(http.StatusUnauthorized, "NOT LOGIN")
		c.Abort()

		return
	}

token:
	if cookie.Token != cookie.Jyauth {
		c.String(http.StatusUnauthorized, "AUTH_ERROR")
		c.Abort()

		return
	}

	auth = cookie
	return
}

// 权限中间件，主要是确认是否登陆成功，设置一个回调函数已在本地服务器中确认
//
// @param callback 	回调函数，用于回传cookie数据
// @param timeout 	token的过期时间
//
func NeedLogin(callback func(c *gin.Context, cookie CookieInfo) bool, timeout int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		now := time.Now().UnixNano()

		for k, v := range c.Request.Header {
			Println(k, v)
		}

		cookie, err := parseCookie(c)
		if nil != err {
			c.String(http.StatusUnauthorized, "NOT_LOGIN"+err.Error())
			c.Abort()

			return
		}

		if cookie.Token != cookie.Jyauth {
			c.String(http.StatusUnauthorized, "AUTH_ERROR")
			c.Abort()

			return
		}

		// 在对应的服务中验证登录是否有效 //
		if !callback(c, cookie) {
			c.String(http.StatusUnauthorized, "NOT_LOGIN")
			c.Abort()

			return
		}

		end := time.Now().UnixNano()
		Println("t--->", end-now)

		c.Next()
	}
}
