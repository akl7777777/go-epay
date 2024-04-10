package epay

import (
	"log"
	"testing"

	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
)

func TestParamsFilter(t *testing.T) {
	asserts := assert.New(t)
	{
		requestParams := map[string]string{
			"pid":          "1",
			"type":         "1",
			"out_trade_no": "1",
			"notify_url":   "1",
			"name":         "1",
			"money":        "1",
			"device":       "1",
			"sign_type":    "MD5",
			"sign":         "",
		}
		expected := map[string]string{
			"pid":          "1",
			"type":         "1",
			"out_trade_no": "1",
			"notify_url":   "1",
			"name":         "1",
			"money":        "1",
			"device":       "1",
		}

		asserts.EqualValues(expected, ParamsFilter(requestParams))
	}
}

func TestParamsSort(t *testing.T) {
	asserts := assert.New(t)
	{
		requestParams := map[string]string{
			"pid":  "pidv",
			"type": "typev",
			//
			"out_trade_no": "out_trade_nov",
			//
			"notify_url": "notify_urlv",
			//
			"name": "namev",
			//
			"money": "moneyv",
			//
			"device": "devicev",
		}

		expectedKeys := []string{"device", "money", "name", "notify_url", "out_trade_no", "pid", "type"}
		expectedValues := []string{"devicev", "moneyv", "namev", "notify_urlv", "out_trade_nov", "pidv", "typev"}

		keys, values := ParamsSort(requestParams)
		asserts.EqualValues(expectedKeys, keys)
		asserts.EqualValues(expectedValues, values)
	}
}

func TestCreateUrlString(t *testing.T) {
	asserts := assert.New(t)
	{
		expectedKeys := []string{"device", "money"}
		expectedValues := []string{"devicev", "moneyv"}

		expectedString := "device=devicev&money=moneyv"
		asserts.EqualValues(expectedString, CreateUrlString(expectedKeys, expectedValues))
	}
	{
		expectedKeys := []string{}
		expectedValues := []string{}

		expectedString := ""
		asserts.EqualValues(expectedString, CreateUrlString(expectedKeys, expectedValues))
	}
}

func TestMD5String(t *testing.T) {
	asserts := assert.New(t)
	{
		urlString := "device=devicev&money=moneyv"
		key := "1234567"
		expected := "3854cc9f022e0fb821bd2e002260245d"
		asserts.EqualValues(expected, MD5String(urlString, key))
	}
}

func TestMapToStruct(t *testing.T) {
	{
		var verifyInfo VerifyRes
		mapstructure.Decode(map[string]string{
			"pid":  "pidv",
			"type": "typev",
			//
			"out_trade_no": "out_trade_nov",
			//
			"notify_url": "notify_urlv",
			//
			"name": "namev",
			//
			"money": "moneyv",
			//
			"device": "devicev",
		}, &verifyInfo)
		log.Println(verifyInfo)
	}
}
