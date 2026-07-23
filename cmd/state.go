package cmd

import (
	"encoding/json"
	evccgen "evcc-cli/internal/gen/evcc"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

func newStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state",
		Short: "Fetch and print evcc state",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newGeneratedClient()
			if err != nil {
				return fmt.Errorf("create api client: %w", err)
			}
			resp, err := client.StateWithResponse(cmd.Context(), &evccgen.StateParams{})
			if err != nil {
				return err
			}

			if cfg.Raw {
				_, err = fmt.Fprintln(cmd.OutOrStdout(), string(resp.Body))
				return err
			}

			if resp.StatusCode() >= 400 {
				return fmt.Errorf("evcc API error: %s: %s", resp.Status(), string(resp.Body))
			}

			if resp.JSON200 == nil {
				return fmt.Errorf("empty state response")
			}

			state := resp.JSON200

			updated := "-"
			if state.HistoryUpdated != nil {
				updated = state.HistoryUpdated.Format("2006-01-02T15:04:05Z07:00")
			}

			interval := 0
			if state.Interval != nil {
				interval = int(*state.Interval)
			}

			gridPower := any(nil)
			if state.Grid != nil {
				gridPower = state.Grid.Power
			}

			homePower := any(nil)
			if state.HomePower != nil {
				homePower = *state.HomePower
			}

			pvPower := any(nil)
			if state.PvPower != nil {
				pvPower = *state.PvPower
			}

			batterySoc := any(nil)
			if state.Battery != nil {
				batterySoc = state.Battery.Soc
			}

			fmt.Fprintf(cmd.OutOrStdout(), "updated: %s\n", updated)
			if interval > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "interval: %ds\n", interval)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "site: grid=%sW home=%sW pv=%sW batterySoc=%s%%\n", watt(gridPower), watt(homePower), watt(pvPower), num(batterySoc))

			if len(state.Loadpoints) == 0 {
				_, err := fmt.Fprintln(cmd.OutOrStdout(), "loadpoints: none")
				return err
			}

			_, _ = fmt.Fprintln(cmd.OutOrStdout(), "loadpoints:")
			for i, lp := range state.Loadpoints {
				name := lp.Title
				if name == "" {
					name = lp.Name
				}
				if name == "" {
					name = fmt.Sprintf("lp-%d", i+1)
				}

				vehicle := lp.VehicleTitle
				if vehicle == "" {
					vehicle = lp.VehicleName
				}
				if vehicle == "" {
					vehicle = "-"
				}

				_, _ = fmt.Fprintf(
					cmd.OutOrStdout(),
					"  %d) %s mode=%s connected=%t charging=%t power=%sW vehicle=%s soc=%s%%\n",
					i+1,
					name,
					string(lp.Mode),
					lp.Connected,
					lp.Charging,
					watt(lp.ChargePower),
					vehicle,
					num(lp.VehicleSoc),
				)
			}

			return nil
		},
	}
	return cmd
}

func num(v any) string {
	switch n := v.(type) {
	case float64:
		return strconv.FormatFloat(n, 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(n), 'f', -1, 32)
	case int:
		return strconv.Itoa(n)
	case int64:
		return strconv.FormatInt(n, 10)
	case json.Number:
		return n.String()
	case nil:
		return "-"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func watt(v any) string {
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f", n)
	case float32:
		return fmt.Sprintf("%.2f", n)
	case int:
		return fmt.Sprintf("%.2f", float64(n))
	case int64:
		return fmt.Sprintf("%.2f", float64(n))
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return n.String()
		}
		return fmt.Sprintf("%.2f", f)
	case nil:
		return "-"
	default:
		return fmt.Sprintf("%v", v)
	}
}
