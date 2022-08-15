package graphql

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/jongregis/uniswapV2_router/graphql/models"
)

type SubGraph struct {
	Client *graphql.Client
}

func NewSubGraph() *SubGraph {
	return &SubGraph{
		Client: graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", nil),
	}
}

func (sg *SubGraph) GetPairs(address string) ([]models.Pair, error) {
	var pairs []models.Pair
	client := graphql.NewClient("https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v2", nil)
	for x := 0; x < 2; x++ {
		query := fmt.Sprintf(`query{pairs(first: 1000, where: {token%v: "%s", reserveETH_gte: 1}){id token0{id symbol} token1{id symbol}}}`, x, strings.ToLower(address))
		res := struct {
			Pairs []models.Pair
		}{}

		err := client.Exec(context.Background(), query, &res, nil)
		if err != nil {
			log.Fatal(err)
		}
		pairs = append(pairs, res.Pairs...)
	}

	return pairs, nil
}
