package stock

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/ztock/ztock/internal/config"
)

func TestNewSinaStock(t *testing.T) {
	assert := assert.New(t)

	index := config.ShangHaiIndexType
	number := "600000"
	cfg := newSinaStock(&config.Config{
		Index:  index,
		Number: number,
	})

	assert.NotNil(cfg)
	assert.Equal(cfg.index, index)
	assert.Equal(cfg.number, number)
}

type testcase struct {
	name         string
	index        config.IndexType
	number       string
	listResponse string
	socket       Stock
	test         func(test *testcase) error
}

func TestGet(t *testing.T) {
	assert := assert.New(t)

	index := config.ShangHaiIndexType
	number := "600000"
	listResponse := `var hq_str_sh600000="test,10.880,10.900,11.120,11.240,10.880,11.110,11.120,113524449,1259911417.000,212600,11.110,273600,11.100,137200,11.090,284200,11.080,147440,11.070,322322,11.120,220500,11.130,313900,11.140,694690,11.150,142367,11.160,2021-03-15,15:00:00,00,";`
	listResponseErr := "test,10.880,10.900,11.120,11.240,10.880,11.110,11.120,113524449,1259911417.000,212600,11.110,273600,11.100,137200,11.090,284200,11.080,147440,11.070,322322,11.120,220500,11.130,313900,11.140,694690,11.150,142367,11.160,2021-03-15,15:00:00,00,"
	date, err := time.Parse(sinaTimeLayout, "2021-03-15 15:00:00")
	if err != nil {
		t.Errorf("time parse failed: err=%#v", err)
	}
	stock := Stock{
		Name:                 "test",
		Number:               number,
		OpeningPrice:         "10.880",
		PreviousClosingPrice: "10.900",
		CurrentPrice:         "11.120",
		HighPrice:            "11.240",
		LowPrice:             "10.880",
		Date:                 date,
	}

	tests := []testcase{
		{
			name:         "Get_Success",
			index:        index,
			number:       number,
			listResponse: listResponse,
			socket:       stock,
			test: func(testcase *testcase) error {
				s := newSinaStock(&config.Config{
					Index:  testcase.index,
					Number: testcase.number,
				})

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s//%s", defaultSinaProtocol, defaultSinaHost),
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(http.StatusOK, testcase.listResponse)
						return resp, nil
					})

				data, err := s.Get()
				if err != nil {
					return err
				}

				assert.Equal(data.Name, testcase.socket.Name)
				assert.Equal(data.Number, testcase.socket.Number)
				assert.Equal(data.OpeningPrice, testcase.socket.OpeningPrice)
				assert.Equal(data.PreviousClosingPrice, testcase.socket.PreviousClosingPrice)
				assert.Equal(data.CurrentPrice, testcase.socket.CurrentPrice)
				assert.Equal(data.HighPrice, testcase.socket.HighPrice)
				assert.Equal(data.LowPrice, testcase.socket.LowPrice)
				assert.Equal(data.Date, testcase.socket.Date)
				return nil
			},
		},
		{
			name:         "Get_List_Platform_Response_Error",
			index:        index,
			number:       number,
			listResponse: listResponseErr,
			socket:       stock,
			test: func(testcase *testcase) error {
				s := newSinaStock(&config.Config{
					Index:  testcase.index,
					Number: testcase.number,
				})

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s//%s", defaultSinaProtocol, defaultSinaHost),
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(http.StatusOK, testcase.listResponse)
						return resp, nil
					})

				_, err := s.Get()
				assert.EqualError(err, "Sina platform returns wrong data")
				return nil
			},
		},
		{
			name:         "Get_List_Platform_Response_404",
			index:        index,
			number:       number,
			listResponse: listResponse,
			socket:       stock,
			test: func(testcase *testcase) error {
				s := newSinaStock(&config.Config{
					Index:  testcase.index,
					Number: testcase.number,
				})

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s//%s", defaultSinaProtocol, defaultSinaHost),
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(http.StatusNotFound, testcase.listResponse)
						return resp, nil
					})

				_, err := s.Get()
				assert.EqualError(err, fmt.Sprintf("Ztock Client Error, status: %d", http.StatusNotFound))
				return nil
			},
		},
		{
			name:         "Get_List_Platform_Response_500",
			index:        index,
			number:       number,
			listResponse: listResponse,
			socket:       stock,
			test: func(testcase *testcase) error {
				s := newSinaStock(&config.Config{
					Index:  testcase.index,
					Number: testcase.number,
				})

				httpmock.Activate()
				defer httpmock.DeactivateAndReset()

				httpmock.RegisterResponder(http.MethodGet, fmt.Sprintf("%s//%s", defaultSinaProtocol, defaultSinaHost),
					func(req *http.Request) (*http.Response, error) {
						resp := httpmock.NewStringResponse(http.StatusInternalServerError, testcase.listResponse)
						return resp, nil
					})

				_, err := s.Get()
				assert.EqualError(err, fmt.Sprintf("Sina Server Error, status: %d", http.StatusInternalServerError))
				return nil
			},
		},
	}

	for _, testcase := range tests {
		err := testcase.test(&testcase)
		if err != nil {
			t.Errorf("%s failed: err=%q", testcase.name, err)
		}
	}
}
