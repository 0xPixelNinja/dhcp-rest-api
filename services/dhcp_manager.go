package services

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"log"

	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/0xPixelNinja/dhcp-rest-api/models"
)

// ListHosts parses the DHCP configuration file and returns a list of hosts.
func ListHosts() ([]models.Host, error) {
	content, err := ioutil.ReadFile(config.AppConfig.DhcpConfPath)
	if err != nil {
		log.Printf("Error reading DHCP config file: %v", err)
		return nil, fmt.Errorf("failed to read DHCP config: %w", err)
	}

	var hosts []models.Host
	// Regex to find host blocks: host <name> { ... }
	hostBlockRegex := regexp.MustCompile(`(?s)host\s+([^\s]+)\s*\{([^}]+)\}`)
	// Regex to parse individual lines within a host block
	hardwareEthernetRegex := regexp.MustCompile(`hardware ethernet\s+([^;]+);`)
	optionRoutersRegex := regexp.MustCompile(`option routers\s+([^;]+);`)
	optionSubnetMaskRegex := regexp.MustCompile(`option subnet-mask\s+([^;]+);`)
	fixedAddressRegex := regexp.MustCompile(`fixed-address\s+([^;]+);`)
	optionDomainNameServersRegex := regexp.MustCompile(`option domain-name-servers\s+([^;]+);`)

	matches := hostBlockRegex.FindAllStringSubmatch(string(content), -1)
	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		hostName := match[1]
		blockContent := match[2]

		host := models.Host{Name: hostName}

		if m := hardwareEthernetRegex.FindStringSubmatch(blockContent); len(m) > 1 {
			host.HardwareEthernet = strings.TrimSpace(m[1])
		}
		if m := optionRoutersRegex.FindStringSubmatch(blockContent); len(m) > 1 {
			host.OptionRouters = strings.TrimSpace(m[1])
		}
		if m := optionSubnetMaskRegex.FindStringSubmatch(blockContent); len(m) > 1 {
			host.OptionSubnetMask = strings.TrimSpace(m[1])
		}
		if m := fixedAddressRegex.FindStringSubmatch(blockContent); len(m) > 1 {
			host.FixedAddress = strings.TrimSpace(m[1])
		}
		if m := optionDomainNameServersRegex.FindStringSubmatch(blockContent); len(m) > 1 {
			host.OptionDomainNameServers = strings.TrimSpace(m[1])
		}
		hosts = append(hosts, host)
	}
	return hosts, nil
}

// AddHost adds a new host entry to the DHCP configuration file.
func AddHost(host models.Host) error {
	hostBlock := fmt.Sprintf("\nhost %s {\n    hardware ethernet %s;\n    option routers %s;\n    option subnet-mask %s;\n    fixed-address %s;\n    option domain-name-servers %s;\n}\n",
		host.Name, host.HardwareEthernet, host.OptionRouters, host.OptionSubnetMask, host.FixedAddress, host.OptionDomainNameServers)

	f, err := os.OpenFile(config.AppConfig.DhcpConfPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening DHCP config file for append: %v", err)
		return fmt.Errorf("failed to open DHCP config for append: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString(hostBlock); err != nil {
		log.Printf("Error writing to DHCP config file: %v", err)
		return fmt.Errorf("failed to write to DHCP config: %w", err)
	}
	return nil
}

// UpdateHost updates an existing host entry in the DHCP configuration file.
func UpdateHost(name string, updates models.HostUpdate) error {
	content, err := ioutil.ReadFile(config.AppConfig.DhcpConfPath)
	if err != nil {
		log.Printf("Error reading DHCP config file: %v", err)
		return fmt.Errorf("failed to read DHCP config for update: %w", err)
	}

	// Regex to find the specific host block
	// Need to escape the name if it contains regex special characters, though less common for hostnames.
	hostBlockRegex := regexp.MustCompile(fmt.Sprintf(`(?s)(host\s+%s\s*\{)([^}]+)(\})`, regexp.QuoteMeta(name)))

	match := hostBlockRegex.FindStringSubmatchIndex(string(content))
	if match == nil {
		log.Printf("Host %s not found for update.", name)
		return fmt.Errorf("host %s not found", name)
	}

	// Extract current fields from the block
	blockContent := string(content[match[4]:match[5]])
	currentHost := models.Host{Name: name} // Start with the original name

	if m := regexp.MustCompile(`hardware ethernet\s+([^;]+);`).FindStringSubmatch(blockContent); len(m) > 1 {
		currentHost.HardwareEthernet = strings.TrimSpace(m[1])
	}
	if m := regexp.MustCompile(`option routers\s+([^;]+);`).FindStringSubmatch(blockContent); len(m) > 1 {
		currentHost.OptionRouters = strings.TrimSpace(m[1])
	}
	if m := regexp.MustCompile(`option subnet-mask\s+([^;]+);`).FindStringSubmatch(blockContent); len(m) > 1 {
		currentHost.OptionSubnetMask = strings.TrimSpace(m[1])
	}
	if m := regexp.MustCompile(`fixed-address\s+([^;]+);`).FindStringSubmatch(blockContent); len(m) > 1 {
		currentHost.FixedAddress = strings.TrimSpace(m[1])
	}
	if m := regexp.MustCompile(`option domain-name-servers\s+([^;]+);`).FindStringSubmatch(blockContent); len(m) > 1 {
		currentHost.OptionDomainNameServers = strings.TrimSpace(m[1])
	}

	// Apply updates
	updatedName := name
	if updates.Name != nil && *updates.Name != "" {
		updatedName = *updates.Name
		currentHost.Name = *updates.Name // Update the name in the host struct as well
	}
	if updates.HardwareEthernet != nil {
		currentHost.HardwareEthernet = *updates.HardwareEthernet
	}
	if updates.OptionRouters != nil {
		currentHost.OptionRouters = *updates.OptionRouters
	}
	if updates.OptionSubnetMask != nil {
		currentHost.OptionSubnetMask = *updates.OptionSubnetMask
	}
	if updates.FixedAddress != nil {
		currentHost.FixedAddress = *updates.FixedAddress
	}
	if updates.OptionDomainNameServers != nil {
		currentHost.OptionDomainNameServers = *updates.OptionDomainNameServers
	}

	// Construct the new host block
	newHostBlockContent := fmt.Sprintf("\n    hardware ethernet %s;\n    option routers %s;\n    option subnet-mask %s;\n    fixed-address %s;\n    option domain-name-servers %s;\n",
		currentHost.HardwareEthernet, currentHost.OptionRouters, currentHost.OptionSubnetMask, currentHost.FixedAddress, currentHost.OptionDomainNameServers)

	newHostBlock := fmt.Sprintf("host %s { %s}", updatedName, newHostBlockContent)

	// Replace the old block with the new block
	newFileContent := string(content[:match[0]]) + newHostBlock + string(content[match[1]+(match[3]-match[2]):])
	// Corrected replacement logic to use the full match indices
	// The hostBlockRegex above has 3 capture groups: (host <name> {), (block_content), (})
	// match[0] and match[1] are start and end of the entire matched block.
	// We want to replace the whole thing.
	newFileContent = string(content[:match[0]]) + newHostBlock + string(content[match[1]:])

	if err := ioutil.WriteFile(config.AppConfig.DhcpConfPath, []byte(newFileContent), 0644); err != nil {
		log.Printf("Error writing updated DHCP config file: %v", err)
		return fmt.Errorf("failed to write updated DHCP config: %w", err)
	}
	return nil
}

// DeleteHost removes a host entry from the DHCP configuration file.
func DeleteHost(name string) error {
	content, err := ioutil.ReadFile(config.AppConfig.DhcpConfPath)
	if err != nil {
		log.Printf("Error reading DHCP config file: %v", err)
		return fmt.Errorf("failed to read DHCP config for delete: %w", err)
	}

	// Regex to find the host block, including potential leading/trailing whitespace for clean removal
	// (?s) allows . to match newlines. Adjusted to remove the entire block including surrounding whitespace lines if possible.
	hostBlockRegex := regexp.MustCompile(fmt.Sprintf(`(?s)\s*host\s+%s\s*\{[^}]*\}\s*`, regexp.QuoteMeta(name)))

	if !hostBlockRegex.Match(content) {
		log.Printf("Host %s not found for deletion.", name)
		// In the Python version, this case returns True (idempotent delete)
		// For Go, we can return nil for success or a specific error/bool if preferred.
		// Returning nil to match Python's behavior of not erroring if not found.
		return nil
	}

	newContent := hostBlockRegex.ReplaceAll(content, []byte{})

	// Trim leading/trailing whitespace from the new content if the file becomes empty
	// or to clean up if the removed block was at the very start/end.
	newContent = []byte(strings.TrimSpace(string(newContent)))
	if len(newContent) > 0 {
		newContent = append(newContent, '\n') // Add a trailing newline if content remains
	}

	if err := ioutil.WriteFile(config.AppConfig.DhcpConfPath, newContent, 0644); err != nil {
		log.Printf("Error writing updated DHCP config file after deletion: %v", err)
		return fmt.Errorf("failed to write DHCP config after delete: %w", err)
	}
	return nil
}

// Helper to read file and split into lines - not used currently but can be useful
func _readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// Helper to write lines back to a file - not used currently
func _writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	return writer.Flush()
}
