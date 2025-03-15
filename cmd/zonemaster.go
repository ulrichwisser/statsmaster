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
	"github.com/apex/log"
	"encoding/json"
	"math/rand"
	"net/http"
	"bytes"
	_ "github.com/go-sql-driver/mysql"
)

func submit2zonemaster(domain string) string {
	url := "https://zonemaster.net/api/start_domain_test"
	requestData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      rand.Intn(89999) + 10000,
		"method":  "start_domain_test",
		"params": map[string]interface{}{
			"language": "en",
			"domain":   domain,
			"profile":  "default",
		},
	}
	jsonData, _ := json.Marshal(requestData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Errorf("Error submitting to %s zonemaster.net: %s", domain, err)
		return ""
	}

	var response map[string]interface{}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		log.Errorf("Error decoding received data after test submission %s: %s", domain, err)
		return ""
	}
	return response["result"].(string)
}

func getTestProgress(testID string) (int, error) {
	url := "https://zonemaster.net/api"
	requestData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      rand.Intn(89999) + 10000,
		"method":  "test_progress",
		"params": map[string]string{
			"test_id": testID,
		},
	}
	jsonData, _ := json.Marshal(requestData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return -1, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return -1, err
	}
	return int(response["result"].(float64)), nil
}

func getTestResults(testid string) (map[string]interface{}, error) {
	url := "https://zonemaster.net/api/get_test_results"
	requestData := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      rand.Intn(89999) + 10000,
		"method":  "get_test_results",
		"params": map[string]string{
			"id":       testid,
			"language": "en",
		},
	}
	jsonData, _ := json.Marshal(requestData)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return response["result"].(map[string]interface{}), nil
}
