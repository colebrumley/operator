package main

import (
	"flag"

	"fmt"

	machinery "github.com/RichardKnop/machinery/v1"
	machcfg "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	log "github.com/Sirupsen/logrus"
	"github.com/colebrumley/operator"
	"github.com/olebedev/config"
)

var (
	cfg    *config.Config
	server *machinery.Server
)

func init() {
	var err error
	// Load config from local yml/json file, env, or flags
	cfg, err = operator.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	flag.Usage = func() {
		fmt.Println("Description:    Operator Client - Submit tasks to Operator via the CLI")
		fmt.Println("Version:        " + operator.Version)
		fmt.Println("Built on:       " + operator.BuildTime)
		fmt.Println("Options:")
		m, _ := cfg.Map("")
		operator.PrettyPrintFlagMap(m, []string{})
	}
	cfg.Env().Flag()
}

func main() {
	operator.ConfigLogging(cfg)

	c, _ := config.RenderJson(cfg.Root)
	log.Debugln("Loaded config: ", c)

	// Using UString here since the defaults always have these covered
	_, err := machinery.NewServer(&machcfg.Config{
		Broker:        cfg.UString("broker.url"),
		ResultBackend: cfg.UString("results.url"),
		Exchange:      cfg.UString("broker.exchange.name"),
		ExchangeType:  cfg.UString("broker.exchange.type"),
		DefaultQueue:  cfg.UString("broker.queue"),
		BindingKey:    cfg.UString("broker.bindingkey"),
	})
	errors.Fail(err, "Could not initialize server")
}
