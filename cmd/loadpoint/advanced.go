package loadpoint

import (
	"fmt"
	"strconv"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

var phaseValues = []string{"0", "1", "3"}

func newPhasesCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "phases",
		Short: "Allowed phases commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get configured phases",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"phasesConfigured",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:       "set [0|1|3]",
		Short:     "Set configured phases",
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: phaseValues,
		RunE: func(cmd *cobra.Command, args []string) error {
			return setLoadpointPhases(cmd, deps, args[0])
		},
	})

	return cmd
}

func newPriorityCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "priority",
		Short: "Loadpoint priority commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get loadpoint priority",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"priority",
				formatFloatAsNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [value]",
		Short: "Set loadpoint priority",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			priority, err := strconv.Atoi(args[0])
			if err != nil || priority < 0 {
				return fmt.Errorf("invalid priority %q: must be an integer >= 0", args[0])
			}
			return setLoadpointPriority(cmd, deps, priority)
		},
	})

	return cmd
}

func newSmartCostLimitCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smart-cost-limit",
		Short: "Smart charging cost limit commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get smart charging cost limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"smartCostLimit",
				formatNullableNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [cost]",
		Short: "Set smart charging cost limit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cost, err := parseNonNegativeFloat(args[0], "cost")
			if err != nil {
				return err
			}
			return setLoadpointSmartCostLimit(cmd, deps, cost)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete smart charging cost limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteLoadpointSmartCostLimit(cmd, deps)
		},
	})

	return cmd
}

func newSmartFeedInPriorityLimitCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "smart-feed-in-priority-limit",
		Short: "Smart feed-in priority limit commands",
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get smart feed-in priority limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"smartFeedInPriorityLimit",
				formatNullableNumber,
			)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "set [cost]",
		Short: "Set smart feed-in priority limit",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cost, err := parseNonNegativeFloat(args[0], "cost")
			if err != nil {
				return err
			}
			return setLoadpointSmartFeedInPriorityLimit(cmd, deps, cost)
		},
	})

	cmd.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "Delete smart feed-in priority limit",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return deleteLoadpointSmartFeedInPriorityLimit(cmd, deps)
		},
	})

	return cmd
}

func setLoadpointPhases(cmd *cobra.Command, deps dependencies, phases string) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointPhasesWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Phases(phases))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d phases set to %s\n", loadpoint, phases)
	return err
}

func setLoadpointPriority(cmd *cobra.Command, deps dependencies, priority int) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointPriorityWithResponse(cmd.Context(), evccgen.Id(loadpoint), priority)
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d priority set to %d\n", loadpoint, priority)
	return err
}

func deleteLoadpointSmartCostLimit(cmd *cobra.Command, deps dependencies) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.DeleteLoadpointSmartCostLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d smart cost limit deleted\n", loadpoint)
	return err
}

func setLoadpointSmartCostLimit(cmd *cobra.Command, deps dependencies, cost float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointSmartCostLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.CostLimit(cost))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d smart cost limit set to %g\n", loadpoint, cost)
	return err
}

func deleteLoadpointSmartFeedInPriorityLimit(cmd *cobra.Command, deps dependencies) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.DeleteLoadpointSmartFeedInPriorityLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d smart feed-in priority limit deleted\n", loadpoint)
	return err
}

func setLoadpointSmartFeedInPriorityLimit(cmd *cobra.Command, deps dependencies, cost float64) error {
	loadpoint, err := getLoadpoint(cmd)
	if err != nil {
		return err
	}

	client, err := deps.clientFactory()
	if err != nil {
		return fmt.Errorf("create api client: %w", err)
	}

	resp, err := client.SetLoadpointSmartFeedInPriorityLimitWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.CostLimit(cost))
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

	_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d smart feed-in priority limit set to %g\n", loadpoint, cost)
	return err
}

func formatNullableNumber(value any) string {
	if value == nil {
		return "unset"
	}

	number, ok := value.(float64)
	if !ok {
		return fmt.Sprintf("%v", value)
	}

	return formatFloatAsNumber(number)
}
