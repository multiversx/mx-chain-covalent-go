# covalent-indexer-go

Covalent indexer acts as a second proxy layer over Elrond's Proxy by providing avro encoded hyper blocks.

## How it works

1. Covalent proxy hyperBlock request is triggered(e.g.: `http://127.0.0.1:7952/hyperblock/by-nonce/37`).
2. Covalent proxy will interrogate the backing Elrond proxy via REST API to fetch requested hyperBlock using query
   parameters defined in the config(`cmd/proxy/config.toml`).
3. Received data from Elrond proxy will be converted to defined avro schema (`schema/block.elrond.avsc`)
4. Covalent proxy response will have the following format (compatible with Elrond proxy API responses):

```  
{
   "data": {},
   "error": "",
   "code": "successful"
}
```

Data field contains a byte array corresponding to the encoded avro schema.

## How to use

1. Go to `cmd/proxy`
2. Build `go build`
3. Run `./proxy`

In `cmd/proxy/config.toml` one can find:

1. Elrond & Covalent proxy configuration(e.g. `port`, `elrondProxyUrl`, etc.)
2. `hyperBlockQueryOptions` used to format hyperBlock query for Elrond proxy. E.g.: following Covalent
   request: `localhost:port/hyperblock/by-nonce/4`, having `withAlteredAccounts = true` and `tokens = all` will trigger
   the following request : `elrondProxy:port/hyperblock/by-nonce/4?withAlteredAccounts=true&tokens=all`

_Please note that altered-accounts endpoints will only work if the backing observers of the Elrond Proxy have support
for historical balances (--operation-mode historical-balances when starting the node)_

## Endpoints

- `/hyperblock/by-nonce/:nonce` (GET) --> returns a hyperblock by nonce, with transactions included
- `/hyperblock/by-hash/:hash` (GET) --> returns a hyperblock by hash, with transactions included

## Avro schema update

In case you want to modify the existing avro schema, after finishing your changes, you need to re-generate the
corresponding code, by:

- Running `go generate` from `schema/codegen.go`
