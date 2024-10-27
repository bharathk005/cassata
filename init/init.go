package init

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"
	"fmt"

	_ "github.com/lib/pq"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func main() {
	// Create Kubernetes client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating in-cluster config: %s", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating Kubernetes client: %s", err)
	}

	// 1. Create JWT secret if it doesn't exist
	createSecret(clientset, os.Getenv("JWT_SECRET_NAME"), os.Getenv("NAMESPACE"), map[string]string{"JWT_SECRET_KEY": generateRandomString(32)})

	// 2. Create database tables if they don't exist
	createDatabaseTables(clientset)

	log.Println("Initialization completed successfully")
}


func createSecret(clientset *kubernetes.Clientset, secretName string, namespace string, data map[string]string) {
	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err == nil {
		log.Printf("Secret %s already exists", secretName)
		return
	}
	if !errors.IsNotFound(err) {
		log.Fatalf("Error checking secret: %s", err)
	}
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		StringData: data,
		Type: v1.SecretTypeOpaque,
	}
	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating secret: %s", err)
	}
	log.Printf("Created secret %s", secretName)
}

func createDatabaseTables(clientset *kubernetes.Clientset) {
	dsn := os.Getenv("DATABASE_DSN")
	external := os.Getenv("DATABASE_EXTERNAL")
	namespace := os.Getenv("NAMESPACE")
	var db *sql.DB
	var err error
	if external == "true" {
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Error connecting to external database: %s", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error pinging external database: %s", err)
		}
	} else {
		dbPassword := generateRandomString(16)
		dsn = fmt.Sprintf("postgres://postgres:%s@cassata-postgres.%s.svc.cluster.local:5432/cassata", dbPassword, namespace)
		db, err = sql.Open("postgres", dsn)
		if err != nil {
			log.Fatalf("Error connecting to local database: %s", err)
		}
		err = db.Ping()
		if err != nil {
			log.Fatalf("Error pinging local database: %s", err)
		}
	}
	createSecret(clientset, os.Getenv("DATABASE_SECRET_NAME"), namespace, map[string]string{"DATABASE_DSN": dsn})
	_, err = db.Exec(`
		-- Migration script to create tables and add unique constraints
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE IF NOT EXISTS workspaces (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE IF NOT EXISTS permissions (
			id SERIAL PRIMARY KEY,
			resource VARCHAR(255) NOT NULL,
			action VARCHAR(50) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);

		CREATE TABLE IF NOT EXISTS user_workspace (
			user_id INTEGER REFERENCES users(id),
			workspace_id INTEGER REFERENCES workspaces(id),
			PRIMARY KEY (user_id, workspace_id)
		);

		CREATE TABLE IF NOT EXISTS workspace_permissions (
			workspace_id INTEGER REFERENCES workspaces(id),
			permission_id INTEGER REFERENCES permissions(id),
			PRIMARY KEY (workspace_id, permission_id)
		);

		CREATE TABLE IF NOT EXISTS resource_maps (
			id SERIAL PRIMARY KEY,
			provider VARCHAR(255) NOT NULL,
			resource_group VARCHAR(255) NOT NULL,
			resource_type VARCHAR(255) NOT NULL,
			k8s_api_group VARCHAR(255) NOT NULL,
			k8s_api_version VARCHAR(255) NOT NULL,
			k8s_resource VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP WITH TIME ZONE
		);

		-- Add unique constraints to prevent duplicate entries
		ALTER TABLE users ADD CONSTRAINT unique_user_name UNIQUE (name);
		ALTER TABLE workspaces ADD CONSTRAINT unique_workspace_name UNIQUE (name);
		ALTER TABLE permissions ADD CONSTRAINT unique_permission UNIQUE (resource, action);
		ALTER TABLE user_workspace ADD CONSTRAINT unique_user_workspace UNIQUE (user_id, workspace_id);
		ALTER TABLE workspace_permissions ADD CONSTRAINT unique_workspace_permission UNIQUE (workspace_id, permission_id);
		ALTER TABLE resource_maps ADD CONSTRAINT unique_resource_map UNIQUE (provider, resource_group, resource_type);
	`)
	if err != nil {
		log.Fatalf("Error running database migration script: %s", err)
	}

	log.Println("Database migration script executed successfully")
}

func generateRandomString(length int) string {
	source := rand.NewSource(time.Now().UnixNano())
	r := rand.New(source)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[r.Intn(len(chars))]
	}
	return string(b)
}
