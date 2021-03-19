package stock

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"

	"github.com/ztock/ztock/internal/config"
	"github.com/ztock/ztock/pkg/stock/util"
)

const (
	// defaultSinaHost defines default sina host for stock service
	defaultSinaHost = "hq.sinajs.cn"

	// defaultSinaProtocol defines default protocol for stock service
	defaultSinaProtocol = "https:"

	// sinaTimeLayout defines default time layout
	sinaTimeLayout = "2006-01-02 15:04:05"
)

// sinaStock represents the sina stock
type sinaStock struct {
	number string
	index  config.IndexType
}

// newSinaStock a new client for sina stock client
func newSinaStock(cfg *config.Config) sinaStock {
	return sinaStock{
		number: cfg.Number,
		index:  cfg.Index,
	}
}

// Get will get the stock info from platform's API
func (s sinaStock) Get() (*Stock, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s//%s", defaultSinaProtocol, defaultSinaHost), nil)
	if err != nil {
		return nil, err
	}

	// Add query params
	q := req.URL.Query()
	q.Add("list", fmt.Sprintf("%s%s", s.index, s.number))
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode
	if statusCode >= http.StatusInternalServerError {
		return nil, fmt.Errorf("Sina Server Error, status: %d", statusCode)
	} else if statusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("Ztock Client Error, status: %d", statusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Match useful data
	data := regexp.MustCompile(`[\\"\\,]+`).Split(string(body), -1)
	if len(data) < 35 {
		return nil, errors.New("Sina platform returns wrong data")
	}

	// Time parse
	t, err := time.Parse(sinaTimeLayout, fmt.Sprintf("%s %s", data[31], data[32]))
	if err != nil {
		return nil, err
	}

	// Calculate the percent increase/decrease.
	pc, err := util.PercentageChangeString(data[3], data[4])
	if err != nil {
		return nil, err
	}

	return &Stock{
		Name:                 data[1],
		Number:               s.number,
		PercentageChange:     fmt.Sprintf("%.2f%%", pc),
		OpeningPrice:         data[2],
		PreviousClosingPrice: data[3],
		CurrentPrice:         data[4],
		HighPrice:            data[5],
		LowPrice:             data[6],
		Date:                 t,
	}, nil
}
