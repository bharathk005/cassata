package main

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"os"
	"time"

	_ "github.com/lib/pq"
	v1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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
	createJWTSecret(clientset)

	// 2. Create database tables if they don't exist
	createDatabaseTables()

	// 3. Create cluster roles, role bindings, and service account
	createKubernetesResources(clientset)

	log.Println("Initialization completed successfully")
}

func createJWTSecret(clientset *kubernetes.Clientset) {
	secretName := os.Getenv("JWT_SECRET_NAME")
	namespace := os.Getenv("NAMESPACE")

	_, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if err == nil {
		log.Printf("JWT secret %s already exists", secretName)
		return
	}
	if !errors.IsNotFound(err) {
		log.Fatalf("Error checking JWT secret: %s", err)
	}

	jwtSecret := os.Getenv("JWT_SECRET_VALUE")
	if jwtSecret == "" {
		jwtSecret = generateRandomString(32)
	}

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: secretName,
		},
		StringData: map[string]string{
			"JWT_SECRET": jwtSecret,
		},
		Type: v1.SecretTypeOpaque,
	}

	_, err = clientset.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		log.Fatalf("Error creating JWT secret: %s", err)
	}
	log.Printf("Created JWT secret %s", secretName)
}

func createDatabaseTables() {
	dsn := os.Getenv("DB_DSN")

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error connecting to database: %s", err)
	}
	defer db.Close()

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

func createKubernetesResources(clientset *kubernetes.Clientset) {
	namespace := os.Getenv("NAMESPACE")
	serviceAccountName := os.Getenv("SERVICE_ACCOUNT_NAME")
	clusterRoleName := os.Getenv("CLUSTER_ROLE_NAME")
	clusterRoleBindingName := os.Getenv("CLUSTER_ROLE_BINDING_NAME")

	// Create ServiceAccount
	sa := &v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceAccountName,
		},
	}
	_, err := clientset.CoreV1().ServiceAccounts(namespace).Create(context.TODO(), sa, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Fatalf("Error creating ServiceAccount: %s", err)
	}
	log.Printf("ServiceAccount %s created or already exists", serviceAccountName)

	// Create ClusterRole
	clusterRole := &rbacv1.ClusterRole{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleName,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{"rbac.authorization.k8s.io"},
				Resources: []string{"clusterroles"},
				Verbs:     []string{"update", "get", "list"},
			},
		},
	}
	_, err = clientset.RbacV1().ClusterRoles().Create(context.TODO(), clusterRole, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Fatalf("Error creating ClusterRole: %s", err)
	}
	log.Printf("ClusterRole %s created or already exists", clusterRoleName)

	// Create ClusterRoleBinding
	clusterRoleBinding := &rbacv1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: clusterRoleBindingName,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      "ServiceAccount",
				Name:      serviceAccountName,
				Namespace: namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "ClusterRole",
			Name:     clusterRoleName,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
	_, err = clientset.RbacV1().ClusterRoleBindings().Create(context.TODO(), clusterRoleBinding, metav1.CreateOptions{})
	if err != nil && !errors.IsAlreadyExists(err) {
		log.Fatalf("Error creating ClusterRoleBinding: %s", err)
	}
	log.Printf("ClusterRoleBinding %s created or already exists", clusterRoleBindingName)
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
