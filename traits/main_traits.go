package traits

import (
	"device-worker-service/internal/domain"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func FindDataByUUId(devices []*domain.Device, uuid string) (string, string) {

	for _, v := range devices {
		if v.UUID == uuid {
			fmt.Println("")
			return v.Host, v.Port
		}
	}

	return "", ""
}

func ParseDevice(jsonData string) ([]*domain.Device, error) {

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		return nil, err
	}

	devices := make([]*domain.Device, 0)
	for uuid, values := range data {
		hostPort := values.([]interface{})
		host := hostPort[0].(string)
		port := fmt.Sprintf("%v", hostPort[1])

		device := &domain.Device{
			UUID: uuid,
			Host: host,
			Port: port,
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func ExtractPID(logMessage string) (int, error) {
	// Define a regex pattern to match the PID
	regexPattern := `Started server process \[([0-9]+)\]`

	// Compile the regex pattern
	re := regexp.MustCompile(regexPattern)

	// Find the match in the log message
	matches := re.FindStringSubmatch(logMessage)

	// Check if a match was found
	if len(matches) < 2 {
		return 0, fmt.Errorf("PID not found in log message")
	}

	// Extract the PID from the matched group
	pidStr := matches[1]

	// Convert the PID string to an integer
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("Failed to convert PID to integer: %v", err)
	}
	return pid, nil
}

func ParseUrl(url string) (string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", nil
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
