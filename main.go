package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var PowerballCounter prometheus.Counter
var PodGauge prometheus.Gauge

func main() {
	// set up endpoints
	http.HandleFunc("/powerball", HandlePowerball)
	http.Handle("/metrics", promhttp.Handler())

	// set up endpoint counter
	counterOpts := prometheus.CounterOpts{
		Subsystem: "a_kasten_assessment",
		Name:      "powerball_counter",
		Help:      "Count of endpoint Powerball",
	}
	PowerballCounter = prometheus.NewCounter(counterOpts)
	if err := prometheus.Register(PowerballCounter); err != nil {
		panic(fmt.Sprintf("Unable to register powerball counter prometheus: %s", err))
	}

	// set up k8s gauge
	gaugeOpts := prometheus.GaugeOpts{
		Subsystem: "a_kasten_assessment",
		Name:      "k8s_pod_count",
		Help:      "Kubernetes cluster pod count",
	}
	PodGauge = prometheus.NewGauge(gaugeOpts)
	if err := prometheus.Register(PodGauge); err != nil {
		panic(fmt.Sprintf("Unable to register pod gauge prometheus: %s", err))
	}

	go SetupK8sPod(PodGauge)

	// starting server
	fmt.Println("Serving it up on port 8080...")
	http.ListenAndServe(":8080", nil)
}

func HandlePowerball(w http.ResponseWriter, r *http.Request) {
	numbers := make(map[int32]struct{})
	for {
		var randomInt int32
		for {
			randomInt = rand.Int31n(69)
			if _, ok := numbers[randomInt]; !ok {
				break
			}
		}
		numbers[randomInt] = struct{}{}
		if len(numbers) == 5 {
			break
		}
	}
	numberKeys := make([]int, 0, len(numbers))
	for k := range numbers {
		numberKeys = append(numberKeys, int(k))
	}
	sort.Ints(numberKeys)

	powerball := rand.Int31n(69)
	fmt.Fprintf(w, "Your random numbers: %d %d %d %d %d; Powerball number: %d", numberKeys[0], numberKeys[1], numberKeys[2], numberKeys[3], numberKeys[4], powerball)
	PowerballCounter.Inc()
}

func SetupK8sPod(podGauge prometheus.Gauge) {
	config, err := clientcmd.BuildConfigFromFlags("", "./kube_config") // normally get this from {home}/.kube/config but for this assessment have the file locally
	if err != nil {
		panic(fmt.Sprintf("Unable to build k8s config: %s", err))
	}
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Sprintf("Unable to create k8s client: %s", err))
	}

	for {
		pods, err := k8sClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(fmt.Sprintf("Unable to list pods: %s", err))
		}
		podGauge.Set(float64(len(pods.Items)))
		time.Sleep(30 * time.Second)
	}
}
