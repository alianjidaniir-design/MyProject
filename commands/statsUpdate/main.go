package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {

	period := flag.String("period", "daily", "status update period: hourly[daily|weekly]")
	dryRun := flag.Bool("dry-run", false, "run without writing updates")
	flag.Parse()

	start := time.Now()
	fmt.Printf("[status-update] start period=%s dryRun=%t\n", *period, *dryRun)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("[status-update] done is %s\n", time.Since(start).String())

}
