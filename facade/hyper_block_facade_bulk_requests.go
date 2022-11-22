package facade

import (
	"errors"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

func (hbf *hyperBlockFacade) createBatchRequests(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) ([][]string, error) {
	if noncesInterval.Start > noncesInterval.End {
		return nil, errInvalidNoncesInterval
	}

	numBatches := noncesInterval.End - noncesInterval.Start + 1
	batchSize := core.MinUint64(options.BatchSize, numBatches)

	batches := make([][]string, 0, numBatches)
	currBatch := make([]string, 0, batchSize)

	idx := uint64(0)
	for nonce := noncesInterval.Start; nonce <= noncesInterval.End; nonce++ {
		idx++

		request := hbf.getHyperBlockByNonceFullPath(nonce, options.QueryOptions)
		currBatch = append(currBatch, request)

		notEnoughRequestsInBatch := len(currBatch) < int(batchSize)
		notLastBatch := idx != numBatches
		if notEnoughRequestsInBatch && notLastBatch {
			continue
		}

		batches = append(batches, currBatch)
	}

	return batches, nil
}

type avroHyperBlockResponse struct {
	encodedHyperBlock []byte
	err               error
}

func (hbf *hyperBlockFacade) requestBatchesConcurrently(batches [][]string) ([][]byte, error) {
	ch := make(chan *avroHyperBlockResponse)
	responses := make([][]byte, 0)

	for _, batch := range batches {
		go func(batch []string) {
			for _, request := range batch {
				encodedHyperBlock, err := hbf.getHyperBlockAvroBytes(request)
				ch <- &avroHyperBlockResponse{
					encodedHyperBlock: encodedHyperBlock,
					err:               err,
				}
			}

		}(batch)
	}

	for {
		select {
		case r := <-ch:
			//fmt.Printf("%s was fetched\n", r.url)
			if r.err != nil {
				return nil, r.err
			}
			responses = append(responses, r.encodedHyperBlock)
			if len(responses) == len(batches) {
				return responses, nil
			}
		case <-time.After(50 * time.Millisecond):
			return nil, errors.New("timeout")
		}
	}
}
