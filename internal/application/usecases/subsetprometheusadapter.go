package usecases

import (
	"log/slog"
	"os"

	"github.com/hanapedia/metrics-processor/internal/application/usecases/query/hexagon"
	"github.com/hanapedia/metrics-processor/internal/domain"
	"github.com/hanapedia/metrics-processor/internal/infrastructure/prometheus"
	"github.com/hanapedia/metrics-processor/pkg/promql"
)

// SubsetPrometheusQueryAdapter creates proemtheusAdapter with subset of queries
// use this adapter when partial requery is needed
func SubsetPrometheusQueryAdapter(config *domain.Config) *prometheus.PrometheusAdapter {
	prometheusAdapter, err := prometheus.NewPrometheusAdapter(config)
	if err != nil {
		slog.Error("Failed to create new Prometheus adapter", "err", err)
		os.Exit(1)
	}

	/* rateConfigs := []query.RateConfig{ */
	/* 	{Name: "5m", Duration: 5 * time.Minute, IsInstant: false}, */
	/* 	{Name: "1m", Duration: 1 * time.Minute, IsInstant: false}, */
	/* 	{Name: "1m", Duration: 1 * time.Minute, IsInstant: true}, */
	/* } */
	defaultSvc := "service-.*"
	filters := []promql.Filter{
		promql.NewFilter("experiment", "=~", config.K6TestName),
		promql.NewFilter("service", "=~", defaultSvc),
		promql.NewFilter("namespace", "=", config.Namespace),
	}

	// Register non-rate or non-irate queries
	queries := []*promql.Query{
		// adaptive timeout
		hexagon.NewAdaptiveTimeoutQuery(hexagon.Call, filters).SetName("adaptive_call_timeout"), // adaptive call timeout
		hexagon.NewAdaptiveTimeoutQuery(hexagon.Task, filters).SetName("adaptive_task_timeout"), // adaptive task timeout
	}
	for _, query := range queries {
		prometheusAdapter.RegisterQuery(query)
	}

	// Register rate & irate queries
	/* for _, rateConfig := range rateConfigs { */
	/* 	queries := []*promql.Query{ */
	/**/
	/* 		// container metrics */
	/* 		container.CreateCpuUsageQuery(containerFilter, rateConfig). */
	/* 			SetName(rateConfig.AddSuffix("cpu_usage")), */
	/* 		container.CreateCpuThrottleQuery(containerFilter, rateConfig). */
	/* 			SetName(rateConfig.AddSuffix("cpu_throttled")), */
	/* 	} */
	/* 	for _, query := range queries { */
	/* 		prometheusAdapter.RegisterQuery(query) */
	/* 	} */
	/* } */

	return prometheusAdapter
}
