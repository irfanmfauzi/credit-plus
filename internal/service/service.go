package service

import (
	"context"
	"credit-plus/internal/model/request"
	"credit-plus/internal/repository"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	CreateCustomer(context.Context, request.CreateCustomerRequest) error

	CreateContract(context.Context, request.CreateContactRequest) (int, error)
}

type service struct {
	db              *sqlx.DB
	customerRepo    repository.CustomerRepository
	installmentRepo repository.InstallmentRepository
}

var (
	dbname     = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() *service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	if dbname == "" {
		dbname = "credit_plus"
	}

	if password == "" {
		password = "root"
	}

	if username == "" {
		username = "root"
	}

	if port == "" {
		port = "3306"
	}

	if host == "" {
		host = "127.0.0.1"
	}

	// Opening a driver typically will not attempt to connect to the database.
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}

	db.SetConnMaxLifetime(0)
	db.SetMaxIdleConns(50)
	db.SetMaxOpenConns(50)

	customerRepo := repository.NewCustomerRepo(db)
	installmentRepo := repository.NewInstallmentRepo(db)

	dbInstance = &service{
		customerRepo:    &customerRepo,
		installmentRepo: &installmentRepo,
		db:              db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	// Ping the database
	err := s.db.PingContext(ctx)
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
	dbStats := s.db.Stats()
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

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dbname)
	return s.db.Close()
}

func (s *service) CreateCustomer(ctx context.Context, req request.CreateCustomerRequest) error {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = s.customerRepo.InsertCustomer(ctx, tx, req)
	if err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (s *service) CreateContract(ctx context.Context, req request.CreateContactRequest) (int, error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	defer tx.Rollback()

	isAllowed, err := s.customerRepo.CheckLimitTransaction(ctx, tx, req.OtrPrice, req.CustomerID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if !isAllowed {
		return http.StatusBadGateway, fmt.Errorf("Your limit is less than price")
	}

	err = s.customerRepo.UpdateLimitTransaction(ctx, tx, req.OtrPrice, req.CustomerID)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	contractId, err := s.customerRepo.InsertContract(ctx, tx, req)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	installmentAmount := req.OtrPrice + (req.OtrPrice * req.Bunga / 100) + req.Fee

	reqInstallments := make([]request.CreateCustomerInstallmentRequest, req.Tenor)

	for i := range req.Tenor {
		reqInstallments[i] = request.CreateCustomerInstallmentRequest{
			ContractId:        contractId,
			InstallmentAmount: installmentAmount / req.Tenor,
			PaidAmount:        0,
			DueDate:           time.Now().AddDate(0, i+1, 0),
		}
	}

	err = s.installmentRepo.BulkInsertInstallment(ctx, tx, reqInstallments)
	if err != nil {
		slog.Error("CreateContract Error BulkInsertInstallment", "Error", err)
		return http.StatusInternalServerError, err
	}

	tx.Commit()
	return http.StatusCreated, nil
}
