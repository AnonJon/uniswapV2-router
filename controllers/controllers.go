package controllers

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	factory "github.com/jongregis/uniswapV2_router/contracts/factory"
	pairContract "github.com/jongregis/uniswapV2_router/contracts/pair"
	"github.com/jongregis/uniswapV2_router/graphql"
	"github.com/jongregis/uniswapV2_router/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	FACTORY     = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	NIL_ADDRESS = "0x0000000000000000000000000000000000000000"
)

type Controller struct {
	Client *graphql.SubGraph
	DB     *gorm.DB
}

func NewController(db *gorm.DB) *Controller {
	return &Controller{Client: graphql.NewSubGraph(), DB: db}
}

func (c *Controller) QueryPair(pair *models.Pair, client *ethclient.Client) (*common.Address, error) {
	exc, err := factory.NewFactoryCaller(common.HexToAddress(FACTORY), client)
	if err != nil {
		logrus.Error("Failed to instantiate the UniswapV2 Factory contractfrom address")
		return nil, err
	}
	addr, err := exc.GetPair(&bind.CallOpts{}, pair.TokenA, pair.TokenB)
	if err != nil {
		return nil, err
	}
	if addr.String() == NIL_ADDRESS {
		return nil, fmt.Errorf("pair not found")
	}

	return &addr, nil
}

func (c *Controller) QueryRate(pair *models.Pair, client *ethclient.Client) (*float64, error) {
	var r float64
	addr, err := c.QueryPair(pair, client)
	if err != nil {
		return nil, err
	}
	exc, err := pairContract.NewPairCaller(*addr, client)
	if err != nil {
		logrus.Error("Failed to instantiate the UniswapV2 Pair contract from address")
		return nil, err
	}
	reserves, err := exc.GetReserves(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	token0, err := exc.Token0(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	r0f := WeiToEther(reserves.Reserve0)
	r1f := WeiToEther(reserves.Reserve1)

	r1, _ := r0f.Float64()
	r2, _ := r1f.Float64()

	if pair.TokenA == token0 {
		r = r2 * pair.Amount / (r1 + pair.Amount)
	} else {
		r = r1 * pair.Amount / (r2 + pair.Amount)
	}

	return &r, nil
}

// GetAllPools returns the list of pools from the UniswapV2 Factory contract
func (c *Controller) GetAllPools(pair *models.Pair, client *ethclient.Client) ([][]*models.Path, error) {
	var contracts []models.Pool
	var possibleHops []common.Address
	var possiblePaths [][]*models.Path
	if err := c.DB.Where("token0 = ?", pair.TokenA).Or("token1 = ?", pair.TokenA).Or("token0 = ?", pair.TokenB).Or("token1 = ?", pair.TokenB).Find(&contracts).Error; err != nil {
		return nil, err
	}
	tokenBsymb := contracts[0].Token1_Symbol
	for _, x := range contracts {
		if x.Token0 == pair.TokenA && x.Token1 == pair.TokenB || x.Token0 == pair.TokenB && x.Token1 == pair.TokenA {
			var w []*models.Path
			w = append(w, &models.Path{Address: x.Address, Symbols: []string{x.Token0_Symbol, x.Token1_Symbol}})
			possiblePaths = append(possiblePaths, w)
			continue
		}
		if x.Token0 == pair.TokenA {
			var addr models.Pool
			var w []*models.Path
			w = append(w, &models.Path{Address: x.Address, Symbols: []string{x.Token0_Symbol, x.Token1_Symbol}})
			possibleHops = append(possibleHops, x.Token1)
			if err := c.DB.Where("token0 = ?", x.Token1).Where("token1 = ?", pair.TokenB).Or("token0 = ?", pair.TokenB).Where("token1 = ?", x.Token1).Find(&addr).Error; err != nil {
				return nil, err
			}
			w = append(w, &models.Path{Address: addr.Address, Symbols: []string{x.Token1_Symbol, tokenBsymb}})
			possiblePaths = append(possiblePaths, w)
			continue
		}
		if x.Token1 == pair.TokenA {
			var addr models.Pool
			var w []*models.Path
			w = append(w, &models.Path{Address: x.Address, Symbols: []string{x.Token1_Symbol, x.Token0_Symbol}})
			possibleHops = append(possibleHops, x.Token0)
			if err := c.DB.Where("token0 = ?", x.Token0).Where("token1 = ?", pair.TokenB).Or("token0 = ?", pair.TokenB).Where("token1 = ?", x.Token0).Find(&addr).Error; err != nil {
				return nil, err
			}
			w = append(w, &models.Path{Address: addr.Address, Symbols: []string{x.Token0_Symbol, tokenBsymb}})
			possiblePaths = append(possiblePaths, w)
			continue
		}
	}

	return possiblePaths, nil
}

// CalculateAllRoutes return the list of all possible routes from the UniswapV2 Pair contracts
func (c *Controller) CalculateAllRoutes(pair *models.Pair, client *ethclient.Client) (*models.Quote, error) {
	var rates [][]float64
	var quotes []*models.Quote

	routes, err := c.GetAllPools(pair, client)
	if err != nil {
		return nil, err
	}
	for _, x := range routes {
		newRate := pair.Amount
		var ratesi []float64
		var quote models.Quote
		for _, y := range x {
			var tokenA common.Address
			var tokenB common.Address
			exc, err := pairContract.NewPairCaller(y.Address, client)
			if err != nil {
				logrus.Error("Failed to instantiate the UniswapV2 Pair contract from address")
				return nil, err
			}
			addr1, err := exc.Token0(&bind.CallOpts{})
			if err != nil {
				return nil, err
			}
			addr2, err := exc.Token1(&bind.CallOpts{})
			if err != nil {
				return nil, err
			}
			if addr1 == pair.TokenA || addr2 == pair.TokenA {
				if addr1 == pair.TokenA {
					tokenA = addr2
				} else if addr2 == pair.TokenA {
					tokenA = addr1
				}
				p := &models.Pair{TokenA: pair.TokenA, TokenB: tokenA, Amount: newRate}
				quote.Route = append(quote.Route, p)
				quote.Path = append(quote.Path, y)
				rate, err := c.QueryRate(p, client)
				if err != nil {
					return nil, err
				}
				newRate = *rate
				ratesi = append(ratesi, *rate)
				if err != nil {
					return nil, err
				}
				continue
			}
			if addr1 == pair.TokenB || addr2 == pair.TokenB {
				if addr1 == pair.TokenB {
					tokenB = addr2
				} else if addr2 == pair.TokenB {
					tokenB = addr1
				}
				p := &models.Pair{TokenA: tokenB, TokenB: pair.TokenB, Amount: newRate}
				quote.Route = append(quote.Route, p)
				quote.Path = append(quote.Path, y)
				rate, err := c.QueryRate(p, client)
				if err != nil {
					return nil, err
				}
				newRate = *rate
				ratesi = append(ratesi, *rate)
				if err != nil {

					return nil, err
				}
				continue
			}
		}
		rates = append(rates, ratesi)
		quote.Rate = newRate
		quotes = append(quotes, &quote)
	}

	return GetBestRoute(quotes, rates), nil
}

func GetBestRoute(quotes []*models.Quote, rates [][]float64) *models.Quote {
	var best float64
	var q int
	for i, x := range quotes {
		log.Printf("Quote %d: %v | Rates: %v", i, x.Rate, rates[i])
		if x.Rate > best {
			best = x.Rate
			q = i
		}
		for _, y := range x.Path {
			log.Printf("\tPath Address: %v | Pair: %v", y.Address, y.Symbols)
		}
		log.Printf("\n")
	}
	return quotes[q]
}

func WeiToEther(wei *big.Int) *big.Float {
	f := new(big.Float)
	f.SetPrec(236)
	f.SetMode(big.ToNearestEven)
	fWei := new(big.Float)
	fWei.SetPrec(236)
	fWei.SetMode(big.ToNearestEven)
	return f.Quo(fWei.SetInt(wei), big.NewFloat(params.Ether))
}
