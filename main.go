package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"

	"github.com/drtoful/gifttt/gifttt"
)

func main() {
	var (
		dbPath     = flag.String("db", "gifttt.db", "path to the database store")
		rulePath   = flag.String("ruledir", "./", "path to rule files")
		cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
		apiBind    = flag.String("ip", "", "ip to bind the api server to")
		apiPort    = flag.String("port", "4200", "port for api server")
	)
	flag.Parse()

	// activate cpu profiling
	if *cpuprofile != "" {
		log.Println("main: Starting CPU profiling '%s'", *cpuprofile)
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if err := gifttt.StoreInit(*dbPath); err != nil {
		log.Fatal(err)
	}

	// start the servers
	rm := gifttt.NewRuleManager(*rulePath)
	api := gifttt.NewAPIServer(*apiBind, *apiPort)

	go rm.Run()
	go api.Run()

	// all listeners are started in the background as
	// gofunc's so we wait here for an interupt signal
	// to stop the service gracefully
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			pprof.StopCPUProfile() // explicitly stop profiling
			log.Fatalf("main: Signal (%d) received, stopping\n", s)
		}
	}
}
