/**
 * @Author: wanglin
 * @Author: wanglin@vspn.com
 * @Date: 2021/10/28 17:22
 * @Desc: 最近联系人（会话）管理相关枚举定义
 */

package recentcontact

// SessionType 会话类型
type SessionType int

const (
	SessionTypeC2C SessionType = 1 // C2C 会话
	SessionTypeG2C SessionType = 2 // G2C 会话
)
