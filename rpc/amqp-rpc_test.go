// Copyright 2014 ISRG.  All rights reserved
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package rpc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/letsencrypt/boulder/core"
	"github.com/letsencrypt/boulder/probs"
	"github.com/letsencrypt/boulder/test"
)

func TestWrapError(t *testing.T) {
	testCases := []error{
		core.InternalServerError("foo"),
		core.NotSupportedError("foo"),
		core.MalformedRequestError("foo"),
		core.UnauthorizedError("foo"),
		core.NotFoundError("foo"),
		core.SyntaxError("foo"),
		core.SignatureValidationError("foo"),
		core.CertificateIssuanceError("foo"),
		core.NoSuchRegistrationError("foo"),
		core.RateLimitedError("foo"),
		core.TooManyRPCRequestsError("foo"),
		errors.New("foo"),
	}
	for _, c := range testCases {
		wrapped := wrapError(c)
		test.AssertEquals(t, wrapped.Type, reflect.TypeOf(c).Name())
		test.AssertEquals(t, wrapped.Value, "foo")
		unwrapped := unwrapError(wrapped)
		test.AssertEquals(t, wrapped.Type, reflect.TypeOf(unwrapped).Name())
		test.AssertEquals(t, unwrapped.Error(), "foo")
	}

	complicated := []struct {
		given    error
		expected error
	}{
		{
			&probs.ProblemDetails{
				Type:       probs.ConnectionProblem,
				Detail:     "whoops",
				HTTPStatus: 417,
			},
			&probs.ProblemDetails{
				Type:       probs.ConnectionProblem,
				Detail:     "whoops",
				HTTPStatus: 417,
			},
		},
		{
			&probs.ProblemDetails{Type: "invalid", Detail: "hm"},
			errors.New("hm"),
		},
		{
			errors.New(""),
			errors.New(""),
		},
	}
	for i, tc := range complicated {
		actual := unwrapError(wrapError(tc.given))
		if !reflect.DeepEqual(tc.expected, actual) {
			t.Errorf("rpc error wrapping case %d: want %#v, got %#v", i, tc.expected, actual)
		}

	}
}
