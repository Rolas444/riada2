package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/riada2/config"
	_ "github.com/riada2/docs" // Importa los documentos de Swagger generados
	"github.com/riada2/internal/core/domain"
	"github.com/riada2/internal/core/ports"
	"github.com/riada2/internal/handlers"
	"github.com/riada2/internal/repository"
	"github.com/riada2/internal/router"
	"github.com/riada2/internal/services"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// @title           Riada2 API
// @version         1.0
// @description     Esta es la API para el proyecto Riada2, con autenticación y gestión de usuarios.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Soporte de API
// @contact.url    http://www.example.com/support
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:3001
// @BasePath  /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Cargar configuración
	cfg, err := config.LoadConfig("./.env")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	// Crear la base de datos si no existe
	err = createDatabaseIfNotExists(cfg)
	if err != nil {
		log.Fatalf("database setup failed: %v", err)
	}

	// Conectar a la base de datos
	db, err := gorm.Open(postgres.Open(cfg.DBSource), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatalf("could not connect to db: %v", err)
	}

	// Migrar el esquema
	err = db.AutoMigrate(&domain.User{}, &domain.Person{}, &domain.Address{}, &domain.Phone{}, &domain.Membership{})
	if err != nil {
		log.Fatalf("could not migrate db: %v", err)
	}

	// Crear manualmente las restricciones de clave foránea para evitar dependencias circulares
	// durante la migración inicial.
	if !db.Migrator().HasConstraint(&domain.User{}, "Person") {
		db.Migrator().CreateConstraint(&domain.User{}, "Person")
	}
	if !db.Migrator().HasConstraint(&domain.Person{}, "User") {
		db.Migrator().CreateConstraint(&domain.Person{}, "User")
	}

	// Inyección de dependencias (unión de piezas)
	userRepo := repository.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo, cfg.JWTSecret)
	userHandler := handlers.NewUserHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, cfg)

	personRepo := repository.NewGormPersonRepository(db)
	personService := services.NewPersonService(personRepo)
	personHandler := handlers.NewPersonHandler(personService)

	addressRepo := repository.NewGormAddressRepository(db)
	addressService := services.NewAddressService(addressRepo, personRepo)
	addressHandler := handlers.NewAddressHandler(addressService)

	phoneRepo := repository.NewGormPhoneRepository(db)
	phoneService := services.NewPhoneService(phoneRepo, personRepo)
	phoneHandler := handlers.NewPhoneHandler(phoneService)

	membershipRepo := repository.NewGormMembershipRepository(db)
	membershipService := services.NewMembershipService(membershipRepo)
	membershipHandler := handlers.NewMembershipHandler(membershipService, personService)

	createDefaultAdmin(db, userRepo, cfg)

	// Configuración de Fiber
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		// En un entorno de producción, deberías restringir esto a tu dominio de frontend.
		// Ejemplo: AllowOrigins: "http://localhost:5173, http://mi-frontend.com",
		AllowOrigins: "*",
		// AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowHeaders: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))
	app.Use(logger.New())

	router.SetupRoutes(app, authHandler, userHandler, personHandler, addressHandler, phoneHandler, membershipHandler, cfg)

	log.Fatal(app.Listen(fmt.Sprintf(":%s", cfg.AppPort)))
}

// createDatabaseIfNotExists se conecta a la base de datos 'postgres' por defecto
// para verificar si la base de datos de la aplicación existe, y la crea si no.
func createDatabaseIfNotExists(cfg *config.Config) error {
	var dbName string
	// Extraer el nombre de la base de datos del DSN
	parts := strings.Split(cfg.DBSource, " ")
	for _, part := range parts {
		if strings.HasPrefix(part, "dbname=") {
			dbName = strings.TrimPrefix(part, "dbname=")
			break
		}
	}

	if dbName == "" {
		return errors.New("dbname not found in DSN")
	}

	// Crear un DSN para la base de datos 'postgres'
	// Reemplazamos dbname=... por dbname=postgres
	tempDSN := strings.Replace(cfg.DBSource, "dbname="+dbName, "dbname=postgres", 1)

	// Conectarse a la base de datos 'postgres'
	db, err := gorm.Open(postgres.Open(tempDSN), &gorm.Config{
		// Silenciar el logger para esta conexión temporal
		Logger: gormlogger.Default.LogMode(gormlogger.Silent),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to 'postgres' database to check existence: %w", err)
	}

	// Obtener el sql.DB subyacente para poder cerrarlo
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB for temporary connection: %w", err)
	}
	defer sqlDB.Close()

	// Verificar si la base de datos de la aplicación existe
	var exists bool
	err = db.Raw("SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = ?)", dbName).Scan(&exists).Error
	if err != nil {
		return fmt.Errorf("failed to check if database '%s' exists: %w", dbName, err)
	}

	if !exists {
		log.Printf("Database '%s' not found, creating...", dbName)
		// Usar Exec para crear la base de datos.
		err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", dbName)).Error
		if err != nil {
			return fmt.Errorf("failed to create database '%s': %w", dbName, err)
		}
		log.Printf("Database '%s' created successfully.", dbName)
	}

	return nil
}

func createDefaultAdmin(db *gorm.DB, userRepo ports.UserRepository, cfg *config.Config) {
	var userCount int64
	db.Model(&domain.User{}).Count(&userCount)

	if userCount == 0 {
		log.Println("No users found. Creating default admin user...")

		if cfg.DefaultAdminUser == "" || cfg.DefaultAdminPassword == "" {
			log.Println("DEFAULT_ADMIN_USER or DEFAULT_ADMIN_PASSWORD not set in .env. Skipping creation.")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cfg.DefaultAdminPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("Failed to hash default admin password: %v", err)
		}

		admin := &domain.User{
			Username:     cfg.DefaultAdminUser,
			PasswordHash: string(hashedPassword),
			Role:         domain.AdminRole,
		}

		if err := userRepo.Save(admin); err != nil {
			log.Fatalf("Failed to create default admin user: %v", err)
		}
		log.Println("Default admin user created successfully.")
	}
}
