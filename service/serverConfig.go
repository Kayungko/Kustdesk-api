package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/model"
	"gorm.io/gorm"
)

type ServerConfigService struct{}

// GetConfigByCode 根据配置码获取配置
func (s *ServerConfigService) GetConfigByCode(code string) (*model.ClientServerConfig, error) {
	// 查找配置码
	var configCode model.ConfigCode
	result := global.DB.Preload("ServerConfig").Where("code = ? AND status = ?", code, model.COMMON_STATUS_ENABLE).First(&configCode)
	if result.Error != nil {
		return nil, fmt.Errorf("invalid config code")
	}

	// 检查过期时间
	if configCode.ExpiresAt != nil && time.Now().After(*configCode.ExpiresAt) {
		return nil, fmt.Errorf("config code has expired")
	}

	// 检查使用次数限制
	if configCode.MaxUsage != nil && configCode.UsageCount >= *configCode.MaxUsage {
		return nil, fmt.Errorf("config code usage limit exceeded")
	}

	// 检查服务器配置是否启用
	if configCode.ServerConfig == nil || configCode.ServerConfig.Status != model.COMMON_STATUS_ENABLE || !*configCode.ServerConfig.IsEnabled {
		return nil, fmt.Errorf("server config is disabled")
	}

	// 更新使用次数
	global.DB.Model(&configCode).Update("usage_count", configCode.UsageCount+1)

	// 返回客户端配置
	clientConfig := &model.ClientServerConfig{
		Name:        configCode.ServerConfig.Name,
		Region:      configCode.ServerConfig.Region,
		IdServer:    configCode.ServerConfig.IdServer,
		RelayServer: configCode.ServerConfig.RelayServer,
		ApiServer:   configCode.ServerConfig.ApiServer,
		Key:         configCode.ServerConfig.Key,
	}

	return clientConfig, nil
}

// RecordConfigCodeUsage 记录配置码使用
func (s *ServerConfigService) RecordConfigCodeUsage(code, clientIP, userAgent string) error {
	// 查找配置码ID
	var configCode model.ConfigCode
	if err := global.DB.Select("id").Where("code = ?", code).First(&configCode).Error; err != nil {
		return err
	}

	// 创建使用记录
	usage := model.ConfigCodeUsage{
		ConfigCodeId: configCode.Id,
		ClientIP:     clientIP,
		UserAgent:    userAgent,
		UsedAt:       time.Now(),
	}

	return global.DB.Create(&usage).Error
}

// List 获取服务器配置列表
func (s *ServerConfigService) List(query *admin.ServerConfigListQuery) (*model.ServerConfigList, error) {
	var configs []*model.ServerConfig
	var total int64

	db := global.DB.Model(&model.ServerConfig{})

	// 添加查询条件
	if query.Name != "" {
		db = db.Where("name LIKE ?", "%"+query.Name+"%")
	}
	if query.Region != "" {
		db = db.Where("region = ?", query.Region)
	}
	if query.IsEnabled != nil {
		db = db.Where("is_enabled = ?", *query.IsEnabled)
	}
	if query.IsDefault != nil {
		db = db.Where("is_default = ?", *query.IsDefault)
	}

	// 计算总数
	db.Count(&total)

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("priority DESC, created_at DESC").
		Limit(int(query.PageSize)).
		Offset(int(offset)).
		Find(&configs).Error

	if err != nil {
		return nil, err
	}

	return &model.ServerConfigList{
		ServerConfigs: configs,
		Pagination: model.Pagination{
			Page:     query.Page,
			PageSize: query.PageSize,
			Total:    total,
		},
	}, nil
}

// Create 创建服务器配置
func (s *ServerConfigService) Create(form *admin.ServerConfigForm) (*model.ServerConfig, error) {
	config := &model.ServerConfig{
		Name:        form.Name,
		Description: form.Description,
		Region:      form.Region,
		IdServer:    form.IdServer,
		RelayServer: form.RelayServer,
		ApiServer:   form.ApiServer,
		Key:         form.Key,
		Priority:    form.Priority,
		Status:      model.COMMON_STATUS_ENABLE,
	}

	// 设置默认值
	if form.IsEnabled != nil {
		config.IsEnabled = form.IsEnabled
	} else {
		enabled := false
		config.IsEnabled = &enabled
	}

	if form.IsDefault != nil {
		config.IsDefault = form.IsDefault
	} else {
		isDefault := false
		config.IsDefault = &isDefault
	}

	// 如果设置为默认配置，需要取消其他默认配置
	if config.IsDefault != nil && *config.IsDefault {
		if err := s.unsetAllDefault(); err != nil {
			return nil, err
		}
	}

	err := global.DB.Create(config).Error
	return config, err
}

// Update 更新服务器配置
func (s *ServerConfigService) Update(id uint, form *admin.ServerConfigForm) (*model.ServerConfig, error) {
	var config model.ServerConfig
	if err := global.DB.First(&config, id).Error; err != nil {
		return nil, err
	}

	// 如果设置为默认配置，需要取消其他默认配置
	if form.IsDefault != nil && *form.IsDefault {
		if err := s.unsetAllDefault(); err != nil {
			return nil, err
		}
	}

	// 更新字段
	config.Name = form.Name
	config.Description = form.Description
	config.Region = form.Region
	config.IdServer = form.IdServer
	config.RelayServer = form.RelayServer
	config.ApiServer = form.ApiServer
	config.Key = form.Key
	config.Priority = form.Priority

	if form.IsEnabled != nil {
		config.IsEnabled = form.IsEnabled
	}
	if form.IsDefault != nil {
		config.IsDefault = form.IsDefault
	}

	err := global.DB.Save(&config).Error
	return &config, err
}

// Delete 删除服务器配置
func (s *ServerConfigService) Delete(id uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 检查是否有关联的配置码
		var count int64
		if err := tx.Model(&model.ConfigCode{}).Where("server_config_id = ?", id).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			return fmt.Errorf("cannot delete server config with existing config codes")
		}

		// 删除配置
		return tx.Delete(&model.ServerConfig{}, id).Error
	})
}

// Detail 获取服务器配置详情
func (s *ServerConfigService) Detail(id uint) (*model.ServerConfig, error) {
	var config model.ServerConfig
	err := global.DB.First(&config, id).Error
	return &config, err
}

// SetDefault 设置默认配置
func (s *ServerConfigService) SetDefault(id uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 取消所有默认配置
		if err := s.unsetAllDefaultWithTx(tx); err != nil {
			return err
		}

		// 设置新的默认配置
		isDefault := true
		return tx.Model(&model.ServerConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
			"is_default": &isDefault,
		}).Error
	})
}

// GetDefault 获取默认配置
func (s *ServerConfigService) GetDefault() (*model.ServerConfig, error) {
	var config model.ServerConfig
	err := global.DB.Where("is_default = ? AND status = ?", true, model.COMMON_STATUS_ENABLE).First(&config).Error
	return &config, err
}

// GenerateConfigCode 生成配置码
func (s *ServerConfigService) GenerateConfigCode(form *admin.ConfigCodeForm, createdBy uint) (*model.ConfigCode, error) {
	// 检查服务器配置是否存在
	var serverConfig model.ServerConfig
	if err := global.DB.First(&serverConfig, form.ServerConfigId).Error; err != nil {
		return nil, fmt.Errorf("server config not found")
	}

	// 生成唯一配置码
	code, err := s.generateUniqueCode()
	if err != nil {
		return nil, err
	}

	configCode := &model.ConfigCode{
		Code:           code,
		ServerConfigId: form.ServerConfigId,
		ExpiresAt:      form.ExpiresAt,
		MaxUsage:       form.MaxUsage,
		CreatedBy:      createdBy,
		Status:         model.COMMON_STATUS_ENABLE,
	}

	err = global.DB.Create(configCode).Error
	return configCode, err
}

// BatchGenerateConfigCode 批量生成配置码
func (s *ServerConfigService) BatchGenerateConfigCode(form *admin.ConfigCodeBatchForm, createdBy uint) ([]*model.ConfigCode, error) {
	// 检查服务器配置是否存在
	var serverConfig model.ServerConfig
	if err := global.DB.First(&serverConfig, form.ServerConfigId).Error; err != nil {
		return nil, fmt.Errorf("server config not found")
	}

	var codes []*model.ConfigCode
	for i := 0; i < form.Count; i++ {
		code, err := s.generateUniqueCode()
		if err != nil {
			return nil, err
		}

		configCode := &model.ConfigCode{
			Code:           code,
			ServerConfigId: form.ServerConfigId,
			ExpiresAt:      form.ExpiresAt,
			MaxUsage:       form.MaxUsage,
			CreatedBy:      createdBy,
			Status:         model.COMMON_STATUS_ENABLE,
		}

		if err := global.DB.Create(configCode).Error; err != nil {
			return nil, err
		}

		codes = append(codes, configCode)
	}

	return codes, nil
}

// ConfigCodeList 获取配置码列表
func (s *ServerConfigService) ConfigCodeList(query *admin.ConfigCodeListQuery) (*model.ConfigCodeList, error) {
	var codes []*model.ConfigCode
	var total int64

	db := global.DB.Model(&model.ConfigCode{}).
		Preload("ServerConfig").
		Preload("Creator")

	// 添加查询条件
	if query.ServerConfigId > 0 {
		db = db.Where("server_config_id = ?", query.ServerConfigId)
	}
	if query.Code != "" {
		db = db.Where("code LIKE ?", "%"+query.Code+"%")
	}
	if query.CreatedBy > 0 {
		db = db.Where("created_by = ?", query.CreatedBy)
	}

	// 计算总数
	db.Count(&total)

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("created_at DESC").
		Limit(int(query.PageSize)).
		Offset(int(offset)).
		Find(&codes).Error

	if err != nil {
		return nil, err
	}

	return &model.ConfigCodeList{
		ConfigCodes: codes,
		Pagination: model.Pagination{
			Page:     query.Page,
			PageSize: query.PageSize,
			Total:    total,
		},
	}, nil
}

// DeleteConfigCode 删除配置码
func (s *ServerConfigService) DeleteConfigCode(id uint) error {
	return global.DB.Transaction(func(tx *gorm.DB) error {
		// 删除使用记录
		if err := tx.Where("config_code_id = ?", id).Delete(&model.ConfigCodeUsage{}).Error; err != nil {
			return err
		}
		// 删除配置码
		return tx.Delete(&model.ConfigCode{}, id).Error
	})
}

// GetConfigCodeStats 获取配置码统计
func (s *ServerConfigService) GetConfigCodeStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总配置码数
	var totalCodes int64
	global.DB.Model(&model.ConfigCode{}).Count(&totalCodes)
	stats["total_codes"] = totalCodes

	// 活跃配置码数（未过期且未达到使用限制）
	var activeCodes int64
	now := time.Now()
	global.DB.Model(&model.ConfigCode{}).Where(
		"status = ? AND (expires_at IS NULL OR expires_at > ?) AND (max_usage IS NULL OR usage_count < max_usage)",
		model.COMMON_STATUS_ENABLE, now).Count(&activeCodes)
	stats["active_codes"] = activeCodes

	// 过期配置码数
	var expiredCodes int64
	global.DB.Model(&model.ConfigCode{}).Where("expires_at IS NOT NULL AND expires_at <= ?", now).Count(&expiredCodes)
	stats["expired_codes"] = expiredCodes

	// 总使用次数
	var totalUsage int64
	global.DB.Model(&model.ConfigCode{}).Select("COALESCE(SUM(usage_count), 0)").Scan(&totalUsage)
	stats["total_usage"] = totalUsage

	// 今日使用次数
	today := time.Now().Format("2006-01-02")
	var todayUsage int64
	global.DB.Model(&model.ConfigCodeUsage{}).Where("DATE(used_at) = ?", today).Count(&todayUsage)
	stats["today_usage"] = todayUsage

	return stats, nil
}

// 辅助方法

// unsetAllDefault 取消所有默认配置
func (s *ServerConfigService) unsetAllDefault() error {
	isDefault := false
	return global.DB.Model(&model.ServerConfig{}).
		Where("is_default = ?", true).
		Update("is_default", &isDefault).Error
}

// unsetAllDefaultWithTx 在事务中取消所有默认配置
func (s *ServerConfigService) unsetAllDefaultWithTx(tx *gorm.DB) error {
	isDefault := false
	return tx.Model(&model.ServerConfig{}).
		Where("is_default = ?", true).
		Update("is_default", &isDefault).Error
}

// generateUniqueCode 生成唯一配置码
func (s *ServerConfigService) generateUniqueCode() (string, error) {
	for i := 0; i < 10; i++ { // 最多尝试10次
		// 生成随机字符串
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			return "", err
		}
		
		// 生成格式化的配置码
		timestamp := time.Now().Format("20060102")
		code := fmt.Sprintf("KUST-%s-%s", timestamp, hex.EncodeToString(bytes)[:16])

		// 检查是否已存在
		var count int64
		if err := global.DB.Model(&model.ConfigCode{}).Where("code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}

	return "", fmt.Errorf("failed to generate unique code after 10 attempts")
}
