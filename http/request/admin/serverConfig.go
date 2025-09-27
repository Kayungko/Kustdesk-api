package admin

import (
	"time"
)

// ServerConfigForm 服务器配置表单
type ServerConfigForm struct {
	Name        string `json:"name" binding:"required,min=1,max=100" label:"配置名称"`
	Description string `json:"description" binding:"max=500" label:"配置描述"`
	Region      string `json:"region" binding:"max=50" label:"地域标识"`
	IdServer    string `json:"id_server" binding:"required,min=1,max=255" label:"ID服务器地址"`
	RelayServer string `json:"relay_server" binding:"max=255" label:"中继服务器地址"`
	ApiServer   string `json:"api_server" binding:"max=255" label:"API服务器地址"`
	Key         string `json:"key" binding:"max=500" label:"服务器密钥"`
	IsEnabled   *bool  `json:"is_enabled" label:"是否启用"`
	IsDefault   *bool  `json:"is_default" label:"是否为默认配置"`
	Priority    int    `json:"priority" label:"优先级"`
}

// ServerConfigListQuery 服务器配置列表查询
type ServerConfigListQuery struct {
	Pagination
	Name      string `form:"name" label:"配置名称"`
	Region    string `form:"region" label:"地域"`
	IsEnabled *bool  `form:"is_enabled" label:"启用状态"`
	IsDefault *bool  `form:"is_default" label:"默认配置"`
}

// ConfigCodeForm 配置码生成表单
type ConfigCodeForm struct {
	ServerConfigId uint       `json:"server_config_id" binding:"required,min=1" label:"服务器配置ID"`
	ExpiresAt      *time.Time `json:"expires_at" label:"过期时间"`
	MaxUsage       *int       `json:"max_usage" binding:"omitempty,min=1" label:"最大使用次数"`
}

// ConfigCodeListQuery 配置码列表查询
type ConfigCodeListQuery struct {
	Pagination
	ServerConfigId uint   `form:"server_config_id" label:"服务器配置ID"`
	Code           string `form:"code" label:"配置码"`
	CreatedBy      uint   `form:"created_by" label:"创建者"`
}

// ConfigCodeBatchForm 批量生成配置码表单
type ConfigCodeBatchForm struct {
	ServerConfigId uint       `json:"server_config_id" binding:"required,min=1" label:"服务器配置ID"`
	Count          int        `json:"count" binding:"required,min=1,max=100" label:"生成数量"`
	ExpiresAt      *time.Time `json:"expires_at" label:"过期时间"`
	MaxUsage       *int       `json:"max_usage" binding:"omitempty,min=1" label:"最大使用次数"`
}

// SetDefaultConfigForm 设置默认配置表单
type SetDefaultConfigForm struct {
	ServerConfigId uint `json:"server_config_id" binding:"required,min=1" label:"服务器配置ID"`
}
