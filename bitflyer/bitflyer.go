package bitflyer

import "log"

import (
	"fmt"
	"sync"
	"context"
	"strconv"
	"encoding/json"
)

type Bitflyer struct {
	client  *Client

	cancel  context.CancelFunc
	mtx     *sync.Mutex
}

func NewBitflyer(api_key string, secret_key string, b_ctx context.Context) (*Bitflyer, error) {
	ctx, cancel := context.WithCancel(b_ctx)
	client := NewClient(api_key, secret_key)
	client.RunPool(ctx)

	self := &Bitflyer {
		client:client,
		cancel:cancel,
		mtx:new(sync.Mutex),
	}
	ok, err := self.checkStatus()
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("store not open.")
	}
	return self, nil
}

func (self *Bitflyer) GetRates(targets []string) (map[string]*Rate, error) {
	self.lock()
	defer self.unlock()

	if targets == nil {
		return nil, fmt.Errorf("not set target.")
	}

	rates := make(map[string]*Rate)
	for _, target := range targets {
		param := "product_code=" + target
		ret, err := self.request2PublicPool("GET", "/v1/getticker", param, nil)
		if err != nil {
			return nil, err
		}
		var r *Rate
		if err := json.Unmarshal(ret, &r); err != nil {
			return nil, err
		}
		if err := r.parseFix(); err != nil {
			return nil, err
		}

		rates[target] = r
	}

	return rates, nil
}

func (self *Bitflyer) LimitOrder(p_code string, o_type string, size float64, price float64) (string, error) {
	self.lock()
	defer self.unlock()

	return self.spotOrder(p_code, o_type, MODE_LIMIT, size, price)
}

func (self *Bitflyer) MarketOrder(p_code string, o_type string, size float64) (string, error) {
	self.lock()
	defer self.unlock()

	return self.spotOrder(p_code, o_type, MODE_MARKET, size, float64(0))
}

func (self *Bitflyer) spotOrder(p_code string, o_type string, mode string, size float64, price float64) (string, error) {
	price_str := strconv.FormatFloat(price, 'f', -1, 64)
	size_str := strconv.FormatFloat(size, 'f', -1, 64)

	var order string
	if mode == MODE_MARKET {
		order = (`{"product_code" : "` + p_code + `", "child_order_type" : "` +
				mode + `", "side":"` + o_type + `", "price": ` + price_str + `,
				"size": ` + size_str + `}`)
	}
	if mode == MODE_LIMIT {
		order = (`{"product_code" : "` + p_code + `", "child_order_type" : "` +
				mode + `", "side":"` + o_type + `"size": ` + size_str + `}`)
	}
	if order == "" {
		return "", fmt.Errorf("undefined mode '%s'", mode)
	}

	ret, err := self.request2PrivatePool("POST", "/v1/me/sendchildorder", "", []byte(order))
	if err != nil {
		return "", err
	}
	var r respOrder
	if err := json.Unmarshal(ret, &r); err != nil {
		return "", err
	}
	if r.Id == "" {
		return "", fmt.Errorf("cannot get response id.")
	}
	return r.Id, nil
}

func (self *Bitflyer) GetOpenOrders(p_code string) ([]*Order, error) {
	self.lock()
	defer self.unlock()

	return self.getOrders(p_code, ORDER_STATE_OPEN)
}

func (self *Bitflyer) GetClosedOrders(p_code string) ([]*Order, error) {
	self.lock()
	defer self.unlock()

	return self.getOrders(p_code, ORDER_STATE_FIXED)
}

func (self *Bitflyer) getOrders(p_code string, state string) ([]*Order, error) {
	param := "product_code=" + p_code + "&child_order_state=" + state
	ret, err := self.request2PrivatePool("GET", "/v1/me/getchildorders", param, nil)
	if err != nil {
		return nil, err
	}
	log.Println(string(ret))
	o := []*Order{}
	if err := json.Unmarshal(ret, &o); err != nil {
		return nil, err
	}
	return o, nil
}

func (self *Bitflyer) GetBalance(c_code string) (*Balance, error) {
	self.lock()
	defer self.unlock()

	ret, err := self.request2PrivatePool("GET", "/v1/me/getbalance", "", nil)
	if err != nil {
		return nil, err
	}

	bs := []*Balance{}
	if err := json.Unmarshal(ret, &bs); err != nil {
		return nil, err
	}
	for _, b := range bs {
		if b.Code != c_code {
			continue
		}
		return b, nil
	}
	return nil, fmt.Errorf("cannot find '%s'", c_code)
}

func (self *Bitflyer) Close() error {
	self.lock()
	defer self.unlock()

	self.cancel()
	return nil
}

func (self *Bitflyer) checkStatus() (bool, error) {
	self.lock()
	defer self.unlock()

	ret, err := self.request2PublicPool("GET", "v1/gethealth", "", nil)
	if err != nil {
		return false, err
	}
	var s *storeStatus
	if err := json.Unmarshal(ret, &s); err != nil {
		return false, err
	}

	if s.Status != STATUS_OPEN {
		return false, nil
	}
	return true, nil
}

func (self *Bitflyer) request2PublicPool(method string, path string, param string, body []byte) ([]byte, error) {
	req, err := self.client.NewRequest(method, path, param, body)
	if err != nil {
		return nil, err
	}
	return self.client.PostPublicPool(req)
}

func (self *Bitflyer) request2PrivatePool(method string, path string, param string, body []byte) ([]byte, error) {
	req, err := self.client.NewRequest(method, path, param, body)
	if err != nil {
		return nil, err
	}
	return self.client.PostPrivatePool(req)
}

func (self *Bitflyer) lock() {
	self.mtx.Lock()
}

func (self *Bitflyer) unlock() {
	self.mtx.Unlock()
}
