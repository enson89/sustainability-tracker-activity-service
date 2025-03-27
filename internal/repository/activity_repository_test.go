package repository_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/enson89/sustainability-tracker-activity-service/internal/repository"
	"github.com/jmoiron/sqlx"
)

func TestCreateActivity(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error opening stub database connection: %v", err)
	}
	defer db.Close()
	sqlxDB := sqlx.NewDb(db, "postgres")
	repo := repository.NewActivityRepository(sqlxDB)

	activity := &model.Activity{
		UserID:      1,
		Type:        "Test",
		Amount:      10.5,
		Description: "Testing activity",
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO activities (user_id, type, amount, description, created_at) VALUES ($1, $2, $3, $4, now()) RETURNING id, created_at")).
		WithArgs(activity.UserID, activity.Type, activity.Amount, activity.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, time.Now()))

	err = repo.CreateActivity(activity)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if activity.ID != 1 {
		t.Errorf("expected activity ID to be 1, got %v", activity.ID)
	}
}
