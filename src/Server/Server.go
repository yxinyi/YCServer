package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
)

var g_server_name = flag.String("ServerName", "", "servername")

func main() {
	flag.Parse()

	if *g_server_name == "single" {
		go func() {
			http.ListenAndServe("0.0.0.0:6060", nil)
		}()
		MainLoop()
	}

}
