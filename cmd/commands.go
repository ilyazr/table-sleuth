package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"table-sleuth/processor"
)

var rootCmd = &cobra.Command{
	Use:   "table-sleuth",
	Short: "Table Sleuth is a tool, that helps you find all DB tables used in a project",
	Long: `Table Sleuth is a tool, that helps you find all DB tables used in a particular Spring Boot project. It can do the following:
	1. Create a mapping of table -> services (meaning, which services use the exact table)
	2. Create a mapping of service -> tables used (meaning, which tables use the exact service).
It can work with both a single project and a directory of projects`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Table Sleuth",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var processProjectsCmd = &cobra.Command{
	Use:     "s2t",
	Short:   "Show what which tables are used in a service",
	Example: `table-sleuth s2t -p /home/projects/spring-boot-app1 -p /home/projects/spring-boot-app2`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := cmd.Flags().GetStringArray("project")
		if err != nil {
			log.Fatalf("error:", err.Error())
		}
		result := processor.ProcessProjects(v)
		printResult(result)
	},
}

var processDirOfProjectsCmd = &cobra.Command{
	Use:     "t2s",
	Short:   "Show what service use particular tables (table to service mapping)",
	Example: `table-sleuth s2t -d /home/spring-boot-projects`,
	Run: func(cmd *cobra.Command, args []string) {
		v, err := cmd.Flags().GetString("dir")
		if err != nil {
			log.Fatalf(err.Error())
		}
		result := processor.ProcessDirOfProjects(v)
		printResult(result)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)

	processProjectsCmd.Flags().StringArrayP("project", "p", []string{}, "Set project (or projects) for processing")
	rootCmd.AddCommand(processProjectsCmd)

	processDirOfProjectsCmd.Flags().StringP("dir", "d", "", "Set dir of projects for processing")
	rootCmd.AddCommand(processDirOfProjectsCmd)

}

func printResult(result interface{}) {
	if result == nil {
		log.Fatal("No Java project found")
	}

	jsonString, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(jsonString))
}

func Execute() error {
	return rootCmd.Execute()
}
