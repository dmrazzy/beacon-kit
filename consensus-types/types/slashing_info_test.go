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

package types_test

import (
	"io"
	"testing"

	"github.com/berachain/beacon-kit/consensus-types/types"
	"github.com/berachain/beacon-kit/primitives/encoding/ssz"
	"github.com/berachain/beacon-kit/primitives/math"
	"github.com/stretchr/testify/require"
)

func generateSlashingInfo() *types.SlashingInfo {
	return &types.SlashingInfo{
		Slot:  12345,
		Index: 67890,
	}
}

func TestSlashingInfo_MarshalSSZ_UnmarshalSSZ(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name     string
		data     *types.SlashingInfo
		expected *types.SlashingInfo
		err      error
	}{
		{
			name:     "Valid SlashingInfo",
			data:     generateSlashingInfo(),
			expected: generateSlashingInfo(),
			err:      nil,
		},
		{
			name: "Empty SlashingInfo",
			data: &types.SlashingInfo{
				Slot:  0,
				Index: 0,
			},
			expected: &types.SlashingInfo{
				Slot:  0,
				Index: 0,
			},
			err: nil,
		},
		{
			name:     "Invalid Buffer Size",
			data:     generateSlashingInfo(),
			expected: nil,
			err:      io.ErrUnexpectedEOF,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			data, err := tc.data.MarshalSSZ()
			require.NoError(t, err)
			require.NotNil(t, data)

			unmarshalled := new(types.SlashingInfo)
			if tc.name == "Invalid Buffer Size" {
				err = ssz.Unmarshal(data[:8], unmarshalled)
				require.ErrorIs(t, err, tc.err)
			} else {
				err = ssz.Unmarshal(data, unmarshalled)
				require.NoError(t, err)
				require.Equal(t, tc.expected, unmarshalled)

				var buf []byte
				buf, err = tc.data.MarshalSSZTo(buf)
				require.NoError(t, err)

				// The two byte slices should be equal
				require.Equal(t, data, buf)
			}
		})
	}
}

func TestSlashingInfo_GetTree(t *testing.T) {
	t.Parallel()
	data := generateSlashingInfo()

	tree, err := data.GetTree()
	require.NoError(t, err)
	require.NotNil(t, tree)

	expectedRoot := data.HashTreeRoot()
	actualRoot := tree.Hash()
	require.Equal(t, string(expectedRoot[:]), string(actualRoot))
}

func TestSlashingInfo_SetSlot(t *testing.T) {
	t.Parallel()
	data := generateSlashingInfo()

	newSlot := math.Slot(67890)
	data.SetSlot(newSlot)

	require.Equal(t, newSlot, data.GetSlot())
}

func TestSlashingInfo_SetIndex(t *testing.T) {
	t.Parallel()
	data := generateSlashingInfo()

	newIndex := math.U64(12345)
	data.SetIndex(newIndex)

	require.Equal(t, newIndex, data.GetIndex())
}
