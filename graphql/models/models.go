package models

import (
	"github.com/ethereum/go-ethereum/common"
)

type Pair struct {
	Id     common.Address
	Token0 *Token
	Token1 *Token
}
type Token struct {
	Id     common.Address
	Symbol string
}
