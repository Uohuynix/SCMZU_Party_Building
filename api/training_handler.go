package api

import (
	"SCMZU_Party_Building/entity"
	"SCMZU_Party_Building/service"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// parseDateString 解析日期字符串，支持多种格式
func parseDateString(dateStr string) *time.Time {
	if dateStr == "" || dateStr == "null" || dateStr == "undefined" {
		return nil
	}

	// 去除前后空格
	dateStr = strings.TrimSpace(dateStr)

	// 尝试的日期格式列表（按优先级排序）
	dateFormats := []string{
		"2006-01-02",                // 简单日期格式（最常见）
		"2006-01-02T15:04:05",       // ISO 8601 格式（无时区）
		"2006-01-02 15:04:05",       // 带时间的格式
		"2006-01-02T15:04:05Z",      // ISO 8601 带 Z
		"2006-01-02T15:04:05Z07:00", // ISO 8601 带时区
	}

	for _, format := range dateFormats {
		if parsed, err := time.Parse(format, dateStr); err == nil {
			return &parsed
		}
	}

	// 如果所有格式都失败，返回 nil
	return nil
}

type TrainingHandler struct {
	trainingService *service.TrainingService
}

func NewTrainingHandler(trainingService *service.TrainingService) *TrainingHandler {
	return &TrainingHandler{trainingService: trainingService}
}

func (h *TrainingHandler) ListTrainings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	trainings, total, err := h.trainingService.ListTrainings(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  trainings,
		"total": total,
		"page":  page,
	})
}

func (h *TrainingHandler) GetTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid training ID"})
		return
	}

	training, err := h.trainingService.GetTraining(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training not found"})
		return
	}

	c.JSON(http.StatusOK, training)
}

func (h *TrainingHandler) CreateTraining(c *gin.Context) {
	// 使用结构体接收请求数据，支持可选字段
	var req struct {
		Period        string `json:"period"`
		StudentID     *uint  `json:"student_id"`
		Unit          string `json:"unit"`
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		Score         string `json:"score"`
		CertificateNo string `json:"certificate_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证必填字段
	if req.Period == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Period is required"})
		return
	}
	if req.Score == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Score is required"})
		return
	}

	// 验证 score 是否为有效的枚举值
	validScores := map[string]bool{
		"excellent": true,
		"good":      true,
		"pass":      true,
		"fail":      true,
	}
	if !validScores[req.Score] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score value. Must be one of: excellent, good, pass, fail"})
		return
	}

	// 构建 Training 实体
	training := entity.Training{
		Period:        req.Period,
		Unit:          req.Unit,
		Score:         req.Score,
		CertificateNo: req.CertificateNo,
	}

	// 处理 student_id（如果为 null 或 0，则设置为 nil，允许数据库使用 NULL）
	if req.StudentID != nil && *req.StudentID > 0 {
		training.StudentID = req.StudentID
	} else {
		training.StudentID = nil
	}

	// 解析日期字符串
	training.StartDate = parseDateString(req.StartDate)
	training.EndDate = parseDateString(req.EndDate)

	err := h.trainingService.CreateTraining(&training)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, training)
}

func (h *TrainingHandler) UpdateTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid training ID"})
		return
	}

	// 接收前端发送的所有字段
	var req struct {
		Period        string `json:"period"`
		StudentID     *uint  `json:"student_id"`
		Unit          string `json:"unit"`
		StartDate     string `json:"start_date"`
		EndDate       string `json:"end_date"`
		Score         string `json:"score"`
		CertificateNo string `json:"certificate_no"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 先查出原记录，再按字段合并更新
	existing, err := h.trainingService.GetTraining(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Training not found"})
		return
	}

	// 更新字段（只更新非空字段）
	if req.Period != "" {
		existing.Period = req.Period
	}
	if req.Score != "" {
		// 验证 score 是否为有效的枚举值
		validScores := map[string]bool{
			"excellent": true,
			"good":      true,
			"pass":      true,
			"fail":      true,
		}
		if !validScores[req.Score] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid score value. Must be one of: excellent, good, pass, fail"})
			return
		}
		existing.Score = req.Score
	}
	if req.Unit != "" {
		existing.Unit = req.Unit
	}
	if req.CertificateNo != "" {
		existing.CertificateNo = req.CertificateNo
	}
	if req.StudentID != nil && *req.StudentID > 0 {
		existing.StudentID = req.StudentID
	} else if req.StudentID != nil && *req.StudentID == 0 {
		existing.StudentID = nil
	}

	// 解析日期字符串
	if req.StartDate != "" {
		existing.StartDate = parseDateString(req.StartDate)
	}
	if req.EndDate != "" {
		existing.EndDate = parseDateString(req.EndDate)
	}

	err = h.trainingService.UpdateTraining(existing)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existing)
}

func (h *TrainingHandler) DeleteTraining(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid training ID"})
		return
	}

	err = h.trainingService.DeleteTraining(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Training deleted successfully"})
}
