package dao

import (
	"SCMZU_Party_Building/entity"
	"errors"

	"gorm.io/gorm"
)

type StudentDAO struct {
	BaseDAO
}

func NewStudentDAO(db *gorm.DB) *StudentDAO {
	return &StudentDAO{BaseDAO: BaseDAO{DB: db}}
}

func (dao *StudentDAO) FindAll(page, pageSize int) ([]entity.Student, int64, error) {
	var students []entity.Student
	var total int64

	offset := (page - 1) * pageSize
	err := dao.DB.Model(&entity.Student{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = dao.DB.Offset(offset).Limit(pageSize).Find(&students).Error
	return students, total, err
}

// FindAllWithoutPagination 不分页查询所有学生，用于管理员查看所有数据
func (dao *StudentDAO) FindAllWithoutPagination() ([]entity.Student, error) {
	var students []entity.Student
	err := dao.DB.Order("id ASC").Find(&students).Error
	// 确保返回空数组而不是nil
	if students == nil {
		students = []entity.Student{}
	}
	return students, err
}

// FindAllBranches 获取所有不重复的党支部列表
func (dao *StudentDAO) FindAllBranches() ([]string, error) {
	var branches []string
	err := dao.DB.Model(&entity.Student{}).
		Distinct("branch").
		Where("branch != '' AND branch IS NOT NULL").
		Pluck("branch", &branches).Error
	return branches, err
}

func (dao *StudentDAO) FindByID(id uint) (*entity.Student, error) {
	var student entity.Student
	err := dao.DB.First(&student, id).Error
	return &student, err
}

func (dao *StudentDAO) Create(student *entity.Student) error {
	return dao.DB.Create(student).Error
}

func (dao *StudentDAO) Update(student *entity.Student) error {
	return dao.DB.Save(student).Error
}

func (dao *StudentDAO) Delete(id uint) error {
	return dao.DB.Delete(&entity.Student{}, id).Error
}

func (dao *StudentDAO) FindByStudentNo(studentNo string) (*entity.Student, error) {
	var student entity.Student
	err := dao.DB.Where("student_no = ?", studentNo).First(&student).Error
	if err != nil {
		// 如果未找到记录，返回 nil 方便上层判断“尚未建档”
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

// FindByBranchAndGroup 按支部+小组分页查询，用于教师查看本小组成员
func (dao *StudentDAO) FindByBranchAndGroup(branch, groupName string, page, pageSize int) ([]entity.Student, int64, error) {
	var students []entity.Student
	var total int64

	offset := (page - 1) * pageSize
	query := dao.DB.Model(&entity.Student{})

	if branch != "" {
		query = query.Where("branch = ?", branch)
	}
	if groupName != "" {
		query = query.Where("group_name = ?", groupName)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(offset).Limit(pageSize).Find(&students).Error; err != nil {
		return nil, 0, err
	}

	return students, total, nil
}

func (dao *StudentDAO) FindByType(studentType string, page, pageSize int) ([]entity.Student, int64, error) {
	var students []entity.Student
	var total int64

	offset := (page - 1) * pageSize
	err := dao.DB.Model(&entity.Student{}).Where("type = ?", studentType).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = dao.DB.Where("type = ?", studentType).Offset(offset).Limit(pageSize).Find(&students).Error
	return students, total, err
}

func (dao *StudentDAO) SearchByName(name string, page, pageSize int) ([]entity.Student, int64, error) {
	var students []entity.Student
	var total int64

	offset := (page - 1) * pageSize
	err := dao.DB.Model(&entity.Student{}).Where("name LIKE ?", "%"+name+"%").Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = dao.DB.Where("name LIKE ?", "%"+name+"%").Offset(offset).Limit(pageSize).Find(&students).Error
	return students, total, err
}
