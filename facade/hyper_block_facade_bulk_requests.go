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

const maxRequestsRetrial = 10

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

func (hbf *hyperBlockFacade) getBlocksByNonces(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) ([][]byte, error) {
	if noncesInterval.Start > noncesInterval.End {
		return nil, errInvalidNoncesInterval
	}

	maxGoroutines := 20
	done := make(chan struct{}, maxGoroutines)
	wg := &sync.WaitGroup{}

	results := make([][]byte, noncesInterval.End-noncesInterval.Start+1)
	mutex := sync.Mutex{}
	currIdx := uint32(0)

	var requestError error
	for nonce := noncesInterval.Start; nonce <= noncesInterval.End && requestError != nil; nonce++ {
		done <- struct{}{}
		wg.Add(1)

		request := hbf.getHyperBlockByNonceFullPath(nonce, options.QueryOptions)
		go func(req string, idx uint32) {
			res, err := hbf.getHyperBlockWithRetrials(req, done, wg)

			mutex.Lock()
			defer mutex.Unlock()

			if err != nil {
				requestError = err
				return
			}

			results[idx] = res
		}(request, currIdx)

		currIdx++
	}

	wg.Wait()

	if requestError != nil {
		return nil, fmt.Errorf("one or more errors occurred; last known error: %w", requestError)
	}
	return results, requestError
}

func (hbf *hyperBlockFacade) getHyperBlockWithRetrials(request string, done chan struct{}, wg *sync.WaitGroup) ([]byte, error) {
	defer func() {
		<-done
		wg.Done()
	}()

	ctRetrials := 0
	for ctRetrials < maxRequestsRetrial {
		res, err := hbf.getHyperBlockAvroBytes(request)
		if err == nil {
			return res, nil
		}

		ctRetrials++
		log.Warn("could not get hyperblock; retrying...",
			"request", request,
			"error", err,
			"num retrials", ctRetrials)
	}

	return nil, fmt.Errorf("%w from request = %s after num of retrials = %d", errCouldNotGetHyperBlock, request, maxRequestsRetrial)
}
