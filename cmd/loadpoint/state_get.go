package loadpoint

import (
	"encoding/json"
	"fmt"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

func runLoadpointStateGet[T any](cmd *cobra.Command, deps dependencies, expr string, formatOutput func(T) string) error {
	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	jq := fmt.Sprintf("{value: .loadpoints[%d].%s}", loadpoint-1, expr)
	resp, err := client.StateWithResponse(cmd.Context(), &evccgen.StateParams{Jq: &jq})
	if err != nil {
		return err
	}

	if deps.rawEnabled() {
		_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
		return err
	}

	if len(resp.Body) == 0 {
		return fmt.Errorf("empty response")
	}

	// we can not use because resp.Json200 is not correct when jq is used
	var payload map[string]any
	if err := json.Unmarshal(resp.Body, &payload); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	rawValue, ok := payload["value"]
	if !ok {
		return fmt.Errorf("field 'value' not found")
	}

	value, ok := rawValue.(T)
	if !ok {
		return fmt.Errorf("unexpected type %T for expression %q", rawValue, expr)
	}

	_, err = fmt.Fprintln(cmd.OutOrStdout(), formatOutput(value))
	return err
}
