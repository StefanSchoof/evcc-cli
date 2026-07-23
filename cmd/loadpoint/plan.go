package loadpoint

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

func newPlanCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Charging plan commands",
	}

	cmd.AddCommand(newPlanGetCmd(deps))
	cmd.AddCommand(newPlanEnergyCmd(deps))
	cmd.AddCommand(newPlanPreviewCmd(deps))
	cmd.AddCommand(newPlanSimulateCmd(deps))
	cmd.AddCommand(newPlanStrategyCmd(deps))

	return cmd
}

func newPlanGetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get charging plan",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.GetLoadpointPlanWithResponse(cmd.Context(), evccgen.Id(loadpoint))
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() || resp.JSON200 == nil {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			return printJSON(cmd, resp.JSON200)
		},
	}

	return cmd
}

func newPlanEnergyCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "energy",
		Short: "Energy-based charging plan commands",
	}

	cmd.AddCommand(newPlanEnergyDeleteCmd(deps))
	cmd.AddCommand(newPlanEnergySetCmd(deps))

	return cmd
}

func newPlanEnergyDeleteCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete energy-based charging plan",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.DeleteLoadpointEnergyPlanWithResponse(cmd.Context(), evccgen.Id(loadpoint))
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d energy-based plan deleted\n", loadpoint)
			return err
		},
	}

	return cmd
}

func newPlanEnergySetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [energy-kwh] [timestamp-rfc3339]",
		Short: "Set energy-based charging plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			energy, err := parseEnergy(args[0])
			if err != nil {
				return err
			}
			timestamp, err := parseTimestamp(args[1])
			if err != nil {
				return err
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.SetLoadpointEnergyPlanWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Energy(energy), evccgen.Timestamp(timestamp))
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() || resp.JSON200 == nil {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			return printJSON(cmd, resp.JSON200)
		},
	}

	return cmd
}

func newPlanPreviewCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preview",
		Short: "Preview charging plans",
	}

	cmd.AddCommand(newPlanPreviewRepeatingCmd(deps))

	return cmd
}

func newPlanPreviewRepeatingCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "repeating [soc] [weekdays] [time-hh:mm] [timezone]",
		Short: "Repeating plan preview",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			soc, err := parseSoc(args[0])
			if err != nil {
				return err
			}
			weekdays, err := parseWeekdays(args[1])
			if err != nil {
				return err
			}
			hhmm, err := parseHourMinute(args[2])
			if err != nil {
				return err
			}
			tz := args[3]
			if strings.TrimSpace(tz) == "" {
				return fmt.Errorf("timezone must not be empty")
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.PreviewLoadpointRepeatingPlanWithResponse(
				cmd.Context(),
				evccgen.Id(loadpoint),
				evccgen.Soc(soc),
				evccgen.Weekdays(weekdays),
				evccgen.HourMinuteTime(hhmm),
				evccgen.Timezone(tz),
			)
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() || resp.JSON200 == nil {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			return printJSON(cmd, resp.JSON200)
		},
	}

	return cmd
}

func newPlanSimulateCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "simulate",
		Short: "Simulate charging plans",
	}

	cmd.AddCommand(newPlanSimulateSocCmd(deps))
	cmd.AddCommand(newPlanSimulateEnergyCmd(deps))

	return cmd
}

func newPlanSimulateSocCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soc [soc] [timestamp-rfc3339]",
		Short: "Simulate soc-based charging plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			soc, err := parseSoc(args[0])
			if err != nil {
				return err
			}
			timestamp, err := parseTimestamp(args[1])
			if err != nil {
				return err
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.PreviewLoadpointSocPlanWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Soc(soc), evccgen.Timestamp(timestamp))
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() || resp.JSON200 == nil {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			return printJSON(cmd, resp.JSON200)
		},
	}

	return cmd
}

func newPlanSimulateEnergyCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "energy [energy-kwh] [timestamp-rfc3339]",
		Short: "Simulate energy-based charging plan",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			energy, err := parseEnergy(args[0])
			if err != nil {
				return err
			}
			timestamp, err := parseTimestamp(args[1])
			if err != nil {
				return err
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.PreviewLoadpointEnergyPlanWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Energy(energy), evccgen.Timestamp(timestamp))
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() || resp.JSON200 == nil {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			return printJSON(cmd, resp.JSON200)
		},
	}

	return cmd
}

func newPlanStrategyCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "strategy",
		Short: "Plan strategy commands",
	}

	cmd.AddCommand(newPlanStrategySetCmd(deps))

	return cmd
}

func newPlanStrategySetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [continuous:true|false] [precondition-seconds]",
		Short: "Set plan strategy",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			continuous, err := strconv.ParseBool(args[0])
			if err != nil {
				return fmt.Errorf("invalid continuous value %q: must be true or false", args[0])
			}
			precondition, err := strconv.Atoi(args[1])
			if err != nil || precondition < 0 {
				return fmt.Errorf("invalid precondition %q: must be >= 0", args[1])
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			body := evccgen.SetLoadpointPlanStrategyJSONRequestBody{
				Continuous:   continuous,
				Precondition: float32(precondition),
			}
			resp, err := client.SetLoadpointPlanStrategyWithResponse(cmd.Context(), evccgen.Id(loadpoint), body)
			if err != nil {
				return err
			}
			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), strings.TrimSpace(string(resp.Body)))
			}

			if deps.rawEnabled() || resp.JSON200 == nil {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			return printJSON(cmd, resp.JSON200)
		},
	}

	return cmd
}

func parseSoc(input string) (float32, error) {
	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid soc %q: must be a number between 0 and 100", input)
	}
	if value < 0 || value > 100 {
		return 0, fmt.Errorf("invalid soc %.2f: must be between 0 and 100", value)
	}
	return float32(value), nil
}

func parseEnergy(input string) (float32, error) {
	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid energy %q: must be a positive number", input)
	}
	if value <= 0 {
		return 0, fmt.Errorf("invalid energy %.2f: must be > 0", value)
	}
	return float32(value), nil
}

func parseTimestamp(input string) (time.Time, error) {
	ts, err := time.Parse(time.RFC3339, input)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid timestamp %q: expected RFC3339, e.g. 2026-07-22T04:00:00Z", input)
	}
	return ts, nil
}

func parseWeekdays(input string) ([]int, error) {
	parts := strings.Split(input, ",")
	if len(parts) == 0 {
		return nil, fmt.Errorf("invalid weekdays %q: expected comma-separated values in range 1..7", input)
	}

	weekdays := make([]int, 0, len(parts))
	seen := map[int]bool{}
	for _, part := range parts {
		value, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil || value < 1 || value > 7 {
			return nil, fmt.Errorf("invalid weekday %q: allowed values are 1..7", strings.TrimSpace(part))
		}
		if !seen[value] {
			seen[value] = true
			weekdays = append(weekdays, value)
		}
	}

	return weekdays, nil
}

func parseHourMinute(input string) (string, error) {
	if _, err := time.Parse("15:04", input); err != nil {
		return "", fmt.Errorf("invalid time %q: expected HH:MM", input)
	}
	return input, nil
}

func printJSON(cmd *cobra.Command, v any) error {
	out, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(cmd.OutOrStdout(), string(out))
	return err
}
