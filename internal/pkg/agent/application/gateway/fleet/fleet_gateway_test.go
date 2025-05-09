// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License 2.0;
// you may not use this file except in compliance with the Elastic License 2.0.

package fleet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/elastic/elastic-agent/internal/pkg/agent/application/coordinator"
	"github.com/elastic/elastic-agent/internal/pkg/agent/application/upgrade/details"
	"github.com/elastic/elastic-agent/internal/pkg/agent/errors"
	"github.com/elastic/elastic-agent/internal/pkg/agent/storage"
	"github.com/elastic/elastic-agent/internal/pkg/agent/storage/store"
	"github.com/elastic/elastic-agent/internal/pkg/fleetapi"
	"github.com/elastic/elastic-agent/internal/pkg/fleetapi/acker/noop"
	"github.com/elastic/elastic-agent/internal/pkg/scheduler"
	agentclient "github.com/elastic/elastic-agent/pkg/control/v2/client"
	"github.com/elastic/elastic-agent/pkg/core/logger"
	"github.com/elastic/elastic-agent/pkg/core/logger/loggertest"
)

type clientCallbackFunc func(headers http.Header, body io.Reader) (*http.Response, error)

type testingClient struct {
	sync.Mutex
	callback clientCallbackFunc
	received chan struct{}
}

func (t *testingClient) Send(
	_ context.Context,
	_ string,
	_ string,
	_ url.Values,
	headers http.Header,
	body io.Reader,
) (*http.Response, error) {
	t.Lock()
	defer t.Unlock()
	defer func() { t.received <- struct{}{} }()
	return t.callback(headers, body)
}

func (t *testingClient) URI() string {
	return "http://localhost"
}

func (t *testingClient) Answer(fn clientCallbackFunc) <-chan struct{} {
	t.Lock()
	defer t.Unlock()
	t.callback = fn
	return t.received
}

func newTestingClient() *testingClient {
	return &testingClient{received: make(chan struct{}, 1)}
}

type withGatewayFunc func(*testing.T, coordinator.FleetGateway, *testingClient, *scheduler.Stepper)

func withGateway(agentInfo agentInfo, settings *fleetGatewaySettings, fn withGatewayFunc) func(t *testing.T) {
	return func(t *testing.T) {
		scheduler := scheduler.NewStepper()
		client := newTestingClient()

		log, _ := logger.New("fleet_gateway", false)

		stateStore := newStateStore(t, log)

		gateway, err := newFleetGatewayWithScheduler(
			log,
			settings,
			agentInfo,
			client,
			scheduler,
			noop.New(),
			emptyStateFetcher,
			stateStore,
		)

		require.NoError(t, err)

		fn(t, gateway, client, scheduler)
	}
}

func ackSeq(channels ...<-chan struct{}) func() {
	return func() {
		for _, c := range channels {
			<-c
		}
	}
}

func wrapStrToResp(code int, body string) *http.Response {
	return &http.Response{
		Status:        fmt.Sprintf("%d %s", code, http.StatusText(code)),
		StatusCode:    code,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Body:          io.NopCloser(bytes.NewBufferString(body)),
		ContentLength: int64(len(body)),
		Header:        make(http.Header),
	}
}

func TestFleetGateway(t *testing.T) {
	agentInfo := &testAgentInfo{}
	settings := &fleetGatewaySettings{
		Duration: 5 * time.Second,
		Backoff:  &backoffSettings{Init: 1 * time.Second, Max: 5 * time.Second},
	}

	t.Run("send no event and receive no action", withGateway(agentInfo, settings, func(
		t *testing.T,
		gateway coordinator.FleetGateway,
		client *testingClient,
		scheduler *scheduler.Stepper,
	) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		waitFn := ackSeq(
			client.Answer(func(headers http.Header, body io.Reader) (*http.Response, error) {
				resp := wrapStrToResp(http.StatusOK, `{ "actions": [] }`)
				return resp, nil
			}),
		)

		errCh := runFleetGateway(ctx, gateway)

		// Synchronize scheduler and acking of calls from the worker go routine.
		scheduler.Next()
		waitFn()

		cancel()
		err := <-errCh
		require.NoError(t, err)
		select {
		case actions := <-gateway.Actions():
			t.Errorf("Expected no actions, got %v", actions)
		default:
		}
	}))

	t.Run("Successfully connects and receives a series of actions", withGateway(agentInfo, settings, func(
		t *testing.T,
		gateway coordinator.FleetGateway,
		client *testingClient,
		scheduler *scheduler.Stepper,
	) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		waitFn := ackSeq(
			client.Answer(func(headers http.Header, body io.Reader) (*http.Response, error) {
				// TODO: assert no events
				resp := wrapStrToResp(http.StatusOK, `
	{
		"actions": [
			{
				"type": "POLICY_CHANGE",
				"id": "id1",
				"data": {
					"policy": {
						"id": "policy-id"
					}
				}
			},
			{
				"type": "ANOTHER_ACTION",
				"id": "id2"
			}
		]
	}
	`)
				return resp, nil
			}),
		)

		errCh := runFleetGateway(ctx, gateway)

		scheduler.Next()
		waitFn()

		cancel()
		err := <-errCh
		require.NoError(t, err)
		select {
		case actions := <-gateway.Actions():
			require.Len(t, actions, 2)
		default:
			t.Errorf("Expected to receive actions")
		}
	}))

	// Test the normal time based execution.
	t.Run("Periodically communicates with Fleet", func(t *testing.T) {
		scheduler := scheduler.NewPeriodic(150 * time.Millisecond)
		client := newTestingClient()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		log, _ := logger.New("tst", false)
		stateStore := newStateStore(t, log)

		gateway, err := newFleetGatewayWithScheduler(
			log,
			settings,
			agentInfo,
			client,
			scheduler,
			noop.New(),
			emptyStateFetcher,
			stateStore,
		)
		require.NoError(t, err)

		waitFn := ackSeq(
			client.Answer(func(headers http.Header, body io.Reader) (*http.Response, error) {
				resp := wrapStrToResp(http.StatusOK, `{ "actions": [] }`)
				return resp, nil
			}),
		)

		errCh := runFleetGateway(ctx, gateway)

		func() {
			var count int
			for {
				waitFn()
				count++
				if count == 4 {
					return
				}
			}
		}()

		cancel()
		err = <-errCh
		require.NoError(t, err)
	})

	t.Run("Test the wait loop is interruptible", func(t *testing.T) {
		// 20mins is the double of the base timeout values for golang test suites.
		// If we cannot interrupt we will timeout.
		d := 20 * time.Minute
		scheduler := scheduler.NewPeriodic(d)
		client := newTestingClient()

		ctx, cancel := context.WithCancel(context.Background())

		log, _ := logger.New("tst", false)
		stateStore := newStateStore(t, log)

		gateway, err := newFleetGatewayWithScheduler(
			log,
			&fleetGatewaySettings{
				Duration: d,
				Backoff:  &backoffSettings{Init: 1 * time.Second, Max: 30 * time.Second},
			},
			agentInfo,
			client,
			scheduler,
			noop.New(),
			emptyStateFetcher,
			stateStore,
		)
		require.NoError(t, err)

		ch2 := client.Answer(func(headers http.Header, body io.Reader) (*http.Response, error) {
			resp := wrapStrToResp(http.StatusOK, `{ "actions": [] }`)
			return resp, nil
		})

		errCh := runFleetGateway(ctx, gateway)

		// Make sure that all API calls to the checkin API are successful, the following will happen:
		// block on the first call.
		<-ch2

		go func() {
			// drain the channel
			for range ch2 {
			}
		}()

		// 1. Gateway will check the API on boot.
		// 2. WaitTick() will block for 20 minutes.
		// 3. Stop will should unblock the wait.
		cancel()
		err = <-errCh
		require.NoError(t, err)
	})

	t.Run("Sends upgrade details", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		scheduler := scheduler.NewStepper()
		client := newTestingClient()

		log, _ := loggertest.New("fleet_gateway")

		stateStore := newStateStore(t, log)

		upgradeDetails := &details.Details{
			TargetVersion: "8.12.0",
			State:         "UPG_WATCHING",
			ActionID:      "foobarbaz",
		}
		stateFetcher := func() coordinator.State {
			return coordinator.State{
				UpgradeDetails: upgradeDetails,
			}
		}

		gateway, err := newFleetGatewayWithScheduler(
			log,
			settings,
			agentInfo,
			client,
			scheduler,
			noop.New(),
			stateFetcher,
			stateStore,
		)

		require.NoError(t, err)

		waitFn := ackSeq(
			client.Answer(func(headers http.Header, body io.Reader) (*http.Response, error) {
				data, err := io.ReadAll(body)
				require.NoError(t, err)

				var checkinRequest fleetapi.CheckinRequest
				err = json.Unmarshal(data, &checkinRequest)
				require.NoError(t, err)

				require.NotNil(t, checkinRequest.UpgradeDetails)
				require.Equal(t, upgradeDetails, checkinRequest.UpgradeDetails)

				resp := wrapStrToResp(http.StatusOK, `{ "actions": [] }`)
				return resp, nil
			}),
		)

		errCh := runFleetGateway(ctx, gateway)

		// Synchronize scheduler and acking of calls from the worker go routine.
		scheduler.Next()
		waitFn()

		cancel()
		err = <-errCh
		require.NoError(t, err)
		select {
		case actions := <-gateway.Actions():
			t.Errorf("Expected no actions, got %v", actions)
		default:
		}
	})
}

func TestRetriesOnFailures(t *testing.T) {
	agentInfo := &testAgentInfo{}
	settings := &fleetGatewaySettings{
		Duration: 5 * time.Second,
		Backoff:  &backoffSettings{Init: 100 * time.Millisecond, Max: 5 * time.Second},
	}

	t.Run("When the gateway fails to communicate with the checkin API we will retry",
		withGateway(agentInfo, settings, func(
			t *testing.T,
			gateway coordinator.FleetGateway,
			client *testingClient,
			scheduler *scheduler.Stepper,
		) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			fail := func(_ http.Header, _ io.Reader) (*http.Response, error) {
				return wrapStrToResp(http.StatusInternalServerError, "something is bad"), nil
			}
			clientWaitFn := client.Answer(fail)

			errCh := runFleetGateway(ctx, gateway)

			// Initial tick is done out of bound so we can block on channels.
			scheduler.Next()

			// Simulate a 500 errors for the next 3 calls.
			<-clientWaitFn
			<-clientWaitFn
			<-clientWaitFn

			// API recover
			waitFn := ackSeq(
				client.Answer(func(_ http.Header, body io.Reader) (*http.Response, error) {
					resp := wrapStrToResp(http.StatusOK, `{ "actions": [] }`)
					return resp, nil
				}),
			)

			waitFn()

			cancel()
			err := <-errCh
			require.NoError(t, err)
			select {
			case actions := <-gateway.Actions():
				t.Errorf("Expected no actions, got %v", actions)
			default:
			}
		}))

	t.Run("The retry loop is interruptible",
		withGateway(agentInfo, &fleetGatewaySettings{
			Duration: 0 * time.Second,
			Backoff:  &backoffSettings{Init: 10 * time.Minute, Max: 20 * time.Minute},
		}, func(
			t *testing.T,
			gateway coordinator.FleetGateway,
			client *testingClient,
			scheduler *scheduler.Stepper,
		) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			fail := func(_ http.Header, _ io.Reader) (*http.Response, error) {
				return wrapStrToResp(http.StatusInternalServerError, "something is bad"), nil
			}
			waitChan := client.Answer(fail)

			errCh := runFleetGateway(ctx, gateway)

			// Initial tick is done out of bound so we can block on channels.
			scheduler.Next()

			// Fail to enter retry loop, all other calls will fails and will force to wait on big initial
			// delay.
			<-waitChan

			cancel()
			err := <-errCh
			require.NoError(t, err)
		}))
}

type testAgentInfo struct{}

func (testAgentInfo) AgentID() string { return "agent-secret" }

func emptyStateFetcher() coordinator.State {
	return coordinator.State{}
}

func runFleetGateway(ctx context.Context, g coordinator.FleetGateway) <-chan error {
	done := make(chan bool)
	errCh := make(chan error, 1)
	go func() {
		err := g.Run(ctx)
		close(done)
		if err != nil && !errors.Is(err, context.Canceled) {
			errCh <- err
		} else {
			errCh <- nil
		}
	}()
	go func() {
		for {
			select {
			case <-done:
				return
			case <-g.Errors():
				// ignore errors here
			}
		}
	}()
	return errCh
}

func newStateStore(t *testing.T, log *logger.Logger) *store.StateStore {
	dir := t.TempDir()

	filename := filepath.Join(dir, "state.enc")
	diskStore, err := storage.NewDiskStore(filename)
	require.NoError(t, err)
	stateStore, err := store.NewStateStore(log, diskStore)
	require.NoError(t, err)

	return stateStore
}

func TestAgentStateToString(t *testing.T) {
	testcases := []struct {
		agentState         agentclient.State
		expectedFleetState string
	}{
		{
			agentState:         agentclient.Healthy,
			expectedFleetState: fleetStateOnline,
		},
		{
			agentState:         agentclient.Failed,
			expectedFleetState: fleetStateError,
		},
		{
			agentState:         agentclient.Starting,
			expectedFleetState: fleetStateStarting,
		},
		// everything else maps to degraded
		{
			agentState:         agentclient.Configuring,
			expectedFleetState: fleetStateOnline,
		},
		{
			agentState:         agentclient.Degraded,
			expectedFleetState: fleetStateDegraded,
		},
		{
			agentState:         agentclient.Stopping,
			expectedFleetState: fleetStateOnline,
		},
		{
			agentState:         agentclient.Stopped,
			expectedFleetState: fleetStateOnline,
		},
		{
			agentState:         agentclient.Upgrading,
			expectedFleetState: fleetStateOnline,
		},
		{
			agentState:         agentclient.Rollback,
			expectedFleetState: fleetStateDegraded,
		},
		{
			// Unknown states should map to degraded.
			agentState:         agentclient.Rollback + 1,
			expectedFleetState: fleetStateDegraded,
		},
	}

	for _, tc := range testcases {
		t.Run(fmt.Sprintf("%s -> %s", tc.agentState, tc.expectedFleetState), func(t *testing.T) {
			actualFleetState := agentStateToString(tc.agentState)
			assert.Equal(t, tc.expectedFleetState, actualFleetState)
		})
	}
}
