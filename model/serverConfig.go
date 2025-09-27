package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"time"
)

// ServerConfig 服务器配置模型
type ServerConfig struct {
	IdModel
	Name        string     `json:"name" gorm:"not null;comment:配置名称"`
	Description string     `json:"description" gorm:"comment:配置描述"`
	Region      string     `json:"region" gorm:"comment:地域标识"`
	IdServer    string     `json:"id_server" gorm:"not null;comment:ID服务器地址"`
	RelayServer string     `json:"relay_server" gorm:"comment:中继服务器地址"`
	ApiServer   string     `json:"api_server" gorm:"comment:API服务器地址"`
	Key         string     `json:"key" gorm:"comment:服务器密钥"`
	IsEnabled   *bool      `json:"is_enabled" gorm:"default:false;not null;comment:是否启用"`
	IsDefault   *bool      `json:"is_default" gorm:"default:false;not null;comment:是否为默认配置"`
	Priority    int        `json:"priority" gorm:"default:0;comment:优先级，数字越大优先级越高"`
	Status      StatusCode `json:"status" gorm:"default:1;not null;comment:状态"`
	TimeModel
}

// ConfigCode 配置码模型
type ConfigCode struct {
	IdModel
	Code           string     `json:"code" gorm:"uniqueIndex;not null;comment:配置码"`
	ServerConfigId uint       `json:"server_config_id" gorm:"not null;index;comment:关联的服务器配置ID"`
	ServerConfig   *ServerConfig `json:"server_config,omitempty" gorm:"foreignKey:ServerConfigId"`
	ExpiresAt      *time.Time `json:"expires_at" gorm:"comment:过期时间，null表示永不过期"`
	UsageCount     int        `json:"usage_count" gorm:"default:0;comment:使用次数"`
	MaxUsage       *int       `json:"max_usage" gorm:"comment:最大使用次数，null表示无限制"`
	CreatedBy      uint       `json:"created_by" gorm:"not null;index;comment:创建者用户ID"`
	Creator        *User      `json:"creator,omitempty" gorm:"foreignKey:CreatedBy"`
	Status         StatusCode `json:"status" gorm:"default:1;not null;comment:状态"`
	TimeModel
}

// ConfigCodeUsage 配置码使用记录
type ConfigCodeUsage struct {
	IdModel
	ConfigCodeId uint       `json:"config_code_id" gorm:"not null;index;comment:配置码ID"`
	ConfigCode   *ConfigCode `json:"config_code,omitempty" gorm:"foreignKey:ConfigCodeId"`
	ClientIP     string     `json:"client_ip" gorm:"comment:客户端IP"`
	UserAgent    string     `json:"user_agent" gorm:"comment:用户代理"`
	UsedAt       time.Time  `json:"used_at" gorm:"not null;comment:使用时间"`
	TimeModel
}

// ServerConfigList 服务器配置列表
type ServerConfigList struct {
	ServerConfigs []*ServerConfig `json:"list"`
	Pagination
}

// ConfigCodeList 配置码列表  
type ConfigCodeList struct {
	ConfigCodes []*ConfigCode `json:"list"`
	Pagination
}

// ClientServerConfig 客户端服务器配置（用于API返回）
type ClientServerConfig struct {
	Name        string `json:"name"`
	Region      string `json:"region"`
	IdServer    string `json:"id_server"`
	RelayServer string `json:"relay_server"`
	ApiServer   string `json:"api_server"`
	Key         string `json:"key"`
}

// EncryptedConfigData 加密的配置数据结构
type EncryptedConfigData struct {
	Config    ClientServerConfig `json:"config"`
	Timestamp time.Time          `json:"timestamp"`
	Version   string             `json:"version"`
}

// GenerateConfigCode 生成配置码
func (sc *ServerConfig) GenerateConfigCode(secretKey string) (string, error) {
	// 构建客户端配置数据
	configData := EncryptedConfigData{
		Config: ClientServerConfig{
			Name:        sc.Name,
			Region:      sc.Region,
			IdServer:    sc.IdServer,
			RelayServer: sc.RelayServer,
			ApiServer:   sc.ApiServer,
			Key:         sc.Key,
		},
		Timestamp: time.Now(),
		Version:   "1.0",
	}

	// 序列化为JSON
	jsonData, err := json.Marshal(configData)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config data: %v", err)
	}

	// 加密数据
	encryptedData, err := encryptAES(jsonData, secretKey)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt config data: %v", err)
	}

	// 生成配置码格式: KUST-{server_id}-{timestamp}-{encrypted_data}
	timestamp := time.Now().Format("20060102")
	configCode := fmt.Sprintf("KUST-%d-%s-%s", sc.Id, timestamp, base64.URLEncoding.EncodeToString(encryptedData))

	return configCode, nil
}

// DecryptConfigCode 解密配置码
func DecryptConfigCode(code, secretKey string) (*ClientServerConfig, error) {
	// 解析配置码格式
	// 提取加密数据部分（跳过前缀）
	var encryptedPart string
	
	// 简单解析，实际可能需要更复杂的解析逻辑
	if len(code) > 20 && code[:5] == "KUST-" {
		// 找到最后一个'-'后的部分作为加密数据
		lastDash := -1
		dashCount := 0
		for i, c := range code {
			if c == '-' {
				dashCount++
				if dashCount == 3 {
					lastDash = i
					break
				}
			}
		}
		if lastDash > 0 && lastDash < len(code)-1 {
			encryptedPart = code[lastDash+1:]
		}
	}

	if encryptedPart == "" {
		return nil, fmt.Errorf("invalid config code format")
	}

	// Base64解码
	encryptedData, err := base64.URLEncoding.DecodeString(encryptedPart)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config code: %v", err)
	}

	// AES解密
	decryptedData, err := decryptAES(encryptedData, secretKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt config data: %v", err)
	}

	// 反序列化JSON
	var configData EncryptedConfigData
	if err := json.Unmarshal(decryptedData, &configData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config data: %v", err)
	}

	// 验证时间戳（可选，避免过期配置）
	if time.Since(configData.Timestamp) > 365*24*time.Hour {
		return nil, fmt.Errorf("config code has expired")
	}

	return &configData.Config, nil
}

// AES加密函数
func encryptAES(data []byte, secretKey string) ([]byte, error) {
	// 使用SHA256生成32字节密钥
	key := sha256.Sum256([]byte(secretKey))
	
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	// 使用GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 生成随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密并附加nonce
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// AES解密函数
func decryptAES(data []byte, secretKey string) ([]byte, error) {
	// 使用SHA256生成32字节密钥
	key := sha256.Sum256([]byte(secretKey))
	
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// TableName 指定表名
func (ServerConfig) TableName() string {
	return "server_configs"
}

func (ConfigCode) TableName() string {
	return "config_codes"
}

func (ConfigCodeUsage) TableName() string {
	return "config_code_usages"
}
