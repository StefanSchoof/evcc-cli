package loadpoint

import (
	"fmt"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"
	evccstate "evcc-cli/internal/gen/evccstate"

	"github.com/spf13/cobra"
)

var modeValues = []string{
	string(evccstate.Off),
	string(evccstate.Now),
	string(evccstate.Minpv),
	string(evccstate.Pv),
}

func newModeCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mode",
		Short: "Get or set the charge mode",
	}

	cmd.AddCommand(newModeGetCmd(deps))
	cmd.AddCommand(newModeSetCmd(deps))

	return cmd
}

func newModeGetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get charge mode of a loadpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"mode",
				func(mode string) string { return mode },
			)
		},
	}
	return cmd
}

func newModeSetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:       "set [off|now|minpv|pv]",
		Short:     "Set charge mode",
		Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		ValidArgs: modeValues,
		RunE: func(cmd *cobra.Command, args []string) error {
			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			mode := args[0]
			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.SetLoadpointModeWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.Mode(mode))
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
				return fmt.Errorf("failed to read mode from response for loadpoint %d", loadpoint)
			}
			setMode := string(*resp.JSON200)

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d mode set to %s\n", loadpoint, setMode)
			return err
		},
	}
	return cmd
}
