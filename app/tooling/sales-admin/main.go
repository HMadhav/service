// This program performs administrative tasks for the garage sale service.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/HMadhav/service/business/data/schema"
	"github.com/HMadhav/service/business/data/store/user"
	"github.com/HMadhav/service/business/sys/auth"
	"github.com/HMadhav/service/business/sys/database"
	"github.com/HMadhav/service/foundation/keystore"
	"github.com/HMadhav/service/foundation/logger"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

func main() {

	// Construct the application logger.
	log, err := logger.New("ADMIN")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	cfg := database.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost",
		Name:         "postgres",
		MaxIdleConns: 0,
		MaxOpenConns: 0,
		DisableTLS:   true,
	}

	if err := migrate(cfg); err != nil {
		fmt.Println("migrating database: %w", err)
		os.Exit(1)
	}

	if err := seed(cfg); err != nil {
		fmt.Println("seeding database: %w", err)
		os.Exit(1)
	}

	if err := seed(cfg); err != nil {
		fmt.Println("seeding database: %w", err)
		os.Exit(1)
	}

	//Hard coded the data according to database
	// userID := "6ba7b810-9dad-11d1-80b4-00c04fd430c8"
	// kid := "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"

	// if err := GenToken(log, cfg, userID, kid); err != nil {
	// 	fmt.Println("genToken: %w", err)
	// 	os.Exit(1)
	// }

}

// Migrate creates the schema in the database.
func migrate(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Migrate(ctx, db); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	fmt.Println("migrations complete")
	return nil
}

// Seed loads test data into the database.
func seed(cfg database.Config) error {
	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := schema.Seed(ctx, db); err != nil {
		return fmt.Errorf("seed database: %w", err)
	}

	fmt.Println("seed data complete")
	return nil
}

// GenToken generates a JWT for the specified user.
func GenToken(log *zap.SugaredLogger, cfg database.Config, userID string, kid string) error {
	if userID == "" || kid == "" {
		fmt.Println("help: gentoken <user_id> <kid>")
		return fmt.Errorf("gentoken issue")
	}

	db, err := database.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	store := user.NewStore(log, db)

	// The call to retrieve a user requires an Admin role by the caller.
	claims := auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Subject: userID,
		},
		Roles: []string{auth.RoleAdmin},
	}

	usr, err := store.QueryByID(ctx, claims, userID)
	if err != nil {
		return fmt.Errorf("retrieve user: %w", err)
	}

	// Construct a key store based on the key files stored in
	// the specified directory.
	keysFolder := "zarf/keys/"
	ks, err := keystore.NewFS(os.DirFS(keysFolder))
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	// Init the auth package.
	activeKID := "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"
	a, err := auth.New(activeKID, ks)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// Generating a token requires defining a set of claims. In this applications
	// case, we only care about defining the subject and the user in question and
	// the roles they have on the database. This token will expire in a year.
	//
	// iss (issuer): Issuer of the JWT
	// sub (subject): Subject of the JWT (the user)
	// aud (audience): Recipient for which the JWT is intended
	// exp (expiration time): Time after which the JWT expires
	// nbf (not before time): Time before which the JWT must not be accepted for processing
	// iat (issued at time): Time at which the JWT was issued; can be used to determine age of the JWT
	// jti (JWT ID): Unique identifier; can be used to prevent the JWT from being replayed (allows a token to be used only once)
	claims = auth.Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   usr.ID,
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().UTC().Unix(),
		},
		Roles: usr.Roles,
	}

	// This will generate a JWT with the claims embedded in them. The database
	// with need to be configured with the information found in the public key
	// file to validate these claims. Dgraph does not support key rotate at
	// this time.
	token, err := a.GenerateToken(claims)
	if err != nil {
		return fmt.Errorf("generating token: %w", err)
	}

	fmt.Printf("-----BEGIN TOKEN-----\n%s\n-----END TOKEN-----\n", token)
	return nil
}
