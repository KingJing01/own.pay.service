package controllers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"paymentservice/types"
	"strings"

	"github.com/astaxie/beego"
	"github.com/pingplusplus/pingpp-go/pingpp"
	"github.com/pingplusplus/pingpp-go/pingpp/charge"
)

type PayserviceController struct {
	beego.Controller
}

//获取支付对象charge
// @Title GetCharge
// @Description 支付对象charge
// @Param   data     body     types.ChargeInput        "请求数据"
// @Success 200 {string} success
// @Failure 400 Invalid email supplied
// @Failure 404 User not found
// @router /GetCharge [post]
func (p *PayserviceController) GetCharge() {
	cinput := &types.ChargeInput{}
	lresult := types.GetChargeResult{}

	json.Unmarshal(p.Ctx.Input.RequestBody, cinput)

	privateKey, err := ioutil.ReadFile("my_rsa_private_key.pem")
	if err != nil {
		lresult.Result = false
		lresult.Message = err.Error()
		p.Data["json"] = lresult
		p.ServeJSON()
		return
	}

	pingpp.Key = "sk_test_ijz50CTy5SK4evvb18ifXjfP"
	pingpp.AccountPrivateKey = string(privateKey)
	// pingpp.AccountPrivateKey = "-----BEGIN RSA PRIVATE KEY-----\n"
	// pingpp.AccountPrivateKey += "MIICXQIBAAKBgQCuY98HS04o2PeIV02irshJAoxUHpzQ2SbCGbxnvikixmmxMZT1\n"
	// pingpp.AccountPrivateKey += "bgaXMjrWJ9ZBehnJaNEiZb0BJN2VQFES7KsrFtM54SVhGIYfG42nNEvzUxC4Pfw+\n"
	// pingpp.AccountPrivateKey += "H9591N0bBG41XZi0u+ewQC8rVe44lFmdd6cvOwL8newDOXZVgq2ZW9f1+QIDAQAB\n"
	// pingpp.AccountPrivateKey += "AoGAe0gLEOMBnArV4sKlFY6t9D8i6QzDGzmIFsuOz2A1QGo3qZY9kct4SUavJVoA\n"
	// pingpp.AccountPrivateKey += "M0WYGTRKNCDsHnTrWGmhZtl8nsvRdjj4IG2XF/ncOiHm5kp16U8itw/+oTTDwNO+\n"
	// pingpp.AccountPrivateKey += "OQEkAxPQF0PEn5Mt/ieLmM02+nBVxB7b1+nDT/f8sTFCCAECQQDk0uy+ndMfs+4n\n"
	// pingpp.AccountPrivateKey += "rtK8obRHy51GLLUaDRQCTKq1ES9Kj4QAJ3IJX73jg4ut5gKCKcAsS0lkDnuaWD66\n"
	// pingpp.AccountPrivateKey += "PDJfWT2BAkEAwxn2PL2s59w0uhJfYTbkq9S/4OxxqBseMRcqZGNOxbo6iWgPVfx2\n"
	// pingpp.AccountPrivateKey += "8ON1BzKRPZEbVnXmJXQ9vFdQOD5th2PkeQJBAKucoR9ogFFzgXZTgAsmf22lAIQD\n"
	// pingpp.AccountPrivateKey += "vaMXEd2ToCeCBuS1c7sl2jm7i09ZdeVq7pCuPUk7AYS/8+VSr2C/CsxFwoECQQC+\n"
	// pingpp.AccountPrivateKey += "TLB1hr0EWzHC3PDTretV/2o5ReeGhQzp7SKYUJUhIAjAxhNPV7XcOMCJiLVKTCNS\n"
	// pingpp.AccountPrivateKey += "LiWSGtOsxa2lbp7/FFxhAkAn/sW72iBZLeQ9J/V13jiOIvVik24lRyKGViKP4AaJ\n"
	// pingpp.AccountPrivateKey += "89BX5s7iNlESd+mB2YTglUJu5opCmxV2zHbP/muFgnBg\n"
	// pingpp.AccountPrivateKey += "-----END RSA PRIVATE KEY-----\n"

	ip := p.Ctx.Request.Host
	ip = ip[0:strings.LastIndex(ip, ":")]
	params := &pingpp.ChargeParams{
		Order_no:  cinput.Order_no,
		App:       pingpp.App{Id: cinput.AppId},
		Channel:   cinput.Channel,
		Amount:    cinput.Amount,
		Currency:  cinput.Currency,
		Client_ip: ip,
		Subject:   cinput.Subject,
		Body:      cinput.Body,
	}

	ch, e := charge.New(params)

	if e != nil {
		lresult.Result = false
		lresult.Message = e.Error()
		p.Data["json"] = lresult
		p.ServeJSON()
		return
	} else {
		lresult.Result = true
		lresult.Charge = *ch
		p.Data["json"] = lresult
		p.ServeJSON()
		return
	}
}

// @router /PayCharge [post]
func (p *PayserviceController) PayCharge() {
	//ch, err := charge.New(&params)
}

// @router /ChargeWebhooks [post]
func (p *PayserviceController) ChargeWebhooks() {

	//示例 - 签名在头部信息的 x-pingplusplus-signature 字段
	signed := p.Ctx.Request.Header.Get("X-Pingplusplus-Signature")
	if signed == "" {
		signed = p.Ctx.Request.Header.Get("x-pingplusplus-signature")
	}

	//示例 - 待验签的数据
	buf := p.Ctx.Input.RequestBody
	event := string(buf)

	// 请从 https://dashboard.pingxx.com 获取「Ping++ 公钥」
	publicKey, err := ioutil.ReadFile("pingpp_rsa_public_key.pem")
	if err != nil {
		fmt.Errorf("read failure: %v", err)
	}

	//base64解码再验证
	decodeStr, _ := base64.StdEncoding.DecodeString(signed)
	errs := pingpp.Verify([]byte(event), publicKey, decodeStr)
	if errs != nil {
		fmt.Println(errs)
		return
	} else {
		fmt.Println("success")
	}

	webhook, err := pingpp.ParseWebhooks(buf)
	//fmt.Println(webhook.Type)
	if err != nil {
		p.Ctx.Output.SetStatus(http.StatusInternalServerError)
		return
	}

	if webhook.Type == "charge.succeeded" {
		// TODO your code for charge
		p.Ctx.Output.SetStatus(http.StatusOK)
	} else if webhook.Type == "refund.succeeded" {
		// TODO your code for refund
		p.Ctx.Output.SetStatus(http.StatusOK)
	} else {
		p.Ctx.Output.SetStatus(http.StatusInternalServerError)
	}

	///////////////////////////////////////////////
	//cinput := &types.WebhooksEvent{}
	lresult := types.GetChargeResult{}
	// res, _ := simplejson.NewJson(p.Ctx.Input.RequestBody)
	// DataInfoJson := res.Get("data").Get("object")
	// test := DataInfoJson.Get("app_id").MustString()
	// fmt.Println(test)
	// //json.Unmarshal([]byte(DataInfoJson), dataInfo)
	// //json.Unmarshal(p.Ctx.Input.RequestBody, cinput)
	// p.Ctx.Output.SetStatus(200)
	lresult.Result = false

	p.Data["json"] = lresult
	p.ServeJSON()
	return
}

// func WebhooksVerify(dataString string, signatureString string, PublicKey publicKey) {

// }
