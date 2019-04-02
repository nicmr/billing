package main

import (
	"github.com/Altemista/altemista-billing/pkg/query"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"log"
	"net/http"
)

var (
	// Session is safe for concurrent use after initialization,
	// as it will not be mutated by the SDK after creation
	sess = createSessionOrFatal()
)

func createSessionOrFatal() *(session.Session) {
	sess, err := session.NewSession()
	if err != nil {
		log.Fatal("Unable to initialize aws session: ", err)
	}
	return sess
}

func costs(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	log.Println(start, end)

	svc := costexplorer.New(sess)

	// TODO: Validate start and end inputs here, and throw http.StatusBadRequest if doesn't match pattern
	output, err := query.CostsBetween(svc, start, end)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(output.String()))
}

func main() {
	http.HandleFunc("/costs", costs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
