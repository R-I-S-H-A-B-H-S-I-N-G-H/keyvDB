package db

import (
	"math"
	"sync"

	"github.com/google/btree"
	"github.com/mmcloughlin/geohash"
)

// Point struct
type Point struct {
	Lat  float64
	Lng  float64
	Name string
}

// GeoItem for B-tree storage
type GeoItem struct {
	Hash  string
	Point Point
}

// Less function (B-tree sorting by hash)
func (a GeoItem) Less(b btree.Item) bool {
	return a.Hash < b.(GeoItem).Hash
}

// GeoDB struct
type GeoDB struct {
	mu    sync.RWMutex
	store map[string]*btree.BTree // Each key has a separate B-tree
}

// NewGeoDB initializes a new GeoDB
func NewGeoDB() *GeoDB {
	return &GeoDB{
		store: make(map[string]*btree.BTree),
	}
}

// GEOADD: Adds a location to the DB (Redis-style)
func (db *GeoDB) GeoAdd(key string, lon, lat float64, member string) {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Ensure key exists
	if _, exists := db.store[key]; !exists {
		db.store[key] = btree.New(2)
	}

	// Encode geohash
	hash := geohash.Encode(lat, lon)

	// Insert into B-tree
	db.store[key].ReplaceOrInsert(GeoItem{
		Hash:  hash,
		Point: Point{Lat: lat, Lng: lon, Name: member},
	})
}

// GEODIST: Compute Haversine distance in meters
func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000 // Radius of Earth in meters
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c // Distance in meters
}

// GEODIST command (distance between two locations)
func (db *GeoDB) GeoDist(key, member1, member2 string) float64 {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Ensure key exists
	tree, exists := db.store[key]
	if !exists {
		return -1 // Key doesn't exist
	}

	var point1, point2 *Point

	// Find members
	tree.Ascend(func(i btree.Item) bool {
		item := i.(GeoItem)
		if item.Point.Name == member1 {
			point1 = &item.Point
		} else if item.Point.Name == member2 {
			point2 = &item.Point
		}
		return point1 == nil || point2 == nil // Stop when both are found
	})

	if point1 == nil || point2 == nil {
		return -1 // One of the points not found
	}

	// Compute Haversine distance
	return haversine(point1.Lat, point1.Lng, point2.Lat, point2.Lng)
}
