package types

import (
	"fmt"
	"time"

	pingpp "github.com/pingplusplus/pingpp-go/pingpp"
)

type ChargeInput struct {
	Order_no  string `json:"Order_no"`
	AppId     string `json:"AppId"`
	Channel   string `json:"Channel"`
	Amount    uint64 `json:"Amount"`
	Currency  string `json:"Currency"`
	Client_ip string `json:"Client_ip"`
	Subject   string `json:"Subject"`
	Body      string `json:"Body"`
}

type OperResult struct {
	Result  bool
	Message string
}

type GetChargeResult struct {
	OperResult
	Charge pingpp.Charge
}

type UnixTime time.Time

type DataInfo struct {
	AppId          string  `json:"app_id"`
	Object         string  `json:"object"`
	AppDisplayName string  `json:"app_display_name"`
	Created        int64   `json:"created"`
	SummaryFrom    int64   `json:"summary_from"`
	SummaryTo      int64   `json:"summary_to"`
	ChargesAmount  float32 `json:"charges_amount"`
	ChargesCount   float32 `json:"charges_count"`
}

type ObjData struct {
	Object DataInfo `json:"aaa"`
}

type WebhooksEvent struct {
	Id              string  `json:"id"`
	Created         int64   `json:"created"`
	Livemode        bool    `json:"livemode"`
	Stype           string  `json:"type"`
	Object          string  `json:"object"`
	PendingWebhooks int     `json:"pending_webhooks"`
	Request         string  `json:"request"`
	Data            ObjData `json:"data"`
}

// MarshalJSON implements json.Marshaler.
func (t UnixTime) MarshalJSON() ([]byte, error) {
	//do your serializing here
	stamp := fmt.Sprintf("%d", time.Time(t).Unix())
	return []byte(stamp), nil
}
