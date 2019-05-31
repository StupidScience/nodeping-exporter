package main

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace = "nodeping"
)

// Collector type for prometheus.Collector interface implementation
type Collector struct {
	CheckUp       prometheus.GaugeVec
	CheckDuration prometheus.GaugeVec

	totalScrapes  prometheus.Counter
	failedScrapes prometheus.Counter

	np NodePing
	sync.Mutex
}

// Describe for prometheus.Collector interface implementation
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

// Collect for prometheus.Collector interface implementation
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.Lock()
	defer c.Unlock()

	c.totalScrapes.Inc()

	checks, err := c.np.GetAllChecks()
	if err != nil {
		c.failedScrapes.Inc()
		return
	}
	wg := &sync.WaitGroup{}
	for _, check := range checks {
		wg.Add(1)
		go func(check Check, ch chan<- prometheus.Metric, wg *sync.WaitGroup) {
			defer wg.Done()
			cs, err := c.np.GetCheckStats(check.ID)
			if err != nil {
				c.failedScrapes.Inc()
			}
			if cs.Success {
				c.CheckUp.WithLabelValues(check.Label, cs.Target, cs.Type).Set(1)
			} else {
				c.CheckUp.WithLabelValues(check.Label, cs.Target, cs.Type).Set(0)
			}
			c.CheckDuration.WithLabelValues(check.Label, cs.Target, cs.Type).Set(float64(cs.Duration) / 1000)

			ch <- c.CheckUp.WithLabelValues(check.Label, cs.Target, cs.Type)
			ch <- c.CheckDuration.WithLabelValues(check.Label, cs.Target, cs.Type)
		}(check, ch, wg)
	}

	wg.Wait()
	ch <- c.totalScrapes
	ch <- c.failedScrapes
}

// NewCollector create new collector struct
func NewCollector(np NodePing) (*Collector, error) {
	c := Collector{
		np: np,
	}

	c.totalScrapes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "exporter_scrapes_total",
		Help:      "Count of total scrapes",
	})

	c.failedScrapes = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Name:      "exporter_failed_balance_scrapes_total",
		Help:      "Count of failed balance scrapes",
	})

	c.CheckUp = *prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "check_up",
			Help:      "Current check status: 1 - for up, 0 - for down",
		},
		[]string{
			"label",
			"target",
			"type",
		},
	)

	c.CheckDuration = *prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "check_duration_seconds",
			Help:      "Last check duration in seconds",
		},
		[]string{
			"label",
			"target",
			"type",
		},
	)

	return &c, nil
}
