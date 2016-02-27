package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Viperoo/golog"
)

var configfile = flag.String("c", "gitlab.toml", "Configuration file")
var logger log.Logger
var debug = flag.Bool("d", false, "Debug mode")
var logfile = flag.String("l", "issue-to-gitlab.log", "Log file")
var apiVerions = "/api/v3/"

func main() {

	/*
	* Parse flags
	 */
	flag.Parse()
	/*
	* Set logger level
	 */
	setLogger()

	ReadConfig(*configfile)

	reader := bufio.NewReader(os.Stdin)

	listProjects()

	fmt.Print("Enter ID of project: ")
	project, _ := reader.ReadString('\n')

	addIssue(project)

	for confirm("\n\nDo you want add other issue to this project? Y/N ") {
		addIssue(project)
	}

	fmt.Println("\nExiting...")
}

func setLogger() {
	file, err := os.OpenFile(*logfile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Print("Error log is not wrtiable.")
		os.Exit(1)
	}
	var multi io.Writer
	if *debug == true {
		multi = io.MultiWriter(file, os.Stdout)
	} else {
		multi = io.MultiWriter(file)
	}

	logger, _ = log.NewLogger(multi,
		log.TIME_FORMAT_SEC,
		log.LOG_FORMAT_SIMPLE,
		log.LogLevel_Debug)
}

func confirm(message string) bool {
	fmt.Print(message)

	reader := bufio.NewReader(os.Stdin)
	c, _ := reader.ReadByte()

	if c == []byte("Y")[0] || c == []byte("y")[0] {
		return true
	} else {
		return false
	}

}
