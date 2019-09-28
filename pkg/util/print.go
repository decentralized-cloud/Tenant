// Package util implements different utilities required by the tenant service
package util

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// These are the colors we'll use in pretty printing output
const noFormat = "\033[0m"
const cSkyyblue = "\033[38;5;117m"
const green = "\033[38;5;28m"
const red = "\033[38;5;1m"

// PrintInfo prints a 'message' with cSkyyblue color text
// message: Message to be printed
func PrintInfo(message string) {
	fmt.Printf(cSkyyblue + message + noFormat + "\n")
}

// PrintSuccess prints 'message' with green color text
// message: Message to be printed
func PrintSuccess(message string) {
	fmt.Printf(green + "✔ " + message + noFormat + "\n")
}

// PrintError prints 'message' with red color text
// message: Message to be printed
func PrintError(message string) {
	fmt.Fprintf(os.Stderr, red+"✘ "+message+noFormat+"\n")
}

// PrintYAML takes input object, marshal it using YAML format and prints it into standard output
// obj: Mandatory. Reference to the object to be printed
func PrintYAML(obj interface{}) {
	marshal, err := yaml.Marshal(&obj)
	if err != nil {
		PrintIfError(err)

		return
	}

	_, err = os.Stdout.Write(marshal)
	PrintIfError(err)
}

// PrintIfError print the error message if the provided err has valid error value
// err: Optinal. The error to be printed
func PrintIfError(err error) {
	if err != nil {
		PrintError(err.Error())
	}
}
