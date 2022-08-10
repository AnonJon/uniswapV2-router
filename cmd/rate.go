package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jongregis/uniswapV2_router/models"
	"github.com/spf13/cobra"
)

var rateCmd = &cobra.Command{
	Use:   "rate",
	Short: "Get rate from a pair",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("finding rate...")
		var reply *float64

		tokenIn, _ := cmd.Flags().GetString("tokenIn")
		tokenOut, _ := cmd.Flags().GetString("tokenOut")
		amount, _ := cmd.Flags().GetString("amount")
		s, err := strconv.ParseFloat(amount, 32)
		if err != nil {
			fmt.Println(err)
		}
		client, err := NewClient()
		if err != nil {
			log.Fatal("error", err)
		}
		pair := &models.Pair{TokenA: common.HexToAddress(tokenIn), TokenB: common.HexToAddress(tokenOut), Amount: s}
		err = client.RPCclient.Call("Handler.GetRate", pair, &reply)
		if err != nil {
			log.Fatal("error", err)
		}
		fmt.Printf("\n\tReturn Rate: %v\n", *reply)
		fmt.Println("")
	},
}

func init() {
	rateCmd.Flags().StringVarP(&TokenIn, "tokenIn", "i", "", "token selling (required)")
	rateCmd.Flags().StringVarP(&TokenOut, "tokenOut", "o", "", "token receiving (required)")
	rateCmd.Flags().StringVarP(&Amount, "amount", "a", "", "token amount selling (required)")
	rateCmd.MarkFlagRequired("tokenIn")
	rateCmd.MarkFlagRequired("tokenOut")
	rateCmd.MarkFlagRequired("amount")
	rootCmd.AddCommand(rateCmd)
}
