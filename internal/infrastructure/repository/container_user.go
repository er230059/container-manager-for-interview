package repository

import (
	"container-manager/internal/errors"
	"context"
	"database/sql"
)

type ContainerUserRepository struct {
	db *sql.DB
}

func NewContainerUserRepository(db *sql.DB) *ContainerUserRepository {
	return &ContainerUserRepository{db: db}
}

func (r *ContainerUserRepository) Create(ctx context.Context, containerID string, userID int64) error {
	query := "INSERT INTO container_user (container_id, user_id) VALUES ($1, $2)"
	_, err := r.db.ExecContext(ctx, query, containerID, userID)
	return err
}

func (r *ContainerUserRepository) Delete(ctx context.Context, containerID string) error {
	query := "DELETE FROM container_user WHERE container_id = $1"
	_, err := r.db.ExecContext(ctx, query, containerID)
	return err
}

func (r *ContainerUserRepository) GetUserIDByContainerID(ctx context.Context, containerID string) (int64, error) {
	query := "SELECT user_id FROM container_user WHERE container_id = $1"
	var userID int64
	err := r.db.QueryRowContext(ctx, query, containerID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.ContainerNotFound
		}
		return 0, err
	}
	return userID, nil
}

func (r *ContainerUserRepository) GetContainerIDsByUserID(ctx context.Context, userID int64) ([]string, error) {
	query := "SELECT container_id FROM container_user WHERE user_id = $1"
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var containerIDs []string
	for rows.Next() {
		var containerID string
		if err := rows.Scan(&containerID); err != nil {
			return nil, err
		}
		containerIDs = append(containerIDs, containerID)
	}
	return containerIDs, rows.Err()
}
