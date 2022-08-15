package models

import (
	"github.com/ethereum/go-ethereum/common"
	uuid "github.com/satori/go.uuid"
)

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

type Pool struct {
	Id            uuid.UUID `gorm:"type:uuid;primary_key;" json:"id,omitempty"`
	PoolNumber    int
	Token0        common.Address
	Token1        common.Address
	Address       common.Address
	Token0_Symbol string
	Token1_Symbol string
}
