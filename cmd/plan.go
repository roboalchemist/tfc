package cmd

import (
	"fmt"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "View plan details and logs",
}

var planShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show plan details",
	Args:  cobra.ExactArgs(1),
	RunE:  runPlanShow,
}

var planLogCmd = &cobra.Command{
	Use:   "log [id]",
	Short: "Stream plan log output",
	Args:  cobra.ExactArgs(1),
	RunE:  runPlanLog,
}

func init() {
	planCmd.AddCommand(planShowCmd, planLogCmd)
	rootCmd.AddCommand(planCmd)
}

type planAttrs struct {
	Status              string `json:"status"`
	LogReadURL          string `json:"log-read-url"`
	ResourceAdditions   int    `json:"resource-additions"`
	ResourceChanges     int    `json:"resource-changes"`
	ResourceDestructions int   `json:"resource-destructions"`
	ResourceImports     int    `json:"resource-imports"`
	HasChanges          bool   `json:"has-changes"`
	ExecutionDetails    struct {
		Mode string `json:"mode"`
	} `json:"execution-details"`
}

func runPlanShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	planID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/plans/"+planID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a planAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type planDetail struct {
		ID    string    `json:"id"`
		Attrs planAttrs `json:"attributes"`
	}
	data := planDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Status", a.Status},
			{"Has Changes", boolStr(a.HasChanges)},
			{"Additions", itoa(a.ResourceAdditions)},
			{"Changes", itoa(a.ResourceChanges)},
			{"Destructions", itoa(a.ResourceDestructions)},
			{"Imports", itoa(a.ResourceImports)},
		},
	}

	return output.RenderTable(td, data, opts)
}

func runPlanLog(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	planID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/plans/"+planID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a planAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	if a.LogReadURL == "" {
		return output.NewNotFoundError(fmt.Sprintf("no log URL available for plan %s (status: %s)", planID, a.Status))
	}

	body, err := fetchLogURL(a.LogReadURL)
	if err != nil {
		return output.NewAPIError(err.Error())
	}
	defer body.Close()

	return output.RenderStream(body, GetOutputOptions())
}
