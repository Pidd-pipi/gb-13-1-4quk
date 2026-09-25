package repository

import (
	"gorm.io/gorm"

	"github.com/campusbooks/campusbooks/internal/model"
)

// UserRepository handles user persistence.
type UserRepository struct{ db *gorm.DB }

// NewUserRepository creates a UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository { return &UserRepository{db: db} }

// Create inserts a user.
func (r *UserRepository) Create(u *model.User) error { return translate(r.db.Create(u).Error) }

// FindByID locates a user by id.
func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var u model.User
	if err := translate(r.db.First(&u, id).Error); err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByEmail locates a user by email.
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var u model.User
	if err := translate(r.db.Where("email = ?", email).First(&u).Error); err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByStudentNo locates a user by student number.
func (r *UserRepository) FindByStudentNo(studentNo string) (*model.User, error) {
	var u model.User
	if err := translate(r.db.Where("student_no = ?", studentNo).First(&u).Error); err != nil {
		return nil, err
	}
	return &u, nil
}

// FindByAccount locates a user by email or student number (login helper).
func (r *UserRepository) FindByAccount(account string) (*model.User, error) {
	var u model.User
	err := translate(r.db.Where("email = ? OR student_no = ?", account, account).First(&u).Error)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// Update persists a user.
func (r *UserRepository) Update(u *model.User) error { return translate(r.db.Save(u).Error) }

// List returns users with pagination and optional role filter.
func (r *UserRepository) List(role string, offset, limit int) ([]model.User, int64, error) {
	var items []model.User
	var total int64
	q := r.db.Model(&model.User{})
	if role != "" {
		q = q.Where("role = ?", role)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("id ASC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// CountBooksBySeller counts books sold/on-sale for a seller (stats helper).
func (r *UserRepository) CountBooksBySeller(sellerID uint) (total, onSale, sold int64, err error) {
	if err = r.db.Model(&model.Book{}).Where("seller_id = ?", sellerID).Count(&total).Error; err != nil {
		return
	}
	if err = r.db.Model(&model.Book{}).Where("seller_id = ? AND status = ?", sellerID, "on_sale").Count(&onSale).Error; err != nil {
		return
	}
	if err = r.db.Model(&model.Book{}).Where("seller_id = ? AND status = ?", sellerID, "sold").Count(&sold).Error; err != nil {
		return
	}
	return
}
