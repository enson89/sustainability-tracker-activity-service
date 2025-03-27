package repository

import (
	"errors"

	"github.com/enson89/sustainability-tracker-activity-service/internal/model"
	"github.com/jmoiron/sqlx"
)

type ActivityRepository interface {
	CreateActivity(activity *model.Activity) error
	GetActivityByID(id int64) (*model.Activity, error)
	UpdateActivity(activity *model.Activity) error
	DeleteActivity(id int64) error
}

type activityRepository struct {
	db *sqlx.DB
}

// NewActivityRepository returns a new ActivityRepository.
func NewActivityRepository(db *sqlx.DB) ActivityRepository {
	return &activityRepository{db: db}
}

func (r *activityRepository) CreateActivity(activity *model.Activity) error {
	query := `INSERT INTO activities (user_id, type, amount, description, created_at)
	          VALUES ($1, $2, $3, $4, now()) RETURNING id, created_at`
	return r.db.QueryRow(query, activity.UserID, activity.Type, activity.Amount, activity.Description).
		Scan(&activity.ID, &activity.CreatedAt)
}

func (r *activityRepository) GetActivityByID(id int64) (*model.Activity, error) {
	activity := &model.Activity{}
	query := `SELECT id, user_id, type, amount, description, created_at FROM activities WHERE id=$1`
	err := r.db.Get(activity, query, id)
	if err != nil {
		return nil, err
	}
	return activity, nil
}

func (r *activityRepository) UpdateActivity(activity *model.Activity) error {
	query := `UPDATE activities SET type=$1, amount=$2, description=$3 WHERE id=$4`
	res, err := r.db.Exec(query, activity.Type, activity.Amount, activity.Description, activity.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("activity not found")
	}
	return nil
}

func (r *activityRepository) DeleteActivity(id int64) error {
	query := `DELETE FROM activities WHERE id=$1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("activity not found")
	}
	return nil
}

func NewPostgresDB(databaseURL string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	return db, nil
}
