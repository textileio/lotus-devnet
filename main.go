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

	err := logging.SetLogLevel("main", "INFO")
	if err != nil {
		panic(err)
	}

	pflag.Int("numminers", 1, "Number of miners in devnet")
	pflag.Int("speed", ldevnet.DefaultDurationMs, "Chain speed block creation in ms")
	pflag.Bool("bigsectors", true, "Use big sectors")
	pflag.String("ipfsaddr", "", "IPFS multiaddr to make Lotus use an IPFS node")
	pflag.Parse()

	config.SetEnvPrefix("TEXLOTUSDEVNET")
	config.AutomaticEnv()
	config.BindPFlags(pflag.CommandLine)

	speed := config.GetInt("speed")
	numMiners := config.GetInt("numminers")
	bigSectors := config.GetBool("bigsectors")
	ipfsAddr := config.GetString("ipfsaddr")

	log.Infof("Starting devnet with speed = %v, numminers = %v, bigsectors = %v, ipfsaddr = %v", speed, numMiners, bigSectors, ipfsAddr)

	_, err = ldevnet.New(numMiners, time.Millisecond*time.Duration(speed), bigSectors, ipfsAddr)
	if err != nil {
		panic(err)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	log.Info("Closing...")
}
