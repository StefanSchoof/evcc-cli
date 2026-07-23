package loadpoint

import (
	"fmt"
	"strconv"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

func newLimitCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "limit",
		Short: "Session limit commands",
	}

	cmd.AddCommand(newLimitEnergyCmd(deps))
	cmd.AddCommand(newLimitSocCmd(deps))

	return cmd
}

func newLimitEnergyCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "energy",
		Short: "Energy limit commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get session energy limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"limitEnergy",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [kwh]",
		Short: "Set session energy limit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			energy, err := parseNonNegativeFloat(args[0], "energy")
			if err != nil {
				return err
			}
			return setLoadpointEnergyLimit(cmd, deps, energy)
		},
	})

	return cmd
}

func newLimitSocCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "soc",
		Short: "SoC limit commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get session SoC limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"limitSoc",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [soc]",
		Short: "Set session SoC limit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			soc, err := parseNonNegativeFloat(args[0], "soc")
			if err != nil {
				return err
			}
			if soc > 100 {
				return fmt.Errorf("invalid soc %.2f: must be between 0 and 100", soc)
			}
			return setLoadpointSocLimit(cmd, deps, soc)
		},
	})

	return cmd
}

func newCurrentCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current",
		Short: "Current limit commands",
	}

	cmd.AddCommand(newCurrentMaxCmd(deps))
	cmd.AddCommand(newCurrentMinCmd(deps))

	return cmd
}

func newCurrentMaxCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "max",
		Short: "Maximum current",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get maximum current",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"maxCurrent",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [ampere]",
		Short: "Set maximum current",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			current, err := parseNonNegativeFloat(args[0], "current")
			if err != nil {
				return err
			}
			return setLoadpointMaxCurrent(cmd, deps, current)
		},
	})

	return cmd
}

func newCurrentMinCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "min",
		Short: "Minimum current",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get minimum current",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"minCurrent",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [ampere]",
		Short: "Set minimum current",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			current, err := parseNonNegativeFloat(args[0], "current")
			if err != nil {
				return err
			}
			return setLoadpointMinCurrent(cmd, deps, current)
		},
	})

	return cmd
}

func newTempCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "temp",
		Short: "Temperature limit commands",
	}

	cmd.AddCommand(newTempMinCmd(deps))

	return cmd
}

func newTempMinCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "min",
		Short: "Minimum temperature",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get minimum temperature",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"ui.minTemp",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [celsius]",
		Short: "Set minimum temperature",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			temp, err := parseNonNegativeFloat(args[0], "temperature")
			if err != nil {
				return err
			}
			return setLoadpointMinTemp(cmd, deps, temp)
		},
	})

	return cmd
}

func setLoadpointEnergyLimit(cmd *cobra.Command, deps dependencies, energy float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointEnergyLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Energy(energy))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d energy limit set to %g kWh\n", loadpoint, energy)
	return err
}

func setLoadpointSocLimit(cmd *cobra.Command, deps dependencies, soc float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointSocLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Soc(soc))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d SoC limit set to %g%%\n", loadpoint, soc)
	return err
}

func setLoadpointMaxCurrent(cmd *cobra.Command, deps dependencies, current float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointMaxCurrentWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Current(current))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d max current set to %g A\n", loadpoint, current)
	return err
}

func setLoadpointMinCurrent(cmd *cobra.Command, deps dependencies, current float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointMinCurrentWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Current(current))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d min current set to %g A\n", loadpoint, current)
	return err
}

func setLoadpointMinTemp(cmd *cobra.Command, deps dependencies, temp float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointMinTempWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Temp(temp))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d min temperature set to %g C\n", loadpoint, temp)
	return err
}

func parseNonNegativeFloat(input string, field string) (float64, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil || value < 0 {
		return 0, fmt.Errorf("invalid %s %q: must be a number >= 0", field, input)
	}
	return value, nil
}
