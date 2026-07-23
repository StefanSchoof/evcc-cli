package loadpoint

import (
	"context"
	"fmt"

	evccgen "evcc-cli/internal/gen/evcc"

	"github.com/spf13/cobra"
)

type ClientFactory func() (*evccgen.ClientWithResponses, error)
type StateFetcher func(ctx context.Context, jq string) ([]byte, error)
type RawEnabled func() bool
type APIErrorFormatter func(status string, body []byte) error

type dependencies struct {
	clientFactory ClientFactory
	rawEnabled    RawEnabled
}

func NewCmd(clientFactory ClientFactory, rawEnabled RawEnabled) *cobra.Command {
	deps := dependencies{
		clientFactory: clientFactory,
		rawEnabled:    rawEnabled,
	}

	cmd := &cobra.Command{
		Use:   "loadpoint",
		Short: "Loadpoint related commands",
	}

	cmd.PersistentFlags().IntP("loadpoint", "l", 1, "loadpoint index starting at 1")
	cmd.AddCommand(newModeCmd(deps))
	cmd.AddCommand(newBatteryCmd(deps))
	cmd.AddCommand(newVehicleCmd(deps))
	cmd.AddCommand(newPlanCmd(deps))
	cmd.AddCommand(newThresholdCmd(deps))
	cmd.AddCommand(newLimitCmd(deps))
	cmd.AddCommand(newCurrentCmd(deps))
	cmd.AddCommand(newTempCmd(deps))
	cmd.AddCommand(newPhasesCmd(deps))
	cmd.AddCommand(newPriorityCmd(deps))
	cmd.AddCommand(newSmartCostLimitCmd(deps))
	cmd.AddCommand(newSmartFeedInPriorityLimitCmd(deps))

	return cmd
}

func getLoadpoint(cmd *cobra.Command) (int, error) {
	loadpoint, err := cmd.Flags().GetInt("loadpoint")
	if err != nil {
		return 0, err
	}
	if loadpoint < 1 {
		return 0, fmt.Errorf("--loadpoint must be >= 1")
	}
	return loadpoint, nil
}
