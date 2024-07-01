package database

import (
	"context"
	"fmt"
	"life-gamifying/internal/models"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service interface {
	RDB() *redis.Client
	DB() *gorm.DB
	RDHealth() map[string]string
	PHealth() map[string]string
	Close()
}

type service struct {
	pgdb *gorm.DB
	rdb  *redis.Client
}

var (
	raddress  = os.Getenv("RDB_ADDRESS")
	rport     = os.Getenv("RDB_PORT")
	rpassword = os.Getenv("RDB_PASSWORD")
	rdatabase = os.Getenv("RDB_DATABASE")
)

func New() Service {
	var pgdb *gorm.DB
	var rdb *redis.Client
	num, err := strconv.Atoi(rdatabase)
	if err != nil {
		log.Fatalf(fmt.Sprintf("database incorrect %v", err))
	}

	fullAddress := fmt.Sprintf("%s:%s", raddress, rport)

	rdb = redis.NewClient(&redis.Options{
		Addr:     fullAddress,
		Password: rpassword,
		DB:       num,
		// Note: It's important to add this for a secure connection. Most cloud services that offer Redis should already have this configured in their services.
		// For manual setup, please refer to the Redis documentation: https://redis.io/docs/latest/operate/oss_and_stack/management/security/encryption/
		// TLSConfig: &tls.Config{
		// 	MinVersion:   tls.VersionTLS12,
		// },
	})

	pgdb, err = gorm.Open(postgres.Open(os.Getenv("PDB_URL")), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatalf(fmt.Sprintf("database down: %v", err))
	}

	// Auto create record for level thresholds
	// This is a one-time operation to create the level thresholds record in the database.
	// The level thresholds are used to calculate the experience points required to reach the next level.
	// The level thresholds are stored in the database and are used to calculate the experience points required to reach the next level.

	err = pgdb.AutoMigrate(
		&models.User{},
		&models.Admin{},
		&models.Daily{},
		&models.Habit{},
		&models.Inventory{},
		&models.InventoryItem{},
		&models.Item{},
		&models.Player{},
		&models.Quest{},
		// &models.Rank{},
		&models.Reward{},
		&models.RewardItem{},
		&models.LevelThresholds{},
	)

	if err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// TODO: Find formula for exp to next level
	// Check if the level thresholds record exists in the database.
	// var levelThresholds models.LevelThresholds
	// err = pgdb.First(&levelThresholds).Error

	// // If the level thresholds record does not exist, create it.
	// if err != nil {
	// 	baseExp := 25.0
	// 	// Create the level thresholds record with default values.
	// 	for i := 1; i <= 100; i++ {
	// 		levelThresholds.Level = uint(i)
	// 		if i > 1 {
	// 			levelThresholds.ExpToNext = uint(baseExp * math.Log2(float64(i)))
	// 		} else {
	// 			levelThresholds.ExpToNext = 0
	// 		}
	// 		err = pgdb.Create(&levelThresholds).Error
	// 		if err != nil {
	// 			log.Fatalf("failed to create level thresholds record: %v", err)
	// 		}
	// 	}
	// }

	// Add the Redis and PostgreSQL clients to the service struct
	s := &service{rdb: rdb, pgdb: pgdb}

	return s
}

func (s *service) RDB() *redis.Client {
	return s.rdb
}

func (s *service) DB() *gorm.DB {
	return s.pgdb
}

// Health returns the health status and statistics of the Redis server.
func (s *service) RDHealth() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Default is now 5s
	defer cancel()

	stats := make(map[string]string)

	// Check Redis health and populate the stats map
	stats = s.checkRedisHealth(ctx, stats)

	return stats
}

// checkRedisHealth checks the health of the Redis server and adds the relevant statistics to the stats map.
func (s *service) checkRedisHealth(ctx context.Context, stats map[string]string) map[string]string {
	// Ping the Redis server to check its availability.
	pong, err := s.rdb.Ping(ctx).Result()
	// Note: By extracting and simplifying like this, `log.Fatalf(fmt.Sprintf("db down: %v", err))`
	// can be changed into a standard error instead of a fatal error.
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	// Redis is up
	stats["redis_status"] = "up"
	stats["redis_message"] = "It's healthy"
	stats["redis_ping_response"] = pong

	// Retrieve Redis server information.
	info, err := s.rdb.Info(ctx).Result()
	if err != nil {
		stats["redis_message"] = fmt.Sprintf("Failed to retrieve Redis info: %v", err)
		return stats
	}

	// Parse the Redis info response.
	redisInfo := parseRedisInfo(info)

	// Get the pool stats of the Redis client.
	poolStats := s.rdb.PoolStats()

	// Prepare the stats map with Redis server information and pool statistics.
	// Note: The "stats" map in the code uses string keys and values,
	// which is suitable for structuring and serializing the data for the frontend (e.g., JSON, XML, HTMX).
	// Using string types allows for easy conversion and compatibility with various data formats,
	// making it convenient to create health stats for monitoring or other purposes.
	// Also note that any raw "memory" (e.g., used_memory) value here is in bytes and can be converted to megabytes or gigabytes as a float64.
	stats["redis_version"] = redisInfo["redis_version"]
	stats["redis_mode"] = redisInfo["redis_mode"]
	stats["redis_connected_clients"] = redisInfo["connected_clients"]
	stats["redis_used_memory"] = redisInfo["used_memory"]
	stats["redis_used_memory_peak"] = redisInfo["used_memory_peak"]
	stats["redis_uptime_in_seconds"] = redisInfo["uptime_in_seconds"]
	stats["redis_hits_connections"] = strconv.FormatUint(uint64(poolStats.Hits), 10)
	stats["redis_misses_connections"] = strconv.FormatUint(uint64(poolStats.Misses), 10)
	stats["redis_timeouts_connections"] = strconv.FormatUint(uint64(poolStats.Timeouts), 10)
	stats["redis_total_connections"] = strconv.FormatUint(uint64(poolStats.TotalConns), 10)
	stats["redis_idle_connections"] = strconv.FormatUint(uint64(poolStats.IdleConns), 10)
	stats["redis_stale_connections"] = strconv.FormatUint(uint64(poolStats.StaleConns), 10)
	stats["redis_max_memory"] = redisInfo["maxmemory"]

	// Calculate the number of active connections.
	// Note: We use math.Max to ensure that activeConns is always non-negative,
	// avoiding the need for an explicit check for negative values.
	// This prevents a potential underflow situation.
	activeConns := uint64(math.Max(float64(poolStats.TotalConns-poolStats.IdleConns), 0))
	stats["redis_active_connections"] = strconv.FormatUint(activeConns, 10)

	// Calculate the pool size percentage.
	poolSize := s.rdb.Options().PoolSize
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	poolSizePercentage := float64(connectedClients) / float64(poolSize) * 100
	stats["redis_pool_size_percentage"] = fmt.Sprintf("%.2f%%", poolSizePercentage)

	// Evaluate Redis stats and update the stats map with relevant messages.
	return s.evaluateRedisStats(redisInfo, stats)
}

// evaluateRedisStats evaluates the Redis server statistics and updates the stats map with relevant messages.
func (s *service) evaluateRedisStats(redisInfo, stats map[string]string) map[string]string {
	poolSize := s.rdb.Options().PoolSize
	poolStats := s.rdb.PoolStats()
	connectedClients, _ := strconv.Atoi(redisInfo["connected_clients"])
	highConnectionThreshold := int(float64(poolSize) * 0.8)

	// Check if the number of connected clients is high.
	if connectedClients > highConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of connected clients"
	}

	// Check if the number of stale connections exceeds a threshold.
	minStaleConnectionsThreshold := 500
	if int(poolStats.StaleConns) > minStaleConnectionsThreshold {
		stats["redis_message"] = fmt.Sprintf("Redis has %d stale connections.", poolStats.StaleConns)
	}

	// Check if Redis is using a significant amount of memory.
	usedMemory, _ := strconv.ParseInt(redisInfo["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(redisInfo["maxmemory"], 10, 64)
	if maxMemory > 0 {
		usedMemoryPercentage := float64(usedMemory) / float64(maxMemory) * 100
		if usedMemoryPercentage >= 90 {
			stats["redis_message"] = "Redis is using a significant amount of memory"
		}
	}

	// Check if Redis has been recently restarted.
	uptimeInSeconds, _ := strconv.ParseInt(redisInfo["uptime_in_seconds"], 10, 64)
	if uptimeInSeconds < 3600 {
		stats["redis_message"] = "Redis has been recently restarted"
	}

	// Check if the number of idle connections is high.
	idleConns := int(poolStats.IdleConns)
	highIdleConnectionThreshold := int(float64(poolSize) * 0.7)
	if idleConns > highIdleConnectionThreshold {
		stats["redis_message"] = "Redis has a high number of idle connections"
	}

	// Check if the connection pool utilization is high.
	poolUtilization := float64(poolStats.TotalConns-poolStats.IdleConns) / float64(poolSize) * 100
	highPoolUtilizationThreshold := 90.0
	if poolUtilization > highPoolUtilizationThreshold {
		stats["redis_message"] = "Redis connection pool utilization is high"
	}

	return stats
}

// parseRedisInfo parses the Redis info response and returns a map of key-value pairs.
func parseRedisInfo(info string) map[string]string {
	result := make(map[string]string)
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			result[key] = value
		}
	}
	return result
}

func (s *service) Close() {
	sqlDB, _ := s.pgdb.DB()
	sqlDB.Close()
	s.rdb.Close()
}

// POSTGRES HEALTH
// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) PHealth() map[string]string {
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Use the DB method to get a *sql.DB connection handle
	sqlDB, err := s.pgdb.DB()
	if err != nil {
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf(fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		return stats
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get database stats (like open connections, in use, idle, etc.)
	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	// Evaluate stats to provide a health message
	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats
}
