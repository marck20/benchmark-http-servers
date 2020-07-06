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
	uFlag     = flag.String("u", ":80", "Target URL for a benchmark of HTTP GET requests. Using a preceding 'http://' is mandatory. Example: http://127.0.0.1:5000 ")
	nFlag     = flag.Int("n", 1, "Number of requests to perform for the benchmarking session. On default performs a single request.")
	cFlag     = flag.Int("c", 1, "Number of concurrent requests. On default performs a single concurrent request.")
	rFlag     = flag.Int("r", 0, "If greater than 0, uses Client.DoRedirect() function from httpserver. Follows up to 'r' redirections. Default is 0.")
	aliveFlag = flag.Bool("k", false, "If true, KeepAlive header is set on HTTP requests. This feature reuses the same HTTP session to perform multiple requests.")

	jobs chan struct{}

	latencies chan int64

	client *fasthttp.Client
)

// Sends HTTP request to urlFlag and mesure its latency
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
		latencies <- elapsed
	}
}

func worker(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	for _ = range jobs {
		request()
	}
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

	//Init vars
	client = &fasthttp.Client{}
	jobs = make(chan struct{}, *nFlag)
	latencies = make(chan int64, *cFlag)
	var wg sync.WaitGroup

	//Creating 'cFlag' workers
	for w := 0; w < *cFlag; w++ {
		wg.Add(1) //Increment waitgrupo counter
		go worker(w, &wg)
	}

	startTime := time.Now()

	// Creating 'nFlag' jobs. Each one which will
	// perform a single HTTP request to 'uFlag'.
	for j := 0; j < *nFlag; j++ {
		jobs <- struct{}{}
	}
	close(jobs)

	sumLatency := int64(0)
	minLatency := int64(-1)
	maxLatency := int64(-1)
	success := 0
	for i := 0; i < *nFlag; i++ {
		latency := <-latencies
		if minLatency == -1 || latency < minLatency {
			minLatency = latency
		}
		if maxLatency == -1 || latency > maxLatency {
			maxLatency = latency
		}
		sumLatency += latency
		success++
	}
	if minLatency == -1 {
		minLatency = 0
	}
	if maxLatency == -1 {
		maxLatency = 0
	}

	wg.Wait() //Block execution until the end of all go routines.
	endTime := time.Now()
	diff := endTime.Sub(startTime)

	close(latencies)

	avgLatencyMs := float64(sumLatency/int64(*nFlag)) / float64(time.Millisecond)
	minLatencyMs := float64(minLatency) / float64(time.Millisecond)
	maxLatencyMs := float64(maxLatency) / float64(time.Millisecond)

	fmt.Println("[+] Benchmark done")
	fmt.Println("[*] Requests: ", success, "/", *nFlag, "  Error: ", 100*(*nFlag-success) / *nFlag, "%")
	fmt.Printf("[*] TPS: %.03f reqs/second\n", float64(*nFlag)/diff.Seconds())
	fmt.Printf("[*] Time: %0.06f seconds\n", diff.Seconds())
	fmt.Printf("[*] Average latency: %0.06f ms   Min latency: %0.06f ms   Max latency: %0.06f ms\n", avgLatencyMs, minLatencyMs, maxLatencyMs)
}
