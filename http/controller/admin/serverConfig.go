package admin

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lejianwen/rustdesk-api/v2/global"
	"github.com/lejianwen/rustdesk-api/v2/http/request/admin"
	"github.com/lejianwen/rustdesk-api/v2/http/response"
	adResp "github.com/lejianwen/rustdesk-api/v2/http/response/admin"
	"github.com/lejianwen/rustdesk-api/v2/service"
)

type ServerConfig struct{}

// List 服务器配置列表
// @Tags 服务器配置
// @Summary 服务器配置列表
// @Description 获取服务器配置列表
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param name query string false "配置名称"
// @Param region query string false "地域"
// @Param is_enabled query bool false "启用状态"
// @Param is_default query bool false "默认配置"
// @Success 200 {object} response.Response{data=admin.ServerConfigListResponse}
// @Failure 500 {object} response.Response
// @Router /admin/server-config/list [get]
// @Security token
func (ct *ServerConfig) List(c *gin.Context) {
	var query admin.ServerConfigListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	list, err := service.AllService.ServerConfigService.List(&query)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	// 转换为响应格式
	var respList []*adResp.ServerConfigResponse
	for _, config := range list.ServerConfigs {
		respList = append(respList, &adResp.ServerConfigResponse{
			Id:          config.Id,
			Name:        config.Name,
			Description: config.Description,
			Region:      config.Region,
			IdServer:    config.IdServer,
			RelayServer: config.RelayServer,
			ApiServer:   config.ApiServer,
			Key:         config.Key,
			IsEnabled:   config.IsEnabled != nil && *config.IsEnabled,
			IsDefault:   config.IsDefault != nil && *config.IsDefault,
			Priority:    config.Priority,
			Status:      int(config.Status),
			CreatedAt:   config.CreatedAt.Time,
			UpdatedAt:   config.UpdatedAt.Time,
		})
	}

	resp := adResp.ServerConfigListResponse{
		List:       respList,
		Pagination: list.Pagination,
	}

	response.Success(c, resp)
}

// Create 创建服务器配置
// @Tags 服务器配置
// @Summary 创建服务器配置
// @Description 创建新的服务器配置
// @Accept json
// @Produce json
// @Param body body admin.ServerConfigForm true "服务器配置信息"
// @Success 200 {object} response.Response{data=admin.ServerConfigResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/server-config/create [post]
// @Security token
func (ct *ServerConfig) Create(c *gin.Context) {
	var form admin.ServerConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	config, err := service.AllService.ServerConfigService.Create(&form)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	resp := &adResp.ServerConfigResponse{
		Id:          config.Id,
		Name:        config.Name,
		Description: config.Description,
		Region:      config.Region,
		IdServer:    config.IdServer,
		RelayServer: config.RelayServer,
		ApiServer:   config.ApiServer,
		Key:         config.Key,
		IsEnabled:   config.IsEnabled != nil && *config.IsEnabled,
		IsDefault:   config.IsDefault != nil && *config.IsDefault,
		Priority:    config.Priority,
		Status:      int(config.Status),
		CreatedAt:   config.CreatedAt.Time,
		UpdatedAt:   config.UpdatedAt.Time,
	}

	response.Success(c, resp)
}

// Update 更新服务器配置
// @Tags 服务器配置
// @Summary 更新服务器配置
// @Description 更新指定的服务器配置
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Param body body admin.ServerConfigForm true "服务器配置信息"
// @Success 200 {object} response.Response{data=admin.ServerConfigResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/server-config/update/{id} [put]
// @Security token
func (ct *ServerConfig) Update(c *gin.Context) {
	id := c.Param("id")
	iid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Fail(c, 400, "Invalid ID")
		return
	}

	var form admin.ServerConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	config, err := service.AllService.ServerConfigService.Update(uint(iid), &form)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	resp := &adResp.ServerConfigResponse{
		Id:          config.Id,
		Name:        config.Name,
		Description: config.Description,
		Region:      config.Region,
		IdServer:    config.IdServer,
		RelayServer: config.RelayServer,
		ApiServer:   config.ApiServer,
		Key:         config.Key,
		IsEnabled:   config.IsEnabled != nil && *config.IsEnabled,
		IsDefault:   config.IsDefault != nil && *config.IsDefault,
		Priority:    config.Priority,
		Status:      int(config.Status),
		CreatedAt:   config.CreatedAt.Time,
		UpdatedAt:   config.UpdatedAt.Time,
	}

	response.Success(c, resp)
}

// Delete 删除服务器配置
// @Tags 服务器配置
// @Summary 删除服务器配置
// @Description 删除指定的服务器配置
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/server-config/delete/{id} [delete]
// @Security token
func (ct *ServerConfig) Delete(c *gin.Context) {
	id := c.Param("id")
	iid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Fail(c, 400, "Invalid ID")
		return
	}

	err = service.AllService.ServerConfigService.Delete(uint(iid))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// Detail 服务器配置详情
// @Tags 服务器配置
// @Summary 服务器配置详情
// @Description 获取指定服务器配置的详细信息
// @Accept json
// @Produce json
// @Param id path int true "配置ID"
// @Success 200 {object} response.Response{data=admin.ServerConfigResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /admin/server-config/detail/{id} [get]
// @Security token
func (ct *ServerConfig) Detail(c *gin.Context) {
	id := c.Param("id")
	iid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Fail(c, 400, "Invalid ID")
		return
	}

	config, err := service.AllService.ServerConfigService.Detail(uint(iid))
	if err != nil {
		response.Fail(c, 404, "Server config not found")
		return
	}

	resp := &adResp.ServerConfigResponse{
		Id:          config.Id,
		Name:        config.Name,
		Description: config.Description,
		Region:      config.Region,
		IdServer:    config.IdServer,
		RelayServer: config.RelayServer,
		ApiServer:   config.ApiServer,
		Key:         config.Key,
		IsEnabled:   config.IsEnabled != nil && *config.IsEnabled,
		IsDefault:   config.IsDefault != nil && *config.IsDefault,
		Priority:    config.Priority,
		Status:      int(config.Status),
		CreatedAt:   config.CreatedAt.Time,
		UpdatedAt:   config.UpdatedAt.Time,
	}

	response.Success(c, resp)
}

// SetDefault 设置默认配置
// @Tags 服务器配置
// @Summary 设置默认配置
// @Description 设置指定的服务器配置为默认配置
// @Accept json
// @Produce json
// @Param body body admin.SetDefaultConfigForm true "配置ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/server-config/set-default [post]
// @Security token
func (ct *ServerConfig) SetDefault(c *gin.Context) {
	var form admin.SetDefaultConfigForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	err := service.AllService.ServerConfigService.SetDefault(form.ServerConfigId)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GenerateConfigCode 生成配置码
// @Tags 配置码
// @Summary 生成配置码
// @Description 为指定服务器配置生成配置码
// @Accept json
// @Produce json
// @Param body body admin.ConfigCodeForm true "配置码生成参数"
// @Success 200 {object} response.Response{data=admin.ConfigCodeGenerateResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config-code/generate [post]
// @Security token
func (ct *ServerConfig) GenerateConfigCode(c *gin.Context) {
	var form admin.ConfigCodeForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 获取当前用户ID（假设从JWT中获取）
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, 401, "Unauthorized")
		return
	}
	userId := userIdInterface.(uint)

	configCode, err := service.AllService.ServerConfigService.GenerateConfigCode(&form, userId)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	resp := &adResp.ConfigCodeGenerateResponse{
		Id:        configCode.Id,
		Code:      configCode.Code,
		ExpiresAt: configCode.ExpiresAt,
		MaxUsage:  configCode.MaxUsage,
	}

	response.Success(c, resp)
}

// BatchGenerateConfigCode 批量生成配置码
// @Tags 配置码
// @Summary 批量生成配置码
// @Description 批量生成多个配置码
// @Accept json
// @Produce json
// @Param body body admin.ConfigCodeBatchForm true "批量生成参数"
// @Success 200 {object} response.Response{data=admin.ConfigCodeBatchGenerateResponse}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config-code/batch-generate [post]
// @Security token
func (ct *ServerConfig) BatchGenerateConfigCode(c *gin.Context) {
	var form admin.ConfigCodeBatchForm
	if err := c.ShouldBindJSON(&form); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	// 获取当前用户ID
	userIdInterface, exists := c.Get("user_id")
	if !exists {
		response.Fail(c, 401, "Unauthorized")
		return
	}
	userId := userIdInterface.(uint)

	codes, err := service.AllService.ServerConfigService.BatchGenerateConfigCode(&form, userId)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	var respCodes []*adResp.ConfigCodeGenerateResponse
	for _, code := range codes {
		respCodes = append(respCodes, &adResp.ConfigCodeGenerateResponse{
			Id:        code.Id,
			Code:      code.Code,
			ExpiresAt: code.ExpiresAt,
			MaxUsage:  code.MaxUsage,
		})
	}

	resp := &adResp.ConfigCodeBatchGenerateResponse{
		Count:       len(codes),
		ConfigCodes: respCodes,
	}

	response.Success(c, resp)
}

// ConfigCodeList 配置码列表
// @Tags 配置码
// @Summary 配置码列表
// @Description 获取配置码列表
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Param server_config_id query int false "服务器配置ID"
// @Param code query string false "配置码"
// @Param created_by query int false "创建者ID"
// @Success 200 {object} response.Response{data=admin.ConfigCodeListResponse}
// @Failure 500 {object} response.Response
// @Router /admin/config-code/list [get]
// @Security token
func (ct *ServerConfig) ConfigCodeList(c *gin.Context) {
	var query admin.ConfigCodeListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, 400, err.Error())
		return
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 20
	}

	list, err := service.AllService.ServerConfigService.ConfigCodeList(&query)
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	// 转换为响应格式
	var respList []*adResp.ConfigCodeResponse
	for _, code := range list.ConfigCodes {
		respCode := &adResp.ConfigCodeResponse{
			Id:             code.Id,
			Code:           code.Code,
			ServerConfigId: code.ServerConfigId,
			ExpiresAt:      code.ExpiresAt,
			UsageCount:     code.UsageCount,
			MaxUsage:       code.MaxUsage,
			CreatedBy:      code.CreatedBy,
			Status:         int(code.Status),
			CreatedAt:      code.CreatedAt.Time,
			UpdatedAt:      code.UpdatedAt.Time,
		}

		// 添加服务器配置信息
		if code.ServerConfig != nil {
			respCode.ServerConfig = &adResp.ServerConfigResponse{
				Id:          code.ServerConfig.Id,
				Name:        code.ServerConfig.Name,
				Description: code.ServerConfig.Description,
				Region:      code.ServerConfig.Region,
			}
		}

		// 添加创建者信息
		if code.Creator != nil {
			respCode.Creator = &adResp.UserResponse{
				Id:       code.Creator.Id,
				Username: code.Creator.Username,
				Nickname: code.Creator.Nickname,
			}
		}

		respList = append(respList, respCode)
	}

	resp := adResp.ConfigCodeListResponse{
		List:       respList,
		Pagination: list.Pagination,
	}

	response.Success(c, resp)
}

// DeleteConfigCode 删除配置码
// @Tags 配置码
// @Summary 删除配置码
// @Description 删除指定的配置码
// @Accept json
// @Produce json
// @Param id path int true "配置码ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /admin/config-code/delete/{id} [delete]
// @Security token
func (ct *ServerConfig) DeleteConfigCode(c *gin.Context) {
	id := c.Param("id")
	iid, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		response.Fail(c, 400, "Invalid ID")
		return
	}

	err = service.AllService.ServerConfigService.DeleteConfigCode(uint(iid))
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.Success(c, nil)
}

// GetConfigCodeStats 获取配置码统计
// @Tags 配置码
// @Summary 配置码统计
// @Description 获取配置码使用统计信息
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=admin.ConfigCodeStatsResponse}
// @Failure 500 {object} response.Response
// @Router /admin/config-code/stats [get]
// @Security token
func (ct *ServerConfig) GetConfigCodeStats(c *gin.Context) {
	stats, err := service.AllService.ServerConfigService.GetConfigCodeStats()
	if err != nil {
		response.Fail(c, 500, err.Error())
		return
	}

	response.Success(c, stats)
}
