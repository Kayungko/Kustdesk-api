package admin

import "github.com/lejianwen/rustdesk-api/v2/model"

type LoginPayload struct {
	Username   string   `json:"username"`
	Email      string   `json:"email"`
	Avatar     string   `json:"avatar"`
	Token      string   `json:"token"`
	RouteNames []string `json:"route_names"`
	Nickname   string   `json:"nickname"`
}

func (lp *LoginPayload) FromUser(user *model.User) {
	lp.Username = user.Username
	lp.Email = user.Email
	lp.Avatar = user.Avatar
	lp.Nickname = user.Nickname
}

type UserOauthItem struct {
	Op     string `json:"op"`
	Status int    `json:"status"`
}

// GroupUsersResponse 用户分组响应
type GroupUsersResponse struct {
	Groups []*model.Group `json:"groups"`
	Users  []*model.User  `json:"users"`
}

// DeviceLimitInfo 设备限制信息
type DeviceLimitInfo struct {
	Limit        int  `json:"limit"`
	IsPersonal   bool `json:"is_personal"`
	CurrentCount int  `json:"current_count"`
	AvailableSlots int `json:"available_slots"`
}

// UserDevicesResponse 用户设备响应
type UserDevicesResponse struct {
	Devices     []*model.UserToken `json:"devices"`
	DeviceLimit DeviceLimitInfo    `json:"device_limit"`
}

// BatchDisableExpiredResponse 批量禁用过期账户响应
type BatchDisableExpiredResponse struct {
	DisabledCount int    `json:"disabled_count"`
	Message       string `json:"message"`
}
