package integration_tests

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"container-manager/internal/application"
	"container-manager/internal/domain/infrastructure"
	"container-manager/internal/infrastructure/repository"
	"container-manager/internal/server"
	"container-manager/internal/server/handler"
	"container-manager/internal/server/middleware"
	"container-manager/pkg/config"
	"container-manager/pkg/postgres"

	"github.com/bwmarrin/snowflake"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

var (
	testDB *sql.DB
	cfg    config.Config
)

func TestMain(m *testing.M) {
	var err error

	cfg, err = config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Printf("Loaded config: %+v", cfg)

	testDB, err = postgres.NewClient(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer testDB.Close()

	if err := testDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	exitCode := m.Run()

	os.Exit(exitCode)
}

func setupTestDB(t *testing.T) {
	t.Helper()
	truncateTables(t)
}

func truncateTables(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	tables := []string{"jobs", "container_user", "users"}

	for _, table := range tables {
		_, err := testDB.ExecContext(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			t.Logf("Failed to truncate table %s: %v", table, err)
		}
	}
}

func setupServer(t *testing.T, runtime infrastructure.ContainerRuntime) *gin.Engine {
	idNode, err := snowflake.NewNode(1)
	require.NoError(t, err)

	fileStorage := repository.NewLocalFileStorage(os.TempDir())

	userRepo := repository.NewUserRepository(testDB)
	containerUserRepo := repository.NewContainerUserRepository(testDB)
	jobRepo := repository.NewJobRepository(testDB)

	jwtSecret := cfg.Server.JWTSecret
	userService := application.NewUserService(userRepo, idNode, jwtSecret)
	fileService := application.NewFileService(fileStorage)
	containerService := application.NewContainerService(runtime, containerUserRepo, jobRepo)
	jobService := application.NewJobService(jobRepo)

	authMiddleware := middleware.NewAuthMiddleware(jwtSecret)
	userHandler := handler.NewUserHandler(userService)
	containerHandler := handler.NewContainerHandler(containerService)
	fileHandler := handler.NewFileHandler(fileService)
	jobHandler := handler.NewJobHandler(jobService)

	r := gin.Default()
	gin.DisableConsoleColor()
	gin.SetMode(gin.TestMode)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowHeaders = []string{"Authorization", "Content-Type", "Accept"}
	r.Use(cors.New(corsConfig))

	server.RegisterRoutes(r, userHandler, containerHandler, fileHandler, jobHandler, authMiddleware)

	return r
}
