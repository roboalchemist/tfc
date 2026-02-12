package cmd

import (
	"fmt"

	"github.com/roboalchemist/tfc/pkg/jsonapi"
	"github.com/roboalchemist/tfc/pkg/output"
	"github.com/spf13/cobra"
)

var (
	flagSVWorkspace string
	flagSVPageSize  int
)

var stateVersionCmd = &cobra.Command{
	Use:     "state-version",
	Aliases: []string{"sv"},
	Short:   "Manage state versions",
}

var stateVersionListCmd = &cobra.Command{
	Use:   "list",
	Short: "List state versions for a workspace",
	RunE:  runStateVersionList,
}

var stateVersionShowCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Show state version details",
	Args:  cobra.ExactArgs(1),
	RunE:  runStateVersionShow,
}

var stateVersionCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new state version",
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("state-version create")
	},
}

var stateVersionDownloadCmd = &cobra.Command{
	Use:   "download [id]",
	Short: "Download the state file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return notImplemented("state-version download")
	},
}

func init() {
	stateVersionListCmd.Flags().StringVar(&flagSVWorkspace, "workspace", "", "Workspace ID (required)")
	stateVersionListCmd.Flags().IntVar(&flagSVPageSize, "page-size", 20, "Results per page")

	stateVersionCreateCmd.Flags().String("workspace", "", "Workspace ID (required)")
	stateVersionCreateCmd.Flags().String("file", "", "Path to state file")
	stateVersionCreateCmd.Flags().String("serial", "", "State serial number")
	stateVersionCreateCmd.Flags().String("md5", "", "MD5 hash of the state file")
	stateVersionCreateCmd.Flags().String("lineage", "", "State lineage")

	stateVersionCmd.AddCommand(
		stateVersionListCmd,
		stateVersionShowCmd,
		stateVersionCreateCmd,
		stateVersionDownloadCmd,
	)
	rootCmd.AddCommand(stateVersionCmd)
}

type svAttrs struct {
	Serial            int    `json:"serial"`
	CreatedAt         string `json:"created-at"`
	Size              int    `json:"size"`
	ResourcesProcessed bool  `json:"resources-processed"`
	Modules           struct {
		Root struct {
			Resources []interface{} `json:"resources"`
		} `json:"root"`
	} `json:"modules"`
	HostedStateDownloadURL string `json:"hosted-state-download-url"`
}

func runStateVersionList(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	if flagSVWorkspace == "" {
		return output.NewUsageError("--workspace is required")
	}

	path := fmt.Sprintf("/workspaces/%s/state-versions?page[size]=%d", flagSVWorkspace, flagSVPageSize)

	var doc jsonapi.Document
	if err := client.Get(path, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	resources, err := jsonapi.ParseList(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	opts := GetOutputOptions()

	type svJSON struct {
		ID    string  `json:"id"`
		Attrs svAttrs `json:"attributes"`
	}
	var jsonData []svJSON
	td := output.TableData{
		Headers: []string{"ID", "SERIAL", "SIZE", "CREATED"},
	}

	for _, r := range resources {
		var a svAttrs
		jsonapi.UnmarshalAttributes(&r, &a)
		td.Rows = append(td.Rows, []string{
			r.ID, itoa(a.Serial), itoa(a.Size), shortDate(a.CreatedAt),
		})
		jsonData = append(jsonData, svJSON{ID: r.ID, Attrs: a})
	}

	return output.RenderTable(td, jsonData, opts)
}

func runStateVersionShow(cmd *cobra.Command, args []string) error {
	client, err := newClient()
	if err != nil {
		return err
	}

	svID := args[0]
	var doc jsonapi.Document
	if err := client.Get("/state-versions/"+svID, &doc); err != nil {
		return output.NewAPIError(err.Error())
	}

	res, err := jsonapi.ParseSingle(&doc)
	if err != nil {
		return output.NewAPIError(err.Error())
	}

	var a svAttrs
	jsonapi.UnmarshalAttributes(res, &a)

	opts := GetOutputOptions()

	type svDetail struct {
		ID    string  `json:"id"`
		Attrs svAttrs `json:"attributes"`
	}
	data := svDetail{ID: res.ID, Attrs: a}

	td := output.TableData{
		Headers: []string{"FIELD", "VALUE"},
		Rows: [][]string{
			{"ID", res.ID},
			{"Serial", itoa(a.Serial)},
			{"Size", itoa(a.Size)},
			{"Resources Processed", boolStr(a.ResourcesProcessed)},
			{"Created", a.CreatedAt},
		},
	}

	return output.RenderTable(td, data, opts)
}
