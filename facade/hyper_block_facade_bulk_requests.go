package facade

import (
	"fmt"
	"sync"
	"time"

	"github.com/ElrondNetwork/covalent-indexer-go/api"
	"github.com/ElrondNetwork/covalent-indexer-go/cmd/proxy/config"
	"github.com/ElrondNetwork/elrond-go-core/core"
)

const (
	maxRequestsRetrial = 10
	waitTimeRetrialsMs = 50
)

func (hbf *hyperBlockFacade) getHyperBlocksByNonces(noncesInterval *api.Interval, options config.HyperBlocksQueryOptions) ([][]byte, error) {
	if noncesInterval.Start > noncesInterval.End {
		return nil, errInvalidNoncesInterval
	}
	if options.BatchSize == 0 {
		return nil, errInvalidBatchSize
	}

	expectedNumOfResults := noncesInterval.End - noncesInterval.Start + 1
	maxGoroutines := core.MinUint64(uint64(options.BatchSize), expectedNumOfResults)
	done := make(chan struct{}, maxGoroutines)
	wg := &sync.WaitGroup{}

	results := make([][]byte, expectedNumOfResults)
	mutex := sync.Mutex{}
	currIdx := uint32(0)

	var requestError error
	for nonce := noncesInterval.Start; nonce <= noncesInterval.End && requestError == nil; nonce++ {
		done <- struct{}{}
		wg.Add(1)

		request := hbf.getHyperBlockByNonceFullPath(nonce, options.QueryOptions)
		go func(req string, idx uint32) {
			res, err := hbf.getHyperBlockWithRetrials(req)

			mutex.Lock()
			defer func() {
				<-done
				wg.Done()
				mutex.Unlock()
			}()

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
	if len(results) != int(expectedNumOfResults) {
		return nil, fmt.Errorf("%w, expected to return %d, only got %d", errCouldNotGetAllHyperBlocks, expectedNumOfResults, len(results))
	}

	return results, requestError
}

func (hbf *hyperBlockFacade) getHyperBlockWithRetrials(request string) ([]byte, error) {
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

		time.Sleep(waitTimeRetrialsMs * time.Millisecond)
	}

	return nil, fmt.Errorf("%w from request = %s after num of retrials = %d", errCouldNotGetHyperBlock, request, maxRequestsRetrial)
}
