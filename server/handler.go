package server

import (
	"log"
	"net/rpc"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jongregis/uniswapV2_router/controllers"
	"github.com/jongregis/uniswapV2_router/models"
	"github.com/jongregis/uniswapV2_router/store"
)

type Handler struct {
	Client *ethclient.Client
}

func NewHandler() *Handler {
	client, err := store.NewETHClient(os.Getenv("ETH_HTTPS"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	h := &Handler{Client: client}
	if err = rpc.Register(h); err != nil {
		panic(err)
	}
	rpc.HandleHTTP()
	return h
}

func (rh *Handler) GetRate(payload *models.Pair, reply *float64) error {
	rate, err := controllers.QueryRate(payload, rh.Client)
	if err != nil {
		return err
	}

	*reply = *rate
	return nil
}

func (rh *Handler) GetQuote(payload *models.Pair, quote *models.Quote) error {
	bestQuote, err := controllers.CalculateAllRoutes(payload, rh.Client)
	if err != nil {
		return err
	}

	*quote = *bestQuote
	return nil
}
