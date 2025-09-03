package model

type UserToken struct {
	IdModel
	UserId     uint   `json:"user_id" gorm:"default:0;not null;index"`
	DeviceUuid string `json:"device_uuid" gorm:"default:'';omitempty;"`
	DeviceId   string `json:"device_id" gorm:"default:'';omitempty;"`
	Token      string `json:"token" gorm:"default:'';not null;index"`
	ExpiredAt  int64  `json:"expired_at" gorm:"default:0;not null;"`
	
	// 新增字段：设备详细信息
	DeviceName    string `json:"device_name" gorm:"default:''"`       // 设备名称
	DeviceType    string `json:"device_type" gorm:"default:''"`       // 设备类型 (web, mobile, desktop)
	DeviceOS      string `json:"device_os" gorm:"default:''"`         // 操作系统
	DeviceIP      string `json:"device_ip" gorm:"default:''"`         // 设备IP地址
	LastActiveAt  int64  `json:"last_active_at" gorm:"default:0"`    // 最后活跃时间
	
	TimeModel
}

type UserTokenList struct {
	UserTokens []UserToken `json:"list"`
	Pagination
}
