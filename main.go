package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

var URL string = os.Getenv("URL")
var FREQUENCY string = os.Getenv("FREQUENCY")

var UP int64 = 0

func incrementUP() {
	UP++
}

var DOWN int64 = 0

func incrementDown() {
	DOWN++
}

var LATENCY time.Duration = 0

func setLatency(l time.Duration) {
	LATENCY = l
}

var LASTSTATUSCODE string = ""

func checkURL() (int, time.Duration, error) {
	start := time.Now()
	client := &http.Client{Timeout: 5 * time.Second} // Set a timeout of 5 seconds

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	elapsed := time.Since(start)
	return resp.StatusCode, elapsed, nil
}

func testEndpoint() {
	max := 48000
	c := 0

	for max > c {
		statusCode, latency, _ := checkURL()
		if statusCode == 200 {
			incrementUP()
		} else {
			incrementDown()
		}
		setLatency(latency)

		f, _ := strconv.Atoi(FREQUENCY)
		time.Sleep(time.Duration(f) * time.Second)
	}
}

func respondMetrics(w http.ResponseWriter, r *http.Request) {

	aval := UP / (UP + DOWN) * 100

	var avLabel string = `{job="gac-checker", up="` + fmt.Sprint(UP) + `", down="` + fmt.Sprint(DOWN) + `"}`
	var data string = "# HELP availability_metric The availability of the url checker" + "\n"
	data += "availability_metric" + avLabel + " " + fmt.Sprint(aval) + "\n"

	data += "# HELP latency_metric the latency of the url checker" + "\n"
	data += "latency_metric" + `{job="gac-checker"}` + " " + fmt.Sprint(LATENCY.Milliseconds()) + "\n"

	w.Write([]byte(data))
}

func main() {
	if FREQUENCY == "" {
		FREQUENCY = "60"
	}
	fmt.Println("Running GAC")
	fmt.Println("URL: " + URL)
	fmt.Println("FREQUENCY: " + FREQUENCY)

	go testEndpoint()
	http.HandleFunc("/metrics", respondMetrics)

	http.ListenAndServe(":8080", nil)
}
