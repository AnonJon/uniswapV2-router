package server

import (
	"log"
	"net/rpc"
	"os"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jongregis/uniswapV2_router/controllers"
	"github.com/jongregis/uniswapV2_router/database"
	"github.com/jongregis/uniswapV2_router/models"
	"github.com/jongregis/uniswapV2_router/services"
	"github.com/jongregis/uniswapV2_router/store"
	"gorm.io/gorm"
)

type Handler struct {
	Client     *ethclient.Client
	Controller *controllers.Controller
	DB         *gorm.DB
}

func NewHandler() *Handler {
	client, err := store.NewETHClient(os.Getenv("ETH_HTTPS"))
	if err != nil {
		log.Fatalf(err.Error())
	}
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatalf(err.Error())
	}

	h := &Handler{
		Client:     client,
		Controller: controllers.NewController(db),
		DB:         db}
	if os.Getenv("BACKFILL") == "true" {
		if err := services.Backill(db, client); err != nil {
			log.Fatalf(err.Error())
		}
	}
	if err = rpc.Register(h); err != nil {
		panic(err)
	}
	rpc.HandleHTTP()
	return h
}

func (rh *Handler) GetRate(payload *models.Pair, reply *float64) error {
	rate, err := rh.Controller.QueryRate(payload, rh.Client)
	if err != nil {
		return err
	}

	*reply = *rate
	return nil
}

func (rh *Handler) GetQuote(payload *models.Pair, quote *models.Quote) error {
	bestQuote, err := rh.Controller.CalculateAllRoutes(payload, rh.Client)
	if err != nil {
		return err
	}

	*quote = *bestQuote
	return nil
}
