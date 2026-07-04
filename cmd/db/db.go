package db

import (
	"github.com/raymondsugiarto/reputation-be/cmd/db/migrate"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/database"
	"github.com/spf13/cobra"
)

// DBCmd represents the db command. Subcommands:
//
//	app db migrate up [step]   — apply / step forward pending migrations
//	app db migrate down [step] — roll back migrations
//	app db seed                — apply db/migrations/seed.sql (idempotent)
//
// `db seed` is wired through migrate.RunSeed so the same seed file
// runs whether you boot the production server (which calls
// MigrateUpAll → Seed) or invoke the CLI directly.
var DBCmd = &cobra.Command{
	Use:   "db",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		switch args[0] {
		case "migrate":
			schema, _ := cmd.Flags().GetString("schema")
			migrate.Migration(args, schema)
		case "seed":
			// Seed requires a live DB connection. Reuse the same
			// initialisation path as the production server so the
			// schema/search_path state matches.
			database.InitForCLI()
			migrate.RunSeed()
		default:
			cmd.Help()
		}
	},
}
