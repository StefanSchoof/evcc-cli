package loadpoint

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

func newThresholdCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threshold",
		Short: "Set start/stop thresholds for solar mode",
	}

	cmd.AddCommand(newThresholdEnableCmd(deps))
	cmd.AddCommand(newThresholdDisableCmd(deps))

	return cmd
}

func newThresholdEnableCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable threshold configuration",
	}

	cmd.AddCommand(newThresholdEnableDelayCmd(deps))
	cmd.AddCommand(newThresholdEnablePowerCmd(deps))

	return cmd
}

func newThresholdDisableCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable threshold configuration",
	}

	cmd.AddCommand(newThresholdDisableDelayCmd(deps))
	cmd.AddCommand(newThresholdDisablePowerCmd(deps))

	return cmd
}

func newThresholdEnableDelayCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delay",
		Short: "Enable delay threshold",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get enable delay threshold",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"enableDelay",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [seconds]",
		Short: "Set enable delay threshold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			delay, err := parseDelaySeconds(args[0])
			if err != nil {
				return err
			}
			return setThresholdEnableDelay(cmd, deps, delay)
		},
	})

	return cmd
}

func newThresholdDisableDelayCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delay",
		Short: "Disable delay threshold",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get disable delay threshold",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"disableDelay",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [seconds]",
		Short: "Set disable delay threshold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			delay, err := parseDelaySeconds(args[0])
			if err != nil {
				return err
			}
			return setThresholdDisableDelay(cmd, deps, delay)
		},
	})

	return cmd
}

func newThresholdEnablePowerCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power",
		Short: "Enable power threshold",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get enable power threshold",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"enableThreshold",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [watts]",
		Short: "Set enable power threshold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			threshold, err := parsePowerWatts(args[0])
			if err != nil {
				return err
			}
			return setThresholdEnablePower(cmd, deps, threshold)
		},
	})

	return cmd
}

func newThresholdDisablePowerCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power",
		Short: "Disable power threshold",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get disable power threshold",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"disableThreshold",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [watts]",
		Short: "Set disable power threshold",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			threshold, err := parsePowerWatts(args[0])
			if err != nil {
				return err
			}
			return setThresholdDisablePower(cmd, deps, threshold)
		},
	})

	return cmd
}

func setThresholdEnableDelay(cmd *cobra.Command, deps dependencies, delay int) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointEnableDelayWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Delay(delay))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d enable delay set to %d s\n", loadpoint, delay)
	return err
}

func setThresholdDisableDelay(cmd *cobra.Command, deps dependencies, delay int) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointDisableDelayWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Delay(delay))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d disable delay set to %d s\n", loadpoint, delay)
	return err
}

func setThresholdEnablePower(cmd *cobra.Command, deps dependencies, threshold float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointEnableThresholdWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Threshold(threshold))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d enable threshold set to %g W\n", loadpoint, threshold)
	return err
}

func setThresholdDisablePower(cmd *cobra.Command, deps dependencies, threshold float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointDisableThresholdWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Threshold(threshold))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d disable threshold set to %g W\n", loadpoint, threshold)
	return err
}

func parseDelaySeconds(input string) (int, error) {
	delay, err := strconv.Atoi(input)
	if err != nil || delay < 0 {
		return 0, fmt.Errorf("invalid delay %q: must be an integer >= 0", input)
	}
	return delay, nil
}

func parsePowerWatts(input string) (float64, error) {
	threshold, err := strconv.ParseFloat(input, 32)
	if err != nil || threshold < 0 {
		return 0, fmt.Errorf("invalid power %q: must be a number >= 0", input)
	}
	return threshold, nil
}

func formatFloatAsNumber(value float64) string {
	if math.Trunc(value) == value {
		return fmt.Sprintf("%d", int(value))
	}
	return fmt.Sprintf("%v", value)
}
