package model

import (
	"SpaxFiz/LaydownLight/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

const EmAccountDataUrl = "http://dcfm.eastmoney.com/EM_mutiSvcExpandInterface/api/js/get"

type EmAccount struct {
	Data []*EmAccountData
}

var reqParam = map[string]interface{}{
	"type":  "GPKHData",
	"token": "894050c76af8597a853f5b408b759f5d",
	"st":    "SDATE",
	"sr":    "1",
	"p":     "1",
	"ps":    "200",
}

type EmAccountData struct {
	Date                       string             `json:"SDATE" ld:"date"`
	NewInvestor                util.TolerantFloat `json:"XZSL" ld:"new_investor"`
	NewInvestorSamePeriodBasis util.TolerantFloat `json:"XZTB" ld:"new_investor_same_period_basis"`
	NewInvestorSequentialBasis util.TolerantFloat `json:"XZHB" ld:"new_investor_sequential_basis"`
	EndInvestor                util.TolerantFloat `json:"QMSL" ld:"end_investor"`
	ShIndex                    util.TolerantFloat `json:"SZZS" ld:"sh_index"`
	MarketTotalValue           util.TolerantFloat `json:"HSZSZ" ld:"market_total_value"`
	MarketAvgValue             util.TolerantFloat `json:"HJZSZ" ld:"market_avg_value"`
}

func (e *EmAccount) Fetch() error {
	target := util.UrlParamConstruct(EmAccountDataUrl, reqParam)
	resp, err := http.Get(target)
	if err != nil {
		log.Errorf("error occur when request EmAccountDataUrl. err=%s", err.Error())
		return err
	}

	if resp.StatusCode == http.StatusOK {
		if err := util.ParseBodyWithTrimAsJSON(resp.Body, nil, &e.Data); err != nil {
			return err
		}
	}

	return nil
}

func (e *EmAccount) Render() (interface{}, error) {
	return util.RewriteDataStructJSON(e.Data)
}
