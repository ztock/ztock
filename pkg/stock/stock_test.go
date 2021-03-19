package stock

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ztock/ztock/internal/config"
)

func TestNewStockContext(t *testing.T) {
	assert := assert.New(t)

	index := config.ShangHaiIndexType
	number := "600000"
	s := NewStockContext(context.Background(), config.SinaPlatformType, &config.Config{
		Index:  index,
		Number: number,
	})

	assert.NotNil(s)
	assert.NotNil(reflect.ValueOf(s).Elem())
}
