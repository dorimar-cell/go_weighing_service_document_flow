package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"go_weighing_service/internal/config"
)

func main() {
	days := flag.Int("days", 180, "delete audit records older than N days")
	dryRun := flag.Bool("dry-run", false, "only count records without deleting")
	flag.Parse()

	cfg := config.Load()
	ctx := context.Background()

	pool, err := config.ConnectDB(ctx, cfg)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	tables := []string{"weighing_audit", "dictionary_audit"}
	for _, table := range tables {
		var count int
		if err := pool.QueryRow(ctx,
			fmt.Sprintf("SELECT count(*) FROM %s WHERE created_at < now() - make_interval(days => $1)", table),
			*days,
		).Scan(&count); err != nil {
			log.Fatalf("count %s: %v", table, err)
		}

		if *dryRun {
			fmt.Printf("%s: %d records older than %d days\n", table, count, *days)
			continue
		}

		tag, err := pool.Exec(ctx,
			fmt.Sprintf("DELETE FROM %s WHERE created_at < now() - make_interval(days => $1)", table),
			*days,
		)
		if err != nil {
			log.Fatalf("delete %s: %v", table, err)
		}
		fmt.Printf("%s: deleted %d records older than %d days\n", table, tag.RowsAffected(), *days)
	}

	if *dryRun {
		os.Exit(0)
	}
}
