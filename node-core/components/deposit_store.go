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

package components

import (
	"path/filepath"

	"cosmossdk.io/depinject"
	"github.com/berachain/beacon-kit/config"
	"github.com/berachain/beacon-kit/log/phuslu"
	"github.com/berachain/beacon-kit/storage/deposit"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cast"
)

// DepositStoreInput is the input for the dep inject framework.
type DepositStoreInput struct {
	depinject.In
	Logger  *phuslu.Logger
	AppOpts config.AppOptions
}

// ProvideDepositStore is a function that provides the module to the
// application.
func ProvideDepositStore(in DepositStoreInput) (deposit.StoreManager, error) {
	var (
		rootDir = cast.ToString(in.AppOpts.Get(flags.FlagHome))
		dataDir = filepath.Join(rootDir, "data")
		nameV1  = "deposits"
	)

	dbV1, err := dbm.NewDB(nameV1, dbm.PebbleDBBackend, dataDir)
	if err != nil {
		return nil, err
	}

	return deposit.NewStore(
		dbV1,
		in.Logger.With("service", "deposit-store"),
	), nil
}
