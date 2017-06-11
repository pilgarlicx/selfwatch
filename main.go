package main

import (
	"flag"
	"log"

	"github.com/leafo/selfwatch/selfwatch"
)

var (
	configFname string
	debugOutput bool
)

func init() {
	flag.StringVar(&configFname, "config", selfwatch.DefaultConfigFname, "Path to json config file")
	flag.BoolVar(&debugOutput, "dump", false, "Print extra debug information")
}

func main() {
	flag.Parse()
	config := selfwatch.LoadConfig(configFname)

	storage, err := selfwatch.NewWatchStorage(config.DbName)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !storage.SchemaExists() {
		storage.CreateSchema()
	}

	command := flag.Arg(0)
	if command == "" {
		command = "start"
	}

	switch command {
	case "summary":
		err := storage.DailyCounts(7)
		if err != nil {
			log.Fatal(err.Error())
		}

	case "start":
		recorder := selfwatch.NewRecorder()
		storage.BindRecorder(recorder, config.SyncDelay)

		if config.RemoteUrl != "" {
			remote := selfwatch.RemoteSync{
				Url:     config.RemoteUrl,
				Storage: storage,
			}

			if config.RemoteFlushDelay > 0 {
				remote.FlushEvery(config.RemoteFlushDelay)
			}
		}

		recorder.Bind()
	default:
		log.Fatal("Unknown command:", command)
	}

}
