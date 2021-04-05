package bitflyer

import (
	"time"
)

type storeStatus struct {
	Status     string    `json:"status"`
}

type Rate struct {
	RawAsk     float64 `json:"best_ask"`
	RawBid     float64 `json:"best_bid"`
	RawProduct string  `json:"product_code"`
	RawTime    string  `json:"timestamp"`
	RawVolume  float64 `json:"volume"`
	RawLast    float64 `json:"last"`

	time       time.Time
}

func (self *Rate) Ask() float64 {
	return self.RawAsk
}

func (self *Rate) Bid() float64 {
	return self.RawBid
}

func (self *Rate) ProductCode() string {
	return self.RawProduct
}

func (self *Rate) Time() time.Time {
	return self.time
}

func (self *Rate) Volume() float64 {
	return self.RawVolume
}

func (self *Rate) Last() float64 {
	return self.RawLast
}

func (self *Rate) parseFix() error {
	t, err := time.Parse("2006-01-02T15:04:05.999999999", self.RawTime)
	if err != nil {
		return err
	}
	self.time = t

	return nil
}
