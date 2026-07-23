package loadpoint

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

func newBatteryCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "battery",
		Short: "Battery related loadpoint commands",
	}

	cmd.AddCommand(newBatteryBoostCmd(deps))

	return cmd
}

func newBatteryBoostCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "boost",
		Short: "Control battery boost",
	}

	cmd.AddCommand(newBatteryBoostGetCmd(deps))
	cmd.AddCommand(newBatteryBoostSocLimitCmd(deps))
	cmd.AddCommand(newBatteryBoostEnableCmd(deps))
	cmd.AddCommand(newBatteryBoostDisableCmd(deps))

	return cmd
}

func newBatteryBoostSocLimitCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soc-limit",
		Short: "Get or set battery boost SoC limit",
	}

	cmd.AddCommand(newBatteryBoostSocLimitGetCmd(deps))
	cmd.AddCommand(newBatteryBoostSocLimitSetCmd(deps))

	return cmd
}

func newBatteryBoostSocLimitGetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get battery boost SoC limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"batteryBoostLimit",
				func(limit float64) string {
					if math.Trunc(limit) == limit {
						return fmt.Sprintf("%d", int(limit))
					}
					return fmt.Sprintf("%v", limit)
				},
			)
		},
	}

	return cmd
}

func newBatteryBoostSocLimitSetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set [soc]",
		Short: "Set battery boost SoC limit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			soc, err := strconv.ParseFloat(args[0], 32)
			if err != nil {
				return fmt.Errorf("invalid SoC %q: must be a number between 0 and 100", args[0])
			}
			if soc < 0 || soc > 100 {
				return fmt.Errorf("invalid SoC %.2f: must be between 0 and 100", soc)
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.SetLoadpointBatteryBoostLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Soc(soc))
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

			if resp.JSON200 == nil {
				return fmt.Errorf("failed to read battery boost SoC limit from response for loadpoint %d", loadpoint)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d battery boost SoC limit set to %d\n", loadpoint, int(*resp.JSON200))
			return err
		},
	}

	return cmd
}

func newBatteryBoostGetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get battery boost status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"batteryBoost",
				func(enabled bool) string {
					if enabled {
						return "enabled"
					}
					return "disabled"
				},
			)
		},
	}

	return cmd
}

func newBatteryBoostEnableCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "enable",
		Short: "Enable battery boost",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setBatteryBoost(cmd, deps, true)
		},
	}

	return cmd
}

func newBatteryBoostDisableCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disable",
		Short: "Disable battery boost",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return setBatteryBoost(cmd, deps, false)
		},
	}

	return cmd
}

func setBatteryBoost(cmd *cobra.Command, deps dependencies, enabled bool) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	value := evccgen.SetLoadpointBatteryBoostParamsEnableFalse
	status := "disabled"
	if enabled {
		value = evccgen.SetLoadpointBatteryBoostParamsEnableTrue
		status = "enabled"
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointBatteryBoostWithResponse(cmd.Context(), evccgen.Id(loadpoint), value)
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

	if resp.JSON200 == nil {
		return fmt.Errorf("failed to read battery boost state from response for loadpoint %d", loadpoint)
	}

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d battery boost %s\n", loadpoint, status)
	return err
}
