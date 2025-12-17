package repository

import (
	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/repository"
	"context"
	"database/sql"
)

var _ repository.ContainerRepository = (*ContainerDatabase)(nil)

type ContainerDatabase struct {
	db *sql.DB
}

func NewContainerDatabase(db *sql.DB) repository.ContainerRepository {
	return &ContainerDatabase{db: db}
}

func (d *ContainerDatabase) Create(ctx context.Context, container *entity.Container) error {
	query := "INSERT INTO containers (id, image, user_id) VALUES ($1, $2, $3)"
	_, err := d.db.ExecContext(ctx, query, container.ID, container.Image, container.UserID)
	return err
}

func (d *ContainerDatabase) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM containers WHERE id = $1"
	_, err := d.db.ExecContext(ctx, query, id)
	return err
}
