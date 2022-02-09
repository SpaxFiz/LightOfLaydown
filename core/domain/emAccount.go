// @description new investor mom trending

package domain

import (
	JSON "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
	"github.com/spaxfiz/unjuanable/util"
	"strconv"
	"time"
)

const emAccountDataURL = "https://datacenter-web.eastmoney.com/api/data/v1/get"

type EmAccount struct {
	Data []*emAccountData
}

var reqParam = map[string]interface{}{
	"reportName":  "RPT_STOCK_OPEN_DATA",
	"columns":     "ALL",
	"pageSize":    "500",
	"sortColumns": "STATISTICS_DATE",
	"sortTypes":   "1",
	"source":      "WEB",
	"client":      "WEB",
	"p":           "1",
	"pageNo":      "1",
	"pageNum":     "1",
	"pageNumber":  "1",
	"_":           strconv.Itoa(int(time.Now().UnixMilli())),
}

type emAccountData struct {
	Date             string             `json:"STATISTICS_DATE" ld:"date"`
	NewInvestor      util.TolerantFloat `json:"ADD_INVESTOR" ld:"new_investor"`
	NewInvestorYOY   util.TolerantFloat `json:"ADD_INVESTOR_YOY" ld:"new_investor_yoy"`
	NewInvestorMOM   util.TolerantFloat `json:"ADD_INVESTOR_QOQ" ld:"new_investor_mom"`
	EndInvestor      util.TolerantFloat `json:"END_INVESTOR" ld:"end_investor"`
	ShIndex          util.TolerantFloat `json:"CLOSE_PRICE" ld:"sh_index"`
	MarketTotalValue util.TolerantFloat `json:"TOTAL_MARKET_CAP" ld:"market_total_value"`
	MarketAvgValue   util.TolerantFloat `json:"AVERAGE_MARKET_CAP" ld:"market_avg_value"`
}

func (e *EmAccount) Fetch() error {
	var raw string
	if err := util.EasyGet(emAccountDataURL, reqParam, &util.TrimBound{Pre: "[", Sur: "]"}, &raw, nil); err != nil {
		log.Errorf("error occur when request emAccountDataURL. err=%s", err.Error())
		return err
	}
	JSON.UnmarshalFromString(raw, &e.Data)
	return nil
}

func (e *EmAccount) Render() (interface{}, error) {
	return util.RewriteStructJSON(e.Data)
}
