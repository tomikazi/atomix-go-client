// Copyright 2020-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package atomix

import (
	"context"
	primitiveapi "github.com/atomix/atomix-api/go/atomix/primitive"
	"github.com/atomix/atomix-go-client/pkg/atomix/counter"
	"github.com/atomix/atomix-go-client/pkg/atomix/election"
	"github.com/atomix/atomix-go-client/pkg/atomix/indexedmap"
	"github.com/atomix/atomix-go-client/pkg/atomix/list"
	"github.com/atomix/atomix-go-client/pkg/atomix/lock"
	_map "github.com/atomix/atomix-go-client/pkg/atomix/map"
	"github.com/atomix/atomix-go-client/pkg/atomix/primitive"
	"github.com/atomix/atomix-go-client/pkg/atomix/set"
	"github.com/atomix/atomix-go-client/pkg/atomix/value"
)

// GetCounter gets the Counter instance of the given name
func GetCounter(ctx context.Context, name string, opts ...primitive.Option) (counter.Counter, error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return counter.New(ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetElection gets the Election instance of the given name
func GetElection(ctx context.Context, name string, opts ...primitive.Option) (election.Election, error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return election.New(ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetIndexedMap gets the IndexedMap instance of the given name
func GetIndexedMap[K, V any](ctx context.Context, name string, opts ...primitive.Option) (indexedmap.IndexedMap[K, V], error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return indexedmap.New[K, V](ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetList gets the List instance of the given name
func GetList[E any](ctx context.Context, name string, opts ...primitive.Option) (list.List[E], error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return list.New[E](ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetLock gets the Lock instance of the given name
func GetLock(ctx context.Context, name string, opts ...primitive.Option) (lock.Lock, error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return lock.New(ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetMap gets the Map instance of the given name
func GetMap[K, V any](ctx context.Context, name string, opts ...primitive.Option) (_map.Map[K, V], error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return _map.New[K, V](ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetSet gets the Set instance of the given name
func GetSet[E any](ctx context.Context, name string, opts ...primitive.Option) (set.Set[E], error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return set.New[E](ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

// GetValue gets the Value instance of the given name
func GetValue[V any](ctx context.Context, name string, opts ...primitive.Option) (value.Value[V], error) {
	client := getClient()
	conn, err := client.connect(ctx, newPrimitiveID(counter.Type, name))
	if err != nil {
		return nil, err
	}
	return value.New[V](ctx, name, conn, getPrimitiveOpts(client.options, opts...)...)
}

func newPrimitiveID(t primitive.Type, name string) primitiveapi.PrimitiveId {
	return primitiveapi.PrimitiveId{
		Type: t.String(),
		Name: name,
	}
}

func getPrimitiveOpts(clientOpts clientOptions, primitiveOpts ...primitive.Option) []primitive.Option {
	return append([]primitive.Option{primitive.WithSessionID(clientOpts.clientID)}, primitiveOpts...)
}
