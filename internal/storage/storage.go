package storage

import "github.com/vaib1343/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudent() ([]types.Student, error)
}
