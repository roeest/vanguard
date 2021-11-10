package vanguard_test

import (
	"github.com/roeest/vanguard"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestETF(t *testing.T) {
	cli := vanguard.NewClient(vanguard.DebugOption)
	const expectedSymbol = "VOO"
	etf, err := cli.GetEtf(expectedSymbol)
	require.NoError(t, err)
	require.Equal(t, expectedSymbol, etf.Symbol)

	h, err := etf.GetHoldings(vanguard.Stock)
	require.NoError(t, err)
	require.NotEmpty(t, h)

	p, err := etf.GetProfile()
	require.NoError(t, err)
	require.True(t, p.ExpenseRatio > 0)

	d, err := etf.GetDiversificationInfo()
	require.NoError(t, err)
	require.NotEmpty(t, d.Sectors)
}
