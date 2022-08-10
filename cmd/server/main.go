package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"github.com/jongregis/uniswapV2_router/server"
	"github.com/sirupsen/logrus"
)

func main() {

	server.NewHandler()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", os.Getenv("PORT")))
	if err != nil {
		logrus.Fatal("listen error: ", err)
	}
	logrus.Info("running rpc server")
	http.Serve(listener, nil)
}
