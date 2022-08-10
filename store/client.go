package store

import (
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
)

func NewETHClient(network string) (*ethclient.Client, error) {
	if len(network) == 0 {
		logrus.Debug("No ETH URL given for network, ignoring")
		return nil, nil
	}
	c, err := ethclient.Dial(network)
	if err != nil {
		log.Fatal(err)
	}
	logrus.Info("Ethereum Client Connected")

	return c, nil
}
