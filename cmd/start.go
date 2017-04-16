// Copyright © 2017 Valerio Gheri <valerio.gheri@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vgheri/gennaker/api"
	"github.com/vgheri/gennaker/engine"
	"github.com/vgheri/gennaker/repository/pg"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start runs the gennaker service",
	Long:  `gennaker start runs the HTTP API server powering gennaker`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Work your own magic here
		fmt.Println("start called")
		repository, err := pg.NewClient(postgresHost, fmt.Sprintf("%d", postgresPort), postgresUsername,
			postgresPassword, postgresDBName, int(postgresMaxConnections))
		if err != nil {
			panic(err)
		}
		deploymentEngine := engine.New(repository, chartsDownloadFolder)
		server, err := api.New(deploymentEngine)
		if err != nil {
			panic(err)
		}
		server.Start(HTTPListenPort)
	},
}

var HTTPListenPort, postgresPort, postgresMaxConnections int32
var postgresHost, postgresUsername, postgresPassword, postgresDBName string
var chartsDownloadFolder string

func init() {
	RootCmd.AddCommand(startCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	startCmd.Flags().Int32VarP(&HTTPListenPort, "http-port", "p", 8080, "Port number for the HTTP server")
	startCmd.Flags().Int32Var(&postgresPort, "pg-port", 5432, "Port number for Postgres")
	startCmd.Flags().Int32Var(&postgresMaxConnections, "db-maxconn", 10, "Max number of connections to Postgres")
	startCmd.Flags().StringVar(&postgresHost, "pg-host", "localhost", "Postgres installation host name")
	startCmd.Flags().StringVar(&postgresDBName, "pg-db", "gennaker", "Postgres database name")
	startCmd.Flags().StringVar(&postgresUsername, "pg-username", "postgres", "Postgres username")
	startCmd.Flags().StringVar(&postgresPassword, "pg-password", "password", "Postgres password")
	startCmd.Flags().StringVarP(&chartsDownloadFolder, "save-dir", "d", "localhost", "Path used to download charts. Must be absolute")
}
