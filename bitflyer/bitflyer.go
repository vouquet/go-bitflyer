package bitflyer

import (
	"fmt"
	"sync"
	"context"
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
