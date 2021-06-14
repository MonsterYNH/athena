package orm

import (
	"time"
)

type Model struct {
	ID        string     `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at"`
}

type User struct {
	Model

	Account  string `gorm:"column:account;unique"`
	Password string `gorm:"column:password"`
	Name     string `gorm:"column:name"`
	Sex      string `gorm:"column:sex"`
	Age      int    `gorm:"column:age"`
	Address  string `gorm:"column:address"`
	Email    string `gorm:"column:email"`
	Phone    string `gorm:"column:phone"`

	Avatar string `gorm:"column:avatar"`
	Status string `grom:"column:status"`
}

func (u User) TableName() string {
	return "users"
}
