package main

import (
	"log"
	"net"
	"os"
	"time"

	"git.epam.com/go-language-global-mentoring-program/internal/cache"
	memorycache "git.epam.com/go-language-global-mentoring-program/internal/cache/memory"
	"git.epam.com/go-language-global-mentoring-program/internal/cache/redis"
	cacheStorage "git.epam.com/go-language-global-mentoring-program/internal/cache/redis/storage"
	"git.epam.com/go-language-global-mentoring-program/internal/grpc/issue"
	"git.epam.com/go-language-global-mentoring-program/internal/grpc/project"
	"git.epam.com/go-language-global-mentoring-program/internal/grpc/user"
	"git.epam.com/go-language-global-mentoring-program/internal/repo"
	"git.epam.com/go-language-global-mentoring-program/internal/repo/cached"
	"git.epam.com/go-language-global-mentoring-program/internal/repo/memory"
	"git.epam.com/go-language-global-mentoring-program/internal/repo/postgres"
	repoStorage "git.epam.com/go-language-global-mentoring-program/internal/repo/postgres/storage"
	"git.epam.com/go-language-global-mentoring-program/pkg/gen"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

const cacheTTL = 15 * time.Minute

func main() {
	_ = godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	//create validators
	userValidator := &user.AppValidator{}
	issueValidator := &issue.AppValidator{}
	projectValidator := &project.AppValidator{}

	store := createStore()

	if os.Getenv("SEED_DATABASE") == "true" {
		seedDatabase(store)
	}

	if cacheStore := createCache(); cacheStore != nil {
		store = cached.NewStore(store, cacheStore, cacheTTL)
	}

	userServer := user.NewServer(userValidator, store)
	issueServer := issue.NewServer(issueValidator, store)
	projectServer := project.NewServer(projectValidator, store)

	gen.RegisterUserServiceServer(grpcServer, userServer)
	gen.RegisterIssueServiceServer(grpcServer, issueServer)
	gen.RegisterProjectServiceServer(grpcServer, projectServer)

	reflection.Register(grpcServer)

	log.Printf("%s", "gRPC server is running on port "+port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("falied to serve gRPC server : %v", err)
	}
}

func createStore() repo.Store {
	storageType := os.Getenv("STORAGE_TYPE")

	switch storageType {
	case "postgres":
		log.Println("Selected storage implementation: POSTGRES (GORM)")

		dbConfig := &repoStorage.ConfigDB{
			Host:     getEnv("DB_HOST", "postgres"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBname:   getEnv("DB_NAME", "issue_tracker"),
			Port:     getEnv("DB_PORT", "5432"),
			SSLmode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		}

		connector := func() (*gorm.DB, error) {
			return repoStorage.ConnectDB(dbConfig)
		}

		store, err := postgres.InitStore(connector)
		if err != nil {
			log.Fatalf("Critical: failed to initialize postgres storage: %v", err)
		}
		return store

	case "memory", "":
		log.Println("Selected storage implementation: IN-MEMORY")
		store, err := memory.NewStore()
		if err != nil {
			log.Fatalf("Critical: failed to create memory store: %v", err)
		}
		return store

	default:
		log.Fatalf("Critical: unknown STORAGE_TYPE %q. Supported values: 'memory', 'postgres'", storageType)
		return nil
	}
}

func createCache() cache.Cache {
	cacheType := os.Getenv("CACHE_TYPE")

	switch cacheType {
	case "redis":
		log.Println("Selected cache implementation: Redis")

		cacheConfig := &cacheStorage.ConfigRedis{
			Host:     getEnv("REDIS_HOST", "redis"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       0,
		}

		client, err := cacheStorage.ConnectRedis(cacheConfig)
		if err != nil {
			log.Fatalf("Critical: failed to connect to redis: %v", err)
		}

		cacheStore, err := redis.NewCache(client)
		if err != nil {
			log.Fatalf("Critical: failed to initialize redis cache: %v", err)
		}

		return cacheStore

	case "memory":
		log.Println("Selected cache implementation: IN-MEMORY")
		return memorycache.NewCache()

	default:
		log.Fatalf("Critical: unknown CACHE_TYPE %q. Supported values: 'redis', 'memory'", cacheType)
		return nil
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
