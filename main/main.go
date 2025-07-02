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
	timescaleURL       = "postgres://postgres:postgres@localhost:5432/shift_db?sslmode=disable"
	postgresURL 	   = "postgres://postgres:postgres@localhost:5433/shift_db?sslmode=disable"
	totalRows   = 100000
	batchSize   = 10000
	concurrency = 10
)

type allocation struct {
	allocationID int64
	opItemID     int64
	marketID     int64
	zoneID       int64
	allocType    string
	startsAt     time.Time
	endsAt       time.Time
}

// Generates all allocations sequentially
func generateAllocations(startID int64, count int) []allocation {
	rand.Seed(time.Now().UnixNano())
	baseTime := time.Now().AddDate(0, 0, -7) // 7 days ago
	types := []string{"b2c", "b2b", "travel"}
	allocs := make([]allocation, count)

	for i := 0; i < count; i++ {
		id := startID + int64(i)
		opItemID := int64(rand.Intn(50) + 1)
		marketID := int64(rand.Intn(5) + 1)
		zoneID := int64(rand.Intn(10) + 1)
		allocType := types[rand.Intn(len(types))]

		randomMinutes := rand.Intn(7 * 24 * 60)
		startsAt := baseTime.Add(time.Duration(randomMinutes) * time.Minute)
		duration := time.Duration(rand.Intn(105)+15) * time.Minute
		endsAt := startsAt.Add(duration)

		allocs[i] = allocation{
			allocationID: id,
			opItemID:     opItemID,
			marketID:     marketID,
			zoneID:       zoneID,
			allocType:    allocType,
			startsAt:     startsAt,
			endsAt:       endsAt,
		}
	}
	return allocs
}

func insertBatch(db *sql.DB, batch []allocation) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`
		INSERT INTO shift_allocations (
			op_item_id, market_id, zone_id, type, starts_at, ends_at
		) VALUES ($1, $2, $3, $4, $5, $6)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, a := range batch {
		_, err := stmt.Exec(a.opItemID, a.marketID, a.zoneID, a.allocType, a.startsAt, a.endsAt)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func main() {
	db, err := sql.Open("postgres", postgresURL)
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
			err := insertBatch(db, batch)
			if err != nil {
				log.Printf("Insert batch failed: %v", err)
			}
		}(allocations[startIndex:endIndex])
	}

	wg.Wait()
	elapsed := time.Since(start)
	fmt.Printf("Inserted %d rows in %s (%.2f rows/sec)\n", totalRows, elapsed, float64(totalRows)/elapsed.Seconds())
}

