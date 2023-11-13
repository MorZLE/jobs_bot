package repository

type Storage interface {
	Set()
	Get()
	Delete()
}
