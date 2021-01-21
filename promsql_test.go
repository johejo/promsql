package promsql

import (
	"database/sql"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestDBStatsCollector(t *testing.T) {
	c := NewDBStatsCollector(new(sql.DB), DBStatsCollectorOpts{
		DriverName: "test",
	})

	registry := prometheus.NewRegistry()
	if err := registry.Register(c); err != nil {
		t.Fatal(err)
	}

	mfs, err := registry.Gather()
	if err != nil {
		t.Fatal(err)
	}

	names := []string{
		"go_test_db_stats_max_open_connections",
		"go_test_db_stats_open_connections",
		"go_test_db_stats_in_use",
		"go_test_db_stats_idle",
		"go_test_db_stats_wait_duration",
		"go_test_db_stats_max_idle_closed",
		"go_test_db_stats_max_idle_time_closed",
		"go_test_db_stats_max_lifetime_closed",
	}
	type result struct {
		found bool
	}
	results := make(map[string]result)
	for _, name := range names {
		results[name] = result{found: false}
	}
	for _, mf := range mfs {
		for _, name := range names {
			if name == *mf.Name {
				results[name] = result{found: true}
				break
			}
		}
	}

	for name, result := range results {
		if !result.found {
			t.Errorf("%s not found", name)
		}
	}
}
