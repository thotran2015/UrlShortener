package store

import (
	"context"
	"fmt"
	"github.com/gocql/gocql"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

// Define struct wrapper around raw Redis client

type StorageService struct {
	redisClient      *redis.Client
	cassandraCluster *gocql.ClusterConfig
}

// Top level declarations for storeService and redis context
var (
	storeService = &StorageService{}
	ctx          = context.Background()
)

// Set expiration for each  key, value pair in redis
const storeExpiration = 12 * time.Hour
const cassandraKeyspace = "url_shortener"
const cassandraTable = "url_mapping"

func InitializeStore() *StorageService {
	// Connect to Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	// Check if we can ping the Redis service
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Println(fmt.Sprintf("Error connecting to Redis: %v\n", err))
	}
	fmt.Printf("Redis started successfully: pong msg: = %s", pong)
	storeService.redisClient = redisClient

	// Connect to Cassandra
	cluster := gocql.NewCluster("localhost:9042")
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("Failed to create a Cassandra session. Error: %v\n", err)
	}
	defer session.Close()
	// Create a keyspace if it doesn't exist
	createKeyspaceQuery := fmt.Sprintf("CREATE KEYSPACE IF NOT EXISTS %s WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': 1}", cassandraKeyspace)
	err = session.Query(createKeyspaceQuery).Exec()
	if err != nil {
		log.Printf("Failed to create a Cassandra keyspace. Error: %v\n", err)
	}
	cluster.Keyspace = cassandraKeyspace
	// Create a table if it doesn't exist
	createTableQuery := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (short_url text PRIMARY KEY, original_url text)", cassandraKeyspace, cassandraTable)
	err = session.Query(createTableQuery).Exec()
	if err != nil {
		log.Printf("Failed to create a Cassandra table. Error: %v\n", err)
	}
	storeService.cassandraCluster = cluster
	return storeService
}

func SaveUrlMapping(shortUrl string, originalUrl string) {
	err := storeService.redisClient.Set(ctx, shortUrl, originalUrl, storeExpiration).Err()
	if err != nil {
		log.Println(fmt.Sprintf("Failed saving key url in redis. Error: %v - shortUrl: %s - originalUrl: %s\n", err, shortUrl, originalUrl))
	}
	session, err := storeService.cassandraCluster.CreateSession()
	if err != nil {
		log.Printf("Failed to create a Cassandra session. Error: %v\n", err)
	}
	defer session.Close()
	insertQuery := fmt.Sprintf("INSERT INTO %s (short_url, original_url) VALUES (?, ?)", cassandraTable)
	err = session.Query(insertQuery, shortUrl, originalUrl).Exec()
	if err != nil {
		log.Printf("Failed to insert a new url mapping into Cassandra. Error: %v\n", err)
	}
}

func GetUrlMapping(shortUrl string) string {
	// Check if the short URL exists in Redis
	res, err := storeService.redisClient.Get(ctx, shortUrl).Result()
	if err != nil {
		log.Printf("Failed to get url from redis. Error: %v\n", err)
	} else {
		return res
	}
	// If not found in Redis, check in Cassandra
	session, err := storeService.cassandraCluster.CreateSession()
	if err != nil {
		log.Printf("Failed to create a Cassandra session. Error: %v\n", err)
	}
	defer session.Close()
	selectQuery := fmt.Sprintf("SELECT original_url FROM %s WHERE short_url = ? LIMIT 1", cassandraTable)
	session.Query(selectQuery, shortUrl).Scan(&res)
	return res
}
