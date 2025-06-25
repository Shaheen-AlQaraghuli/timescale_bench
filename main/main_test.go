package main

import (
	"database/sql"
	"math/rand"
	"testing"
)

func BenchmarkInsertSingleRowParallelTimescale(b *testing.B) {
	db, err := sql.Open("postgres", timescaleURL)
	if err != nil {
		b.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	allocations := generateAllocations(1_000_000, totalRows)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			randomIdx := rand.Intn(len(allocations))
			a := allocations[randomIdx]

			_, err := db.Exec(`
				INSERT INTO shift_allocations (
					op_item_id, market_id, zone_id, type, starts_at, ends_at
				) VALUES ($1, $2, $3, $4, $5, $6)
			`, a.opItemID, a.marketID, a.zoneID, a.allocType, a.startsAt, a.endsAt)

			if err != nil {
				b.Errorf("Insert failed: %v", err)
			}
		}
	})
}

func BenchmarkInsertSingleRowParallelPostgres(b *testing.B) {
	db, err := sql.Open("postgres", postgresURL)
	if err != nil {
		b.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)

	allocations := generateAllocations(1_000_000, totalRows)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			randomIdx := rand.Intn(len(allocations))
			a := allocations[randomIdx]

			_, err := db.Exec(`
				INSERT INTO shift_allocations (
					op_item_id, market_id, zone_id, type, starts_at, ends_at
				) VALUES ($1, $2, $3, $4, $5, $6)
			`, a.opItemID, a.marketID, a.zoneID, a.allocType, a.startsAt, a.endsAt)

			if err != nil {
				b.Errorf("Insert failed: %v", err)
			}
		}
	})
}

func BenchmarkUpdateByIDTimescale(b *testing.B) {
	db, err := sql.Open("postgres", timescaleURL)
	if err != nil {
		b.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	maxID := 1000
	types := []string{"b2c", "b2b", "travel"}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			id := rand.Intn(maxID) + 1 // +1 to avoid 0 if ids start from 1
			newType := types[rand.Intn(len(types))]

			_, err := db.Exec(`
				UPDATE shift_allocations
				SET type = $1
				WHERE id = $2
			`, newType, id)

			if err != nil {
				b.Errorf("Update failed: %v", err)
			}
		}
	})
}

func BenchmarkUpdateByIDPostgres(b *testing.B) {
 db, err := sql.Open("postgres", postgresURL)
 if err != nil {
  b.Fatalf("Failed to connect to DB: %v", err)
 }
 defer db.Close()

 maxID := 1000
 types := []string{"b2c", "b2b", "travel"}

 b.ResetTimer()

 b.RunParallel(func(pb *testing.PB) {
  for pb.Next() {
   id := rand.Intn(maxID) + 1 // +1 to avoid 0 if ids start from 1
   newType := types[rand.Intn(len(types))]

   _, err := db.Exec(`
	UPDATE shift_allocations
	SET type = $1
	WHERE id = $2
   `, newType, id)

   if err != nil {
	b.Errorf("Update failed: %v", err)
   }
  }
 })
}