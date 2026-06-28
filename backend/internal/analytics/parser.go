package analytics

import (
	"fmt"
	"strings"
	"time"
)

func ParseMarketDeployed(factoryAddress string, log BlockscoutLog) (MarketDeployed, bool, error) {
	if log.Decoded == nil {
		return MarketDeployed{}, false, nil
	}
	if !strings.HasPrefix(log.Decoded.MethodCall, "MarketDeployed(") {
		return MarketDeployed{}, false, nil
	}

	params := decodedParamMap(log.Decoded.Parameters)
	event := MarketDeployed{
		FactoryAddress:         factoryAddress,
		MarketIDHash:           params["marketId"],
		MarketAddress:          params["market"],
		CreatorAddress:         params["creator"],
		ResolverAddress:        params["resolver"],
		CollateralTokenAddress: params["collateralToken"],
		CloseTimestamp:         params["closeTimestamp"],
		Question:               params["question"],
		TransactionHash:        log.TransactionHash,
		BlockNumber:            log.BlockNumber,
		LogIndex:               log.Index,
		Raw:                    log.Raw,
	}

	if event.MarketAddress == "" {
		return MarketDeployed{}, true, fmt.Errorf("MarketDeployed log missing market parameter")
	}
	if event.TransactionHash == "" {
		return MarketDeployed{}, true, fmt.Errorf("MarketDeployed log missing transaction hash")
	}
	if log.BlockTimestamp != "" {
		timestamp, err := time.Parse(time.RFC3339Nano, log.BlockTimestamp)
		if err != nil {
			return MarketDeployed{}, true, fmt.Errorf("parse MarketDeployed block timestamp: %w", err)
		}
		event.BlockTimestamp = timestamp
	}

	return event, true, nil
}

func decodedParamMap(params []BlockscoutParameter) map[string]string {
	mapped := make(map[string]string, len(params))
	for _, param := range params {
		mapped[param.Name] = param.Value
	}
	return mapped
}
