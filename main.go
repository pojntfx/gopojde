package main

import (
	"log"

	"github.com/docker/docker/client"
	"github.com/pojntfx/gopojde/pkg/orchestration"
	"github.com/pojntfx/gopojde/pkg/servers"
	"github.com/pojntfx/gopojde/pkg/services"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configFileKey             = "configFile"
	listenAddressKey          = "listenAddress"
	websocketListenAddressKey = "websocketListenAddress"
)

func main() {
	// Create command
	cmd := &cobra.Command{
		Use:   "gopojde-backend",
		Short: "Experimental Go implementation of https://github.com/pojntfx/pojde.",
		Long: `Experimental Go implementation of https://github.com/pojntfx/pojde.

For more information, please visit https://github.com/pojntfx/gopojde.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Bind config file
			if !(viper.GetString(configFileKey) == "") {
				viper.SetConfigFile(viper.GetString(configFileKey))

				if err := viper.ReadInConfig(); err != nil {
					return err
				}
			}

			// Connect to Docker
			docker, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				log.Fatal("could not connect to Docker:", err)
			}

			// Create orchestration backends
			instancesManager := orchestration.NewInstancesManager(docker)

			// Create services
			instancesService := services.NewInstancesService(instancesManager)

			// Create servers
			srv := servers.NewGRPCServer(viper.GetString(listenAddressKey), viper.GetString(websocketListenAddressKey), instancesService)

			// Start servers
			log.Printf("gopojde backend listening on %v (gRPC) and %v (gRPC-Web)", viper.GetString(listenAddressKey), viper.GetString(websocketListenAddressKey))

			return srv.ListenAndServe()
		},
	}

	// Bind flags
	cmd.PersistentFlags().StringP(configFileKey, "c", "", "Config file to use")

	cmd.PersistentFlags().String(listenAddressKey, ":15323", "Listen address")
	cmd.PersistentFlags().String(websocketListenAddressKey, ":15324", "Listen address (for the WebSocket proxy)")

	// Bind env variables
	if err := viper.BindPFlags(cmd.PersistentFlags()); err != nil {
		log.Fatal(err)
	}
	viper.SetEnvPrefix("gopojde_backend")
	viper.AutomaticEnv()

	// Run command
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
