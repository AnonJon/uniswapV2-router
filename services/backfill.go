package services

import (
	"errors"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jongregis/uniswapV2_router/contracts/erc20"
	"github.com/jongregis/uniswapV2_router/contracts/factory"
	pairContract "github.com/jongregis/uniswapV2_router/contracts/pair"
	"github.com/jongregis/uniswapV2_router/models"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const (
	FACTORY     = "0x5C69bEe701ef814a2B6a3EDD4B1652CB9cc5aA6f"
	NIL_ADDRESS = "0x0000000000000000000000000000000000000000"
)

func Backill(db *gorm.DB, client *ethclient.Client) error {
	logrus.Info("Backfilling pools...")
	var pool models.Pool
	db.AutoMigrate(&models.Pool{})

	exc, err := factory.NewFactoryCaller(common.HexToAddress(FACTORY), client)
	if err != nil {
		logrus.Error("Failed to instantiate the UniswapV2 Factory contractfrom address")
		return err
	}

	num, err := exc.AllPairsLength(&bind.CallOpts{})
	if err != nil {
		return err
	}
	if err := db.Order("pool_number desc").First(&pool).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
		} else {
			return err
		}
	}
	if pool.PoolNumber != 0 {
		pool.PoolNumber = pool.PoolNumber + 1
	}
	for i := pool.PoolNumber; i < int(num.Int64())+1; i++ {
		addr, err := exc.AllPairs(&bind.CallOpts{}, big.NewInt(int64(i)))
		if err != nil {
			return err
		}
		exc, err := pairContract.NewPairCaller(addr, client)
		if err != nil {
			logrus.Error("Failed to instantiate the UniswapV2 Pair contract from address")
			return err
		}

		t1, err := exc.Token0(&bind.CallOpts{})
		if err != nil {
			return err
		}
		t2, err := exc.Token1(&bind.CallOpts{})
		if err != nil {
			return err
		}

		x0, err := erc20.NewErc20Caller(t1, client)
		if err != nil {
			logrus.Error("Failed to instantiate the ERC20 contract from address")
			return err
		}

		s0, err := x0.Symbol(&bind.CallOpts{})
		if err != nil {
			log.Println("error getting symbol0", err)
			continue
		}

		x1, err := erc20.NewErc20Caller(t2, client)
		if err != nil {
			logrus.Error("Failed to instantiate the ERC20 contract from address")
			return err
		}
		s1, err := x1.Symbol(&bind.CallOpts{})
		if err != nil {
			log.Println("error getting symbol1", err)
			continue
		}
		if err := db.Create(&models.Pool{
			Id:            uuid.NewV4(),
			Address:       addr,
			Token0:        t1,
			Token1:        t2,
			Token0_Symbol: s0,
			Token1_Symbol: s1,
			PoolNumber:    i,
		}).Error; err != nil {
			return err
		}
		log.Printf("saved pool: [%s, %s]", s0, s1)
	}
	logrus.Info("Finished backfilling pools")
	return nil
}
