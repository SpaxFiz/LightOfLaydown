package util

import (
	log "github.com/sirupsen/logrus"
	"net/http"
)

const ua = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/95.0.4638.54 Safari/537.36"

func EasyGet(target string, reqParam map[string]interface{}, trim *TrimBound, setter interface{}) error {
	url := UrlParamConstruct(target, reqParam)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("user-agent", ua)

	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("error occur when request emAccountDataURL. err=%s", err.Error())
		return err
	}

	if resp.StatusCode == http.StatusOK {
		if err := ParseBodyWithTrimAsJSON(resp.Body, trim, &setter); err != nil {
			return err
		}
	}
	return nil
}
