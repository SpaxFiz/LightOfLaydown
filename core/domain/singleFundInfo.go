// @program: unjuanable
// @author: Fizzy
// @created: 2021-12-17

package domain

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spaxfiz/unjuanable/util"
	"regexp"
	"strconv"
	"time"
)

const fundInfoURL = "https://fundgz.1234567.com.cn/js/%v.js?rt=%d"

var (
	p1, _ = regexp.Compile(`"dwjz":"(.+?)"`)
	p2, _ = regexp.Compile(`"name":"(.+?)"`)
)

func SingleFundHandler(ctx *gin.Context) {
	code := ctx.Query("code")
	url := fmt.Sprintf(fundInfoURL, code, time.Now().UnixMilli())
	var data string
	if err := util.EasyGet(url, nil, nil, &data, nil); err != nil {
		ctx.JSON(500, gin.H{"message": err.Error()})
		return
	}
	var price, name string
	if match := p1.FindAllStringSubmatch(data, -1); len(match) > 0 {
		price = match[0][1]
	}
	if match := p2.FindAllStringSubmatch(data, -1); len(match) > 0 {
		name = match[0][1]
	}
	if price != "" && name != "" {
		p, _ := strconv.ParseFloat(price, 64)
		ctx.JSON(200, gin.H{
			"price": p,
			"name":  name,
		})
		return
	}
	ctx.JSON(500, gin.H{"message": "data is null"})
}
