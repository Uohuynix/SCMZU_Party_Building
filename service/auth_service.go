package service

import (
	"SCMZU_Party_Building/dao"
	"SCMZU_Party_Building/entity"
	"SCMZU_Party_Building/util"
	"errors"

	"gorm.io/gorm"
)

type AuthService struct {
	userDao *dao.UserDAO
}

func NewAuthService(userDao *dao.UserDAO) *AuthService {
	return &AuthService{userDao: userDao}
}

func (s *AuthService) Login(username, password string) (*entity.User, error) {
	user, err := s.userDao.FindByUsername(username)
	if err != nil {
		return nil, err
	}

	if !util.ComparePassword(user.Password, password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) Register(username, password, name, role, branch, groupName string) (*entity.User, error) {
	// 检查用户名是否已存在，避免数据库层面直接报错
	if _, err := s.userDao.FindByUsername(username); err == nil {
		return nil, ErrUserExists
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 只有“未找到”才允许继续注册，其他错误直接返回
		return nil, err
	}

	hashedPassword, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &entity.User{
		Username:  username,
		Password:  hashedPassword,
		Name:      name,
		Role:      role,
		Branch:    branch,
		GroupName: groupName,
	}

	err = s.userDao.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *AuthService) GetProfile(userID uint) (*entity.User, error) {
	return s.userDao.FindByID(userID)
}

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrUserExists = errors.New("username already exists")
