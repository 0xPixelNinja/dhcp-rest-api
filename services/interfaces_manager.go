package services

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/0xPixelNinja/dhcp-rest-api/config"
)

var (
	interfacesV4Regex = regexp.MustCompile(`INTERFACESv4="([^"]*)"`)
	interfacesV6Regex = regexp.MustCompile(`INTERFACESv6="([^"]*)"`)
)

// GetInterfaces reads the interfaces configuration file and returns the v4 and v6 interfaces.
func GetInterfaces() (map[string]string, error) {
	content, err := ioutil.ReadFile(config.AppConfig.InterfacesConfPath)
	if err != nil {
		if os.IsNotExist(err) { // If file doesn't exist, treat as empty configuration
			log.Printf("Interfaces config file not found at %s, returning empty config", config.AppConfig.InterfacesConfPath)
			return map[string]string{"v4": "", "v6": ""}, nil
		}
		log.Printf("Error reading interfaces config file: %v", err)
		return nil, fmt.Errorf("failed to read interfaces config: %w", err)
	}

	interfaces := make(map[string]string)
	v4Match := interfacesV4Regex.FindStringSubmatch(string(content))
	if len(v4Match) > 1 {
		interfaces["v4"] = v4Match[1]
	} else {
		interfaces["v4"] = ""
	}

	v6Match := interfacesV6Regex.FindStringSubmatch(string(content))
	if len(v6Match) > 1 {
		interfaces["v6"] = v6Match[1]
	} else {
		interfaces["v6"] = ""
	}

	return interfaces, nil
}

// SaveInterfaces writes the provided interface settings back to the configuration file.
// It reads the file line by line, updates existing INTERFACESv4/v6 lines,
// and preserves all other content.
// This function will return an error if the interfaces configuration file does not exist.
func SaveInterfaces(interfaces map[string]string) error {
	filePath := config.AppConfig.InterfacesConfPath

	// Attempt to open for reading. If it doesn't exist, Python's version would fail.
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Error: interfaces config file %s does not exist. Cannot save.", filePath)
			return fmt.Errorf("interfaces config file '%s' does not exist: %w", filePath, err)
		}
		log.Printf("Error opening interfaces config file for read: %v", err)
		return fmt.Errorf("failed to open interfaces config file '%s' for read: %w", filePath, err)
	}

	scanner := bufio.NewScanner(file)
	var newLines []string

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "INTERFACESv4=") {
			newLines = append(newLines, fmt.Sprintf("INTERFACESv4=\"%s\"", interfaces["v4"]))
		} else if strings.HasPrefix(line, "INTERFACESv6=") {
			newLines = append(newLines, fmt.Sprintf("INTERFACESv6=\"%s\"", interfaces["v6"]))
		} else {
			newLines = append(newLines, line)
		}
	}
	// Important: Close the file *before* checking scanner.Err() and *before* writing.
	if errClose := file.Close(); errClose != nil {
		log.Printf("Warning: error closing interfaces config file after reading: %v", errClose)
		// Depending on severity, you might choose to return errClose here if critical
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning interfaces config file '%s': %v", filePath, err)
		return fmt.Errorf("failed to scan interfaces config file '%s': %w", filePath, err)
	}

	outputContent := strings.Join(newLines, "\n")
	// Ensure a trailing newline if there's content, which is good practice for config files.
	if len(newLines) > 0 && !strings.HasSuffix(outputContent, "\n") {
		outputContent += "\n"
	}

	// Write the potentially modified content back to the file, truncating it.
	if err := ioutil.WriteFile(filePath, []byte(outputContent), 0644); err != nil {
		log.Printf("Error writing updated interfaces config to file '%s': %v", filePath, err)
		return fmt.Errorf("failed to write updated interfaces config to '%s': %w", filePath, err)
	}

	return nil
}

// AddInterface adds an interface to the specified type (v4 or v6) if not already present.
func AddInterface(ifaceType string, ifaceName string) error {
	interfaces, err := GetInterfaces()
	if err != nil {
		return err // Error already logged by GetInterfaces
	}

	key := strings.ToLower(ifaceType)
	if key != "v4" && key != "v6" {
		return fmt.Errorf("invalid interface type: %s, must be v4 or v6", ifaceType)
	}

	currentValue := interfaces[key]
	currentList := strings.Fields(currentValue) // Splits by whitespace

	for _, existingIface := range currentList {
		if existingIface == ifaceName {
			log.Printf("Interface %s already present in %s", ifaceName, key)
			return nil // Already present, no action needed
		}
	}

	currentList = append(currentList, ifaceName)
	interfaces[key] = strings.Join(currentList, " ")

	return SaveInterfaces(interfaces)
}

// DeleteInterface removes an interface from the specified type (v4 or v6) if present.
func DeleteInterface(ifaceType string, ifaceName string) error {
	interfaces, err := GetInterfaces()
	if err != nil {
		return err // Error already logged by GetInterfaces
	}

	key := strings.ToLower(ifaceType)
	if key != "v4" && key != "v6" {
		return fmt.Errorf("invalid interface type: %s, must be v4 or v6", ifaceType)
	}

	currentValue := interfaces[key]
	currentList := strings.Fields(currentValue)
	newList := []string{}
	found := false

	for _, existingIface := range currentList {
		if existingIface == ifaceName {
			found = true
		} else {
			newList = append(newList, existingIface)
		}
	}

	if !found {
		log.Printf("Interface %s not found in %s, no action needed", ifaceName, key)
		return nil // Not found, no action needed
	}

	interfaces[key] = strings.Join(newList, " ")

	return SaveInterfaces(interfaces)
}
