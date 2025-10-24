package metrics

type ErrorMetrics struct {
	Total      int64
	ClientErrs int64
	ServerErrs int64
	Timeouts   int64
}
