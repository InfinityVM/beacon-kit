// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2024, Berachain Foundation. All rights reserved.
// Use of this software is governed by the Business Source License included
// in the LICENSE file of this repository and at www.mariadb.com/bsl11.
//
// ANY USE OF THE LICENSED WORK IN VIOLATION OF THIS LICENSE WILL AUTOMATICALLY
// TERMINATE YOUR RIGHTS UNDER THIS LICENSE FOR THE CURRENT AND ALL OTHER
// VERSIONS OF THE LICENSED WORK.
//
// THIS LICENSE DOES NOT GRANT YOU ANY RIGHT IN ANY TRADEMARK OR LOGO OF
// LICENSOR OR ITS AFFILIATES (PROVIDED THAT YOU MAY USE A TRADEMARK OR LOGO OF
// LICENSOR AS EXPRESSLY REQUIRED BY THIS LICENSE).
//
// TO THE EXTENT PERMITTED BY APPLICABLE LAW, THE LICENSED WORK IS PROVIDED ON
// AN “AS IS” BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package main

import (
	"encoding/json"
	"flag"
	clibuilder "github.com/berachain/beacon-kit/mod/cli/pkg/builder"
	clicomponents "github.com/berachain/beacon-kit/mod/cli/pkg/components"
	spec "github.com/berachain/beacon-kit/mod/config/pkg/spec"
	nodebuilder "github.com/berachain/beacon-kit/mod/node-core/pkg/builder"
	nodecomponents "github.com/berachain/beacon-kit/mod/node-core/pkg/components"
	nodetypes "github.com/berachain/beacon-kit/mod/node-core/pkg/types"
	common "github.com/berachain/beacon-kit/mod/primitives/pkg/common"
	"go.uber.org/automaxprocs/maxprocs"
	"log/slog"
	"os"
)

type Node = nodetypes.Node

// ProvideChainSpecWithInput provides the chain spec based on the input.
func ProvideChainSpecWithInput(chainSpec *spec.ChainSpecInput) func() common.ChainSpec {
	return func() common.ChainSpec {
		return nodecomponents.ProvideChainSpec(chainSpec)
	}
}

// run runs the beacon node.
func run(chainSpec *spec.ChainSpecInput) error {
	// Set the uber max procs
	if _, err := maxprocs.Set(); err != nil {
		return err
	}

	// Build the node using the node-core.
	nb := nodebuilder.New(
		// Set the Runtime Components to the Default.
		nodebuilder.WithComponents[Node, *Logger, *LoggerConfig](
			DefaultComponents(),
		),
	)

	// Build the root command using the builder
	cb := clibuilder.New(
		// Set the Name to the Default.
		clibuilder.WithName[Node, *ExecutionPayload, *Logger](
			"beacond",
		),
		// Set the Description to the Default.
		clibuilder.WithDescription[Node, *ExecutionPayload, *Logger](
			"A basic beacon node, usable most standard networks.",
		),
		// Set the Runtime Components to the Default.
		clibuilder.WithComponents[Node, *ExecutionPayload, *Logger](
			append(
				clicomponents.DefaultClientComponents(),
				// TODO: remove these, and eventually pull cfg and chainspec
				// from built node
				ProvideChainSpecWithInput(chainSpec),
			),
		),
		// Set the NodeBuilderFunc to the NodeBuilder Build.
		clibuilder.WithNodeBuilderFunc[
			Node, *ExecutionPayload, *Logger,
		](nb.Build),
	)

	cmd, err := cb.Build()
	if err != nil {
		return err
	}

	// eventually we want to decouple from cosmos cli, and just pass in a built
	// Node and Cmd to a runner

	// for now, running the cmd will start the node
	return cmd.Run(clicomponents.DefaultNodeHome)
}

// main is the entry point.
func main() {
	// Define the chain-spec flag
	chainSpecPath := flag.String("chain-spec", "", "Path to the chain specification JSON file")
	flag.Parse()
	specType := os.Getenv("CHAIN_SPEC")
	if *chainSpecPath == "" && specType == "" {
		slog.Error("Path to chain-spec must be specified or set env var CHAIN_SPEC")
		os.Exit(1)
	}

	var chainSpecInput *spec.ChainSpecInput
	if *chainSpecPath != "" {
		b, err := os.ReadFile(*chainSpecPath)
		if err != nil {
			slog.Error("Failed to open chain specification file", "path", *chainSpecPath, "err", err)
			os.Exit(1)
		}

		chainSpecInput := new(spec.ChainSpecInput)
		err = json.Unmarshal(b, chainSpecInput)
		if err != nil {
			slog.Error("Failed to unmarshal chain specification", "path", *chainSpecPath, "err", err)
			os.Exit(1)
		}
	}

	if err := run(chainSpecInput); err != nil {
		//nolint:sloglint // todo fix.
		slog.Error("startup failure", "error", err)
		os.Exit(1)
	}
}
