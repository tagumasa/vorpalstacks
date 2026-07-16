package testutil

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iot"
	"github.com/aws/aws-sdk-go-v2/service/iotdataplane"
	"vorpalstacks-sdk-tests/config"
)

// runIoTShadowTests exercises the Thing Shadow data plane (GetThingShadow /
// UpdateThingShadow / DeleteThingShadow / ListNamedShadowsForThing) via the
// iotdataplane SDK client. The shadow handlers are registered under the iot
// service and routed via the /things/{thingName}/shadow REST paths (Unit 3).
func (r *TestRunner) runIoTShadowTests(tc *iotTestContext) []TestResult {
	cfg, err := config.LoadDefaultAWSConfig(config.AWSConfig{
		Endpoint: r.endpoint,
		Region:   r.region,
	})
	if err != nil {
		return []TestResult{{Service: "iot", TestName: "Shadow_Setup", Status: "FAIL", Error: err.Error()}}
	}
	dataClient := iotdataplane.NewFromConfig(cfg)
	ctx := context.Background()

	thingName := uniqueName("shadow-thing")
	defer tc.client.DeleteThing(ctx, &iot.DeleteThingInput{ThingName: aws.String(thingName)})

	if _, err := tc.client.CreateThing(ctx, &iot.CreateThingInput{ThingName: aws.String(thingName)}); err != nil {
		return []TestResult{{Service: "iot", TestName: "Shadow_Setup", Status: "FAIL", Error: fmt.Sprintf("CreateThing failed: %v", err)}}
	}

	var results []TestResult
	add := func(name string, fn func() error) {
		results = append(results, r.RunTest("iot", name, fn))
	}

	// Classic shadow: Update (desired) -> Get (verify version/state) -> Update
	// (reported, optimistic version) -> Delete -> Get NotFound.
	add("Shadow_UpdateThingShadow_Classic", func() error {
		payload, _ := json.Marshal(map[string]interface{}{
			"state": map[string]interface{}{
				"desired": map[string]interface{}{"color": "green"},
			},
		})
		out, err := dataClient.UpdateThingShadow(ctx, &iotdataplane.UpdateThingShadowInput{
			ThingName: aws.String(thingName), Payload: payload,
		})
		if err != nil {
			return fmt.Errorf("UpdateThingShadow failed: %w", err)
		}
		if out.Payload == nil {
			return fmt.Errorf("expected non-nil payload")
		}
		return nil
	})

	add("Shadow_GetThingShadow_Classic", func() error {
		out, err := dataClient.GetThingShadow(ctx, &iotdataplane.GetThingShadowInput{
			ThingName: aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("GetThingShadow failed: %w", err)
		}
		var state struct {
			State struct {
				Desired map[string]interface{} `json:"desired"`
			} `json:"state"`
			Version *int64 `json:"version"`
		}
		if err := json.Unmarshal(out.Payload, &state); err != nil {
			return fmt.Errorf("decode shadow payload failed: %w", err)
		}
		if state.Version == nil || *state.Version < 1 {
			return fmt.Errorf("expected version >= 1, got %v", state.Version)
		}
		if state.State.Desired["color"] != "green" {
			return fmt.Errorf("expected desired.color=green, got %v", state.State.Desired["color"])
		}
		return nil
	})

	add("Shadow_UpdateThingShadow_Reported", func() error {
		payload, _ := json.Marshal(map[string]interface{}{
			"state": map[string]interface{}{
				"reported": map[string]interface{}{"color": "blue"},
			},
		})
		if _, err := dataClient.UpdateThingShadow(ctx, &iotdataplane.UpdateThingShadowInput{
			ThingName: aws.String(thingName), Payload: payload,
		}); err != nil {
			return fmt.Errorf("UpdateThingShadow (reported) failed: %w", err)
		}
		return nil
	})

	// Named shadow lifecycle + ListNamedShadowsForThing.
	namedShadow := "config"
	add("Shadow_UpdateThingShadow_Named", func() error {
		payload, _ := json.Marshal(map[string]interface{}{
			"state": map[string]interface{}{
				"desired": map[string]interface{}{"firmware": "1.2.0"},
			},
		})
		if _, err := dataClient.UpdateThingShadow(ctx, &iotdataplane.UpdateThingShadowInput{
			ThingName: aws.String(thingName), ShadowName: aws.String(namedShadow), Payload: payload,
		}); err != nil {
			return fmt.Errorf("UpdateThingShadow (named) failed: %w", err)
		}
		return nil
	})

	add("Shadow_ListNamedShadowsForThing", func() error {
		out, err := dataClient.ListNamedShadowsForThing(ctx, &iotdataplane.ListNamedShadowsForThingInput{
			ThingName: aws.String(thingName),
		})
		if err != nil {
			return fmt.Errorf("ListNamedShadowsForThing failed: %w", err)
		}
		found := false
		for _, n := range out.Results {
			if n == namedShadow {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("expected named shadow %q in results %v", namedShadow, out.Results)
		}
		return nil
	})

	add("Shadow_DeleteThingShadow_Named", func() error {
		if _, err := dataClient.DeleteThingShadow(ctx, &iotdataplane.DeleteThingShadowInput{
			ThingName: aws.String(thingName), ShadowName: aws.String(namedShadow),
		}); err != nil {
			return fmt.Errorf("DeleteThingShadow (named) failed: %w", err)
		}
		return nil
	})

	add("Shadow_DeleteThingShadow_Classic", func() error {
		if _, err := dataClient.DeleteThingShadow(ctx, &iotdataplane.DeleteThingShadowInput{
			ThingName: aws.String(thingName),
		}); err != nil {
			return fmt.Errorf("DeleteThingShadow (classic) failed: %w", err)
		}
		return nil
	})

	add("Shadow_GetThingShadow_Deleted_NotFound", func() error {
		_, err := dataClient.GetThingShadow(ctx, &iotdataplane.GetThingShadowInput{
			ThingName: aws.String(thingName),
		})
		return expectNotFound(err)
	})

	return results
}
