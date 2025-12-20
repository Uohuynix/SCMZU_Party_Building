package service

import (
	"SCMZU_Party_Building/dao"
	"SCMZU_Party_Building/entity"
)

type StudentService struct {
	studentDao     *dao.StudentDAO
	developmentDao *dao.DevelopmentDAO
	rewardDao      *dao.RewardDAO
}

func NewStudentService(
	studentDao *dao.StudentDAO,
	developmentDao *dao.DevelopmentDAO,
	rewardDao *dao.RewardDAO,
) *StudentService {
	return &StudentService{
		studentDao:     studentDao,
		developmentDao: developmentDao,
		rewardDao:      rewardDao,
	}
}

func (s *StudentService) ListStudents(page, pageSize int) ([]entity.Student, int64, error) {
	return s.studentDao.FindAll(page, pageSize)
}

// ListAllStudents 不分页获取所有学生，用于管理员
func (s *StudentService) ListAllStudents() ([]entity.Student, error) {
	return s.studentDao.FindAllWithoutPagination()
}

// GetAllBranches 获取所有党支部列表
func (s *StudentService) GetAllBranches() ([]string, error) {
	return s.studentDao.FindAllBranches()
}

func (s *StudentService) GetStudent(id uint) (*entity.Student, error) {
	return s.studentDao.FindByID(id)
}

// GetStudentByStudentNo 根据学号获取学生，供“学生本人”权限判断使用
func (s *StudentService) GetStudentByStudentNo(studentNo string) (*entity.Student, error) {
	return s.studentDao.FindByStudentNo(studentNo)
}

// ListStudentsByBranchAndGroup 教师按支部+小组查看成员
func (s *StudentService) ListStudentsByBranchAndGroup(branch, groupName string, page, pageSize int) ([]entity.Student, int64, error) {
	return s.studentDao.FindByBranchAndGroup(branch, groupName, page, pageSize)
}

func (s *StudentService) CreateStudent(student *entity.Student) error {
	return s.studentDao.Create(student)
}

func (s *StudentService) UpdateStudent(student *entity.Student) error {
	return s.studentDao.Update(student)
}

func (s *StudentService) DeleteStudent(id uint) error {
	return s.studentDao.Delete(id)
}

func (s *StudentService) GetDevelopment(studentID uint) (*entity.Development, error) {
	return s.developmentDao.FindByStudentID(studentID)
}

func (s *StudentService) UpdateDevelopment(development *entity.Development) error {
	if development.ID == 0 {
		return s.developmentDao.Create(development)
	}
	return s.developmentDao.Update(development)
}

func (s *StudentService) GetRewards(studentID uint) ([]entity.Reward, error) {
	return s.rewardDao.FindByStudentID(studentID)
}

func (s *StudentService) CreateReward(reward *entity.Reward) error {
	return s.rewardDao.Create(reward)
}

func (s *StudentService) UpdateReward(reward *entity.Reward) error {
	return s.rewardDao.Update(reward)
}

func (s *StudentService) DeleteReward(id uint) error {
	return s.rewardDao.Delete(id)
}
