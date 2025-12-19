package application

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"

	"container-manager/internal/domain/entity"
	"container-manager/internal/domain/infrastructure"

	"github.com/google/uuid"
)

type ContainerService struct {
	runtime           infrastructure.ContainerRuntime
	containerUserRepo infrastructure.ContainerUserRepository
	jobRepo           infrastructure.JobRepository

	singleflightGroup singleflight.Group
	mutexMap          sync.Map
}

func NewContainerService(runtime infrastructure.ContainerRuntime, containerUserRepo infrastructure.ContainerUserRepository, jobRepo infrastructure.JobRepository) *ContainerService {
	return &ContainerService{
		runtime:           runtime,
		containerUserRepo: containerUserRepo,
		jobRepo:           jobRepo,
	}
}

func (s *ContainerService) CreateContainer(ctx context.Context, userID int64, options infrastructure.ContainerCreateOptions) (string, error) {
	payload, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	job := &entity.Job{
		ID:        uuid.New().String(),
		Type:      "container_creation",
		Status:    entity.JobStatusPending,
		Payload:   payload,
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.jobRepo.Create(context.Background(), job); err != nil {
		return "", err
	}

	go s.runCreateContainerJob(job, userID, options)

	return job.ID, nil
}

func (s *ContainerService) runCreateContainerJob(job *entity.Job, userID int64, options infrastructure.ContainerCreateOptions) {
	ctx := context.Background()

	job.Status = entity.JobStatusRunning
	job.UpdatedAt = time.Now()
	if err := s.jobRepo.Update(ctx, job); err != nil {
		log.Printf("failed to update job %s to running: %v", job.ID, err)
		return
	}

	containerID, err := s.runtime.Create(ctx, options)
	if err != nil {
		job.Status = entity.JobStatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		if updateErr := s.jobRepo.Update(ctx, job); updateErr != nil {
			log.Printf("failed to update job %s to failed: %v", job.ID, updateErr)
		}
		return
	}

	err = s.containerUserRepo.Create(ctx, containerID, userID)
	if err != nil {
		s.runtime.Remove(ctx, containerID)
		job.Status = entity.JobStatusFailed
		job.Error = err.Error()
		job.UpdatedAt = time.Now()
		if updateErr := s.jobRepo.Update(ctx, job); updateErr != nil {
			log.Printf("failed to update job %s to failed: %v", job.ID, updateErr)
		}
		return
	}

	result, err := json.Marshal(map[string]string{"container_id": containerID})
	if err != nil {
		job.Status = entity.JobStatusFailed
		job.Error = "failed to marshal result"
		job.UpdatedAt = time.Now()
	} else {
		job.Status = entity.JobStatusCompleted
		job.Result = result
		job.UpdatedAt = time.Now()
	}

	if updateErr := s.jobRepo.Update(ctx, job); updateErr != nil {
		log.Printf("failed to update job %s to completed/failed: %v", job.ID, updateErr)
	}
}

func (s *ContainerService) StartContainer(ctx context.Context, userID int64, id string) error {
	_, err, _ := s.singleflightGroup.Do("start:"+id, func() (any, error) {
		mutex := s.getMutex(id)
		if !mutex.TryLock() {
			return nil, errors.New("conflict operation")
		}
		defer mutex.Unlock()

		containerUserID, err := s.containerUserRepo.GetUserIDByContainerID(ctx, id)
		if err != nil {
			return nil, err
		}
		if containerUserID != userID {
			return nil, errors.New("permission denied")
		}
		err = s.runtime.Start(ctx, id)
		return nil, err
	})
	return err
}

func (s *ContainerService) StopContainer(ctx context.Context, userID int64, id string) error {
	_, err, _ := s.singleflightGroup.Do("stop:"+id, func() (any, error) {
		mutex := s.getMutex(id)
		if !mutex.TryLock() {
			return nil, errors.New("conflict operation")
		}
		defer mutex.Unlock()

		containerUserID, err := s.containerUserRepo.GetUserIDByContainerID(ctx, id)
		if err != nil {
			return nil, err
		}
		if containerUserID != userID {
			return nil, errors.New("permission denied")
		}
		err = s.runtime.Stop(ctx, id)
		return nil, err
	})
	return err
}

func (s *ContainerService) RemoveContainer(ctx context.Context, userID int64, id string) error {
	_, err, _ := s.singleflightGroup.Do("remove:"+id, func() (any, error) {
		mutex := s.getMutex(id)
		if !mutex.TryLock() {
			return nil, errors.New("conflict operation")
		}
		defer s.mutexMap.Delete(id)

		containerUserID, err := s.containerUserRepo.GetUserIDByContainerID(ctx, id)
		if err != nil {
			return nil, err
		}
		if containerUserID != userID {
			return nil, errors.New("permission denied")
		}
		err = s.runtime.Remove(ctx, id)
		if err != nil {
			return nil, err
		}
		err = s.containerUserRepo.Delete(ctx, id)
		return nil, err
	})
	return err
}

func (s *ContainerService) getMutex(id string) *sync.Mutex {
	m, _ := s.mutexMap.LoadOrStore(id, &sync.Mutex{})
	return m.(*sync.Mutex)
}
