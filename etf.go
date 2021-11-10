package vanguard

import (
	"fmt"
)

const (
	profileResource            = "profile"
	diversificationResource    = "diversification"
	allEtfHoldingQueryTemplate = "start=%d&count=%d"
)

type Etf struct {
	c      *client
	Symbol string
}

func newEtf(c *client, symbol string) (*Etf, error) {
	var p Profile

	err := c.getResource(symbol, profileResource, &p)
	if err != nil {
		return nil, err
	}

	return &Etf{
		Symbol: symbol,
		c:      c,
	}, nil
}

func (e *Etf) GetHoldings(h holdingType) ([]Holding, error) {
	r, err := e.getCountHoldings(h, 1, 1)
	if err != nil {
		return nil, err
	}
	result := make([]Holding, 0, r.Size)
	for i := 1; i < r.Size; i += 5000 {
		r, err = e.getCountHoldings(h, i, 5000)
		if err != nil {
			return nil, err
		}
		h, err := r.Fund.toHoldings()
		if err != nil {
			return nil, err
		}
		result = append(result, h...)
	}
	return result, nil
}
func (e *Etf) getCountHoldings(h holdingType, start, count int) (rawFunds, error) {
	var r rawFunds
	err := e.c.getResourceWithQueryParams(e.Symbol, h.resource(), &r, fmt.Sprintf(allEtfHoldingQueryTemplate, start, count))
	return r, err
}

func (e *Etf) GetProfile() (Profile, error) {
	var p struct {
		FundProfile rawProfile `json:"fundProfile"`
	}
	err := e.c.getResource(e.Symbol, profileResource, &p)
	if err != nil {
		return Profile{}, err
	}
	return p.FundProfile.toProfile()
}

// type

func (e *Etf) GetDiversificationInfo() (DiversificationInfo, error) {
	var p struct {
		Sector rawDiversification `json:"sector"`
	}
	err := e.c.getResource(e.Symbol, diversificationResource, &p)
	if err != nil {
		return DiversificationInfo{}, err
	}
	return p.Sector.toDiversificationInfo()
}

// https://api.vanguard.com/rs/ire/01/ind/fund/0968/diversification.json
