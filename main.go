package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func createMarkdownFile(templatePath, outputPath, date, name string, tags []string) {
	content, err := ioutil.ReadFile(templatePath)
	if err != nil {
		fmt.Printf("Error reading the template file: %s\n", err)
		return
	}

	// Replace {{DATE}} placeholder with the provided date value
	contentStr := string(content)
	contentStr = strings.Replace(contentStr, "{{DATE}}", date, -1)

	// Replace {{NAME}} placeholder with the provided name value
	var nameStr string
	if name != "" {
		nameStr = name
	} else {
		nameStr = ""
	}
	contentStr = strings.Replace(contentStr, "{{NAME}}", nameStr, -1)

	// Replace {{TAGS}} placeholder with the provided tags value or "[]" if tags are empty
	var tagsStr string
	if len(tags) > 0 {
		tagsJSON, err := json.Marshal(tags)
		if err != nil {
			fmt.Printf("Error converting tags to JSON: %s\n", err)
			return
		}
		tagsStr = string(tagsJSON)
	} else {
		tagsStr = "[]"
	}
	contentStr = strings.Replace(contentStr, "{{TAGS}}", tagsStr, -1)

	err = ioutil.WriteFile(outputPath, []byte(contentStr), 0644)
	if err != nil {
		fmt.Printf("Error creating the output markdown file: %s\n", err)
		return
	}

	fmt.Printf("Markdown file '%s' created using template '%s'.\n", outputPath, templatePath)
}

func getDefaultOutputFileName(template, name string) string {
	now := time.Now()
	formattedDate := now.Format("2006-01-02")
	formattedTime := now.Format("15:04:02") // 24-hour format (HH:mm:ss)

	if name != "" {
		return fmt.Sprintf("%s_%s_%s_%s", formattedDate, formattedTime, template, name)
	}
	return fmt.Sprintf("docs/journal/%s_%s_%s", formattedDate, formattedTime, template)
}

func getUserInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func main() {
	templateDir := "templates"

	// Parse command-line arguments
	flag.Parse()
	args := flag.Args()

	var template, templateName, name, tagsInput, date string

	// If no arguments are provided, use the default template
	if len(args) == 0 {
		template = "default.md" // Set the default template name to "default"
	} else if len(args) >= 2 {
		// If two or more arguments are provided, use the provided template and name
		template = args[0]
		name = args[1]
		tagsInput = getUserInput("Tags (optional, JSON-formatted array): ")
		date = getUserInput("Date (optional): ")
	} else {
		// If only one argument is provided, prompt for template and name
		template = args[0]
		templateName = getUserInput("Template name (optional): ")
		name = getUserInput("Name (optional): ")
		tagsInput = getUserInput("Tags (optional, JSON-formatted array): ")
		date = getUserInput("Date (optional): ")
	}

	if date == "" {
		// If no explicit date is provided, use the current date
		now := time.Now()
		date = now.Format("2006-01-02")
	}

	var tags []string
	if tagsInput != "" {
		err := json.Unmarshal([]byte(tagsInput), &tags)
		if err != nil {
			fmt.Printf("Error parsing tags as JSON array: %s\n", err)
			return
		}
	}

	if templateName == "" {
		templateName = template
	}

	outputFile := filepath.Join(".", getDefaultOutputFileName(template, name))
	templatePath := filepath.Join(templateDir, templateName)

	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		fmt.Printf("Template '%s' not found.\n", templateName)
		return
	}

	createMarkdownFile(templatePath, outputFile, date, name, tags)
}
