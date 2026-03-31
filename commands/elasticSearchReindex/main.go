package main

import (
	"flag"
	"fmt"
	"time"
)

func main() {
	indesName := flag.String("index", "student", "target ElasticSearch index")
	batchSize := flag.Int("batch", 500, "batch size for reindex")
	flag.Parse()

	start := time.Now()
	fmt.Printf("[es-reindex] start index=%s batch=%d\n", *indesName, *batchSize)
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("[es-reindex] done is %s\n", time.Since(start).String())

}
