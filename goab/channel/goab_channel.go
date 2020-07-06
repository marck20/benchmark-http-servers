package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var (
	uFlag     = flag.String("u", ":8080", "Target URL for a benchmark of HTTP GET requests. Using a preceding 'http://' is mandatory. Example: http://127.0.0.1:8080 ")
	nFlag     = flag.Int("n", 1, "Number of requests to perform for the benchmarking session. On default performs a single request.")
	cFlag     = flag.Int("c", 1, "Number of concurrent requests. On default performs a single concurrent request.")
	rFlag     = flag.Int("r", 0, "If greater than 0, uses Client.DoRedirect() function from httpserver. Follows up to 'r' redirections. Default is 0.")
	aliveFlag = flag.Bool("k", false, "If true, KeepAlive header is set on HTTP requests. This feature reuses the same HTTP session to perform multiple requests.")

	results chan int64

	client *fasthttp.Client
)

// Sends HTTP request to urlFlag
func request() {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.SetRequestURI(*uFlag)

	if *aliveFlag == true {
		req.Header.Set("Connection", "keep-alive")
	} else {
		req.Header.Set("Connection", "close")
	}

	startTime := time.Now()
	err := client.DoRedirects(req, resp, *rFlag)
	elapsed := time.Since(startTime).Microseconds()

	if err == nil {
		results <- elapsed
	}
}

// Worker function is called as a goroutine. First reserves a slot from
// slots channel and performs a single HTTP GET request. If Success
// increments 'success' counter inside a mutual exclusion segment.
// @param slots chan struct{} - Channel to control the goroutines concurrency.
// @param wg *sync.WaitGroup
func worker(slots chan struct{}, wg *sync.WaitGroup) {
	//Decrement workergroup counter at the end of the routine
	defer wg.Done()

	//Reseve a slot for the current worker.
	slots <- struct{}{}

	//Send HTTP GET request to url
	request()

	//Release a slot for a future worker
	<-slots
}

func main() {
	flag.Parse()

	//Set GOMAXPROCS to total cpu threads
	goMaxProcs := os.Getenv("GOMAXPROCS")
	if goMaxProcs == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	//Display info
	fmt.Println("[+] Starting benchmark")
	fmt.Println("[i] URL: ", *uFlag)
	fmt.Println("[i] CPU cores: ", runtime.NumCPU(), "  Number of transacions: ", *nFlag, "  Concurrent workers: ", *cFlag, "  KeepAlive: ", *aliveFlag)

	//Buffer channel of empty structs which represents
	//worker slots. Each worker will consume an available
	//slot in order to perform its task. Empty structs
	//do not allocate memory, its size is 0.
	slots := make(chan struct{}, *cFlag)

	//Buffer channel
	results = make(chan int64, *cFlag)

	client = &fasthttp.Client{}
	var wg sync.WaitGroup
	startTime := time.Now()

	for i := 0; i < *nFlag; i++ {
		wg.Add(1) //Increment waitgrupo counter
		go worker(slots, &wg)
	}

	sumLatency := int64(0)
	minLatency := int64(0)
	maxLatency := int64(0)
	first := true
	success := 0
	for r := 0; r < *nFlag; r++ {
		latency := <-results
		if latency < minLatency || first {
			minLatency = latency
		}
		if latency > maxLatency || first {
			maxLatency = latency
		}
		if first {
			first = false
		}
		sumLatency += latency
		success++
	}

	wg.Wait() //Block execution until the end of all go routines.

	endTime := time.Now()
	diff := endTime.Sub(startTime)

	close(slots)
	close(results)

	avgLatencyMs := float64(sumLatency/int64(*nFlag)) / float64(time.Millisecond)
	minLatencyMs := float64(minLatency) / float64(time.Millisecond)
	maxLatencyMs := float64(maxLatency) / float64(time.Millisecond)

	fmt.Println("[+] Benchmark done")
	fmt.Println("[*] Requests: ", success, "/", *nFlag, "  Error: ", 100*(*nFlag-success) / *nFlag, "%")
	fmt.Printf("[*] TPS: %.03f reqs/second\n", float64(*nFlag)/diff.Seconds())
	fmt.Printf("[*] Time: %0.06f seconds\n", diff.Seconds())
	fmt.Printf("[*] Average latency: %0.06f ms   Min latency: %0.06f ms   Max latency: %0.06f ms\n", avgLatencyMs, minLatencyMs, maxLatencyMs)
}
