package admin

import (
	"time"
	"github.com/lejianwen/rustdesk-api/v2/model"
)

// ServerConfigResponse 服务器配置响应
type ServerConfigResponse struct {
	Id          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Region      string    `json:"region"`
	IdServer    string    `json:"id_server"`
	RelayServer string    `json:"relay_server"`
	ApiServer   string    `json:"api_server"`
	Key         string    `json:"key"`
	IsEnabled   bool      `json:"is_enabled"`
	IsDefault   bool      `json:"is_default"`
	Priority    int       `json:"priority"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ServerConfigListResponse 服务器配置列表响应
type ServerConfigListResponse struct {
	List []*ServerConfigResponse `json:"list"`
	model.Pagination
}

// ConfigCodeResponse 配置码响应
type ConfigCodeResponse struct {
	Id             uint                  `json:"id"`
	Code           string                `json:"code"`
	ServerConfigId uint                  `json:"server_config_id"`
	ServerConfig   *ServerConfigResponse `json:"server_config,omitempty"`
	ExpiresAt      *time.Time            `json:"expires_at"`
	UsageCount     int                   `json:"usage_count"`
	MaxUsage       *int                  `json:"max_usage"`
	CreatedBy      uint                  `json:"created_by"`
	Creator        *UserResponse         `json:"creator,omitempty"`
	Status         int                   `json:"status"`
	CreatedAt      time.Time             `json:"created_at"`
	UpdatedAt      time.Time             `json:"updated_at"`
}

// UserResponse 用户响应（简化版）
type UserResponse struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

// ConfigCodeListResponse 配置码列表响应
type ConfigCodeListResponse struct {
	List []*ConfigCodeResponse `json:"list"`
	model.Pagination
}

// ConfigCodeGenerateResponse 配置码生成响应
type ConfigCodeGenerateResponse struct {
	Code         string    `json:"code"`
	Id           uint      `json:"id"`
	ExpiresAt    *time.Time `json:"expires_at"`
	MaxUsage     *int      `json:"max_usage"`
	DownloadUrl  string    `json:"download_url,omitempty"`
	QrCodeUrl    string    `json:"qr_code_url,omitempty"`
}

// ConfigCodeBatchGenerateResponse 批量生成配置码响应
type ConfigCodeBatchGenerateResponse struct {
	Count        int                           `json:"count"`
	ConfigCodes  []*ConfigCodeGenerateResponse `json:"config_codes"`
	DownloadUrl  string                        `json:"download_url,omitempty"`
}

// ConfigCodeUsageResponse 配置码使用记录响应
type ConfigCodeUsageResponse struct {
	Id           uint                  `json:"id"`
	ConfigCodeId uint                  `json:"config_code_id"`
	ConfigCode   *ConfigCodeResponse   `json:"config_code,omitempty"`
	ClientIP     string                `json:"client_ip"`
	UserAgent    string                `json:"user_agent"`
	UsedAt       time.Time             `json:"used_at"`
	CreatedAt    time.Time             `json:"created_at"`
}

// ConfigCodeStatsResponse 配置码统计响应
type ConfigCodeStatsResponse struct {
	TotalCodes     int64 `json:"total_codes"`
	ActiveCodes    int64 `json:"active_codes"`
	ExpiredCodes   int64 `json:"expired_codes"`
	TotalUsage     int64 `json:"total_usage"`
	TodayUsage     int64 `json:"today_usage"`
	WeekUsage      int64 `json:"week_usage"`
	MonthUsage     int64 `json:"month_usage"`
}

// ServerConfigStatsResponse 服务器配置统计响应  
type ServerConfigStatsResponse struct {
	TotalConfigs   int64 `json:"total_configs"`
	EnabledConfigs int64 `json:"enabled_configs"`
	DefaultConfig  *ServerConfigResponse `json:"default_config,omitempty"`
}
