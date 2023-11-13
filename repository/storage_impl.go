package repository

import "github.com/MorZLE/jobs_bot/config"

func NewRepository(cnf *config.Config) (Storage, error) {
	return &repository{}, nil
}

type repository struct {
}

func (r *repository) Set() {

}

func (r *repository) Get() {

}
func (r *repository) Delete() {

}
