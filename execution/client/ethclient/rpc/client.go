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

package rpc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/berachain/beacon-kit/primitives/encoding/json"
	beaconhttp "github.com/berachain/beacon-kit/primitives/net/http"
	"github.com/berachain/beacon-kit/primitives/net/jwt"
)

var _ Client = (*client)(nil)

type Client interface {
	Start(context.Context)
	Call(ctx context.Context, target any, method string, params ...any) error
	Close() error
}

// client is an Ethereum RPC client that provides a
// convenient way to interact with an Ethereum node.
type client struct {
	// url is the URL of the RPC endpoint.
	url string
	// client is the HTTP client used to make RPC calls.
	client *http.Client
	// reqPool is a sync.Pool for reusing RPC request objects.
	reqPool *sync.Pool
	// jwtSecret is the JWT secret used for authentication.
	jwtSecret *jwt.Secret
	// jwtRefreshInterval is the interval at which the JWT token should be
	// refreshed.
	jwtRefreshInterval time.Duration

	// mu protects header for concurrent access.
	mu sync.RWMutex

	// header is the HTTP header used for RPC requests.
	header http.Header
}

// New create new rpc client with given url.
func NewClient(
	url string,
	secret *jwt.Secret,
	jwtRefreshInterval time.Duration,
) Client {
	rpc := &client{
		url:    url,
		client: http.DefaultClient,
		reqPool: &sync.Pool{
			New: func() any {
				return &Request{
					ID:      1,
					JSONRPC: "2.0",
				}
			},
		},
		jwtSecret:          secret,
		jwtRefreshInterval: jwtRefreshInterval,
		header:             http.Header{"Content-Type": {"application/json"}},
	}

	return rpc
}

// Start starts the rpc client.
func (rpc *client) Start(ctx context.Context) {
	ticker := time.NewTicker(rpc.jwtRefreshInterval)
	defer ticker.Stop()

	if err := rpc.updateHeader(); err != nil {
		panic(err)
	}
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := rpc.updateHeader(); err != nil {
				// TODO: log or something.
				continue
			}
		}
	}
}

// Close closes the RPC client.
func (rpc *client) Close() error {
	rpc.client.CloseIdleConnections()
	return nil
}

// Call calls the given method with the given parameters.
func (rpc *client) Call(
	ctx context.Context,
	target any,
	method string,
	params ...any,
) error {
	result, err := rpc.callRaw(ctx, method, params...)
	if err != nil {
		return err
	}

	if target == nil {
		return nil
	}

	return json.Unmarshal(result, target)
}

// Call returns raw response of method call.
func (rpc *client) callRaw(
	ctx context.Context,
	method string,
	params ...any,
) (json.RawMessage, error) {
	// Pull a request from the pool, we know that it already has the correct
	// JSONRPC version and ID set.
	//nolint:errcheck // this is safe.
	request := rpc.reqPool.Get().(*Request)
	defer rpc.reqPool.Put(request)

	// Update the request with the method and params.
	request.Method = method
	request.Params = params

	body, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		rpc.url,
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	rpc.mu.RLock()
	req.Header = rpc.header.Clone()
	rpc.mu.RUnlock()

	response, err := rpc.client.Do(req)
	if err != nil {
		return nil, err
	}
	if response == nil {
		return nil, ErrNilResponse
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	switch response.StatusCode {
	case http.StatusOK:
		// OK: just proceed (no return)
	case http.StatusUnauthorized:
		return nil, beaconhttp.ErrUnauthorized
	default:
		// Return a default error
		return nil, fmt.Errorf("unexpected status code %d: %s", response.StatusCode, string(data))
	}

	resp := new(Response)
	if err = json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil
}
