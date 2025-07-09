package services

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/0xPixelNinja/dhcp-rest-api/config"
	"github.com/0xPixelNinja/dhcp-rest-api/models"
)

func ListHosts() ([]models.Host, error) {
	content, err := os.ReadFile(config.AppConfig.DhcpConfPath)
	if err != nil {
		log.Printf("Error reading DHCP config file: %v", err)
		return nil, fmt.Errorf("failed to read DHCP config: %w", err)
	}

	var hosts []models.Host
	hostBlockRegex := regexp.MustCompile(`(?s)host\s+([^\s]+)\s*\{([^}]+)\}`)
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

func UpdateHost(name string, updates models.HostUpdate) error {
	content, err := os.ReadFile(config.AppConfig.DhcpConfPath)
	if err != nil {
		log.Printf("Error reading DHCP config file: %v", err)
		return fmt.Errorf("failed to read DHCP config for update: %w", err)
	}

	hostBlockRegex := regexp.MustCompile(fmt.Sprintf(`(?s)(host\s+%s\s*\{)([^}]+)(\})`, regexp.QuoteMeta(name)))

	match := hostBlockRegex.FindStringSubmatchIndex(string(content))
	if match == nil {
		log.Printf("Host %s not found for update.", name)
		return fmt.Errorf("host %s not found", name)
	}

	// Extract current fields from the block
	blockContent := string(content[match[4]:match[5]])
	currentHost := models.Host{Name: name}

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
		currentHost.Name = *updates.Name
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

	// Build the new host block
	newHostBlockContent := fmt.Sprintf("\n    hardware ethernet %s;\n    option routers %s;\n    option subnet-mask %s;\n    fixed-address %s;\n    option domain-name-servers %s;\n",
		currentHost.HardwareEthernet, currentHost.OptionRouters, currentHost.OptionSubnetMask, currentHost.FixedAddress, currentHost.OptionDomainNameServers)

	newHostBlock := fmt.Sprintf("host %s { %s}", updatedName, newHostBlockContent)

	// Replace the entire host block
	newFileContent := string(content[:match[0]]) + newHostBlock + string(content[match[1]:])

	if err := os.WriteFile(config.AppConfig.DhcpConfPath, []byte(newFileContent), 0644); err != nil {
		log.Printf("Error writing updated DHCP config file: %v", err)
		return fmt.Errorf("failed to write updated DHCP config: %w", err)
	}
	return nil
}

func DeleteHost(name string) error {
	content, err := os.ReadFile(config.AppConfig.DhcpConfPath)
	if err != nil {
		log.Printf("Error reading DHCP config file: %v", err)
		return fmt.Errorf("failed to read DHCP config for delete: %w", err)
	}

	hostBlockRegex := regexp.MustCompile(fmt.Sprintf(`(?s)\s*host\s+%s\s*\{[^}]*\}\s*`, regexp.QuoteMeta(name)))

	if !hostBlockRegex.Match(content) {
		log.Printf("Host %s not found for deletion.", name)
		return nil // idempotent delete
	}

	newContent := hostBlockRegex.ReplaceAll(content, []byte{})

	// Clean up whitespace
	newContent = []byte(strings.TrimSpace(string(newContent)))
	if len(newContent) > 0 {
		newContent = append(newContent, '\n')
	}

	if err := os.WriteFile(config.AppConfig.DhcpConfPath, newContent, 0644); err != nil {
		log.Printf("Error writing updated DHCP config file after deletion: %v", err)
		return fmt.Errorf("failed to write DHCP config after delete: %w", err)
	}
	return nil
}
