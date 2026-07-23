package loadpoint

import (
	"fmt"
	"strings"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

func newVehicleCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vehicle",
		Short: "Vehicle assignment and detection commands",
	}

	cmd.AddCommand(newVehicleGetCmd(deps))
	cmd.AddCommand(newVehicleAssignCmd(deps))
	cmd.AddCommand(newVehicleRemoveCmd(deps))
	cmd.AddCommand(newVehicleDetectionCmd(deps))

	return cmd
}

func newVehicleGetCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get assigned vehicle",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLoadpointStateGet(
				cmd,
				deps,
				"vehicleName",
				func(name string) string {
					if strings.TrimSpace(name) == "" {
						return "unassigned"
					}
					return name
				},
			)
		},
	}

	return cmd
}

func newVehicleAssignCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assign [vehicle-name]",
		Short: "Assign a vehicle to the loadpoint",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			vehicleName := strings.TrimSpace(args[0])
			if vehicleName == "" {
				return fmt.Errorf("vehicle name must not be empty")
			}

			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.AssignLoadpointVehicleWithResponse(cmd.Context(), evccgen.Id(loadpoint), evccgen.VehicleName(vehicleName))
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

			result := vehicleName
			if resp.JSON200 != nil && resp.JSON200.Vehicle != nil && strings.TrimSpace(string(*resp.JSON200.Vehicle)) != "" {
				result = string(*resp.JSON200.Vehicle)
			}

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d vehicle assigned to %s\n", loadpoint, result)
			return err
		},
	}

	return cmd
}

func newVehicleRemoveCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove assigned vehicle from the loadpoint",
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

			resp, err := client.RemoveLoadpointVehicleWithResponse(cmd.Context(), evccgen.Id(loadpoint))
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

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d vehicle removed\n", loadpoint)
			return err
		},
	}

	return cmd
}

func newVehicleDetectionCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detection",
		Short: "Vehicle detection commands",
	}

	cmd.AddCommand(newVehicleDetectionStartCmd(deps))

	return cmd
}

func newVehicleDetectionStartCmd(deps dependencies) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start",
		Aliases: []string{"detect-start", "start-detection"},
		Short:   "Start vehicle detection",
		Args:    cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			loadpoint, err := getLoadpoint(cmd)
			if err != nil {
				return err
			}

			client, err := deps.clientFactory()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}

			resp, err := client.StartLoadpointVehicleDetectionWithResponse(cmd.Context(), evccgen.Id(loadpoint))
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

			_, err = fmt.Fprintf(cmd.OutOrStdout(), "loadpoint %d vehicle detection started\n", loadpoint)
			return err
		},
	}

	return cmd
}
