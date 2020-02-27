package nerdgraph

import (
	"encoding/json"
	"errors"
	"fmt"

	prettyjson "github.com/hokaccha/go-prettyjson"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	variables string
)

// Command represents the nerdgraph command.
var Command = &cobra.Command{
	Use:   "nerdgraph",
	Short: "Top-level command for executing raw GraphQL requests to the NerdGraph API.",
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Execute a raw GraphQL query request to the NerdGraph API.",
	Long: `
Execute a raw GraphQL query request to the NerdGraph API.

The query command accepts a single argument in the form of a GraphQL query as a string.
This command accepts an optional flag, --variables, which should be a JSON string where the
keys are the variables to be referenced in the GraphQL query.
`,
	Example: `newrelic nerdgraph query 'query($guid: EntityGuid!) { actor { entity(guid: $guid) { guid name domain entityType } } }' --variables '{"guid": "<GUID>"}'`,
	Args: func(cmd *cobra.Command, args []string) error {
		argsCount := len(args)

		if argsCount < 1 {
			return errors.New("missing graph query argument")
		}

		if argsCount > 1 {
			return errors.New("command expects only 1 argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			var variablesParsed map[string]interface{}

			err := json.Unmarshal([]byte(variables), &variablesParsed)
			if err != nil {
				log.Fatal(err)
			}

			query := args[0]

			result, err := nrClient.NerdGraph.Query(query, variablesParsed)
			if err != nil {
				log.Fatal(err)
			}

			json, err := prettyjson.Marshal(result)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(string(json))
		})
	},
}

func init() {
	Command.AddCommand(queryCmd)
	queryCmd.Flags().StringVar(&variables, "variables", "{}", "(Optional) The variables to pass to the GraphQL query, represented as a JSON string.")
}
