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
	// UrlParameterWithMetaData represents the name of an URL parameter to query meta data for altered accounts with tokens
	UrlParameterWithMetaData = "withMetadata"
)

// HyperBlockQueryOptions holds options for hyperBlock queries
type HyperBlockQueryOptions struct {
	WithLogs            bool
	WithAlteredAccounts bool
	NotarizedAtSource   bool
	Tokens              string
	WithMetaData        bool
}

var options = HyperBlockQueryOptions{
	WithLogs:            true,
	WithAlteredAccounts: true,
	NotarizedAtSource:   true,
	Tokens:              "all",
	WithMetaData:        true,
}
