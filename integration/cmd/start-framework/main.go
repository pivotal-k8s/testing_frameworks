package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"sigs.k8s.io/testing_frameworks/integration"
)

func main() {
	cp := &integration.ControlPlane{}
	if err := cp.Start(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer cp.Stop()

	fmt.Printf(`API Server running. Connect with:

    http:  kubectl --server %s

Press Ctrl+C to exit
`, cp.APIURL())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
