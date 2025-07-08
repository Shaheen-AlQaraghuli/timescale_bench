package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

const (
	timescaleURL = "postgres://postgres:postgres@localhost:5432/shift_db?sslmode=disable"
	totalRows    = 100_000
	concurrency  = 10
)

type allocation struct {
	allocationID    int64
	allocType       string
	durationMinutes int
	date            time.Time
	opItemID        int64
	zoneID          int64
	marketID        int64
}

func generateAllocations(startID int64, count int) []allocation {
	rand.Seed(time.Now().UnixNano())
	baseTime := time.Now().AddDate(0, 0, -7) // Start from 7 days ago
	types := []string{"b2c", "b2b", "travel"}
	allocs := make([]allocation, count)

	for i := 0; i < count; i++ {
		id := startID + int64(i)
		allocType := types[rand.Intn(len(types))]
		opItemID := int64(rand.Intn(50) + 1)
		marketID := int64(rand.Intn(5) + 1)
		zoneID := int64(rand.Intn(10) + 1)

		randomMinutes := rand.Intn(7 * 24 * 60)
		date := baseTime.Add(time.Duration(randomMinutes) * time.Minute).Truncate(time.Minute)
		duration := rand.Intn(106) + 15 // 15 to 120 minutes

		allocs[i] = allocation{
			allocationID:    id,
			allocType:       allocType,
			durationMinutes: duration,
			date:            date,
			opItemID:        opItemID,
			zoneID:          zoneID,
			marketID:        marketID,
		}
	}

	return allocs
}

func insertAllocations(db *sql.DB, allocs []allocation) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
		INSERT INTO shift_allocations (
			allocation_id, type, duration_minutes, date, op_item_id, zone_id, market_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range allocs {
		_, err := stmt.Exec(a.allocationID, a.allocType, a.durationMinutes, a.date, a.opItemID, a.zoneID, a.marketID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func main() {
	db, err := sql.Open("postgres", timescaleURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	fmt.Println("Generating allocations...")
	allocations := generateAllocations(1, totalRows)
	fmt.Println("Generation done.")

	start := time.Now()

	var wg sync.WaitGroup
	chunkSize := totalRows / concurrency

	for i := 0; i < concurrency; i++ {
		startIndex := i * chunkSize
		endIndex := startIndex + chunkSize
		if i == concurrency-1 {
			endIndex = totalRows
		}

		wg.Add(1)
		go func(batch []allocation) {
			defer wg.Done()
			err := insertAllocations(db, batch)
			if err != nil {
				log.Printf("Insert batch failed: %v", err)
			}
		}(allocations[startIndex:endIndex])
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Inserted %d rows in %s (%.2f rows/sec)\n", totalRows, elapsed, float64(totalRows)/elapsed.Seconds())
}
