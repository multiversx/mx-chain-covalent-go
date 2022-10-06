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

// TODO: Remove this once the feat is complete, since:
// - withLogs is implemented, but not available on mainnet
// - withBalances is not yet implemented
var options = HyperBlockQueryOptions{
	WithLogs:     true,
	WithBalances: false,
}
