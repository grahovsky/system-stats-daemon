package stats

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestStats(t *testing.T) {
	t.Parallel()
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
			t.Parallel()
			actVal := SafeParseFloat(tc.sVal)
			if tc.expErr {
				require.NotEqual(t, actVal, tc.expVal)
			} else {
				require.Equal(t, actVal, tc.expVal)
			}
		})
	}
}
