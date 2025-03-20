/**
 * Description：
 * FileName：token.go
 * Author：CJiaの用心
 * Create：2025/3/20 23:16:11
 * Remark：
 */

package config

type TokenConfig struct {
	ApiKeyAuth string `yaml:"apiKeyAuth" json:"apiKeyAuth"`
	Expire     int    `yaml:"expire" json:"expire"`
}
