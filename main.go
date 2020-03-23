package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"github.com/textileio/lotus-devnet/ldevnet"
)

var (
	log    = logging.Logger("main")
	config = viper.New()
)

func main() {
	pflag.Int("numminers", 1, "Number of miners in devnet")
	pflag.Int("speed", 50, "Chain speed block creation in ms")
	pflag.Parse()

	config.SetEnvPrefix("TEXLOTUSDEVNET")
	config.AutomaticEnv()
	config.BindPFlags(pflag.CommandLine)

	speed := config.GetInt("speed")
	numMiners := config.GetInt("numminers")

	_, err := ldevnet.New(numMiners, time.Millisecond*time.Duration(speed))
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	log.Info("Closing...")
}
