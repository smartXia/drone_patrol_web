package services

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"drone-patrol-backend/internal/models"
)

type ErrorCodeService struct{}

func NewErrorCodeService() *ErrorCodeService {
	return &ErrorCodeService{}
}

// 获取错误码列表
func (s *ErrorCodeService) GetErrorCodes() (*models.APIResponse, error) {
	// 读取错误码文档
	errorCodesFile := filepath.Join("..", "..", "public", "docs", "dji-error-codes.md")
	file, err := os.Open(errorCodesFile)
	if err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("读取错误码文件失败: %v", err),
			Data:    []models.ErrorCode{},
		}, err
	}
	defer file.Close()

	var errorCodes []models.ErrorCode
	scanner := bufio.NewScanner(file)

	// 跳过标题行
	scanner.Scan()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Split(line, "\t")
		if len(parts) >= 2 {
			code := strings.TrimSpace(parts[0])
			description := strings.TrimSpace(parts[1])

			// 根据错误码范围确定分类
			codeNum, _ := strconv.Atoi(code)
			category := s.getErrorCategory(codeNum)
			severity := s.getErrorSeverity(description)

			errorCodes = append(errorCodes, models.ErrorCode{
				Code:        code,
				Description: description,
				Category:    category,
				Severity:    severity,
			})
		}
	}

	if err := scanner.Err(); err != nil {
		return &models.APIResponse{
			Code:    1,
			Message: fmt.Sprintf("读取错误码文件失败: %v", err),
			Data:    []models.ErrorCode{},
		}, err
	}

	return &models.APIResponse{
		Code:    0,
		Message: "ok",
		Data:    errorCodes,
	}, nil
}

// 根据错误码确定分类
func (s *ErrorCodeService) getErrorCategory(codeNum int) string {
	switch {
	case codeNum >= 1000 && codeNum < 2000:
		return "系统错误"
	case codeNum >= 2000 && codeNum < 3000:
		return "通信错误"
	case codeNum >= 3000 && codeNum < 4000:
		return "传感器错误"
	case codeNum >= 4000 && codeNum < 5000:
		return "电池错误"
	case codeNum >= 5000 && codeNum < 6000:
		return "电机错误"
	case codeNum >= 6000 && codeNum < 7000:
		return "GPS错误"
	case codeNum >= 7000 && codeNum < 8000:
		return "相机错误"
	case codeNum >= 8000 && codeNum < 9000:
		return "云台错误"
	case codeNum >= 9000 && codeNum < 10000:
		return "遥控器错误"
	default:
		return "未知错误"
	}
}

// 根据错误描述确定严重程度
func (s *ErrorCodeService) getErrorSeverity(description string) string {
	description = strings.ToLower(description)

	switch {
	case strings.Contains(description, "严重") || strings.Contains(description, "致命") || strings.Contains(description, "critical"):
		return "严重"
	case strings.Contains(description, "警告") || strings.Contains(description, "warning"):
		return "警告"
	case strings.Contains(description, "注意") || strings.Contains(description, "notice"):
		return "注意"
	case strings.Contains(description, "信息") || strings.Contains(description, "info"):
		return "信息"
	default:
		return "未知"
	}
}
