package dao

import (
	"SCMZU_Party_Building/entity"

	"gorm.io/gorm"
)

type DevelopmentDAO struct {
	BaseDAO
}

func NewDevelopmentDAO(db *gorm.DB) *DevelopmentDAO {
	return &DevelopmentDAO{BaseDAO: BaseDAO{DB: db}}
}

func (dao *DevelopmentDAO) FindByStudentID(studentID uint) (*entity.Development, error) {
	var development entity.Development
	err := dao.DB.Where("student_id = ?", studentID).Preload("Student").First(&development).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &development, err
}

func (dao *DevelopmentDAO) Create(development *entity.Development) error {
	// 使用 Omit 排除关联和自动时间戳字段，确保字符串字段正确保存
	// GORM 会自动处理 CreatedAt 和 UpdatedAt，所以不需要手动设置
	return dao.DB.Omit("Student").Create(development).Error
}

func (dao *DevelopmentDAO) Update(development *entity.Development) error {
	// 使用 map 更新，明确指定字段类型为字符串，避免 GORM 自动类型转换
	updates := map[string]interface{}{
		"join_league_date":                development.JoinLeagueDate,
		"apply_date":                      development.ApplyDate,
		"party_talk_date":                 development.PartyTalkDate,
		"league_recommend_date":           development.LeagueRecommendDate,
		"activist_date":                   development.ActivistDate,
		"candidate_date":                  development.CandidateDate,
		"superior_report_date":            development.SuperiorReportDate,
		"committee_review_date":           development.CommitteeReviewDate,
		"party_committee_pre_review_date": development.PartyCommitteePreReviewDate,
		"branch_meeting_date":             development.BranchMeetingDate,
		"conversion_apply_date":           development.ConversionApplyDate,
		"conversion_branch_meeting_date":  development.ConversionBranchMeetingDate,
		"probation_date":                  development.ProbationDate,
		"conversion_date":                 development.ConversionDate,
		"transfer_date":                   development.TransferDate,
		"introduction_date":               development.IntroductionDate,
		"status":                          development.Status,
	}
	return dao.DB.Model(development).Updates(updates).Error
}
