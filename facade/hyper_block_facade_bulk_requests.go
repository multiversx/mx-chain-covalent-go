package facade

import (
	"errors"
	"fmt"
	"sync"
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
		currBatch = make([]string, 0, batchSize)
	}

	return batches, nil
}

type avroHyperBlockResponse struct {
	encodedHyperBlock []byte
	err               error
}

func (hbf *hyperBlockFacade) requestBatchesConcurrently(batches [][]string, totalRequests uint64) ([][]byte, error) {
	ch := make(chan *avroHyperBlockResponse)
	responses := make([][]byte, 0)

	for _, batch := range batches {
		go func(batch []string) {
			for _, request := range batch {
				encodedHyperBlock, err := hbf.getHyperBlockAvroBytes(request)
				fmt.Println("sending request: ", request)
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
			if len(responses) == int(totalRequests) {
				return responses, nil
			}
		case <-time.After(500 * time.Millisecond):
			return nil, errors.New("timeout")
		}
	}
}

type encodedHyperBlock struct {
	nonce        uint64
	encodedBytes []byte
}

func (hbf *hyperBlockFacade) getBlocksByNonces(requests []string) ([][]byte, error) {
	maxGoroutines := 20
	done := make(chan struct{}, maxGoroutines)
	wg := &sync.WaitGroup{}

	results := make([][]byte, 0, len(requests))
	mutex := sync.Mutex{}
	for _, request := range requests {
		done <- struct{}{}
		wg.Add(1)
		fmt.Println("sending request: ", request)
		go func(req string) {
			res, err := hbf.getBlockByRequest(req, done, wg)
			if err != nil {
				// TODO treat error
				return
			}

			mutex.Lock()
			results = append(results, res)
			mutex.Unlock()

		}(request)
	}

	wg.Wait()
	return results, nil
}

func (hbf *hyperBlockFacade) getBlockByRequest(request string, done chan struct{}, wg *sync.WaitGroup) ([]byte, error) {
	defer func() {
		<-done
		wg.Done()
	}()
	return hbf.getHyperBlockAvroBytes(request)
}
