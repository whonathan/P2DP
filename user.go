package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func validateUsername(shortUsername string) (string, bool) {
	prefixes := map[string]int{
		"d": 15,
		"j": 14,
		"z": 15,
	}

	if len(shortUsername) != 3 {
		return "", false
	}

	prefix := string(shortUsername[0])
	maxNum, validPrefix := prefixes[prefix]
	if !validPrefix {
		return "", false
	}

	num, err := strconv.Atoi(shortUsername[1:])
	if err != nil || num < 1 || num > maxNum {
		return "", false
	}

	fullUsername := fmt.Sprintf("52260.%s%02d", prefix, num)
	return fullUsername, true
}

func promptUser() (string, string, int, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter username identifier (format: d01-d15, j01-j14, or z01-z15): ")
	shortUsername, _ := reader.ReadString('\n')
	shortUsername = strings.TrimSpace(shortUsername)

	username, valid := validateUsername(shortUsername)
	if !valid {
		return "", "", 0, fmt.Errorf("invalid username format")
	}

	client, err := NewPLNClient(100, username)
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to initialize client: %v", err)
	}

	fmt.Println("Fetching prepaid data...")
	err = client.FetchAndStorePrepaidData(username, "52260")
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to fetch data: %v", err)
	}

	kdrbmData, err := client.GetKDRBMData()
	if err != nil {
		return "", "", 0, fmt.Errorf("failed to get KDRBM data: %v", err)
	}

	if len(kdrbmData) == 0 {
		return "", "", 0, fmt.Errorf("no KDRBM data available")
	}

	fmt.Println("\nAvailable KDRBM data:")
	for _, data := range kdrbmData {
		fmt.Printf("KDRBM: %s, Remaining: %d\n", data.KDRBM, data.Count)
	}

	fmt.Print("\nSelect KDRBM: ")
	kdrbm, _ := reader.ReadString('\n')
	kdrbm = strings.TrimSpace(kdrbm)

	var selectedKDRBM *KDRBMData
	for _, data := range kdrbmData {
		if data.KDRBM == kdrbm {
			selectedKDRBM = &data
			break
		}
	}

	if selectedKDRBM == nil {
		return "", "", 0, fmt.Errorf("invalid KDRBM selection")
	}

	fmt.Printf("\nEnter number of data to send (max %d): ", min(110, selectedKDRBM.Count))
	countStr, _ := reader.ReadString('\n')
	count, err := strconv.Atoi(strings.TrimSpace(countStr))
	if err != nil {
		return "", "", 0, fmt.Errorf("invalid number format")
	}

	if count <= 0 || count > 110 || count > selectedKDRBM.Count {
		return "", "", 0, fmt.Errorf("invalid count: must be between 1 and min(110, available data)")
	}

	return username, kdrbm, count, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
