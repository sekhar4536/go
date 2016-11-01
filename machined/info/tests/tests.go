// Copyright 2015-2016 Platina Systems, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style license described in the
// LICENSE file.

package tests

import (
	"github.com/platinasystems/go/machined/info"
	"github.com/platinasystems/go/machined/info/tests/test_string"
)

func New() []info.Interface {
	return []info.Interface{
		test_string.New(),
	}
}