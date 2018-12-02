// Package controller implements the API controller that handles the connections
// between the clients and the work to do.
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
	"fmt"

	"github.com/buaazp/fasthttprouter"
	"github.com/golang/glog"
	"github.com/valyala/fasthttp"

	"github.com/radar-go/metrics/pkg/config"
)

const jsonContentType = "application/json; charset=utf-8"

// Controller handles all the incoming requests to the API.
type Controller struct {
	Router *fasthttprouter.Router
	cfg    *config.Config
}

// New returns a new Controller object.
func New(cfg *config.Config) *Controller {
	c := &Controller{
		Router: fasthttprouter.New(),
		cfg:    cfg,
	}

	c.register()
	return c
}

// logPaths logs the requested path to the info log.
func logPath(path []byte) {
	glog.Infof("Request path: %s", path)
}

// register defines all the router paths the API implements.
func (c *Controller) register() {
	c.Router.HandleMethodNotAllowed = true
	c.Router.NotFound = c.notFound
	c.Router.MethodNotAllowed = c.methodNotAllowed
	c.Router.PanicHandler = c.panic

	c.Router.GET("/healthcheck", c.healthcheck)
}

// panic handles when the server have a fatal error.
func (c *Controller) panic(ctx *fasthttp.RequestCtx, from interface{}) {
	logPath(ctx.Path())
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	ctx.SetContentType(jsonContentType)
	ctx.SetBodyString(fmt.Sprintf(`{"error": "API fatal error calling %s"}`,
		ctx.Path()))
}

// methodNotAllowed handles the response when a method call is not allowed from
// the client.
func (c *Controller) methodNotAllowed(ctx *fasthttp.RequestCtx) {
	logPath(ctx.Path())
	ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
	ctx.SetContentType(jsonContentType)
	ctx.SetBodyString(fmt.Sprintf(`{"error": "Method not allowed calling %s"}`,
		ctx.Path()))
}

// notFound handles the response when a path have not been found.
func (c *Controller) notFound(ctx *fasthttp.RequestCtx) {
	logPath(ctx.Path())
	ctx.SetStatusCode(fasthttp.StatusNotFound)
	ctx.SetContentType(jsonContentType)
	ctx.SetBodyString(fmt.Sprintf(`{"error": "Path %s not found"}`,
		ctx.Path()))
}

// healthcheck handler.
func (c *Controller) healthcheck(ctx *fasthttp.RequestCtx) {
	logPath(ctx.Path())
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType(jsonContentType)
	ctx.SetBodyString(`{"status": "ok"}`)
}

// internalServerError response
func (c *Controller) internalServerError(ctx *fasthttp.RequestCtx, msg string) {
	ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	c.setError(ctx, msg)
}

// badRequest response.
func (c *Controller) badRequest(ctx *fasthttp.RequestCtx, msg string) {
	ctx.SetStatusCode(fasthttp.StatusBadRequest)
	c.setError(ctx, msg)
}

// unauthorized to call the endpoint response.
func (c *Controller) unauthorized(ctx *fasthttp.RequestCtx, msg string) {
	ctx.SetStatusCode(fasthttp.StatusForbidden)
	c.setError(ctx, msg)
}

// forbidden handles the response when the request is forbidden.
func (c *Controller) forbidden(ctx *fasthttp.RequestCtx, msg string) {
	ctx.SetStatusCode(fasthttp.StatusNotAcceptable)
	c.setError(ctx, msg)
}

// setError sets the error response.
func (c *Controller) setError(ctx *fasthttp.RequestCtx, msg string) {
	ctx.SetContentType(jsonContentType)
	ctx.SetBodyString(fmt.Sprintf(`{"error":"%s"}`, msg))
}
