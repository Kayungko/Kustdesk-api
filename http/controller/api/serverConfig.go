package api

import (
	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type ServerConfig struct{}

// GetConfigByCode 根据配置码获取服务器配置
// @Tags 客户端配置
// @Summary 根据配置码获取服务器配置
// @Description 客户端使用配置码获取服务器配置信息
// @Accept json
// @Produce json
// @Param code path string true "配置码"
// @Success 200 {object} response.Response{data=model.ClientServerConfig}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/config/{code} [get]
func (ct *ServerConfig) GetConfigByCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Fail(c, 400, "Config code is required")
		return
	}

	// 获取客户端信息
	clientIP := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// 获取配置
	config, err := service.AllService.ServerConfigService.GetConfigByCode(code)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	// 记录使用情况
	go func() {
		service.AllService.ServerConfigService.RecordConfigCodeUsage(code, clientIP, userAgent)
	}()

	response.Success(c, config)
}

// ValidateConfigCode 验证配置码
// @Tags 客户端配置
// @Summary 验证配置码
// @Description 验证配置码是否有效，不返回具体配置信息
// @Accept json
// @Produce json
// @Param code path string true "配置码"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/config/{code}/validate [get]
func (ct *ServerConfig) ValidateConfigCode(c *gin.Context) {
	code := c.Param("code")
	if code == "" {
		response.Fail(c, 400, "Config code is required")
		return
	}

	// 只验证配置码，不记录使用
	_, err := service.AllService.ServerConfigService.GetConfigByCode(code)
	if err != nil {
		response.Fail(c, 404, err.Error())
		return
	}

	response.Success(c, gin.H{
		"valid": true,
		"message": "Config code is valid",
	})
}
