package internal

import (
	"log"
	"os"
	"strings"

	"github.com/gwenwindflower/tbd/shared"

	"gopkg.in/yaml.v2"
)

func WriteYAML(tables shared.SourceTables, buildDir string, database *string, schema *string) {
	yamlData, err := yaml.Marshal(tables)
	if err != nil {
		log.Fatalf("Failed to marshal data into YAML %v\n", err)
	}

	// Convert the YAML data to a string and split it into lines
	yamlLines := strings.Split(string(yamlData), "\n")

	// Remove the first line with sources: because it's redundant with the sources section in the YAML
	yamlLines = yamlLines[1:]

	// Create a new slice to hold the indented lines
	indentedYamlLines := make([]string, len(yamlLines))

	// Indent each line with 6 spaces
	for i, line := range yamlLines {
		indentedYamlLines[i] = "      " + line
	}

	// Join the indented lines back into a single string
	indentedYaml := strings.Join(indentedYamlLines, "\n")

	// Prepend the version and source information to the indented YAML
	finalYaml := "version: 2\n\nsources:\n  - name: " + *database + "\n    schema: " + *schema + "\n    tables:\n" + indentedYaml

	// Write the final YAML content to the specified file
	writeError := os.WriteFile(buildDir+"/_sources.yml", []byte(finalYaml), 0644)
	if writeError != nil {
		log.Fatalf("Failed to write file %v\n", writeError)
	}
}
