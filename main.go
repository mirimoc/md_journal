package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

func getUserInput(prompt string, args ...interface{}) string {
	fmt.Printf(prompt, args...)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func listMarkdownFiles(folderPath string) []string {
	var markdownFiles []string

	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Printf("Warning: Reading the folder %s\n", err)
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			// Skip directories
			continue
		}

		if strings.HasSuffix(file.Name(), ".md") {
			markdownFiles = append(markdownFiles, file.Name())
		}
	}

	return markdownFiles
}

func openMarkdownFile(filePath string) error {
	cmd := exec.Command("xdg-open", filePath) // Try xdg-open for Linux
	if runtime.GOOS == "darwin" {
		cmd = exec.Command("open", filePath) // For macOS, use "open" command
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "start", filePath) // For Windows, use "cmd" and "start" command
	}
	err := cmd.Start()
	return err
}

func main() {
	// Get the absolute path to the executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting the executable path: %s\n", err)
		return
	}

	// Get the directory path of the executable
	exeDir := filepath.Dir(exePath)

	// Update the template directory to be relative to the executable location
	templateDir := filepath.Join(exeDir, "templates")

	// Define the flags for both short and long forms
	var wizardFlag bool
	flag.BoolVar(&wizardFlag, "w", false, "Run in wizard mode (interactive prompts for optionals)")
	flag.BoolVar(&wizardFlag, "wizard", false, "Run in wizard mode (interactive prompts for optionals)")
	flag.Parse()
	args := flag.Args()

	var template, templateName, name, tagsInput, date string

	folderPath := "templates/"
	templates := strings.Join(listMarkdownFiles(folderPath), ", ")
	// Check if the wizard flag is provided
	if wizardFlag {
		// Run in wizard mode with interactive prompts for optionals

		//templateName = getUserInput("Template name (possible templates: %s): ", templates)
		// Prompt the user for template name with the possible templates shown
		templateName = getUserInput("Template name (possible templates: %s): ", templates)
		template = templateName
		name = getUserInput("Name (optional): ")
		tagsInput = getUserInput("Tags (optional, JSON-formatted array): ")
		date = getUserInput("Date (optional): ")
	} else {
		if len(args) == 0 {
			template = "task.md" // Set the default template name to "default"
		} else if len(args) == 1 {
			// If two or more arguments are provided, use the provided template and name
			template = args[0]
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
	// Open the generated markdown file with the default program associated with its file type
	if err := openMarkdownFile(outputFile); err != nil {
		fmt.Printf("Error opening the markdown file: %s\n", err)
	}
}
