package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/samber/lo"
	epay "github.com/star-horizon/go-epay"
)

func main() {
	client := epay.NewClient(&epay.Config{
		PartnerID: "1000",
		Key:       "KEY",
	})
	notify, _ := url.Parse("http://testdomain.natappfree.cc/verify")
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		url, params, err := client.Purchase(&epay.PurchaseArgs{
			Type:           epay.WechatPay,
			ServiceTradeNo: "20160806151343593",
			Name:           "测试商品",
			Money:          "0.01",
			Device:         epay.PC,
			NotifyUrl:      notify,
			ReturnUrl:      notify,
		})
		if err != nil {
			log.Println(err)
			return
		}

		html := "<form id='alipaysubmit' name='alipaysubmit' action='" + url + "' method='POST'>"
		for key, value := range params {
			html += "<input type='hidden' name='" + key + "' value='" + value + "'/>"
		}
		html += "<input type='submit'>POST</form>"

		writer.Header().Set("Content-Type", "text/html")
		writer.Write([]byte(html))
	})
	mux.HandleFunc("/verify", func(writer http.ResponseWriter, request *http.Request) {
		params := lo.Reduce(lo.Keys(request.URL.Query()), func(r map[string]string, t string, i int) map[string]string {
			r[t] = request.URL.Query().Get(t)
			return r
		}, map[string]string{})

		verifyInfo, err := client.Verify(params)
		if err == nil && verifyInfo.VerifyStatus {
			writer.Write([]byte("success"))
		} else {
			writer.Write([]byte("fail"))
		}

		if verifyInfo.TradeStatus == epay.StatusTradeSuccess {
			log.Println(verifyInfo)
		}
	})
	http.ListenAndServe(":8080", mux)
}
