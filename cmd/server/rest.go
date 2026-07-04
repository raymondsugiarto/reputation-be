package cmdserver

import (
	"github.com/raymondsugiarto/reputation-be/cmd/db/migrate"
	"github.com/raymondsugiarto/reputation-be/config"
	"github.com/raymondsugiarto/reputation-be/pkg/infrastructure/server"

	"github.com/spf13/cobra"
)

var RestCmd = &cobra.Command{
	Use:   "api",
	Short: "Start Api Server",
	Long:  `Start the Rest Api Server`,
	Run:   startRest,
}

func startRest(cmd *cobra.Command, args []string) {
	config.GetConfig()
	httpServer := server.NewRest()
	httpServer.Initialize()
}

// production mode
var StartRestCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Api Server Production",
	Long:  `Start the Rest Api Server Production`,
	Run:   startRestProduction,
}

func startRestProduction(cmd *cobra.Command, args []string) {
	migrate.MigrateUpAll()

	config.GetConfig()
	httpServer := server.NewRest()
	httpServer.Initialize()
}
