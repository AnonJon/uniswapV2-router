package models

import "github.com/ethereum/go-ethereum/common"

type Pair struct {
	TokenA common.Address
	TokenB common.Address
	Amount float64
}

type Quote struct {
	Route []*Pair
	Rate  float64
	Path  []*Path
}

type PairContract struct {
	Id     common.Address `json:"id"`
	Token0 *Token         `json:"token0"`
	Token1 *Token         `json:"token1"`
}

type Token struct {
	Id     common.Address `json:"id"`
	Symbol string         `json:"symbol"`
}

type Path struct {
	Address common.Address
	Symbols []string
}
