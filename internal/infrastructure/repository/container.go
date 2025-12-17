package repository

import (
	"container-manager/internal/domain/repository"
	"context"
	"database/sql"
)

var _ repository.ContainerUserRepository = (*ContainerUserDatabase)(nil)

type ContainerUserDatabase struct {
	db *sql.DB
}

func NewContainerDatabase(db *sql.DB) repository.ContainerUserRepository {
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
