/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

// Package goplugin - golang plugin in docker runner
package goplugin

import (
	"context"
	"net/rpc"
	"os/exec"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/rpctypes"
	"github.com/pkg/errors"
)

// Options of the GoPlugin
type Options struct {
	// Listen  is address `GoPlugin` listens on and provides RPC interface for runner(s)
	Listen string
}

// RunnerOptions - set of options to control internal isolated code runner(s)
type RunnerOptions struct {
	// Listen is address the runner listens on and provides RPC interface for the `GoPlugin`
	Listen string
	// CodeStoragePath is path to directory where the runner caches code
	CodeStoragePath string
}

// GoPlugin is a logic runner of code written in golang and compiled as go plugins
type GoPlugin struct {
	Cfg             *configuration.LogicRunner
	MessageBus      core.MessageBus
	ArtifactManager core.ArtifactManager
	runner          *exec.Cmd
	client          *rpc.Client
	clientMutex     sync.Mutex
}

// NewGoPlugin returns a new started GoPlugin
func NewGoPlugin(conf *configuration.LogicRunner, eb core.MessageBus, am core.ArtifactManager) (*GoPlugin, error) {
	gp := GoPlugin{
		Cfg:             conf,
		MessageBus:      eb,
		ArtifactManager: am,
	}

	return &gp, nil
}

// Stop stops runner(s) and RPC service
func (gp *GoPlugin) Stop() error {
	return nil
}

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) Downstream() (*rpc.Client, error) {
	if gp.client != nil {
		return gp.client, nil
	}

	gp.clientMutex.Lock()
	defer gp.clientMutex.Unlock()

	if gp.client != nil {
		return gp.client, nil
	}

	client, err := rpc.Dial(gp.Cfg.GoPlugin.RunnerProtocol, gp.Cfg.GoPlugin.RunnerListen)
	if err != nil {
		return nil, errors.Wrapf(
			err, "couldn't dial '%s' over %s",
			gp.Cfg.GoPlugin.RunnerListen, gp.Cfg.GoPlugin.RunnerProtocol,
		)
	}

	gp.client = client
	return gp.client, nil
}

// Downstream returns a connection to `ginsider`
func (gp *GoPlugin) clientClose() {
	gp.clientMutex.Lock()
	defer gp.clientMutex.Unlock()
	if gp.client == nil {
		return
	}
	gp.client.Close()
	gp.client = nil
}

func (gp *GoPlugin) CallRPC(ctx context.Context, method string, req interface{}, res interface{}) error {
	start := time.Now()
	for {
		client, err := gp.Downstream()
		if err != nil {
			inslogger.FromContext(ctx).Warnf("problem with rpc connection %s", err.Error())
			continue
		}

		select {
		case call := <-client.Go(method, req, res, nil).Done:
			inslogger.FromContext(ctx).Debugf("%s done work, time spend in here - %s", method, time.Since(start))
			if call.Error == rpc.ErrShutdown {
				gp.clientClose()
				inslogger.FromContext(ctx).Warnf("problem with API call %s", call.Error.Error())
			} else {
				return call.Error
			}
		case <-time.After(timeout):
			inslogger.FromContext(ctx).Warnf("problem with API call timeout")
		}
	}
}

const timeout = time.Minute * 10

// CallMethod runs a method on an object in controlled environment
func (gp *GoPlugin) CallMethod(
	ctx context.Context, callContext *core.LogicCallContext,
	code core.RecordRef, data []byte,
	method string, args core.Arguments,
) (
	[]byte, core.Arguments, error,
) {
	res := rpctypes.DownCallMethodResp{}
	req := rpctypes.DownCallMethodReq{
		Context:   callContext,
		Code:      code,
		Data:      data,
		Method:    method,
		Arguments: args,
	}

	err := gp.CallRPC(ctx, "RPC.CallMethod", req, &res)
	return res.Data, res.Ret, err
}

// CallConstructor runs a constructor of a contract in controlled environment
func (gp *GoPlugin) CallConstructor(
	ctx context.Context, callContext *core.LogicCallContext,
	code core.RecordRef, name string, args core.Arguments,
) (
	[]byte, error,
) {
	res := rpctypes.DownCallConstructorResp{}
	req := rpctypes.DownCallConstructorReq{
		Context:   callContext,
		Code:      code,
		Name:      name,
		Arguments: args,
	}

	err := gp.CallRPC(ctx, "RPC.CallMethod", req, &res)
	return res.Ret, err
}
