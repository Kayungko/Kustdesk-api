package admin

import (
	"github.com/gin-gonic/gin"
	"rustdesk-api/global"
	"rustdesk-api/http/response"
	"rustdesk-api/model"
	"time"
)

type SystemController struct{}

// SystemConfig 系统配置结构体
type SystemConfig struct {
	MaxConcurrentDevices int    `json:"max_concurrent_devices"`
	Register            bool   `json:"register"`
	RegisterStatus      int    `json:"register_status"`
	CaptchaThreshold    int    `json:"captcha_threshold"`
	BanThreshold        int    `json:"ban_threshold"`
	DisablePwdLogin     bool   `json:"disable_pwd_login"`
	WebClient           int    `json:"web_client"`
	WebSso              bool   `json:"web_sso"`
	TokenExpire         string `json:"token_expire"`
}

// SystemStatus 系统状态结构体
type SystemStatus struct {
	ServerTime    time.Time `json:"server_time"`
	Uptime        string    `json:"uptime"`
	Version       string    `json:"version"`
	GoVersion     string    `json:"go_version"`
	DatabaseType  string    `json:"database_type"`
	RedisStatus   bool      `json:"redis_status"`
}

// UserStatistics 用户统计信息
type UserStatistics struct {
	TotalUsers    int64 `json:"total_users"`
	ActiveUsers   int64 `json:"active_users"`
	OnlineUsers   int64 `json:"online_users"`
	DisabledUsers int64 `json:"disabled_users"`
}

// DeviceStatistics 设备统计信息
type DeviceStatistics struct {
	TotalDevices   int64 `json:"total_devices"`
	OnlineDevices  int64 `json:"online_devices"`
	ExpiredDevices int64 `json:"expired_devices"`
	ActiveDevices  int64 `json:"active_devices"`
}

// GetConfig 获取系统配置
func (sc *SystemController) GetConfig(c *gin.Context) {
	cfg := global.Config
	
	systemConfig := SystemConfig{
		MaxConcurrentDevices: cfg.App.MaxConcurrentDevices,
		Register:            cfg.App.Register,
		RegisterStatus:      cfg.App.RegisterStatus,
		CaptchaThreshold:    cfg.App.CaptchaThreshold,
		BanThreshold:        cfg.App.BanThreshold,
		DisablePwdLogin:     cfg.App.DisablePwdLogin,
		WebClient:           cfg.App.WebClient,
		WebSso:              cfg.App.WebSso,
		TokenExpire:         cfg.App.TokenExpire.String(),
	}
	
	response.Success(c, systemConfig)
}

// UpdateConfig 更新系统配置
func (sc *SystemController) UpdateConfig(c *gin.Context) {
	var req SystemConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidateFail(c, err.Error())
		return
	}
	
	// 更新配置（注意：这里只是更新内存中的配置，实际项目中可能需要写入配置文件）
	cfg := &global.Config
	cfg.App.MaxConcurrentDevices = req.MaxConcurrentDevices
	cfg.App.Register = req.Register
	cfg.App.RegisterStatus = req.RegisterStatus
	cfg.App.CaptchaThreshold = req.CaptchaThreshold
	cfg.App.BanThreshold = req.BanThreshold
	cfg.App.DisablePwdLogin = req.DisablePwdLogin
	cfg.App.WebClient = req.WebClient
	cfg.App.WebSso = req.WebSso
	
	// 解析TokenExpire
	if duration, err := time.ParseDuration(req.TokenExpire); err == nil {
		cfg.App.TokenExpire = duration
	}
	
	response.Success(c, "配置更新成功")
}

// GetStatus 获取系统状态
func (sc *SystemController) GetStatus(c *gin.Context) {
	status := SystemStatus{
		ServerTime:   time.Now(),
		Uptime:       "运行中", // 这里可以计算实际运行时间
		Version:      "1.0.0",  // 从版本文件或构建信息获取
		GoVersion:    "1.21",   // 从runtime获取
		DatabaseType: global.Config.Gorm.Type,
		RedisStatus:  true,     // 检查Redis连接状态
	}
	
	response.Success(c, status)
}

// GetUserStatistics 获取用户统计信息
func (sc *SystemController) GetUserStatistics(c *gin.Context) {
	var stats UserStatistics
	
	// 获取总用户数
	global.DB.Model(&model.User{}).Count(&stats.TotalUsers)
	
	// 获取活跃用户数（状态为1的用户）
	global.DB.Model(&model.User{}).Where("status = ?", 1).Count(&stats.ActiveUsers)
	
	// 获取在线用户数（有活跃Token的用户）
	now := time.Now().Unix()
	global.DB.Model(&model.UserToken{}).Where("expired_at > ?", now).
		Distinct("user_id").Count(&stats.OnlineUsers)
	
	// 获取禁用用户数
	global.DB.Model(&model.User{}).Where("status = ?", 2).Count(&stats.DisabledUsers)
	
	response.Success(c, stats)
}

// GetDeviceStatistics 获取设备统计信息
func (sc *SystemController) GetDeviceStatistics(c *gin.Context) {
	var stats DeviceStatistics
	
	now := time.Now().Unix()
	
	// 获取总设备数（Token数）
	global.DB.Model(&model.UserToken{}).Count(&stats.TotalDevices)
	
	// 获取在线设备数（未过期的Token）
	global.DB.Model(&model.UserToken{}).Where("expired_at > ?", now).Count(&stats.OnlineDevices)
	
	// 获取过期设备数
	global.DB.Model(&model.UserToken{}).Where("expired_at <= ?", now).Count(&stats.ExpiredDevices)
	
	// 获取活跃设备数（最近5分钟内有活动的设备）
	fiveMinutesAgo := now - 300
	global.DB.Model(&model.UserToken{}).
		Where("expired_at > ? AND last_active_at > ?", now, fiveMinutesAgo).
		Count(&stats.ActiveDevices)
	
	response.Success(c, stats)
}
