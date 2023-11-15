package repository

import (
	"context"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/model"
	"github.com/egorgasay/dockerdb/v3"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

func NewRepository(cnf *config.Config) (Storage, error) {
	ctx := context.TODO()
	config := dockerdb.EmptyConfig().DBName("postgres").DBUser("postgres").
		DBPassword("postgres").StandardDBPort("5432").
		Vendor(dockerdb.Postgres15).SQL().PullImage()

	vdb, err := dockerdb.New(ctx, config.Build())
	if err != nil {
		return nil, fmt.Errorf("failed to connect docker: %w", err)
	}
	db, err := gorm.Open(postgres.Open(vdb.GetSQLConnStr()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.Debug().AutoMigrate(&model.Student{}, &model.Employee{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *gorm.DB
}

func (r *repository) Set(student model.Student) error {
	if err := r.db.Create(&student).Error; err != nil {
		return fmt.Errorf("error create user: %w", err)
	}
	return nil
}

func (r *repository) Get(id int64) (model.Student, error) {
	var student model.Student
	if err := r.db.Where("tgid = ?", id).First(&student).Error; err != nil {
		return model.Student{}, fmt.Errorf("error get user: %w", err)
	}
	return student, nil
}

func (r *repository) Delete(id int64) error {
	var student model.Student
	if err := r.db.Where("tgid = ?", id).Delete(&student).Error; err != nil {
		return fmt.Errorf("error delete user: %w", err)
	}
	return nil
}

func (r *repository) Close() {
	sqlDB, _ := r.db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}
