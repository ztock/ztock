package util

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPercentageChange(t *testing.T) {
	old, new := 20, 60
	pct := PercentageChange(old, new)

	if int(pct) != 200.0 {
		t.Fatalf("%f is wrong percent!", pct)
	}
}

func TestPercentageChangeFloat(t *testing.T) {
	old, new := 20.0, 60.0
	pct := PercentageChangeFloat(old, new)

	if int(pct) != 200.0 {
		t.Fatalf("%f is wrong percent!", pct)
	}
}

type testcase struct {
	name string
	old  string
	new  string
	test func(test *testcase) error
}

func TestPercentageChangeString(t *testing.T) {
	assert := assert.New(t)

	tests := []testcase{
		{
			name: "Calculate_Success",
			old:  "20.0",
			new:  "60.0",
			test: func(testcase *testcase) error {
				pct, err := PercentageChangeString(testcase.old, testcase.new)

				assert.Nil(err)
				assert.Equal(int(pct), 200)
				return nil
			},
		},
		{
			name: "String_Error",
			old:  "old",
			new:  "new",
			test: func(testcase *testcase) error {
				_, err := PercentageChangeString(testcase.old, testcase.new)

				fmt.Println(err)
				assert.EqualError(err, `strconv.ParseFloat: parsing "old": invalid syntax`)
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
