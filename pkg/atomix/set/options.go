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

package set

import (
	api "github.com/atomix/atomix-api/go/atomix/primitive/set"
	"github.com/atomix/atomix-go-client/pkg/atomix/primitive"
	"github.com/atomix/atomix-go-client/pkg/atomix/primitive/codec"
)

// Option is a set option
type Option[E any] interface {
	primitive.Option
	applyNewSet(options *newSetOptions[E])
}

// newSetOptions is set options
type newSetOptions[E any] struct {
	elementCodec codec.Codec[E]
}

func WithCodec[E any](elementCodec codec.Codec[E]) Option[E] {
	return codecOption[E]{
		elementCodec: elementCodec,
	}
}

type codecOption[E any] struct {
	primitive.EmptyOption
	elementCodec codec.Codec[E]
}

func (o codecOption[E]) applyNewSet(options *newSetOptions[E]) {
	options.elementCodec = o.elementCodec
}

// WatchOption is an option for set Watch calls
type WatchOption interface {
	beforeWatch(request *api.EventsRequest)
	afterWatch(response *api.EventsResponse)
}

// WithReplay returns a Watch option to replay entries
func WithReplay() WatchOption {
	return replayOption{}
}

type replayOption struct{}

func (o replayOption) beforeWatch(request *api.EventsRequest) {
	request.Replay = true
}

func (o replayOption) afterWatch(response *api.EventsResponse) {

}
