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
func (r *RedisRepository) DisplayAllKeys(database Database) {
	r.SetRedisDatabaseConnecion(database)
	keys, err := r.clientMap[database].Keys(r.ctx, "*").Result()
	if err != nil {
		r.logger.Fatalw("Error fetching keys from Redis: %v", err)
		return
	}
	r.logger.Infow("Keys in Redis:")
	for _, key := range keys {
		if strings.Contains(key, "LOGGER") {
			continue
		} else if strings.Contains(key, "BUFFER") {
			continue
		} else {
			r.logger.Infow(key)
		}
	}
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
	repo.DisplayAllKeys(APPL_DB)
	repo.DisplayAllKeys(CONFIG_DB)
	result, err := repo.GetLldpLocalChassis()
	if err!=nil{
		fmt.Println("ERROR")
	}

	fmt.Println(result)

}
