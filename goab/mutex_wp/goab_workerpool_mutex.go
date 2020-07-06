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

	mu         sync.Mutex // guards number of success requests and latency statistics
	success    int
	sumLatency int64
	minLatency int64
	maxLatency int64

	client *fasthttp.Client
)

// Sends HTTP request to urlFlag and mesure its latency
func request() {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()

	req.SetRequestURI(*uFlag)

	if *aliveFlag == true {
		req.Header.Set("Connection", "keep-alive")
	} else {
		req.Header.Set("Connection", "close")
		//req.SetConnectionClose()
	}

	startTime := time.Now()
	err := client.DoRedirects(req, resp, *rFlag)
	latency := time.Since(startTime).Microseconds()

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)

	if err == nil {
		mu.Lock() //Starts mutual exclusion segment
		success++
		sumLatency += latency
		if latency < minLatency || minLatency == int64(-1) {
			minLatency = latency
		}
		if latency > maxLatency || maxLatency == int64(-1) {
			maxLatency = latency
		}
		mu.Unlock() //Ends mutual exclusion segment
	} else {
		fmt.Println("request error: ", err)
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
	fmt.Println("[+] Starting Mutex Workpool Benchmark")
	fmt.Println("[i] URL: ", *uFlag)
	fmt.Println("[i] CPU cores: ", runtime.NumCPU(), "  Number of transacions: ", *nFlag, "  Concurrent workers: ", *cFlag, "  KeepAlive: ", *aliveFlag)

	//Init vars
	minLatency = -1
	maxLatency = -1
	client = &fasthttp.Client{}
	jobs = make(chan struct{}, *nFlag)
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

	wg.Wait() //Block execution until the end of all go routines.
	endTime := time.Now()
	diff := endTime.Sub(startTime)

	avgLatencyMs := float64(sumLatency/int64(*nFlag)) / float64(time.Millisecond)
	minLatencyMs := float64(minLatency) / float64(time.Millisecond)
	maxLatencyMs := float64(maxLatency) / float64(time.Millisecond)

	fmt.Println("[+] Benchmark done")
	fmt.Println("[*] Requests: ", success, "/", *nFlag, "  Error: ", 100*(*nFlag-success) / *nFlag, "%")
	fmt.Printf("[*] TPS: %.03f reqs/second\n", float64(*nFlag)/diff.Seconds())
	fmt.Printf("[*] Time: %0.06f seconds\n", diff.Seconds())
	fmt.Printf("[*] Average latency: %0.06f ms   Min latency: %0.06f ms   Max latency: %0.06f ms\n", avgLatencyMs, minLatencyMs, maxLatencyMs)
}
