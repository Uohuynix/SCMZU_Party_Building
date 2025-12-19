package entity

import (
	"time"
)

type Training struct {
	ID            uint       `gorm:"primaryKey" json:"id"`
	Period        string     `gorm:"size:50" json:"period"`
	StudentID     *uint      `gorm:"default:NULL" json:"student_id"`
	Unit          string     `gorm:"size:100" json:"unit"`
	StartDate     *time.Time `gorm:"type:date" json:"start_date"`
	EndDate       *time.Time `gorm:"type:date" json:"end_date"`
	Score         string     `gorm:"type:ENUM('excellent','good','pass','fail')" json:"score"`
	CertificateNo string     `gorm:"size:50" json:"certificate_no"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Student       Student    `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}
