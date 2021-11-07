package accrual

import (
	"encoding/json"
	"fmt"
	"github.com/im-tollu/go-musthave-diploma-tpl/model"
	"log"
	"net/http"
	"time"
)

type Client struct {
	http.Client

	Addr string
}

func NewClient(addr string) *Client {
	client := http.Client{}
	client.Timeout = 1 * time.Second

	return &Client{
		Client: client,
		Addr:   addr,
	}
}

func (c *Client) GetOrderAccruals(orderNr int64) (model.OrderAccrual, error) {
	url := fmt.Sprintf("%s/api/orders/%d", c.Addr, orderNr)
	log.Printf("requesting accrual: %s", url)

	response, errReq := c.Get(url)
	if errReq != nil {
		return model.OrderAccrual{}, fmt.Errorf("cannot request accrual server: %w", errReq)
	}
	defer response.Body.Close()

	if contentType := response.Header.Get("Content-Type"); contentType != "application/json" {
		return model.OrderAccrual{}, fmt.Errorf("unexpected content type [%s]", contentType)
	}

	accrualJSON := OrderAccrualJSON{}

	dec := json.NewDecoder(response.Body)
	if err := dec.Decode(&accrualJSON); err != nil {
		return model.OrderAccrual{}, fmt.Errorf("cannot parse accrual service response: %w", err)
	}

	orderAccrual, errConv := accrualJSON.ToOrderAccrual()
	if errConv != nil {
		return model.OrderAccrual{}, fmt.Errorf("cannot convert response to Accrual: %w", errConv)
	}

	return orderAccrual, nil
}
