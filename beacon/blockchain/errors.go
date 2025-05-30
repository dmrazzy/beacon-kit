// SPDX-License-Identifier: BUSL-1.1
//
// Copyright (C) 2025, Berachain Foundation. All rights reserved.
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

package blockchain

import "github.com/berachain/beacon-kit/errors"

var (
	// ErrTooManyConsensusTxs is an error for consensus blocks having more than MaxConsensusTxsCount txs.
	ErrTooManyConsensusTxs = errors.New("too many consensus txs")
	// ErrUnexpectedBlockSlot is an error for consensus blocks with non consecutive slots.
	ErrUnexpectedBlockSlot = errors.New("unexpected block slot")
	// ErrNilBlk is an error for when the beacon block is nil.
	ErrNilBlk = errors.New("nil beacon block")
	// ErrNilBlob is an error for when the BlobSidecars is nil.
	ErrNilBlob = errors.New("nil blob")
	// ErrVersionMismatch is an error for when the fork for the block timestamp does not match the fork
	// for the ABCI timestamp.
	ErrVersionMismatch = errors.New("ABCI fork version mismatch")
	// ErrDataNotAvailable indicates that the required data is not available.
	ErrDataNotAvailable = errors.New("data not available")
	// ErrSidecarCommitmentMismatch indicates that the BeaconBlockBody commitments do not match the sidecars.
	ErrSidecarCommitmentMismatch = errors.New("sidecars commitments mismatch")
	// ErrSidecarSignatureMismatch indicates that the sidecar signature is invalid.
	ErrSidecarSignatureMismatch = errors.New("sidecar signature mismatch")
)
