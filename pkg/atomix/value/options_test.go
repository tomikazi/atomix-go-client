// SPDX-FileCopyrightText: 2019-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package value

import (
	api "github.com/atomix/atomix-api/go/atomix/primitive/value"
	"github.com/atomix/atomix-go-framework/pkg/atomix/meta"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOptions(t *testing.T) {
	request := &api.SetRequest{}
	IfMatch(meta.ObjectMeta{Revision: 1}).beforeSet(request)
	assert.Equal(t, meta.Revision(1), meta.Revision(request.Preconditions[0].GetMetadata().Revision.Num))
}
