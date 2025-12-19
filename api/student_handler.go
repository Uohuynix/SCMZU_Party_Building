package api

import (
	"SCMZU_Party_Building/entity"
	"SCMZU_Party_Building/service"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type StudentHandler struct {
	studentService *service.StudentService
	authService    *service.AuthService
}

func NewStudentHandler(studentService *service.StudentService, authService *service.AuthService) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
		authService:    authService,
	}
}

func (h *StudentHandler) ListStudents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	usernameVal, _ := c.Get("username")
	username, _ := usernameVal.(string)
	userIDVal, _ := c.Get("userID")
	userID, _ := userIDVal.(uint)

	var (
		students []entity.Student
		total    int64
		err      error
	)

	switch role {
	case "student":
		// 学生：只返回自己的记录（约定用户名=学号）
		student, errGet := h.studentService.GetStudentByStudentNo(username)
		if errGet != nil {
			c.JSON(http.StatusOK, gin.H{
				"data":  []entity.Student{},
				"total": 0,
				"page":  page,
			})
			return
		}
		students = []entity.Student{*student}
		total = 1
	case "teacher":
		// 教师：只能查看本支部、本小组成员
		user, errProfile := h.authService.GetProfile(userID)
		if errProfile != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load teacher profile"})
			return
		}
		students, total, err = h.studentService.ListStudentsByBranchAndGroup(user.Branch, user.GroupName, page, pageSize)
	default: // admin 或其他：查看全部
		students, total, err = h.studentService.ListStudents(page, pageSize)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  students,
		"total": total,
		"page":  page,
	})
}

func (h *StudentHandler) GetStudentDetail(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	student, err := h.studentService.GetStudent(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：学生只能看自己；教师只能看本小组；管理员可看所有
	if !h.canAccessStudent(c, student) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to access this student"})
		return
	}

	c.JSON(http.StatusOK, student)
}

func (h *StudentHandler) CreateStudent(c *gin.Context) {
	var student entity.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	// 学生：允许为自己创建一条党员档案，但不允许为他人创建
	if role == "student" {
		usernameVal, _ := c.Get("username")
		username, _ := usernameVal.(string)

		// 学生用户名约定为学号，只能为自己建档
		if username == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "username is required for student role"})
			return
		}

		// 如果已经存在档案，则直接返回该档案，视为“已创建”，避免重复报错
		if existing, _ := h.studentService.GetStudentByStudentNo(username); existing != nil {
			c.JSON(http.StatusOK, existing)
			return
		}

		// 强制把学号设置为当前登录学生，避免伪造他人档案
		student.StudentNo = username
		// 性别字段是 ENUM('male','female')，为空会导致 MySQL 截断错误，给一个安全默认值
		if student.Gender != "male" && student.Gender != "female" {
			student.Gender = "male"
		}
		// 党员类型为空时默认 masses，避免空值导致后续统计异常
		if student.Type == "" {
			student.Type = "masses"
		}
	} else if role != "teacher" && role != "admin" {
		// 其他未知角色不允许创建
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to create student record"})
		return
	}

	err := h.studentService.CreateStudent(&student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, student)
}

func (h *StudentHandler) UpdateStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var req entity.Student
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existing, err := h.studentService.GetStudent(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：学生只能修改自己的档案；教师只能修改本小组；管理员无限制
	if !h.canAccessStudent(c, existing) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to update this student"})
		return
	}

	// 合并更新：以数据库中的 existing 为基础，只更新前端允许修改的字段，
	// 避免因为未传字段而把学号等必填字段清空导致保存失败。
	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Gender != "" {
		existing.Gender = req.Gender
	}
	if req.Ethnicity != "" {
		existing.Ethnicity = req.Ethnicity
	}
	if req.BirthDate != nil {
		existing.BirthDate = req.BirthDate
	}
	if req.Education != "" {
		existing.Education = req.Education
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.IDCard != "" {
		existing.IDCard = req.IDCard
	}
	if req.MajorClass != "" {
		existing.MajorClass = req.MajorClass
	}
	if req.Type != "" {
		existing.Type = req.Type
	}
	// 照片 URL 允许前端一起修改（例如重置为某个地址）
	if req.PhotoURL != "" {
		existing.PhotoURL = req.PhotoURL
	}

	// 保持支部/小组信息不被学生随意修改；教师/管理员可以通过单独接口修改
	if roleVal, ok := c.Get("role"); ok && roleVal == "student" {
		existing.Branch = existing.Branch
		existing.GroupName = existing.GroupName
	} else {
		if req.Branch != "" {
			existing.Branch = req.Branch
		}
		if req.GroupName != "" {
			existing.GroupName = req.GroupName
		}
	}

	err = h.studentService.UpdateStudent(existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *StudentHandler) DeleteStudent(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// 只有管理员可以删除学员
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admin can delete students"})
		return
	}

	err = h.studentService.DeleteStudent(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

func (h *StudentHandler) GetDevelopment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	development, err := h.studentService.GetDevelopment(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, development)
}

func (h *StudentHandler) UpdateDevelopment(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var development entity.Development
	if err := c.ShouldBindJSON(&development); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	development.StudentID = uint(studentID)
	err = h.studentService.UpdateDevelopment(&development)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, development)
}

func (h *StudentHandler) GetRewards(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	rewards, err := h.studentService.GetRewards(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rewards)
}

func (h *StudentHandler) CreateReward(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var reward entity.Reward
	if err := c.ShouldBindJSON(&reward); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	reward.StudentID = uint(studentID)
	err = h.studentService.CreateReward(&reward)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, reward)
}

// UploadPhoto 上传或更新学生照片，保存文件到本地并更新 photo_url
func (h *StudentHandler) UploadPhoto(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	student, err := h.studentService.GetStudent(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：沿用 canAccessStudent
	if !h.canAccessStudent(c, student) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to update this student photo"})
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "photo file is required"})
		return
	}

	// 照片保存目录：behind/uploads/photos
	saveDir := "uploads/photos"
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	ext := filepath.Ext(file.Filename)
	if ext == "" {
		ext = ".jpg"
	}
	filename := strconv.Itoa(id) + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ext
	fullPath := filepath.Join(saveDir, filename)

	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save photo"})
		return
	}

	// 为前端提供访问路径，这里简单使用相对路径，前端可拼接服务器地址
	student.PhotoURL = "/" + fullPath
	if err := h.studentService.UpdateStudent(student); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update student photo url"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"photo_url": student.PhotoURL,
	})
}

// canAccessStudent 根据当前登录用户角色，判断是否有权访问/修改指定学生
func (h *StudentHandler) canAccessStudent(c *gin.Context, student *entity.Student) bool {
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	usernameVal, _ := c.Get("username")
	username, _ := usernameVal.(string)
	userIDVal, _ := c.Get("userID")
	userID, _ := userIDVal.(uint)

	switch role {
	case "admin":
		return true
	case "student":
		// 约定学生用户名=学号
		return student.StudentNo == username
	case "teacher":
		user, err := h.authService.GetProfile(userID)
		if err != nil {
			return false
		}
		// 必须在相同支部+小组内
		if user.Branch == "" || user.GroupName == "" {
			return false
		}
		return student.Branch == user.Branch && student.GroupName == user.GroupName
	default:
		return false
	}
}
