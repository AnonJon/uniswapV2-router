package cmd

import (
	"fmt"
	"net/rpc"
	"os"

	"github.com/spf13/cobra"
)

var TokenIn string
var TokenOut string
var Amount string

var rootCmd = &cobra.Command{
	Use:   "uniswapV2_calculator",
	Short: "Find the best path to trade between two tokens",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

type Client struct {
	RPCclient *rpc.Client
}

func NewClient() (*Client, error) {
	client, err := rpc.DialHTTP("tcp", fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT")))
	if err != nil {
		return nil, err
	}
	return &Client{
		RPCclient: client,
	}, nil
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.Flags().StringVarP(&TokenIn, "tokenIn", "i", "", "token selling (required)")
	rootCmd.Flags().StringVarP(&TokenOut, "tokenOut", "o", "", "token receiving (required)")
	rootCmd.Flags().StringVarP(&Amount, "amount", "a", "", "token amount selling (required)")
}
