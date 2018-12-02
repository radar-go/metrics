package controller

/* Copyright (C) 2018 Radar team (see AUTHORS)

   This file is part of metrics.

   metrics is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   metrics is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with metrics. If not, see <http://www.gnu.org/licenses/>.
*/

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/radar-go/metrics/helper"
	"github.com/radar-go/metrics/pkg/config"
)

const (
	unexpectedControllerResponse = "Unexpected controller response"
	unexpectedStatusCode         = "Unexpected status code"
)

func TestController(t *testing.T) {
	cfg, err := config.New("./testdata")
	assert.NoError(t, err, "Unexpected error reading the configs.")
	assert.NotNil(t, cfg, "Unexpected error creating the configs object.")

	ctx, c := setup(t, cfg)
	c.panic(ctx, "test")
	assert.Equal(t, 500, ctx.Response.StatusCode(), unexpectedStatusCode)
	assert.Equal(t, []byte(`{"error": "API fatal error calling /"}`),
		ctx.Response.Body(), unexpectedControllerResponse)

	c.methodNotAllowed(ctx)
	assert.Equal(t, 405, ctx.Response.StatusCode(), unexpectedStatusCode)
	assert.Equal(t, []byte(`{"error": "Method not allowed calling /"}`),
		ctx.Response.Body(), unexpectedControllerResponse)

	c.notFound(ctx)
	assert.Equal(t, 404, ctx.Response.StatusCode(), unexpectedStatusCode)
	assert.Equal(t, []byte(`{"error": "Path / not found"}`),
		ctx.Response.Body(), unexpectedControllerResponse)

	c.internalServerError(ctx, "Internal server error")
	assert.Equal(t, 500, ctx.Response.StatusCode(), unexpectedStatusCode)
	assert.Equal(t, []byte(`{"error":"Internal server error"}`),
		ctx.Response.Body(), unexpectedControllerResponse)

	c.badRequest(ctx, "Bad request")
	assert.Equal(t, 400, ctx.Response.StatusCode(), unexpectedStatusCode)
	assert.Equal(t, []byte(`{"error":"Bad request"}`),
		ctx.Response.Body(), unexpectedControllerResponse)

	c.unauthorized(ctx, "User not authorized to do this request")
	assert.Equal(t, 403, ctx.Response.StatusCode(), unexpectedControllerResponse)
	assert.Equal(t, []byte(`{"error":"User not authorized to do this request"}`),
		ctx.Response.Body(), unexpectedControllerResponse)

	c.forbidden(ctx, "Request forbidden")
	assert.Equal(t, 406, ctx.Response.StatusCode(), unexpectedStatusCode)
	assert.Equal(t, []byte(`{"error":"Request forbidden"}`),
		ctx.Response.Body(), unexpectedControllerResponse)
}

func TestEndpoints(t *testing.T) {
	tests := map[string]struct {
		path        string
		code        int
		method      string
		data        string
		contentType string
		headers     map[string]string
	}{
		"healthcheck": {
			path:   "/healthcheck",
			code:   200,
			method: "GET",
		},
	}

	cfg, err := config.New("./testdata")
	assert.NoError(t, err, "Unexpected error reading the configs.")
	assert.NotNil(t, cfg, "Unexpected error creating the configs object.")

	ctx, c := setup(t, cfg)
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ctx.Request.Reset()
			ctx.Request.SetRequestURI(tc.path)
			ctx.Request.Header.SetMethod(tc.method)
			ctx.Request.Header.SetContentType(tc.contentType)
			if len(tc.headers) > 0 {
				for k, v := range tc.headers {
					ctx.Request.Header.Add(k, v)
				}
			}

			if len(tc.data) > 0 {
				ctx.Request.SetBodyString(tc.data)
			}

			c.Router.Handler(ctx)
			if tc.code != ctx.Response.StatusCode() {
				t.Errorf("Expected %d, Got %d", tc.code, ctx.Response.StatusCode())
			}

			helper.SaveGoldenData(t, name, ctx.Response.Body())
			expected := helper.GetGoldenData(t, name)
			assert.Equal(t, expected, ctx.Response.Body(), "Unexpected response")
		})
	}
}

func setup(t *testing.T, cfg *config.Config) (*fasthttp.RequestCtx, *Controller) {
	t.Helper()
	ctx := &fasthttp.RequestCtx{}
	c := New(cfg)

	return ctx, c
}
