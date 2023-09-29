package stats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	type testCase struct {
		sVal   string
		expVal float64
		expErr bool
	}
	testCases := []testCase{
		{
			sVal:   "1.23",
			expVal: 1.23,
		},
		{
			sVal:   "1,23",
			expVal: 1.23,
		},
		{
			sVal:   "0",
			expVal: 0.0,
		},
		{
			sVal:   "0..0",
			expVal: 0.0,
		},
		{
			sVal:   "unknown",
			expVal: 0.0,
		},
		{
			sVal:   "unknown",
			expVal: 1.0,
			expErr: true,
		},
		{
			sVal:   "1:0",
			expVal: 1.0,
			expErr: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run("fix", func(t *testing.T) {
			actVal := SafeParseFloat(tc.sVal)
			if tc.expErr {
				require.NotEqual(t, actVal, tc.expVal)
			} else {
				require.Equal(t, actVal, tc.expVal)
			}
		})
	}
}
