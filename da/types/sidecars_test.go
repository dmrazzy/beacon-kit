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
// AN "AS IS" BASIS. LICENSOR HEREBY DISCLAIMS ALL WARRANTIES AND CONDITIONS,
// EXPRESS OR IMPLIED, INCLUDING (WITHOUT LIMITATION) WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE, NON-INFRINGEMENT, AND
// TITLE.

package types_test

import (
	"strconv"
	"testing"

	ctypes "github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/da/types"
	byteslib "github.com/berachain/beacon-kit/primitives/bytes"
	"github.com/berachain/beacon-kit/primitives/common"
	"github.com/berachain/beacon-kit/primitives/crypto"
	"github.com/berachain/beacon-kit/primitives/eip4844"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func TestEmptySidecarMarshalling(t *testing.T) {
	t.Parallel()
	inclusionProof := make([]common.Root, 0)
	// Create an empty BlobSidecar
	for i := 1; i <= ctypes.KZGInclusionProofDepth; i++ {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		proof, errBytes := byteslib.ToBytes32(it)
		require.NoError(t, errBytes)
		inclusionProof = append(inclusionProof, common.Root(proof))
	}

	sidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.SignedBeaconBlockHeader{
			Header:    &ctypes.BeaconBlockHeader{},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		[48]byte{},
		inclusionProof,
	)

	// Marshal the empty sidecar
	marshalled, err := sidecar.MarshalSSZ()
	require.NoError(
		t,
		err,
		"Marshalling empty sidecar should not produce an error",
	)
	require.NotNil(
		t,
		marshalled,
		"Marshalling empty sidecar should produce a result",
	)

	// Unmarshal the empty sidecar
	unmarshalled := new(types.BlobSidecar)
	err = ssz.Unmarshal(marshalled, unmarshalled)
	require.NoError(
		t,
		err,
		"Unmarshalling empty sidecar should not produce an error",
	)

	// Compare the original and unmarshalled empty sidecars
	require.Equal(
		t,
		sidecar,
		unmarshalled,
		"The original and unmarshalled empty sidecars should be equal",
	)
}

func TestValidateBlockRoots(t *testing.T) {
	t.Parallel()
	inclusionProof := make([]common.Root, 0)
	// Create a sample BlobSidecar with valid roots
	for i := 1; i <= ctypes.KZGInclusionProofDepth; i++ {
		it := byteslib.ExtendToSize([]byte(strconv.Itoa(i)), byteslib.B32Size)
		proof, errBytes := byteslib.ToBytes32(it)
		require.NoError(t, errBytes)
		inclusionProof = append(inclusionProof, common.Root(proof))
	}

	validSidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.SignedBeaconBlockHeader{
			Header: &ctypes.BeaconBlockHeader{
				StateRoot: [32]byte{1},
				BodyRoot:  [32]byte{2},
			},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		[48]byte{},
		[48]byte{},
		inclusionProof,
	)

	// Validate the sidecar with valid roots
	sidecars := types.BlobSidecars{
		validSidecar,
	}
	err := sidecars.ValidateBlockRoots()
	require.NoError(
		t,
		err,
		"Validating sidecar with valid roots should not produce an error",
	)

	// Create a sample BlobSidecar with invalid roots
	differentBlockRootSidecar := types.BuildBlobSidecar(
		math.U64(0),
		&ctypes.SignedBeaconBlockHeader{
			Header: &ctypes.BeaconBlockHeader{
				StateRoot: [32]byte{1},
				BodyRoot:  [32]byte{3},
			},
			Signature: crypto.BLSSignature{},
		},
		&eip4844.Blob{},
		eip4844.KZGCommitment{},
		eip4844.KZGProof{},
		inclusionProof,
	)
	// Validate the sidecar with invalid roots
	sidecarsInvalid := types.BlobSidecars{
		validSidecar,
		differentBlockRootSidecar,
	}
	err = sidecarsInvalid.ValidateBlockRoots()
	require.Error(
		t,
		err,
		"Validating sidecar with invalid roots should produce an error",
	)
}

func TestZeroSidecarsInBlobSidecarsIsNotNil(t *testing.T) {
	t.Parallel()
	// This test exists to ensure that proposing a BlobSidecars with 0
	// Sidecars is not considered IsNil().
	sidecars := &types.BlobSidecars{}
	require.NotNil(t, sidecars)
}
