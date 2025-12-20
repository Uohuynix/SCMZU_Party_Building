package api

import (
	"SCMZU_Party_Building/entity"
	"SCMZU_Party_Building/service"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	// 如果pageSize为0或大于10000，视为不分页请求
	noPagination := pageSize == 0 || pageSize > 10000

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
		// 教师：可以查看本支部所有小组成员
		user, errProfile := h.authService.GetProfile(userID)
		if errProfile != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load teacher profile"})
			return
		}
		if user.Branch == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "teacher branch is empty, please contact administrator to set your branch"})
			return
		}
		// 使用ListStudentsByBranchAndGroup，但groupName传空字符串，这样会查询本支部的所有成员
		if noPagination {
			// 教师也可以不分页获取本支部所有成员
			allStudents, errAll := h.studentService.ListAllStudents()
			if errAll != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": errAll.Error()})
				return
			}
			// 过滤出本支部的成员
			var filteredStudents []entity.Student
			for _, s := range allStudents {
				if s.Branch == user.Branch {
					filteredStudents = append(filteredStudents, s)
				}
			}
			students = filteredStudents
			total = int64(len(students))
		} else {
			students, total, err = h.studentService.ListStudentsByBranchAndGroup(user.Branch, "", page, pageSize)
		}
	default: // admin 或其他：查看全部
		if noPagination {
			// 管理员不分页，返回所有数据
			students, err = h.studentService.ListAllStudents()
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			total = int64(len(students))
		} else {
			students, total, err = h.studentService.ListStudents(page, pageSize)
		}
	}

	if err != nil {
		log.Printf("ListStudents error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 确保即使students为空也返回空数组，而不是nil
	if students == nil {
		students = []entity.Student{}
	}

	log.Printf("ListStudents success: role=%s, total=%d, count=%d", role, total, len(students))
	c.JSON(http.StatusOK, gin.H{
		"data":  students,
		"total": total,
		"page":  page,
	})
}

// GetAllBranches 获取所有党支部列表
func (h *StudentHandler) GetAllBranches(c *gin.Context) {
	branches, err := h.studentService.GetAllBranches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": branches})
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

	c.JSON(http.StatusOK, gin.H{"data": student})
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
	} else if role == "teacher" {
		// 教师：如果前端没有传branch，自动设置为教师所在的党支部
		if student.Branch == "" {
			userIDVal, _ := c.Get("userID")
			userID, _ := userIDVal.(uint)
			user, errProfile := h.authService.GetProfile(userID)
			if errProfile != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load teacher profile"})
				return
			}
			if user.Branch == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "teacher branch is empty, please contact administrator to set your branch"})
				return
			}
			student.Branch = user.Branch
		}
		// 确保性别和类型有默认值
		if student.Gender != "male" && student.Gender != "female" {
			student.Gender = "male"
		}
		if student.Type == "" {
			student.Type = "masses"
		}
	} else if role != "admin" {
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

	// 合并更新：以数据库中的 existing 为基础，只更新前端允许修改的字段
	// 注意：这里使用指针检查来判断字段是否被传递，允许空字符串清空字段
	// 但为了避免误操作，对于必填字段（如Name）需要特殊处理

	// 姓名是必填字段，如果传递了空字符串，保持原值
	if req.Name != "" {
		existing.Name = req.Name
	}

	// 性别字段：如果传递了值则更新（包括空字符串，但需要验证有效性）
	if req.Gender != "" {
		// 验证性别值是否有效
		if req.Gender == "male" || req.Gender == "female" {
			existing.Gender = req.Gender
		}
	}

	// 可选字段：允许空字符串清空
	existing.Ethnicity = req.Ethnicity
	existing.Education = req.Education
	existing.Phone = req.Phone
	existing.IDCard = req.IDCard
	existing.MajorClass = req.MajorClass
	existing.PhotoURL = req.PhotoURL

	// 日期字段：如果传递了值则更新
	if req.BirthDate != nil {
		existing.BirthDate = req.BirthDate
	}

	// 党员类型：如果传递了值则更新（需要验证有效性）
	if req.Type != "" {
		validTypes := map[string]bool{
			"masses":       true,
			"activist":     true,
			"candidate":    true,
			"probationary": true,
			"formal":       true,
		}
		if validTypes[req.Type] {
			existing.Type = req.Type
		}
	}

	// 保持支部/小组信息不被学生随意修改；教师/管理员可以通过单独接口修改
	if roleVal, ok := c.Get("role"); ok && roleVal == "student" {
		// 学生不能修改支部和小组信息，保持原值不变
		// 不需要赋值，因为existing已经是数据库中的值
	} else {
		// 教师和管理员可以修改支部和小组信息
		existing.Branch = req.Branch
		existing.GroupName = req.GroupName
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

	// 先获取学生信息，用于权限检查
	student, err := h.studentService.GetStudent(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：学生只能看自己；教师只能看本支部；管理员可看所有
	if !h.canAccessStudent(c, student) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to access this student's development"})
		return
	}

	development, err := h.studentService.GetDevelopment(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": development})
}

func (h *StudentHandler) UpdateDevelopment(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// 先获取学生信息，用于权限检查
	student, err := h.studentService.GetStudent(uint(studentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：学生只能修改自己的档案；教师只能修改本支部；管理员无限制
	if !h.canAccessStudent(c, student) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to update this student's development"})
		return
	}

	// 使用自定义结构体接收请求，避免日期解析问题
	var req struct {
		StudentID                   uint    `json:"student_id"`
		JoinLeagueDate              *string `json:"join_league_date"`
		ApplyDate                   *string `json:"apply_date"`
		PartyTalkDate               *string `json:"party_talk_date"`
		LeagueRecommendDate         *string `json:"league_recommend_date"`
		ActivistDate                *string `json:"activist_date"`
		CandidateDate               *string `json:"candidate_date"`
		SuperiorReportDate          *string `json:"superior_report_date"`
		CommitteeReviewDate         *string `json:"committee_review_date"`
		PartyCommitteePreReviewDate *string `json:"party_committee_pre_review_date"`
		BranchMeetingDate           *string `json:"branch_meeting_date"`
		ConversionApplyDate         *string `json:"conversion_apply_date"`
		ConversionBranchMeetingDate *string `json:"conversion_branch_meeting_date"`
		ProbationDate               *string `json:"probation_date"`
		ConversionDate              *string `json:"conversion_date"`
		TransferDate                *string `json:"transfer_date"`
		IntroductionDate            *string `json:"introduction_date"`
		Status                      string  `json:"status"`
	}

	// 使用 ShouldBindBodyWith 或直接读取 JSON 来避免日期解析问题
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("UpdateDevelopment: JSON binding error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	// 调试日志：打印接收到的数据
	log.Printf("UpdateDevelopment: Received data for studentID=%d, join_league_date=%v", studentID, req.JoinLeagueDate)

	// 先获取现有的发展信息（如果存在）
	existing, err := h.studentService.GetDevelopment(uint(studentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 创建或更新发展信息
	var development entity.Development
	if existing != nil {
		development = *existing
	} else {
		development.StudentID = uint(studentID)
	}

	// 辅助函数：处理字符串字段，将空字符串转换为 nil
	processStringField := func(val *string) *string {
		if val == nil {
			return nil
		}
		trimmed := strings.TrimSpace(*val)
		if trimmed == "" {
			return nil
		}
		return &trimmed
	}

	// 更新字段（只更新请求中提供的字段，处理空字符串为 null）
	if req.JoinLeagueDate != nil {
		development.JoinLeagueDate = processStringField(req.JoinLeagueDate)
	}
	if req.ApplyDate != nil {
		development.ApplyDate = processStringField(req.ApplyDate)
	}
	if req.PartyTalkDate != nil {
		development.PartyTalkDate = processStringField(req.PartyTalkDate)
	}
	if req.LeagueRecommendDate != nil {
		development.LeagueRecommendDate = processStringField(req.LeagueRecommendDate)
	}
	if req.ActivistDate != nil {
		development.ActivistDate = processStringField(req.ActivistDate)
	}
	if req.CandidateDate != nil {
		development.CandidateDate = processStringField(req.CandidateDate)
	}
	if req.SuperiorReportDate != nil {
		development.SuperiorReportDate = processStringField(req.SuperiorReportDate)
	}
	if req.CommitteeReviewDate != nil {
		development.CommitteeReviewDate = processStringField(req.CommitteeReviewDate)
	}
	if req.PartyCommitteePreReviewDate != nil {
		development.PartyCommitteePreReviewDate = processStringField(req.PartyCommitteePreReviewDate)
	}
	if req.BranchMeetingDate != nil {
		development.BranchMeetingDate = processStringField(req.BranchMeetingDate)
	}
	if req.ConversionApplyDate != nil {
		development.ConversionApplyDate = processStringField(req.ConversionApplyDate)
	}
	if req.ConversionBranchMeetingDate != nil {
		development.ConversionBranchMeetingDate = processStringField(req.ConversionBranchMeetingDate)
	}
	if req.ProbationDate != nil {
		development.ProbationDate = processStringField(req.ProbationDate)
	}
	if req.ConversionDate != nil {
		development.ConversionDate = processStringField(req.ConversionDate)
	}
	if req.TransferDate != nil {
		development.TransferDate = processStringField(req.TransferDate)
	}
	if req.IntroductionDate != nil {
		development.IntroductionDate = processStringField(req.IntroductionDate)
	}
	if req.Status != "" {
		development.Status = req.Status
	}

	err = h.studentService.UpdateDevelopment(&development)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": development})
}

func (h *StudentHandler) GetRewards(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// 先获取学生信息，用于权限检查
	student, err := h.studentService.GetStudent(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：学生只能看自己；教师只能看本支部；管理员可看所有
	if !h.canAccessStudent(c, student) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to access this student's rewards"})
		return
	}

	rewards, err := h.studentService.GetRewards(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": rewards})
}

func (h *StudentHandler) CreateReward(c *gin.Context) {
	studentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	// 先获取学生信息，用于权限检查
	student, err := h.studentService.GetStudent(uint(studentID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// 权限控制：学生不能创建奖惩记录；教师只能为本支部学生创建；管理员无限制
	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)
	if role == "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "students cannot create reward records"})
		return
	}
	if !h.canAccessStudent(c, student) {
		c.JSON(http.StatusForbidden, gin.H{"error": "no permission to create reward for this student"})
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

	c.JSON(http.StatusCreated, gin.H{"data": reward})
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
			log.Printf("canAccessStudent: failed to get teacher profile (userID: %d), error: %v", userID, err)
			return false
		}
		// 教师可以查看本支部的所有成员（只需要支部相同）
		if user.Branch == "" {
			log.Printf("canAccessStudent: teacher branch is empty (userID: %d, username: %s)", userID, user.Username)
			return false
		}
		// 如果学生的支部为空，也允许访问（可能是未分配支部的学生）
		if student.Branch == "" {
			// 允许教师访问未分配支部的学生，或者根据业务需求决定是否允许
			// 这里选择允许，因为教师可能负责管理未分配的学生
			log.Printf("canAccessStudent: student branch is empty, allowing access (studentID: %d, teacherBranch: %s)", student.ID, user.Branch)
			return true
		}
		// 检查支部是否匹配（忽略大小写和前后空格）
		teacherBranch := strings.TrimSpace(user.Branch)
		studentBranch := strings.TrimSpace(student.Branch)
		match := strings.EqualFold(teacherBranch, studentBranch)
		if !match {
			log.Printf("canAccessStudent: branch mismatch - teacher: '%s' (userID: %d), student: '%s' (studentID: %d)",
				teacherBranch, userID, studentBranch, student.ID)
		}
		return match
	default:
		return false
	}
}
