package repository

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
	"github.com/MorZLE/jobs_bot/model"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var Category = []string{
	"ИБ",
	"ДОУ",
	"Финансы",
	"Реклама",
	"Логистика",
	"Разработчик",
	"Оператор ИС",
	"Страховое дело",
	"Землеустройство",
	"Банковское дело",
	"Оператор верстки",
	"Издательское дело",
	"Прикладная геодезия",
	"Графический дизайнер",
	"Управление качеством",
	"Экономика и бух учет",
	"Системное администрирование",
	"Другое",
}

func NewRepository(cnf *config.Config) (Storage, error) {
	//ctx := context.TODO()
	//config := dockerdb.EmptyConfig().DBName("dbjobsbot").DBUser("fl0user").
	//	DBPassword("kYLDaq9SdN8f").StandardDBPort("5432").
	//	Vendor(dockerdb.Postgres15).SQL().PullImage()
	//
	//vdb, err := dockerdb.New(ctx, config.Build())
	//if err != nil {
	//	return nil, fmt.Errorf("failed to connect docker: %w", err)
	//}
	db, err := gorm.Open(postgres.Open(cnf.DB), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.Debug().AutoMigrate(&model.Student{}, &model.Employee{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	mCategory := make(map[string][]model.Student)
	for _, c := range Category {
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

	return &repository{db: db, m: mCategory, mIdxResume: make(map[int64]int)}, nil
}

type repository struct {
	db         *gorm.DB
	m          map[string][]model.Student
	mIdxResume map[int64]int
}

func (r *repository) Set(student model.Student) error {
	if err := r.db.Create(&student).Error; err != nil {
		return fmt.Errorf("error create user: %w", err)
	}
	r.mIdxResume[student.Tgid] = len(r.m[student.Category])
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

func (r *repository) Delete(id int64, category string) error {
	var student model.Student
	if err := r.db.Where("tgid = ?", id).Delete(&student).Error; err != nil {
		return fmt.Errorf("error delete user: %w", err)
	}
	users := r.m[category]
	idx := r.mIdxResume[id]
	st := users[idx]
	st.Status = constants.StatusDeleted
	users[idx] = st
	r.m[category] = users
	return nil
}

func (r *repository) Close() {
	sqlDB, _ := r.db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}

func (r *repository) GetOneResume(category string, direction string, count int) (model.Student, error) {
	if _, ok := r.m[category]; !ok {
		return model.Student{}, constants.ErrNotCategory
	}
	if len(r.m[category]) == 0 {
		return model.Student{}, constants.ErrNotResume
	}
	if len(r.m[category]) < count {
		return model.Student{}, constants.ErrNotFound
	}
	switch direction {
	case constants.Offer:
		return r.m[category][count], nil
	default:
		res := r.m[category][count]
		if res.Status == constants.StatusDeleted {
			return model.Student{}, constants.ErrDeleteResume
		}
	}
	if len(r.m[category]) == count+1 {
		return r.m[category][count], constants.ErrLastResume
	}
	return r.m[category][count], nil
}
