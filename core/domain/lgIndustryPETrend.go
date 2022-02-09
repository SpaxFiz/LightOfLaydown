// @program: unjuanable
// @author: Fizzy
// @created: 2021-12-02
// @description: industry PE trend

package domain

import (
	"encoding/gob"
	"github.com/sirupsen/logrus"
	"github.com/spaxfiz/unjuanable/core/storage"
	"github.com/spaxfiz/unjuanable/util"
	"time"
)

const (
	industryURL    = "https://api.jiucaishuo.com/v2/guzhi/newtubiaolinedata"
	cacheKeyPrefix = "IND:C:"
)

var peParam = map[string]interface{}{
	"act_time":    time.Now().UnixMilli(),
	"authtoken":   "",
	"data_source": "xichou",
	"gu_code":     "",
	"pe_category": "pe",
	"type":        "pc",
	"version":     "1.7.7",
	"year":        10,
}

var industryCode = map[string]string{
	//"000928.SH": "中证800能源",
	"000991.SH": "全指医药",
	//"000992.SH": "全指金融",
	"000993.SH": "全指信息",
	//"000989.SH": "中证医疗",
	//"399986.SZ": "中证银行",
	"000922.CSI": "中证红利",
	//"000934.SH": "金融地产",
	"GDAXI.GI":   "德国DAX",
	"H30533.CSI": "中国互联50",
	"SPX.GI":     "标普500",
	"HSI.HI":     "恒指",
	"930771.CSI": "新能源",
	"801120.SI":  "食品饮料",
	"000989.SH":  "全指可选",
}

type resp struct {
	Data struct {
		Tubiao struct {
			Series []struct {
				Data [][2]float64
			}
		}
	}
}

type series struct {
	Data   [][2]float64 `json:"data"`
	Name   string       `json:"name"`
	Type   string       `json:"type"`
	Symbol string       `json:"symbol"`
}

type s []*series

type IndustryPETrend struct {
	Data s
}

func (i *IndustryPETrend) Fetch() error {
	fn := func() (interface{}, error) {
		i.Data = make([]*series, 0, len(industryCode))
		for k, v := range industryCode {
			p := genParam(k)
			r := new(resp)
			if err := util.EasyPost(industryURL, p, nil, r); err != nil {
				logrus.Error(err)
				return nil, err
			}
			i.Data = append(i.Data, &series{
				Data:   r.Data.Tubiao.Series[1].Data,
				Name:   v,
				Type:   "line",
				Symbol: "none",
			})
		}
		return i.Data, nil
	}
	return storage.GetCache().LoadOrDo(cacheKeyPrefix, &i.Data, fn)
}

func (i *IndustryPETrend) Render() (interface{}, error) {
	return i.Data, nil
}

func zip(raw [][2]float64) [][2]interface{} {
	rs := make([][2]interface{}, len(raw))
	for i := range raw {
		t := time.Unix(int64(raw[i][0])/1000, 0).Format(util.DateLayout)
		rs[i] = [2]interface{}{t, raw[i][1]}
	}
	return rs
}

func genParam(code string) map[string]interface{} {
	rs := make(map[string]interface{})
	for k, v := range peParam {
		rs[k] = v
	}
	rs["gu_code"] = code
	rs["act_time"] = time.Now().UnixMilli()
	return rs
}

func init() {
	gob.Register(s{})
}
