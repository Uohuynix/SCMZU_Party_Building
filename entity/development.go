package entity

import (
	"time"
)

type Development struct {
	ID                          uint      `gorm:"primaryKey" json:"id"`
	StudentID                   uint      `json:"student_id"`
	JoinLeagueDate              *string   `gorm:"type:varchar(50)" json:"join_league_date"`                // 入团时间（文本格式）
	ApplyDate                   *string   `gorm:"type:varchar(50)" json:"apply_date"`                      // 递交入党申请书时间（文本格式）
	PartyTalkDate               *string   `gorm:"type:varchar(50)" json:"party_talk_date"`                 // 党组织派人谈话时间（文本格式）
	LeagueRecommendDate         *string   `gorm:"type:varchar(50)" json:"league_recommend_date"`           // 团组织推优时间（文本格式）
	ActivistDate                *string   `gorm:"type:varchar(50)" json:"activist_date"`                   // 确立为入党积极分子时间（文本格式）
	CandidateDate               *string   `gorm:"type:varchar(50)" json:"candidate_date"`                  // 确定发展对象时间（文本格式）
	SuperiorReportDate          *string   `gorm:"type:varchar(50)" json:"superior_report_date"`            // 提交上级党组织备案时间（文本格式）
	CommitteeReviewDate         *string   `gorm:"type:varchar(50)" json:"committee_review_date"`           // 支委会审查时间（文本格式）
	PartyCommitteePreReviewDate *string   `gorm:"type:varchar(50)" json:"party_committee_pre_review_date"` // 党委预审通过时间（文本格式）
	BranchMeetingDate           *string   `gorm:"type:varchar(50)" json:"branch_meeting_date"`             // 支部大会讨论时间（发展）（文本格式）
	ConversionApplyDate         *string   `gorm:"type:varchar(50)" json:"conversion_apply_date"`           // 递交转正申请书时间（文本格式）
	ConversionBranchMeetingDate *string   `gorm:"type:varchar(50)" json:"conversion_branch_meeting_date"`  // 支部大会讨论时间（转正）（文本格式）
	ProbationDate               *string   `gorm:"type:varchar(50)" json:"probation_date"`                  // 预备党员时间（文本格式）
	ConversionDate              *string   `gorm:"type:varchar(50)" json:"conversion_date"`                 // 转正时间（文本格式）
	TransferDate                *string   `gorm:"type:varchar(50)" json:"transfer_date"`                   // 党组织关系转接时间（文本格式）
	IntroductionDate            *string   `gorm:"type:varchar(50)" json:"introduction_date"`               // 党组织关系介绍信时间（文本格式）
	Status                      string    `gorm:"type:ENUM('masses','activist','candidate','probationary','formal');default:'masses'" json:"status"`
	CreatedAt                   time.Time `json:"created_at"`
	UpdatedAt                   time.Time `json:"updated_at"`
	Student                     Student   `gorm:"foreignKey:StudentID" json:"student,omitempty"`
}
