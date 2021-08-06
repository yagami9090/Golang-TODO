package repository

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"todo/entities"
)

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) GormRepository {
	return GormRepository{db: db}
}

func (repo GormRepository) NewTask(task *entities.Task) error {
	return errors.WithMessage(repo.db.Create(task).Error, "gorm insert")
}

func (repo GormRepository) TaskDone(id uint) error {
	return errors.WithMessage(repo.db.Model(&entities.Task{}).Where("id = ?", id).Update("done", true).Error, "gorm update task done")
}

func (repo GormRepository) List() ([]*entities.Task, error) {
	var tasks []*entities.Task
	err := repo.db.Find(&tasks).Error
	return tasks, errors.WithMessage(err, "gorm list task")
}
