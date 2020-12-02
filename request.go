package youyoulibs

import (
	"github.com/gin-gonic/gin"
)

type RequestPage struct {
	Page int `form:"page"`
	Size int `form:"size"`
}

// 从请求中获取page和size，主要是获取请求区间
func PageSize(c *gin.Context) RequestPage {
	var param RequestPage

	err := c.ShouldBindQuery(&param)
	if err != nil {
		param.Page = 0
		param.Size = 20

		return param
	}

	if 0 >= param.Page {
		param.Page = 0
	}

	if 0 >= param.Size {
		param.Size = 20
	}

	if 0 < param.Page {
		param.Page = param.Page - 1
	}

	param.Page = param.Page * param.Size

	return param
}
