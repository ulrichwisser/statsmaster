/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/apex/log"
)

// retrieveCmd represents the retrieve command
var retrieveCmd = &cobra.Command{
	Use:   "retrieve",
	Short: "Retrieve test results from zonemaster.net",
	Long:  `Retrieve test results from zonemaster.net`,
	Run: execRetrieve,
}

func init() {
	rootCmd.AddCommand(retrieveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// retrieveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// retrieveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")


	// Use flags for viper values
	viper.BindPFlags(retrieveCmd.Flags())
}

func execRetrieve(cmd *cobra.Command, args []string) {
	// open DB
	db := openDB()

	rows, err := db.Query("SELECT testid, region, country, sector FROM QUEUE ORDER BY id")
    if err != nil {
		log.Errorf("Error fetching test IDs from queue: %s", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		// get queue info
		var testid string
		var region string
		var country string
		var sector string
		if err := rows.Scan(&testid, &region, &country, &sector); err != nil {
				log.Errorf("Error scanning row: %s", err)
				continue
		}

		// get progress
		progress, err := getTestProgress(testid)
		if err != nil {
			log.Errorf("Error checking progress for %s: %s", testid, err)
			continue
		}

		// jump to next id if progress is not 100%
		if progress < 100 {
			log.Debugf("Progress for TESTID %s is %d", testid, progress)
			continue
		}

		// results should be ready
		results, err := getTestResults(testid)
		if err != nil {
			log.Errorf("Error fetching results for %s: %s", testid, err)
			continue
		}

		// save results to database
		if err := saveResults(db, testid, region, country, sector, results); err != nil {
			log.Errorf("Error saving results for %s: %s", testid, err)
		} else {
			popFromQueue(db, testid)
		}
	}
}
