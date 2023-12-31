package repository

import (
	"errors"
	"fmt"
	"github.com/MorZLE/jobs_bot/config"
	"github.com/MorZLE/jobs_bot/constants"
	"github.com/MorZLE/jobs_bot/logger"
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

var Blacklist []int64
var Admins []int64

func NewRepository(cnf *config.Config) (Storage, error) {
	db, err := gorm.Open(postgres.Open(cnf.DB), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.Debug().AutoMigrate(&model.Student{}, &model.BanUser{}, &model.AdminInvait{}, &model.Admin{})
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
	var bans []model.BanUser
	err = db.Find(&bans).Error
	if err != nil {
		return nil, fmt.Errorf("error get banuser: %w", err)
	}
	for _, ban := range bans {
		Blacklist = append(Blacklist, ban.Tgid)
	}
	storage := &repository{db: db, m: mCategory, mIdxResume: make(map[int64]int)}
	adm, err := storage.GetAdmins()
	if err != nil {
		logger.Error("failed to get admins", err)
	}
	Admins = append(Admins, cnf.Admin)
	for _, admin := range adm {
		Admins = append(Admins, admin.Tgid)
	}
	return storage, nil
}

type repository struct {
	db         *gorm.DB
	m          map[string][]model.Student //Массив студентов по категориям
	mIdxResume map[int64]int              //Индекс студента в массиве
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

func (r *repository) GetOneResume(category string, direction string, count int, wantStatus string) (model.Student, error) {
	valCat, ok := r.m[category]
	if !ok {
		return model.Student{}, constants.ErrNotCategory
	}
	if len(valCat) == 0 {
		return model.Student{}, constants.ErrNotResume
	}
	if len(valCat) < count {
		return model.Student{}, constants.ErrNotFound
	}
	switch direction {
	case constants.Offer:
		return valCat[count], nil
	default:
		res := valCat[count]
		if res.Status == constants.StatusDeleted && len(valCat) == count+1 {
			return model.Student{}, constants.ErrNotFound
		}
		if res.Status != wantStatus {
			if len(valCat) == count+1 {
				return model.Student{}, constants.ErrNotFound
			}
			return model.Student{}, constants.ErrDeleteResume
		}
	}
	if len(valCat) == count+1 {
		return valCat[count], constants.ErrLastResume
	}
	return valCat[count], nil
}

func (r *repository) BanUser(idx int, category string) error {
	users := r.m[category]
	st := users[idx]
	st.Status = constants.StatusPublished
	users[idx] = st
	r.m[category] = users

	err := r.db.Model(&model.Student{}).Where("tgid = ?", st.Tgid).Update("status", constants.StatusBanned).Error
	if err != nil {
		return fmt.Errorf("error banUser user: %w", err)
	}
	banUser := model.BanUser{
		Tgid:     st.Tgid,
		Username: st.Username,
	}
	if err := r.db.Create(&banUser).Error; err != nil {
		return fmt.Errorf("error banUser user: %w", err)
	}
	Blacklist = append(Blacklist, banUser.Tgid)
	return nil
}
func (r *repository) PublishUser(idx int, category string) error {
	users := r.m[category]
	st := users[idx]
	st.Status = constants.StatusPublished
	users[idx] = st
	r.m[category] = users

	err := r.db.Model(&model.Student{}).Where("tgid = ?", st.Tgid).Update("status", constants.StatusPublished).Error
	if err != nil {
		return fmt.Errorf("error PublishUser user: %w", err)
	}
	return nil
}
func (r *repository) DeclineUser(idx int, category string) error {
	users := r.m[category]
	st := users[idx]
	st.Status = constants.StatusPublished
	users[idx] = st
	r.m[category] = users
	err := r.db.Model(&model.Student{}).Where("tgid = ?", st.Tgid).Update("status", constants.StatusRejected).Error

	if err != nil {
		return fmt.Errorf("error DeclineUser user: %w", err)
	}
	return nil
}

func (r *repository) Statistics() (map[string][]model.Student, error) {
	return r.m, nil
}

func (r *repository) UnbanUsername(username string) error {
	var banUser model.BanUser
	err := r.db.Where("username = ?", username).Find(&banUser).Error
	if err != nil {
		return err
	}
	tgid := banUser.Tgid
	if tgid == 0 {
		return constants.ErrNotFound
	}
	err = r.db.Where("username = ?", username).Delete(&model.BanUser{}).Error
	if err != nil {
		return err
	}
	deleteBanUser(tgid)
	return nil
}
func (r *repository) UnbanTgID(tgid int64) error {
	err := r.db.Where("tgid = ?", tgid).Delete(&model.BanUser{}).Error
	if err != nil {
		return err
	}
	deleteBanUser(tgid)
	return nil
}

func (r *repository) ViewBanList() ([]model.BanUser, error) {
	var banUsers []model.BanUser
	err := r.db.Find(&banUsers).Error
	if err != nil {
		return nil, err
	}
	return banUsers, nil
}

func deleteBanUser(tgid int64) {
	for r := range Blacklist {
		if Blacklist[r] == tgid {
			Blacklist = append(Blacklist[:r], Blacklist[r+1:]...)
		}
	}
	return
}

func (r *repository) NewAdminURL(username, url string) error {
	adminURL := model.AdminInvait{
		Username: username,
		Url:      url,
	}
	return r.db.Create(&adminURL).Error
}

func (r *repository) CheckUrlAdmin(username, url string) error {
	var adminURL model.AdminInvait
	err := r.db.Where("username = ?", username).Where("url = ?", url).Find(&adminURL).Error
	if err != nil {
		return err
	}
	return nil
}
func (r *repository) DeleteUrlInvaite(username, url string) error {
	err := r.db.Where("username = ?", username).Where("url = ?", url).Delete(&model.AdminInvait{}).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) CreateAdmin(username string, id int64) error {
	admin := model.Admin{
		Username: username,
		Tgid:     id,
		Lvl:      constants.Junior,
	}
	err := r.db.Create(&admin).Error
	if err != nil {
		return err
	}
	Admins = append(Admins, admin.Tgid)
	return nil
}

func (r *repository) GetAdmins() ([]model.Admin, error) {
	var admins []model.Admin
	err := r.db.Find(&admins).Error
	if err != nil {
		return nil, err
	}
	return admins, nil
}

func (r *repository) DeleteAdmin(username string) error {
	var admin model.Admin
	err := r.db.Where("username = ?", username).Find(&admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return constants.ErrNotFound
		}
		return err
	}
	id := admin.Tgid
	err = r.db.Where("username = ?", username).Delete(&model.Admin{}).Error
	if err != nil {

		return err
	}
	for r := range Admins {
		if Admins[r] == id {
			if r == len(Admins)-1 {
				Admins = Admins[:r]
				return nil
			}
			Admins = append(Admins[:r], Admins[r+1:]...)
		}
	}
	return nil
}
