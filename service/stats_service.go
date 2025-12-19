package service

import (
	"SCMZU_Party_Building/dao"
	"SCMZU_Party_Building/entity"
)

type StatsService struct {
	trainingDao *dao.TrainingDAO
}

func NewStatsService(trainingDao *dao.TrainingDAO) *StatsService {
	return &StatsService{trainingDao: trainingDao}
}

func (s *StatsService) GetTrainingStats() (map[string]int, error) {
	// Implement training statistics logic here
	// For example: count number of people in each grade level
	stats := make(map[string]int)

	// Get all training records
	var trainings []entity.Training
	err := s.trainingDao.DB.Find(&trainings).Error
	if err != nil {
		return nil, err
	}

	// Count number of people for each grade
	for _, training := range trainings {
		if training.Score != "" {
			stats[training.Score]++
		}
	}

	return stats, nil
}
