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

package beacon

import (
	"strings"

	beacontypes "github.com/berachain/beacon-kit/mod/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/mod/node-api/handlers/utils"
)

func (h *Handler[
	BeaconBlockHeaderT, ContextT, _, _,
]) GetBlobSidecars(c ContextT) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetBlobSidecarsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	_, err = utils.SlotFromBlockID(req.BlockID, h.backend)
	if err != nil {
		return nil, err
	}

	// Return a sample blob sidecar
	return beacontypes.BlobSidecarsResponse[BeaconBlockHeaderT]{
		Data: []*beacontypes.BlobSidecarsData[BeaconBlockHeaderT]{
			{
				Index:                       0,
				Blob:                        "0x" + strings.Repeat("00", 2),
				KZGCommitment:               "0x" + strings.Repeat("00", 2),
				KZGProof:                    "0x" + strings.Repeat("00", 2),
				KZGCommitmentInclusionProof: make([]string, 17), // Array of 17 empty strings
				// SignedBlockHeader: &beacontypes.BlockHeader[BeaconBlockHeaderT]{
				// 	Message:   bytes.B48{},
				// 	Signature: bytes.B48{}, // TODO: implement
				// }, // Empty block header
			},
		},
	}, nil
}
