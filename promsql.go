// Package promsql exports *sql.DB stats as prometheus metrics collector.
package promsql

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

// DB represents a part of *sql.DB.
type DB interface {
	// Stats returns database statistics.
	Stats() sql.DBStats
}

var _ DB = (*sql.DB)(nil)

// DBStatsCollectorOpts defines the behavior of a db stats collector
// created with NewDBStatsCollector.
type DBStatsCollectorOpts struct {
	// DriverName holds the name of driver.
	// It will not used for empty strings.
	DriverName string
}

// NewDBStatsCollector returns a collector that exports metrics about the given *sql.DB.
// See https://golang.org/pkg/database/sql/#DBStats for more information on stats.
func NewDBStatsCollector(db DB, opts DBStatsCollectorOpts) prometheus.Collector {
	var fqName func(name string) string
	if opts.DriverName == "" {
		fqName = func(name string) string {
			return "go_db_stats_" + name
		}
	} else {
		fqName = func(name string) string {
			return "go_" + opts.DriverName + "_db_stats_" + name
		}
	}

	return &dbStatsCollector{
		db: db,
		maxOpenConnsDesc: prometheus.NewDesc(
			fqName("max_open_connections"),
			"Maximum number of open connections to the database.",
			nil, nil,
		),
		openConnsDesc: prometheus.NewDesc(
			fqName("open_connections"),
			"The number of established connections both in use and idle.",
			nil, nil,
		),
		inUseDesc: prometheus.NewDesc(
			fqName("in_use"),
			"The number of connections currently in use.",
			nil, nil,
		),
		idleDesc: prometheus.NewDesc(
			fqName("idle"),
			"The number of idle connections.",
			nil, nil,
		),
		waitCountDesc: prometheus.NewDesc(
			fqName("wait_count"),
			"The total number of connections waited for.",
			nil, nil,
		),
		waitDurationDesc: prometheus.NewDesc(
			fqName("wait_duration"),
			"The total time blocked waiting for a new connection.",
			nil, nil,
		),
		maxIdleClosedDesc: prometheus.NewDesc(
			fqName("max_idle_closed"),
			"The total number of connections closed due to SetMaxIdleConns.",
			nil, nil,
		),
		maxIdleTimeClosedDesc: prometheus.NewDesc(
			fqName("max_idle_time_closed"),
			"The total number of connections closed due to SetConnMaxIdleTime.",
			nil, nil,
		),
		maxLifetimeClosedDesc: prometheus.NewDesc(
			fqName("max_lifetime_closed"),
			"The total number of connections closed due to SetConnMaxLifetime.",
			nil, nil,
		),
	}
}

type dbStatsCollector struct {
	db                    DB
	maxOpenConnsDesc      *prometheus.Desc
	openConnsDesc         *prometheus.Desc
	inUseDesc             *prometheus.Desc
	idleDesc              *prometheus.Desc
	waitCountDesc         *prometheus.Desc
	waitDurationDesc      *prometheus.Desc
	maxIdleClosedDesc     *prometheus.Desc
	maxIdleTimeClosedDesc *prometheus.Desc
	maxLifetimeClosedDesc *prometheus.Desc
}

// Describe returns all descriptions of the collector.
func (c *dbStatsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.maxOpenConnsDesc
	ch <- c.inUseDesc
	ch <- c.idleDesc
	ch <- c.waitCountDesc
	ch <- c.waitDurationDesc
	ch <- c.maxIdleClosedDesc
	ch <- c.maxIdleTimeClosedDesc
	ch <- c.maxLifetimeClosedDesc
}

// Collect returns the current state of all metrics of the collector.
func (c *dbStatsCollector) Collect(ch chan<- prometheus.Metric) {
	stats := c.db.Stats()
	ch <- prometheus.MustNewConstMetric(c.maxOpenConnsDesc, prometheus.GaugeValue, float64(stats.MaxOpenConnections))
	ch <- prometheus.MustNewConstMetric(c.openConnsDesc, prometheus.GaugeValue, float64(stats.OpenConnections))
	ch <- prometheus.MustNewConstMetric(c.inUseDesc, prometheus.GaugeValue, float64(stats.InUse))
	ch <- prometheus.MustNewConstMetric(c.idleDesc, prometheus.GaugeValue, float64(stats.Idle))
	ch <- prometheus.MustNewConstMetric(c.waitCountDesc, prometheus.CounterValue, float64(stats.WaitCount))
	ch <- prometheus.MustNewConstMetric(c.waitDurationDesc, prometheus.CounterValue, float64(stats.WaitDuration))
	ch <- prometheus.MustNewConstMetric(c.maxIdleClosedDesc, prometheus.CounterValue, float64(stats.MaxIdleClosed))
	ch <- prometheus.MustNewConstMetric(c.maxIdleTimeClosedDesc, prometheus.CounterValue, float64(stats.MaxIdleTimeClosed))
	ch <- prometheus.MustNewConstMetric(c.maxLifetimeClosedDesc, prometheus.CounterValue, float64(stats.MaxLifetimeClosed))
}
