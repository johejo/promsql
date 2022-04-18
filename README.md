Use [prometheus/client_golang's collector](https://pkg.go.dev/github.com/prometheus/client_golang@v1.12.1/prometheus/collectors#NewDBStatsCollector) instead.

# promsql

[![ci](https://github.com/johejo/promsql/workflows/ci/badge.svg?branch=main)](https://github.com/johejo/promsql/actions?query=workflow%3Aci)
[![Go Reference](https://pkg.go.dev/badge/github.com/johejo/promsql.svg)](https://pkg.go.dev/github.com/johejo/promsql)
[![codecov](https://codecov.io/gh/johejo/promsql/branch/main/graph/badge.svg)](https://codecov.io/gh/johejo/promsql)
[![Go Report Card](https://goreportcard.com/badge/github.com/johejo/promsql)](https://goreportcard.com/report/github.com/johejo/promsql)

Package promsql exports \*sql.DB stats as prometheus metrics collector.

## Example

```go
package promsql_test

import (
	"bufio"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/johejo/promsql"
	_ "github.com/go-sql-driver/mysql"
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
```

## License

MIT

## Author

Mitsuo Heijo
