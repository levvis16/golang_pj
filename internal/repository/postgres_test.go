package repository

import (
	"database/sql"
	"testing"
	"time"

	"subscription-service/internal/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5433 user=postgres password=postgres dbname=subscriptions sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec(`TRUNCATE TABLE subscriptions`)
	if err != nil {
		t.Fatal(err)
	}

	return db
}

func TestCreateSubscription(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSubscriptionRepository(db)

	now := time.Now()
	sub := &models.Subscription{
		ID:          uuid.New(),
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err := repo.Create(sub)
	if err != nil {
		t.Errorf("Failed to create subscription: %v", err)
	}
}

func TestGetSubscriptionByID(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSubscriptionRepository(db)

	id := uuid.New()
	now := time.Now()
	sub := &models.Subscription{
		ID:          id,
		ServiceName: "Spotify",
		Price:       299,
		UserID:      uuid.New(),
		StartDate:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo.Create(sub)

	result, err := repo.GetByID(id)
	if err != nil {
		t.Errorf("Failed to get subscription: %v", err)
	}

	if result == nil {
		t.Error("Subscription not found")
	}
}

func TestUpdateSubscription(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSubscriptionRepository(db)

	id := uuid.New()
	now := time.Now()
	sub := &models.Subscription{
		ID:          id,
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo.Create(sub)

	sub.Price = 1099
	sub.UpdatedAt = time.Now()

	err := repo.Update(sub)
	if err != nil {
		t.Errorf("Failed to update subscription: %v", err)
	}

	updated, _ := repo.GetByID(id)
	if updated.Price != 1099 {
		t.Errorf("Expected price 1099, got %d", updated.Price)
	}
}

func TestDeleteSubscription(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSubscriptionRepository(db)

	id := uuid.New()
	now := time.Now()
	sub := &models.Subscription{
		ID:          id,
		ServiceName: "Netflix",
		Price:       999,
		UserID:      uuid.New(),
		StartDate:   now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	repo.Create(sub)

	err := repo.Delete(id)
	if err != nil {
		t.Errorf("Failed to delete subscription: %v", err)
	}

	result, _ := repo.GetByID(id)
	if result != nil {
		t.Error("Subscription still exists after deletion")
	}
}

func TestGetTotalCostWithFilters(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewSubscriptionRepository(db)

	userID := uuid.New()
	now := time.Now()

	subs := []models.Subscription{
		{
			ID:          uuid.New(),
			ServiceName: "Netflix",
			Price:       999,
			UserID:      userID,
			StartDate:   now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			ServiceName: "Spotify",
			Price:       299,
			UserID:      userID,
			StartDate:   now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		{
			ID:          uuid.New(),
			ServiceName: "Netflix",
			Price:       999,
			UserID:      uuid.New(),
			StartDate:   now,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
	}

	for _, sub := range subs {
		repo.Create(&sub)
	}

	filter := &models.SubscriptionFilter{
		UserID: userID.String(),
	}

	total, err := repo.GetTotalCost(filter)
	if err != nil {
		t.Errorf("Failed to get total cost: %v", err)
	}

	if total != 1298 {
		t.Errorf("Expected total 1298, got %d", total)
	}
}
