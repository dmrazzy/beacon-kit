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

package beacon

import (
	"fmt"

	"github.com/berachain/beacon-kit/node-api/handlers"
	beacontypes "github.com/berachain/beacon-kit/node-api/handlers/beacon/types"
	"github.com/berachain/beacon-kit/node-api/handlers/types"
	"github.com/berachain/beacon-kit/node-api/handlers/utils"
	"github.com/berachain/beacon-kit/primitives/version"
)

func (h *Handler) GetPendingPartialWithdrawals(c handlers.Context) (any, error) {
	req, err := utils.BindAndValidate[beacontypes.GetPendingPartialWithdrawalsRequest](
		c, h.Logger(),
	)
	if err != nil {
		return nil, err
	}

	slot, err := utils.SlotFromStateID(req.StateID, h.backend)
	if err != nil {
		return nil, err
	}

	st, _, err := h.backend.StateAtSlot(slot)
	if err != nil {
		return nil, err
	}

	// Get the fork version.
	forkVersion, err := st.GetFork()
	if err != nil {
		return nil, err
	}

	if version.IsBefore(forkVersion.CurrentVersion, version.Electra()) {
		return nil, fmt.Errorf("%w: Electra fork not active yet", types.ErrInvalidRequest)
	}

	// Get the pending partial withdrawals from the state.
	partialWithdrawals, err := h.backend.PendingPartialWithdrawalsAtState(st)
	if err != nil {
		return nil, err
	}

	return beacontypes.NewPendingPartialWithdrawalsResponse(
		forkVersion.CurrentVersion,
		partialWithdrawals,
	), nil
}
