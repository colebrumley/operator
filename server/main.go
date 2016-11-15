package main

import (
	"flag"
	"os"

	"fmt"

	machinery "github.com/RichardKnop/machinery/v1"
	machcfg "github.com/RichardKnop/machinery/v1/config"
	"github.com/RichardKnop/machinery/v1/errors"
	log "github.com/Sirupsen/logrus"
	"github.com/colebrumley/operator"
	oapi "github.com/colebrumley/operator/api"
	"github.com/colebrumley/operator/tasks"
	"github.com/olebedev/config"
)

var (
	cfg    *config.Config
	server *machinery.Server
	worker *machinery.Worker
)

func init() {
	var err error
	// Load config from local yml/json file, env, or flags
	cfg, err = operator.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	flag.Usage = func() {
		fmt.Println("Description:    Operator Server - Task execution daemon for Operator")
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
	server, err := machinery.NewServer(&machcfg.Config{
		Broker:        cfg.UString("broker.url"),
		ResultBackend: cfg.UString("results.url"),
		Exchange:      cfg.UString("broker.exchange.name"),
		ExchangeType:  cfg.UString("broker.exchange.type"),
		DefaultQueue:  cfg.UString("broker.queue"),
		BindingKey:    cfg.UString("broker.bindingkey"),
	})
	errors.Fail(err, "Could not initialize server")

	// Register tasks
	list, _ := cfg.List("tasks.enabled")
	for _, name := range list {
		n := name.(string)
		log.Info("Registering task handler '" + n + "'")
		server.RegisterTask(n, tasks.TaskList[n].Fn)
	}

	if cfg.UBool("api.enabled", false) {
		// Start the API
		api := oapi.OperatorAPI{
			ListenAddress: cfg.UString("api.addr", ":8080"),
			UseBasicAuth:  cfg.UBool("api.basicauth", false),
			Password:      cfg.UString("api.password"),
			UseTLS:        cfg.UBool("api.usetls", false),
			Version:       1,
		}

		log.Infoln("Launching API")
		go api.Start(server)
	} else {
		log.Infoln("Skipping API")
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker = server.NewWorker(getWorkerName())
	if err = worker.Launch(); err != nil {
		errors.Fail(err, "Could not launch worker")
	}
}

func getWorkerName() string {
	hostname, _ := os.Hostname()
	pid := os.Getpid()
	str := fmt.Sprintf("%s.%v.opworker", hostname, pid)
	return str
}
