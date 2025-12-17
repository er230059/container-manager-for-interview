package repository

import (
	"container-manager/internal/domain/repository"
	"context"
	"database/sql"
	"errors"
)

var _ repository.ContainerUserRepository = (*ContainerUserDatabase)(nil)

type ContainerUserDatabase struct {
	db *sql.DB
}

func NewContainerUserDatabase(db *sql.DB) repository.ContainerUserRepository {
	return &ContainerUserDatabase{db: db}
}

func (d *ContainerUserDatabase) Create(ctx context.Context, containerID string, userID int64) error {
	query := "INSERT INTO container_user (container_id, user_id) VALUES ($1, $2)"
	_, err := d.db.ExecContext(ctx, query, containerID, userID)
	return err
}

func (d *ContainerUserDatabase) Delete(ctx context.Context, containerID string) error {
	query := "DELETE FROM container_user WHERE container_id = $1"
	_, err := d.db.ExecContext(ctx, query, containerID)
	return err
}

func (d *ContainerUserDatabase) GetUserIDByContainerID(ctx context.Context, containerID string) (int64, error) {
	query := "SELECT user_id FROM container_user WHERE container_id = $1"
	var userID int64
	err := d.db.QueryRowContext(ctx, query, containerID).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New("container not found")
		}
		return 0, err
	}
	return userID, nil
}
