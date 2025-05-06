package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zahartd/biathlon_competitions_system/internal/config"
	"github.com/zahartd/biathlon_competitions_system/internal/engine"
	"github.com/zahartd/biathlon_competitions_system/internal/events"
	"github.com/zahartd/biathlon_competitions_system/internal/output"
)

func main() {
	cfgPath := flag.String("config", "", "path to JSON config")
	eventsPath := flag.String("events", "", "path to incoming events")
	outlogPath := flag.String("out", "", "path to output log")
	flag.Parse()

	cfg, err := config.Load(cfgPath)
	if err != nil {
		log.Fatalf("Failed to load configs: %s", err.Error())
		flag.Usage()
		os.Exit(1)
	}

	eventsFile, err := os.Open(*eventsPath)
	if err != nil {
		log.Fatalf("Failed to load events: %s", err.Error())
		flag.Usage()
		os.Exit(1)
	}
	defer eventsFile.Close()

	outlogFile, err := os.OpenFile(
		*outlogPath,
		os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
		0o644,
	)
	if err != nil {
		log.Fatalf("Incorrect output log: %s", err.Error())
		flag.Usage()
		os.Exit(1)
	}
	defer outlogFile.Close()

	eventParser := events.NewParser()
	resultLogger := output.NewLogger(outlogFile)
	eventEngine := engine.NewEngine(cfg, resultLogger)

	scanner := bufio.NewScanner(eventsFile)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Parsed line: %s", line)
		event, err := eventParser.ParseEvent(line)
		if err != nil {
			log.Fatalf("Failed to parse event %s: %s", line, err.Error())
			os.Exit(1)
		}
		log.Printf("Parsed event: %v", event)

		err = eventEngine.ProcessEvent(event)
		if err != nil {
			log.Fatalf("Failed to process event %v: %s", event, err.Error())
			os.Exit(1)
		}
		log.Printf("Processed event: %v", event)
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Failed to reading events: %s", err.Error())
	}

	eventEngine.Finilize()
	rows := eventEngine.GetReport()
	for _, r := range rows {
		fmt.Fprint(os.Stdout, r.Format())
	}
}
