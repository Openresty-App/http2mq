package main

import (
	"flag"

	"os"

	"fmt"

	"github.com/xuyz/http2mq/app"
)

func main() {
	conf := flag.String("conf,c", "conf/http2mq.yaml", "http2mq configuration file")
	flag.Parse()

	a, err := app.NewApp(*conf)
	if err != nil {
		fmt.Printf("create app error %s", err.Error())
		os.Exit(2)
	}
	defer a.Close()

	a.Run()
}
