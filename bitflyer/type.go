package bitflyer

import (
	"time"
)

type storeStatus struct {
	Status     string    `json:"status"`
}

type Rate struct {
	RawAsk     float64   `json:"best_ask"`
	RawBid     float64   `json:"best_bid"`
	RawProduct string    `json:"product_code"`
	T          time.Time `json:"timestamp"`
	RawVolume  float64   `json:"volume"`
	RawLast    float64   `json:"last"`
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
	return self.T
}

func (self *Rate) Volume() float64 {
	return self.RawVolume
}

func (self *Rate) Last() float64 {
	return self.RawLast
}
