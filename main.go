package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	"github.com/HikaruEgashira/mise-osquery-extension/pkg/scanner"
	"github.com/osquery/osquery-go"
	"github.com/osquery/osquery-go/plugin/table"
)

func main() {
	var (
		socket   = flag.String("socket", "", "Path to osquery socket")
		timeout  = flag.Int("timeout", 10, "Timeout in seconds")
		interval = flag.Int("interval", 3, "Interval in seconds")
		verbose  = flag.Bool("verbose", false, "Enable verbose logging")
	)
	flag.Parse()

	// Debug: log all possible socket sources (only if verbose)
	if *verbose {
		log.Printf("Debug: Args=%v, OSQUERY_SOCKET=%s, socket flag=%s",
			flag.Args(), os.Getenv("OSQUERY_SOCKET"), *socket)
	}

	if *socket == "" {
		// Try to get socket from OSQUERY_SOCKET environment variable
		if envSocket := os.Getenv("OSQUERY_SOCKET"); envSocket != "" {
			*socket = envSocket
			if *verbose {
				log.Printf("Using socket from OSQUERY_SOCKET: %s", *socket)
			}
		}
	}

	// When launched by osqueryi --extension, socket is passed as first positional argument
	if *socket == "" && len(flag.Args()) > 0 {
		*socket = flag.Args()[0]
		if *verbose {
			log.Printf("Using socket from positional arg: %s", *socket)
		}
	}

	if *socket == "" {
		log.Fatalln("Usage: version_manager_packages_extension --socket <path>")
	}

	if *verbose {
		log.Printf("Using socket: %s", *socket)
	}

	if *verbose {
		log.Printf("Creating extension manager server...")
	}
	server, err := osquery.NewExtensionManagerServer(
		"mise_packages_extension",
		*socket,
		osquery.ServerTimeout(time.Duration(*timeout)*time.Second),
		osquery.ServerPingInterval(time.Duration(*interval)*time.Second),
	)
	if err != nil {
		log.Fatalf("Error creating extension: %v", err)
	}
	if *verbose {
		log.Printf("Extension manager server created successfully")
	}

	columns := []table.ColumnDefinition{
		table.TextColumn("tool"),
		table.TextColumn("version"),
		table.TextColumn("manager"),
		table.TextColumn("install_path"),
	}

	if *verbose {
		log.Printf("Registering table plugin: mise_packages")
	}
	server.RegisterPlugin(table.NewPlugin("mise_packages", columns, func(ctx context.Context, queryContext table.QueryContext) ([]map[string]string, error) {
		if *verbose {
			log.Printf("Table query called for mise_packages")
		}
		packages, err := scanner.ScanAllManagers()
		if err != nil {
			log.Printf("Error scanning packages: %v", err)
			return []map[string]string{}, nil
		}

		if *verbose {
			log.Printf("Found %d packages", len(packages))
		}
		var results []map[string]string
		for _, pkg := range packages {
			results = append(results, map[string]string{
				"tool":         pkg.Tool,
				"version":      pkg.Version,
				"manager":      pkg.Manager,
				"install_path": pkg.InstallPath,
			})
		}
		return results, nil
	}))
	if *verbose {
		log.Printf("Table plugin registered successfully")
	}

	log.Printf("Starting mise_packages extension on socket: %s", *socket)
	log.Printf("Registered table: mise_packages")
	if err := server.Run(); err != nil {
		log.Fatalf("Error running extension: %v", err)
	}
	log.Printf("Extension stopped")
}
