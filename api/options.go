package api

const (
	// UrlParameterWithLogs represents the name of an URL parameter to query logs per hyperBlock
	UrlParameterWithLogs = "withLogs"
	// UrlParameterNotarizedAtSource represents the name of an URL parameter to query hyper blocks which are notarized at source
	UrlParameterNotarizedAtSource = "notarizedAtSource"
	// UrlParameterWithAlteredAccounts represents the name of an URL parameter to query altered accounts per hyperBlock
	UrlParameterWithAlteredAccounts = "withAlteredAccounts"
	// UrlParameterTokens represents the name of an URL parameter to query altered accounts with tokens
	UrlParameterTokens = "tokens"
)

// Interval defines a [start,end] interval
type Interval struct {
	Start uint64
	End   uint64
}
