package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/roeest/vanguard"
)

func main() {
	symbol := flag.String("symbol", "VOO", "The name of the symbol of the ETF to query")
	flag.Parse()
	cli := vanguard.NewClient()
	etf, err := cli.GetEtf(*symbol)
	if err != nil {
		log.Fatalln("failed getting etf", err)
	}
	p, err := etf.GetProfile()
	if err != nil {
		log.Fatalln("failed getting profile", err)
	}
	fmt.Println(p)
	h, err := etf.GetHoldings(vanguard.Stock)
	if err != nil {
		log.Fatalln("failed getting holdings", err)
	}
	fmt.Println("Total holdings:", len(h))
}
