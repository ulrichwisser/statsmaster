/*
Copyright © 2025 Ulrich Wisser

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/apex/log"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	_ "github.com/go-sql-driver/mysql"
)

var tld2region map[string]string

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "read domain list from file and submit domains to zonemaster.net",
	Long: `read domain list from file and submit domains to zonemaster.net`,
	Run: execSubmit,
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// submitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	submitCmd.Flags().StringP(TLD_FILE, TLD_FILE_SHORT, "", "name of file with the TLD list")
	submitCmd.Flags().StringP(COUNTRIES_FILE, COUNTRIES_FILE_SHORT, "", "name of file with list of domain by country")
	submitCmd.Flags().StringP(REGIONS_FILE, REGIONS_FILE_SHORT, REGIONS_FILE_DEFAULT, "name of file with ICANN regions")

	// Use flags for viper values
	viper.BindPFlags(submitCmd.Flags())
}

func execSubmit(cmd *cobra.Command, args []string) {
	tldFilename := viper.GetString(TLD_FILE)
	countriesFilename := viper.GetString(COUNTRIES_FILE)
	log.Infof("Using TLD file: %s", tldFilename)
	log.Infof("Using countries file: %s", tldFilename)

	if len(tldFilename) == 0 && len(countriesFilename) == 0 {
		cmd.Help();
		log.Fatal("Neiter TLD nor Country file are given.")
	}

	// load in regions information
	getRegions()

    // open database
	db := openDB()

	if len(tldFilename) > 0 {
		log.Infof("Using TLD file: %s", tldFilename)
		handleTldFile(db, tldFilename)
	}
	if len(countriesFilename) > 0 {
		log.Infof("Using countries file: %s", tldFilename)
	}

}

func handleTldFile(db *sql.DB, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Error reading TLD file %s: %s", filename, err)
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		 line= strings.TrimSpace(line)
		
		 // jump over empty lines
		if line == "" {
			continue
		}
		
		// jump over comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		
		// now we do have avalid TLD
		tld := strings.ToLower(line)

		// get region for TLD
		region,_ := tld2region[tld]
		
		//
		log.Debugf("TLD %s  Region %s", tld, region)
		testid := submit2zonemaster(tld)
		if testid != "" {
			pushToQueue(db, testid, region, "", "")
		}
	}

}

func getRegions() {
	file, err := os.Open(viper.GetString(REGIONS_FILE))
	if err != nil {
		log.Fatalf("Error opening regions file: %s", err)
	}
	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	err = json.Unmarshal(byteValue, &tld2region)
	if err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}
}