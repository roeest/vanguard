package vanguard

import (
	"strconv"
	"time"
)

type holdingType string

func (t holdingType) resource() string {
	return holdingResource + string(t)
}

const (
	Stock            = holdingType("stock")
	Bond             = holdingType("bond")
	ShortTermReserve = holdingType("short-term-reserve")

	timeFormat      = "2006-01-02T15:04:05-07:00"
	holdingResource = "portfolio-holding/"
)

type rawProfile struct {
	ExpenseRatio     string `json:"expenseRatio"`
	ExpenseRatioAsOf string `json:"expenseRatioAsOfDate"`
	InceptionDate    string `json:"inceptionDate"`
	LongName         string `json:"longName"`
	Symbol           string `json:"ticker"`
}

func (p *rawProfile) toProfile() (Profile, error) {
	var (
		result = Profile{Symbol: p.Symbol, LongName: p.LongName}
		err    error
	)

	result.ExpenseRatio, err = strconv.ParseFloat(p.ExpenseRatio, 64)
	if err != nil {
		return Profile{}, err
	}
	result.ExpenseRatioAsOf, err = parseTime(p.ExpenseRatioAsOf)
	if err != nil {
		return Profile{}, err
	}
	result.InceptionDate, err = parseTime(p.InceptionDate)
	if err != nil {
		return Profile{}, err
	}
	return result, nil
}

type Profile struct {
	ExpenseRatio     float64
	ExpenseRatioAsOf time.Time
	InceptionDate    time.Time
	LongName         string
	Symbol           string
}

type rawHoldings struct {
	Holdings []rawHolding `json:"entity"`
}

func (r rawHoldings) toHoldings() ([]Holding, error) {
	var (
		result = make([]Holding, 0, len(r.Holdings))
		err    error
	)
	for _, h := range r.Holdings {
		var resultHolding = Holding{Symbol: h.Symbol}
		resultHolding.Shares, err = strconv.Atoi(h.SharesHeld)
		if err != nil {
			return nil, err
		}
		resultHolding.AsOf, err = parseTime(h.AsOf)
		if err != nil {
			return nil, err
		}
		result = append(result, resultHolding)
	}
	return result, nil
}

type rawHolding struct {
	AsOf       string `json:"asOfDate"`
	SharesHeld string `json:"sharesHeld"`
	Symbol     string `json:"ticker"`
}

type rawFunds struct {
	Fund rawHoldings `json:"fund"`
	Size int         `json:"size"`
}

type Holding struct {
	AsOf   time.Time
	Shares int
	Symbol string
}

func parseTime(val string) (time.Time, error) {
	t, err := time.Parse(timeFormat, val)
	return t.UTC(), err
}
