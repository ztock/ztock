package stock

import (
	"context"
	"time"

	"github.com/ztock/ztock/internal/config"
)

// stockSupper is the interface for the stock
type stockSupper interface {
	Get() (*Stock, error)
}

// stockContext represents the socket context
type stockContext struct {
	strategy stockSupper
}

// Stock represents the socket
type Stock struct {
	Name                 string
	Number               string
	PercentageChange     string
	OpeningPrice         string
	PreviousClosingPrice string
	CurrentPrice         string
	HighPrice            string
	LowPrice             string
	Date                 time.Time
}

// NewStockContext creates a new client for stock context client
func NewStockContext(ctx context.Context, platformType config.PlatformType, cfg *config.Config) stockSupper {
	s := new(stockContext)
	switch platformType {
	case config.SinaPlatformType:
		s.strategy = newSinaStock(cfg)
	default:
		s.strategy = newSinaStock(cfg)
	}
	return s
}

// Get will get the stock info from platform's API
func (s stockContext) Get() (*Stock, error) {
	return s.strategy.Get()
}
