package epay

import (
	"net/url"
	"strings"

	"github.com/mitchellh/mapstructure"
)

type PurchaseType string

var (
	Alipay    PurchaseType = "alipay" // Alipay 支付宝
	WechatPay PurchaseType = "wxpay"  // WechatPay 微信
)

type DeviceType string

var (
	PC     DeviceType = "pc"     // PC PC端
	MOBILE DeviceType = "mobile" // MOBILE 移动端
)

type PurchaseArgs struct {
	// 支付类型
	Type PurchaseType
	// 商家订单号
	ServiceTradeNo string
	// 商品名称
	Name string
	// 金额
	Money string
	// 设备类型
	Device    DeviceType
	NotifyUrl *url.URL
	ReturnUrl *url.URL
}

const (
	PurchaseUrl = "/submit.php"
)

// Purchase 生成支付链接和参数
func (c *Client) Purchase(args *PurchaseArgs) (string, map[string]string, error) {
	// see https://payment.moe/doc.html
	requestParams := map[string]string{
		"pid":          c.Config.PartnerID,
		"type":         string(args.Type),
		"out_trade_no": args.ServiceTradeNo,
		"notify_url":   args.NotifyUrl.String(),
		"name":         args.Name,
		"money":        args.Money,
		"device":       string(args.Device),
		"sign_type":    "MD5",
		"return_url":   args.ReturnUrl.String(),
		"sign":         "",
	}

	// 修复：正确拼接路径，避免路径被替换
	u := *c.BaseUrl
	// 确保 BaseUrl 的路径以 / 结尾
	if u.Path != "" && u.Path[len(u.Path)-1] != '/' {
		u.Path += "/"
	}
	// 移除 PurchaseUrl 开头的 /，然后拼接
	u.Path += strings.TrimPrefix(PurchaseUrl, "/")

	return u.String(), GenerateParams(requestParams, c.Config.Key), nil
}

const StatusTradeSuccess = "TRADE_SUCCESS"

type VerifyRes struct {
	// 支付类型
	Type PurchaseType
	// 易支付订单号
	TradeNo string `mapstructure:"trade_no"`
	// 商家订单号
	ServiceTradeNo string `mapstructure:"out_trade_no"`
	// 商品名称
	Name string
	// 金额
	Money string
	// 订单支付状态
	TradeStatus string `mapstructure:"trade_status"`
	// 签名检验
	VerifyStatus bool `mapstructure:"-"`
}

func (c *Client) Verify(params map[string]string) (*VerifyRes, error) {
	sign := params["sign"]
	var verifyRes VerifyRes
	// 从 map 映射到 struct 上
	err := mapstructure.Decode(params, &verifyRes)
	// 验证签名
	verifyRes.VerifyStatus = sign == GenerateParams(params, c.Config.Key)["sign"]
	if err != nil {
		return nil, err
	} else {
		return &verifyRes, nil
	}
}
