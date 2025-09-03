package model

import "time"

type User struct {
	IdModel
	Username string `json:"username" gorm:"default:'';not null;uniqueIndex"`
	Email    string `json:"email" gorm:"default:'';not null;index"`
	// Email	string     	`json:"email" `
	Password string     `json:"-" gorm:"default:'';not null;"`
	Nickname string     `json:"nickname" gorm:"default:'';not null;"`
	Avatar   string     `json:"avatar" gorm:"default:'';not null;"`
	GroupId  uint       `json:"group_id" gorm:"default:0;not null;index"`
	IsAdmin  *bool      `json:"is_admin" gorm:"default:0;not null;"`
	Status   StatusCode `json:"status" gorm:"default:1;not null;"`
	Remark   string     `json:"remark" gorm:"default:'';not null;"`
	
	// 新增字段：账户生效时间段
	AccountStartTime *time.Time `json:"account_start_time" gorm:"default:null"` // 账户生效开始时间
	AccountEndTime   *time.Time `json:"account_end_time" gorm:"default:null"`   // 账户生效结束时间
	
	// 新增字段：个人设备数量限制
	MaxDevices *int `json:"max_devices" gorm:"default:null"` // 个人最大设备数量，null表示使用全局配置
	
	TimeModel
}

// BeforeSave 钩子用于确保 email 字段有合理的默认值
//func (u *User) BeforeSave(tx *gorm.DB) (err error) {
//	// 如果 email 为空，设置为默认值
//	if u.Email == "" {
//		u.Email = fmt.Sprintf("%s@example.com", u.Username)
//	}
//	return nil
//}

type UserList struct {
	Users []*User `json:"list,omitempty"`
	Pagination
}

var UserRouteNames = []string{
	"MyTagList", "MyAddressBookList", "MyInfo", "MyAddressBookCollection", "MyPeer", "MyShareRecordList", "MyLoginLog",
}
var AdminRouteNames = []string{"*"}
