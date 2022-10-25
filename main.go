package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"time"
)

var REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
	Name: "go_app_request_count",
	Help: "Total App HTTP Requests Count",
})

var REQUEST_INPROGRESS = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "go_app_requests_in_progress",
	Help: "Total App HTTP Requests Count in Progress",
})

var REQUEST_RESPOND_TIME = promauto.NewSummaryVec(prometheus.SummaryOpts{
	Name: "go_app_response_latency_seconds",
	Help: "response latency in seconds",
}, []string{"path"})

func routeMiddleware(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start_time := time.Now()
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		next.ServeHTTP(w, r)
		time_taken := time.Since(start_time)
		REQUEST_RESPOND_TIME.WithLabelValues(path).Observe(time_taken.Seconds())

	})

}
func main() {
	// Start the application
	startMyApp()
}

func startMyApp() {
	router := mux.NewRouter()
	router.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		REQUEST_INPROGRESS.Inc()
		//vars := mux.Vars(r)
		//name := vars["name"]
		greetings := fmt.Sprintf("Hello %s :)", "ronal")
		time.Sleep(5 * time.Second)
		rw.Write([]byte(greetings))
		REQUEST_INPROGRESS.Dec()
		REQUEST_COUNT.Inc()
	}).Methods("GET")
	router.Use(routeMiddleware)
	log.Println("Started the application server...")
	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe(":8080", router)
}
