package api

const (
	// UrlParameterWithBalances represents the name of an URL parameter to query balances per hyperBlock
	UrlParameterWithBalances = "withBalances"
	// UrlParameterWithLogs represents the name of an URL parameter to query logs per hyperBlock
	UrlParameterWithLogs = "withLogs"
)

// HyperBlockQueryOptions holds options for hyperBlock queries
type HyperBlockQueryOptions struct {
	WithLogs     bool
	WithBalances bool
}
