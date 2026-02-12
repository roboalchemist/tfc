package cmd

import (
	"fmt"
	"io"
	"net/http"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "View apply details and logs",
}

var applyShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show apply details",
	Args:  cobra.ExactArgs(1),
	RunE:  runApplyShow,
}

var applyLogCmd = &cobra.Command{
	Use:   "log [id]",
	Short: "Stream apply log output",
	Args:  cobra.ExactArgs(1),
	RunE:  runApplyLog,
}

func init() {
	applyCmd.AddCommand(applyShowCmd, applyLogCmd)
	rootCmd.AddCommand(applyCmd)
}

type applyAttrs struct {
	Status              string `json:"status"`
	LogReadURL          string `json:"log-read-url"`
	ResourceAdditions   int    `json:"resource-additions"`
	ResourceChanges     int    `json:"resource-changes"`
	ResourceDestructions int   `json:"resource-destructions"`
	ResourceImports     int    `json:"resource-imports"`
}

func runApplyShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	applyID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/applies/"+applyID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a applyAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type applyDetail struct {
		ID    string     `json:"id"`
		Attrs applyAttrs `json:"attributes"`
	}
	data := applyDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Status", a.Status},
			{"Additions", itoa(a.ResourceAdditions)},
			{"Changes", itoa(a.ResourceChanges)},
			{"Destructions", itoa(a.ResourceDestructions)},
			{"Imports", itoa(a.ResourceImports)},
		},
	}

	return output.RenderTable(td, data, opts)
}

func runApplyLog(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	applyID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/applies/"+applyID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a applyAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	if a.LogReadURL == "" {
		return output.NewNotFoundError(fmt.Sprintf("no log URL available for apply %s (status: %s)", applyID, a.Status))
	}

	body, err := fetchLogURL(a.LogReadURL)
	if err != nil {
		return output.NewAPIError(err.Error())
	}
	defer body.Close()

	return output.RenderStream(body, GetOutputOptions())
}

// fetchLogURL fetches a TFC log URL (these are pre-signed S3 URLs, no auth needed).
func fetchLogURL(url string) (io.ReadCloser, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetch log: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("fetch log: status %d", resp.StatusCode)
	}
	return resp.Body, nil
}
