package repository

import (
	"context"
	"fmt"
	"saidvandeklundert/agrippa/agrippalogger"
	"strings"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// GetDatabaseByType returns the integer value (index) of the database.
func GetDatabaseByType(db Database) int {
	return int(db)
}

func databaseToString(db Database) string {
	switch db {
	case APPL_DB:
		return "0"
	case ASIC_DB:
		return "1"
	case COUNTERS_DB:
		return "2"
	case LOGLEVEL_DB:
		return "3"
	case CONFIG_DB:
		return "4"
	case PFC_WD_DB:
		return "5"
	case STATE_DB:
		return "6"
	default:
		return "Unknown"
	}
}

/*
Repository for interacting with Redis.

Offers convenient methods to retrieve data structures from the Redis
database maintained by Sonic.
*/
type RedisRepository struct {
	ctx       context.Context
	logger    *zap.SugaredLogger
	clientMap map[Database]*redis.Client
}

// Returns a pointer to a RedisRepository for the given table
func NewRedisRepository() *RedisRepository {

	return &RedisRepository{ctx: context.Background(), logger: agrippalogger.GetLogger(), clientMap: make(map[Database]*redis.Client)}
}

// Retrieves the device metadata from the configuration
func (r *RedisRepository) GetDeviceMetadata() (*DeviceMetadata, error) {
	r.SetRedisDatabaseConnecion(CONFIG_DB)
	var metadata DeviceMetadata
	err := r.clientMap[CONFIG_DB].HGetAll(r.ctx, "DEVICE_METADATA|localhost").Scan(&metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to scan device metadata: %w", err)
	}

	return &metadata, nil
}

// Retrieves the device LLDP local chassis data from APPL_DB
func (r *RedisRepository) GetLldpLocalChassis() (*LldpLocalChassis, error) {
	r.SetRedisDatabaseConnecion(APPL_DB)
	var data LldpLocalChassis
	err := r.clientMap[APPL_DB].HGetAll(r.ctx, "LLDP_LOC_CHASSIS").Scan(&data)
	if err != nil {
		return nil, fmt.Errorf("failed to scan device metadata: %w", err)
	}

	return &data, nil
}

/*
Check if a client for the given database is already part of the connection pool.
If this is the case, do nothing. Otherwise, create it
*/
func (r *RedisRepository) SetRedisDatabaseConnecion(database Database) {
	_, exists := r.clientMap[database]
	if !exists {
		r.clientMap[database] = redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:6379",
			DB:   GetDatabaseByType(database),
		})
	}
}

// Log all keys in target database
func (r *RedisRepository) GetAllKeys(database Database) ([]string, error) {
	r.SetRedisDatabaseConnecion(database)
	keys, err := r.clientMap[database].Keys(r.ctx, "*").Result()
	if err != nil {
		r.logger.Fatalw("Error fetching keys from Redis: %v", err)
		return []string{}, err
	}
	r.logger.Info("Keys in Redis database number ", database)
	return keys, nil
}

// Log all keys in target database
func (r *RedisRepository) DisplayAllKeys(database Database) {
	keys, err := r.GetAllKeys(database)
	if err != nil {
		r.logger.Error(err)
	}
	for _, key := range keys {
		r.logger.Infow(key)
	}
}

/*
Example function that sbscribes a Redis client to a pattern and then continously
logs this.
*/
func (r *RedisRepository) SubDebugLogger(database Database) {

	var builder strings.Builder
	builder.WriteString("__keyspace@")
	builder.WriteString(databaseToString(database))
	builder.WriteString("*__:*")
	redis_subscription_pattern := builder.String()
	fmt.Println(redis_subscription_pattern)
	// Subscribtion client:
	pubsub := r.clientMap[database].PSubscribe(r.ctx, redis_subscription_pattern)
	defer pubsub.Close()

	// Wait for confirmation that subscription is created
	_, err := pubsub.Receive(r.ctx)
	if err != nil {
		r.logger.Fatal(err)

	}

	// Start listening for messages
	ch := pubsub.Channel()

	fmt.Println("Listening for Redis changes...")

	// Continuously read messages from the channel
	for msg := range ch {
		r.logger.Info("Channel: %s, Message: %s, Pattern: %s\n", msg.Channel, msg.Payload, msg.PayloadSlice)
	}
}

// Runs a 'HGETALL' to retrieve all fields and values of a hash stored in the given database for target key.
// Returns a map of string keys to string values.
func (r *RedisRepository) GetKeyValue(database Database, key string) *redis.MapStringStringCmd {
	r.SetRedisDatabaseConnecion(database)
	value := r.clientMap[database].HGetAll(r.ctx, key)
	return value
}

func GetRepository() {
	logger := agrippalogger.GetLogger()
	logger.Infow("get repository")
	repo := NewRedisRepository()

	meta_data, err := repo.GetDeviceMetadata()
	if err != nil {
		logger.Fatalw("Error GetDeviceMetadata: %v", err)
		return
	}
	logger.Info("Device metadata", meta_data)

	repo.DisplayAllKeys(CONFIG_DB)
	result, err := repo.GetLldpLocalChassis()
	if err != nil {
		fmt.Println("ERROR")
	}

	fmt.Println(result)
	repo.DisplayAllKeys(APPL_DB)
	repo.DisplayAllKeys(STATE_DB)
	getKeyValueResult := repo.GetKeyValue(CONFIG_DB, "PORT|Ethernet0")
	fmt.Println(getKeyValueResult)

	repo.SubDebugLogger(CONFIG_DB)
}
