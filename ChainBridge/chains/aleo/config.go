/* Copyright (c) 2021 Forte Labs
 * All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains the property of
 * Forte Labs and their suppliers, if any.  The intellectual and technical
 * concepts contained herein are proprietary to Forte Labs and their suppliers
 * and may be covered by U.S. and Foreign Patents, patents in process, and are
 * protected by trade secret or copyright law. Dissemination of this information
 * or reproduction of this material is strictly forbidden unless prior written
 * permission is obtained from Forte Labs.
 */

package aleo

import (
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-utils/core"

	"github.com/ChainSafe/chainbridge-utils/msg"
)

// Chain specific options
var (
	HttpOpt       = "http"
	StartBlockOpt = "startBlock"
	RelayerIdOpt  = "relayerId"
)

// Config encapsulates all necessary parameters in aleo compatible custodians
type Config struct {
	name               string // Human-readable chain name
	id                 msg.ChainId
	endpoint           string // url for the custodian's endpoint
	from               string // address of key to use, not used for the custodian
	keystorePath       string // Location of keyfiles, not used for the custodian
	blockstorePath     string // Location of the blockstore, not used for custodian
	freshStart         bool   // Disables loading from blockstore at start, in this case, loads all transactions on the custodian
	http               bool
	startBlock         *big.Int
	blockConfirmations *big.Int
	relayerId		   string
}

// parseChainConfig uses a core.ChainConfig to construct the corresponding Config
func parseChainConfig(chainCfg *core.ChainConfig) (*Config, error) {
	config := &Config{
		name:               chainCfg.Name,
		id:                 chainCfg.Id,
		endpoint:           chainCfg.Endpoint,
		from:               chainCfg.From,
		keystorePath:       chainCfg.KeystorePath,
		blockstorePath:     chainCfg.BlockstorePath,
		freshStart:         chainCfg.FreshStart,
		http:               false,
		startBlock:         big.NewInt(0),
		blockConfirmations: big.NewInt(0),
		relayerId:          "ab8f33ee-1f93-4cca-a104-a545ec6bec92",
	}

	if HTTP, ok := chainCfg.Opts[HttpOpt]; ok && HTTP == "true" {
		config.http = true
		delete(chainCfg.Opts, HttpOpt)
	} else if HTTP, ok := chainCfg.Opts[HttpOpt]; ok && HTTP == "false" {
		config.http = false
		delete(chainCfg.Opts, HttpOpt)
	}

	if startBlock, ok := chainCfg.Opts[StartBlockOpt]; ok && startBlock != "" {
		block := big.NewInt(0)
		_, pass := block.SetString(startBlock, 10)
		if pass {
			config.startBlock = block
			delete(chainCfg.Opts, StartBlockOpt)
		} else {
			return nil, fmt.Errorf("unable to parse %s", StartBlockOpt)
		}
	}

	if relayerId, ok := chainCfg.Opts[RelayerIdOpt]; ok && relayerId != "" {
		config.relayerId = relayerId
		delete(chainCfg.Opts, RelayerIdOpt)
	}

	if len(chainCfg.Opts) != 0 {
		return nil, fmt.Errorf("unknown Opts Encountered: %#v", chainCfg.Opts)
	}

	return config, nil

}
