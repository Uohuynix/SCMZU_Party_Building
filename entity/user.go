package entity

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Password string `gorm:"size:255;not null" json:"-"`
	Role     string `gorm:"type:ENUM('admin','teacher','student');default:'student'" json:"role"`
	Name     string `gorm:"size:50" json:"name"`
	// 对于教师：所属党支部/党小组；对于学生/管理员可为空
	Branch    string    `gorm:"size:100" json:"branch"`
	GroupName string    `gorm:"size:100" json:"group_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
