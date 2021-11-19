// Copyright 2019-present Open Networking Foundation.
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

package value

import (
	"context"
	api "github.com/atomix/atomix-api/go/atomix/primitive/value/v1"
	"github.com/atomix/atomix-go-client/pkg/atomix/primitive"
	"github.com/atomix/atomix-sdk-go/pkg/errors"
	"github.com/atomix/atomix-sdk-go/pkg/logging"
	"github.com/atomix/atomix-sdk-go/pkg/meta"
	"google.golang.org/grpc"
	"io"
)

var log = logging.GetLogger("atomix", "client", "value")

// Type is the value type
const Type primitive.Type = "Value"

// Client provides an API for creating Values
type Client interface {
	// GetValue gets the Value instance of the given name
	GetValue(ctx context.Context, name string, opts ...Option) (Value, error)
}

// Value provides a simple atomic value
type Value interface {
	primitive.Primitive

	// Set sets the current value and returns the version
	Set(ctx context.Context, value []byte, opts ...SetOption) (meta.ObjectMeta, error)

	// Get gets the current value and version
	Get(ctx context.Context) ([]byte, meta.ObjectMeta, error)

	// Watch watches the value for changes
	Watch(ctx context.Context, ch chan<- Event) error
}

// EventType is the type of a set event
type EventType string

const (
	// EventUpdate indicates the value was updated
	EventUpdate EventType = "update"
)

// Event is a value change event
type Event struct {
	meta.ObjectMeta

	// Type is the change event type
	Type EventType

	// Value is the updated value
	Value []byte
}

// New creates a new Lock primitive for the given partitions
// The value will be created in one of the given partitions.
func New(ctx context.Context, name string, conn *grpc.ClientConn, opts ...Option) (Value, error) {
	options := newValueOptions{}
	for _, opt := range opts {
		if op, ok := opt.(Option); ok {
			op.applyNewValue(&options)
		}
	}
	sessions := api.NewValueManagerClient(conn)
	request := &api.OpenSessionRequest{
		Options: options.sessionOptions,
	}
	response, err := sessions.OpenSession(ctx, request)
	if err != nil {
		return nil, errors.From(err)
	}
	return &value{
		Client:  primitive.NewClient(Type, name, response.SessionID),
		client:  api.NewValueClient(conn),
		session: sessions,
	}, nil
}

// value is the single partition implementation of Lock
type value struct {
	*primitive.Client
	client  api.ValueClient
	session api.ValueManagerClient
}

func (v *value) Set(ctx context.Context, value []byte, opts ...SetOption) (meta.ObjectMeta, error) {
	request := &api.SetRequest{
		Value: api.Object{
			Value: value,
		},
	}
	for i := range opts {
		opts[i].beforeSet(request)
	}
	response, err := v.client.Set(v.GetContext(ctx), request)
	if err != nil {
		return meta.ObjectMeta{}, errors.From(err)
	}
	for i := range opts {
		opts[i].afterSet(response)
	}
	return meta.FromProto(response.Value.ObjectMeta), nil
}

func (v *value) Get(ctx context.Context) ([]byte, meta.ObjectMeta, error) {
	request := &api.GetRequest{}
	response, err := v.client.Get(v.GetContext(ctx), request)
	if err != nil {
		return nil, meta.ObjectMeta{}, errors.From(err)
	}
	return response.Value.Value, meta.FromProto(response.Value.ObjectMeta), nil
}

func (v *value) Watch(ctx context.Context, ch chan<- Event) error {
	request := &api.EventsRequest{}
	stream, err := v.client.Events(v.GetContext(ctx), request)
	if err != nil {
		return errors.From(err)
	}

	openCh := make(chan struct{})
	go func() {
		defer close(ch)
		open := false
		defer func() {
			if !open {
				close(openCh)
			}
		}()
		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				err = errors.From(err)
				if errors.IsCanceled(err) || errors.IsTimeout(err) {
					return
				}
				log.Errorf("Watch failed: %v", err)
				return
			}

			if !open {
				close(openCh)
				open = true
			}
			switch response.Event.Type {
			case api.Event_UPDATE:
				ch <- Event{
					ObjectMeta: meta.FromProto(response.Event.Value.ObjectMeta),
					Type:       EventUpdate,
					Value:      response.Event.Value.Value,
				}
			}
		}
	}()

	select {
	case <-openCh:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (v *value) Close(ctx context.Context) error {
	request := &api.CloseSessionRequest{
		SessionID: v.SessionID(),
	}
	_, err := v.session.CloseSession(ctx, request)
	if err != nil {
		return errors.From(err)
	}
	return nil
}
