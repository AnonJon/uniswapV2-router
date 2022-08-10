package cmd

import (
	"fmt"
	"log"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jongregis/uniswapV2_router/models"
	"github.com/spf13/cobra"
)

// quoteCmd represents the quote command
var quoteCmd = &cobra.Command{
	Use:   "quote",
	Short: "Get the best quote for a trade",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running quote...")
		var reply *models.Quote

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
		err = client.RPCclient.Call("Handler.GetQuote", pair, &reply)
		if err != nil {
			log.Fatal("error", err)
		}

		for _, y := range reply.Path {
			fmt.Printf("\n\t[V2] 100.00%% = %v -- [%v] --> %v\n", y.Symbols[0], y.Address, y.Symbols[1])
		}
		fmt.Printf("\n\tReturn Rate: %v %v\n", reply.Rate, reply.Path[len(reply.Path)-1].Symbols[1])
		fmt.Println("")

	},
}

func init() {
	quoteCmd.Flags().StringVarP(&TokenIn, "tokenIn", "i", "", "token selling (required)")
	quoteCmd.Flags().StringVarP(&TokenOut, "tokenOut", "o", "", "token receiving (required)")
	quoteCmd.Flags().StringVarP(&Amount, "amount", "a", "", "token amount selling (required)")
	quoteCmd.MarkFlagRequired("tokenIn")
	quoteCmd.MarkFlagRequired("tokenOut")
	quoteCmd.MarkFlagRequired("amount")
	rootCmd.AddCommand(quoteCmd)
}
