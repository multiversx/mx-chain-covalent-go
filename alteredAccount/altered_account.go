package alteredAccount

// TODO: Remove this from here and use the one from core, once feat/altered-account from elrond-go & elrond-go-core are finished

// AlteredAccount is the altered account dto response from Elrond proxy
type AlteredAccount struct {
	Address string `json:"address"`
	Balance string `json:"balance,omitempty"`
	Nonce   uint64 `json:"nonce"`
}
