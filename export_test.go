package covalent

import "github.com/ElrondNetwork/covalent-indexer-go/schema"

func (c *covalentIndexer) SendBlockResultToCovalent(result *schema.BlockResult) {
	c.sendBlockResultToCovalent(result)
}
