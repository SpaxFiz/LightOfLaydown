// @program: unjuanable
// @author: Fizzy
// @created: 2021-11-25
// @description: PE trending of Chinese A stock market

package domain

import (
	"encoding/gob"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spaxfiz/unjuanable/core/storage"
	"github.com/spaxfiz/unjuanable/util"
	"sort"
	"sync"
	"time"
)

const PEURL = "https://www.legulegu.com/api/stockdata/index-basic"
const peCacheKey = "PE:C"

type dataGroup string

var marketCodeMapping = map[dataGroup]string{
	"000300.SH": "PE300",
	"000016.SH": "PE50",
	"000905.SH": "PE500",
	"000852.SH": "PE1000",
}

func (d dataGroup) isValid() bool {
	return marketCodeMapping[d] != ""
}

func (d dataGroup) MarshalJSON() ([]byte, error) {
	if !d.isValid() {
		return nil, errors.New("invalid market code")
	}
	return util.JSON.Marshal(marketCodeMapping[d])
}

type PERecord struct {
	AVGPETTM              util.TolerantFloat `json:"ttmPe" ld:"avg_pe_ttm"`         // 等权滚动 pe
	WeightedPETTM         util.TolerantFloat `json:"addTtmPe" ld:"weighted_pe_ttm"` // 加权滚动 pe
	Date                  util.CustomDate    `json:"date" ld:"date"`
	AVGPETTMQuantile      util.TolerantFloat `json:"ttmPeQuantile" ld:"avg_pe_ttm_quantile"`         // 等权 pe 分位数
	WeightedPETTMQuantile util.TolerantFloat `json:"addTtmPeQuantile" ld:"weighted_pe_ttm_quantile"` // 加权 pe 分位数
	ClosePrice            util.TolerantFloat `json:"close" ld:"close_price"`
	MarketCode            dataGroup          `json:"industryCode" ld:"market_code"`
}

type PERecordSeries [][]*PERecord

func (p PERecordSeries) Len() int {
	return len(p)
}

func (p PERecordSeries) Less(i, j int) bool {
	return p[i][0].Date < p[j][0].Date
}

func (p PERecordSeries) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

type PETrend struct {
	Data [][]*PERecord
}

func (P *PETrend) Fetch() error {
	tenYearBefore := time.Now().AddDate(-10, 0, 0).Format(util.DateLayout)
	token, _ := util.GetPECrawlToken()
	fn := func() (interface{}, error) {
		var lock sync.Mutex
		result := make(map[util.CustomDate][]*PERecord)
		final := make(PERecordSeries, 0)

		var sg sync.WaitGroup
		sg.Add(len(marketCodeMapping))
		for item := range marketCodeMapping {
			item := item
			go func() {
				defer sg.Done()
				param := map[string]interface{}{
					"token":     token,
					"indexCode": item,
				}
				current := struct {
					Data []*PERecord `json:"data"`
				}{}
				if err := util.EasyGet(PEURL, param, nil, &current, nil); err != nil {
					log.Errorf("error occur when request PEURL. err=%s", err.Error())
					return
				}
				lock.Lock()
				for _, item := range current.Data {
					if string(item.Date) < tenYearBefore {
						continue
					}
					if result[item.Date] == nil {
						result[item.Date] = make([]*PERecord, 0, 1)
					}
					result[item.Date] = append(result[item.Date], item)
				}
				lock.Unlock()
			}()
		}
		sg.Wait()
		for _, val := range result {
			final = append(final, val)
		}
		sort.Sort(&final)
		return final, nil
	}
	if err := storage.GetCache().LoadOrDo(peCacheKey, &P.Data, fn); err != nil {
		return err
	}
	return nil
}

func (P *PETrend) Render() (interface{}, error) {
	return util.RewriteStructJSON(P.Data)
}

func init() {
	gob.Register(PERecordSeries{})
}
