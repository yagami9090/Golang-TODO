package repository

import (
	"sync"
	"todo/entities"
)

var index uint
var tasks map[uint]*entities.Task = make(map[uint]*entities.Task)
var mutex = sync.Mutex{}

type MemoryRepository struct{}

func (repo MemoryRepository) NewTask(task *entities.Task) error {
	mutex.Lock()
	defer func() {
		index++
		mutex.Unlock()
	}()

	tasks[index] = &entities.Task{
		Title: task.Title,
		Done:  false,
	}
	return nil
}

func (repo MemoryRepository) TaskDone(id uint) error {
	tasks[id].Done = true
	return nil
}

func (repo MemoryRepository) List() (map[uint]*entities.Task, error) {
	return tasks, nil
}
