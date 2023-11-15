package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
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
	category := []string{"Разработчик", "Инфо без-ть", "Системный ад-р", "Банковское дело", "Страховой агент", "Мечтатель"}
	mCategory := make(map[string][]model.Student)
	for _, c := range category {
		var students []model.Student
		err := db.Where("category = ?", c).Find(&students).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				continue
			}
			return nil, fmt.Errorf("error get user: %w", err)
		}
		mCategory[c] = students
	}

	return &repository{db: db, m: mCategory}, nil
}

type repository struct {
	db *gorm.DB
	m  map[string][]model.Student
}

func (r *repository) Set(student model.Student) error {
	if err := r.db.Create(&student).Error; err != nil {
		return fmt.Errorf("error create user: %w", err)
	}
	r.m[student.Category] = append(r.m[student.Category], student)
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

func (r *repository) GetOneResume(category string, count int) (model.Student, error) {
	if _, ok := r.m[category]; !ok {
		return model.Student{}, constants.ErrNotCategory
	}
	if len(r.m[category]) == 0 {
		return model.Student{}, constants.ErrNotResume
	}
	if len(r.m[category]) <= count {
		return model.Student{}, constants.ErrNotFound
	}

	return r.m[category][count], nil
}
