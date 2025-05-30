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

package handlers

import (
	"github.com/berachain/beacon-kit/log"
)

// Route is a route for the node API.
type Route struct {
	Method  string
	Path    string
	Handler handlerFn
}

// DecorateWithLogs adds logging to the route's handler function as soon as
// a request is received and when a response is ready.
func (r *Route) DecorateWithLogs(logger log.Logger) {
	handler := r.Handler
	r.Handler = func(ctx Context) (any, error) {
		logger.Info("received request", "method", r.Method, "path", r.Path)
		res, err := handler(ctx)
		if err != nil {
			logger.Error("error handling request", "error", err)
		}
		logger.Info("request handled")
		return res, err
	}
}

// RouteSet is a set of routes for the node API.
type RouteSet struct {
	BasePath string
	Routes   []*Route
}

// NewRouteSet creates a new route set.
func NewRouteSet(basePath string, routes ...*Route) *RouteSet {
	return &RouteSet{
		BasePath: basePath,
		Routes:   routes,
	}
}
