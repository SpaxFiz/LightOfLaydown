// @description new funds mom trending

package domain

import (
	"SpaxFiz/LaydownLight/util"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sort"
)

const emNewFundDataURL = "http://fund.eastmoney.com/data/FundNewIssue.aspx"

var emNewFundParam = map[string]interface{}{
	"t":     "xcln",
	"sort":  "jzrgq,desc",
	"y":     "",
	"page":  "1,5000",
	"isbuy": "1",
	"v":     "0.4069919776543214",
}

type EmNewFund struct {
	Data []*emNewFundData
}

func (e *EmNewFund) Len() int {
	return len(e.Data)
}

func (e *EmNewFund) Less(i, j int) bool {
	return e.Data[i].Date < e.Data[j].Date
}

func (e *EmNewFund) Swap(i, j int) {
	e.Data[i], e.Data[j] = e.Data[j], e.Data[i]
}

type emNewFundData struct {
	Date  string
	Count map[string]int
}

func (e *EmNewFund) Fetch() error {
	target := util.UrlParamConstruct(emNewFundDataURL, emNewFundParam)
	resp, err := http.Get(target)
	if err != nil {
		log.Errorf("error occur when request emAccountDataURL. err=%s", err.Error())
		return err
	}

	var data [][]string
	if resp.StatusCode == http.StatusOK {
		if err := util.ParseBodyWithTrimAsJSON(resp.Body, &util.TrimBound{Pre: "[", Sur: "]"}, &data); err != nil {
			return err
		}
	}
	return e.transformData(data)
}

func (e *EmNewFund) Render() (interface{}, error) {
	return e.Data, nil
}

func (e *EmNewFund) transformData(rawData [][]string) (err error) {
	defer func() {
		if e := recover(); e != nil {
			log.Errorf("EmNewFund transformData failed. err=%v", e)
			err = e.(error)
		}
	}()
	result := make(map[string]*emNewFundData)
	for i := len(rawData) - 1; i >= 0; i-- {
		d := rawData[i][6][:7]
		fundType := rawData[i][4]
		if result[d] == nil {
			result[d] = &emNewFundData{
				Date:  d,
				Count: make(map[string]int),
			}
		}
		result[d].Count[fundType]++
	}
	e.Data = make([]*emNewFundData, 0, len(result))
	for _, val := range result {
		e.Data = append(e.Data, val)
	}
	sort.Sort(e)
	return nil
}
