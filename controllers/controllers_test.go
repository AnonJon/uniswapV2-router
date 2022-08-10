package controllers_test

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jongregis/uniswapV2_router/controllers"
	"github.com/jongregis/uniswapV2_router/models"
	"github.com/jongregis/uniswapV2_router/server"
	"github.com/stretchr/testify/assert"
)

func TestWeiToEther(t *testing.T) {
	wei := big.NewInt(1100000000000000000)
	expected := float64(1.1)
	actual := controllers.WeiToEther(wei)
	actualFloat, _ := actual.Float64()
	assert.Equal(t, expected, actualFloat)

}

func TestQueryPair(t *testing.T) {
	pair := &models.Pair{
		TokenA: common.HexToAddress("0xb4efd85c19999d84251304bda99e90b92300bd93"),
		TokenB: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	}
	pair2 := &models.Pair{
		TokenA: common.HexToAddress("0x0000000000000000000000000000000000000000"),
		TokenB: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	}
	handler := server.NewHandler()
	addr, err := controllers.QueryPair(pair, handler.Client)
	assert.Nil(t, err)
	assert.Equal(t, "0x70EA56e46266f0137BAc6B75710e3546f47C855D", addr.Hex())

	// Test with invalid pair
	_, err = controllers.QueryPair(pair2, handler.Client)
	assert.NotNil(t, err)
	assert.Equal(t, "pair not found", err.Error())
}

func TestQueryRate(t *testing.T) {
	pair := &models.Pair{
		TokenA: common.HexToAddress("0xb4efd85c19999d84251304bda99e90b92300bd93"),
		TokenB: common.HexToAddress("0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"),
	}
	handler := server.NewHandler()
	rate, err := controllers.QueryRate(pair, handler.Client)
	assert.Nil(t, err)
	assert.Equal(t, float64(0.0011), *rate)
}
