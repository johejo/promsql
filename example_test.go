package promsql_test

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/johejo/promsql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Example() {
	db, err := sql.Open("mysql", "CONNECTION_STRING")
	if err != nil {
		panic(err)
	}
	if err := prometheus.Register(promsql.NewDBStatsCollector(db, promsql.DBStatsCollectorOpts{
		DriverName: "mysql",
	})); err != nil {
		panic(err)
	}

	go func() {
		http.ListenAndServe(":8080", promhttp.Handler())
	}()

	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	s := bufio.NewScanner(resp.Body)

	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "mysql") {
			fmt.Println(line)
		}
	}
}
