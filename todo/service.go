package todo

import "todo/entities"

type Repository interface {
	NewTask(*entities.Task) error
	TaskDone(uint) error
	List() ([]*entities.Task, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return Service{repo: repo}
}

func (srv Service) Add(task entities.Task) error {
	return srv.repo.NewTask(&task)
}

func (srv Service) Done(id uint) error {
	return srv.repo.TaskDone(uint(id))
}

func (srv Service) List() ([]*entities.Task, error) {
	return srv.repo.List()
}
