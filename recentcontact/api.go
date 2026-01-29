/**
 * @Author: wanglin
 * @Author: wanglin@vspn.com
 * @Date: 2021/10/28 16:05
 * @Desc: 最近联系人（会话）管理 API 实现
 */

package recentcontact

import (
	"github.com/d60-Lab/tencent-im/internal/core"
	"github.com/d60-Lab/tencent-im/internal/types"
)

const (
	service              = "recentcontact"
	commandFetchSessions = "get_list"
	commandDeleteSession = "delete"

	// 新增接口命令（2022年新增）
	commandCreateContactGroup = "create_contact_group" // 创建会话分组数据
	commandDelContactGroup    = "del_contact_group"    // 删除会话分组数据
	commandUpdateContactGroup = "update_contact_group" // 更新会话分组数据
	commandSearchContactGroup = "search_contact_group" // 搜索会话分组标记数据
	commandMarkContact        = "mark_contact"         // 创建或更新会话标记数据
	commandGetContactGroup    = "get_contact_group"    // 拉取会话分组标记数据
)

type API interface {
	// FetchSessions 拉取会话列表
	// 支持分页拉取会话列表
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/62118
	FetchSessions(arg *FetchSessionsArg) (ret *FetchSessionsRet, err error)

	// PullSessions 续拉取会话列表
	// 本API是借助"拉取会话列表"API进行扩展实现
	// 支持分页拉取会话列表
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/62118
	PullSessions(arg *PullSessionsArg, fn func(ret *FetchSessionsRet)) (err error)

	// DeleteSession 删除单个会话
	// 删除指定会话，支持同步清理漫游消息。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/62119
	DeleteSession(fromUserId, toUserId string, SessionType SessionType, isClearRamble ...bool) (err error)

	// ========== 新增接口（2022年） ==========

	// CreateContactGroup 创建会话分组数据
	// App 管理员可以创建会话分组。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/81914
	CreateContactGroup(userId string, groupName string, contactIds []string) (err error)

	// DeleteContactGroup 删除会话分组数据
	// App 管理员可以删除会话分组。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/81915
	DeleteContactGroup(userId string, groupNames []string) (err error)

	// UpdateContactGroup 更新会话分组数据
	// App 管理员可以更新会话分组。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/81916
	UpdateContactGroup(userId string, groupName string, newName string, contactIds []string) (err error)

	// SearchContactGroup 搜索会话分组标记数据
	// App 管理员可以搜索会话分组标记数据。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/81917
	SearchContactGroup(userId string, groupName string, markType int) (contacts []Contact, err error)

	// MarkContact 创建或更新会话标记数据
	// App 管理员可以创建或更新会话标记。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/81918
	MarkContact(userId string, contactId string, markType int) (err error)

	// GetContactGroup 拉取会话分组标记数据
	// App 管理员可以拉取会话分组标记数据。
	// 点击查看详细文档:
	// https://cloud.tencent.com/document/product/269/81919
	GetContactGroup(userId string) (groups []ContactGroup, err error)
}

type api struct {
	client core.Client
}

func NewAPI(client core.Client) API {
	return &api{client: client}
}

// FetchSessions 拉取会话列表
// 支持分页拉取会话列表
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/62118
func (a *api) FetchSessions(arg *FetchSessionsArg) (ret *FetchSessionsRet, err error) {
	req := &fetchSessionsReq{
		UserId:        arg.UserId,
		TimeStamp:     arg.TimeStamp,
		StartIndex:    arg.StartIndex,
		TopTimeStamp:  arg.TopTimeStamp,
		TopStartIndex: arg.TopStartIndex,
	}

	if arg.IsAllowTopSession {
		req.AssistFlags += 1 << 0
	}

	if arg.IsReturnEmptySession {
		req.AssistFlags += 1 << 1
	}

	if arg.IsAllowTopSessionPaging {
		req.AssistFlags += 1 << 2
	}

	resp := &fetchSessionsResp{}

	if err = a.client.Post(service, commandFetchSessions, req, resp); err != nil {
		return
	}

	ret = &FetchSessionsRet{
		TimeStamp:     resp.TimeStamp,
		StartIndex:    resp.StartIndex,
		TopTimeStamp:  resp.TopTimeStamp,
		TopStartIndex: resp.TopStartIndex,
		List:          resp.Sessions,
		HasMore:       resp.CompleteFlag == 0,
	}

	return
}

// PullSessions 续拉取会话列表
// 本API是借助"拉取会话列表"API进行扩展实现
// 支持分页拉取会话列表
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/62118
func (a *api) PullSessions(arg *PullSessionsArg, fn func(ret *FetchSessionsRet)) (err error) {
	var (
		ret *FetchSessionsRet
		req = &FetchSessionsArg{
			UserId:                  arg.UserId,
			IsAllowTopSession:       arg.IsAllowTopSession,
			IsReturnEmptySession:    arg.IsReturnEmptySession,
			IsAllowTopSessionPaging: arg.IsAllowTopSessionPaging,
		}
	)

	for ret == nil || ret.HasMore {
		ret, err = a.FetchSessions(req)
		if err != nil {
			return
		}

		fn(ret)

		if ret.HasMore {
			req.TimeStamp = ret.TimeStamp
			req.StartIndex = ret.StartIndex
			req.TopTimeStamp = ret.TopTimeStamp
			req.TopStartIndex = ret.TopStartIndex
		}
	}

	return
}

// DeleteSession 删除单个会话
// 删除指定会话，支持同步清理漫游消息。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/62119
func (a *api) DeleteSession(fromUserId, toUserId string, SessionType SessionType, isClearRamble ...bool) (err error) {
	req := &deleteSessionReq{
		FromUserId: fromUserId,
		ToUserId:   toUserId,
		Type:       SessionType,
	}

	if len(isClearRamble) > 0 && isClearRamble[0] {
		req.ClearRamble = 1
	}

	if err = a.client.Post(service, commandDeleteSession, req, &types.ActionBaseResp{}); err != nil {
		return
	}

	return
}

// ========== 新增接口实现（2022年） ==========

// CreateContactGroup 创建会话分组数据
// App 管理员可以创建会话分组。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/81914
func (a *api) CreateContactGroup(userId string, groupName string, contactIds []string) (err error) {
	req := &createContactGroupReq{
		UserId:     userId,
		GroupName:  groupName,
		ContactIds: contactIds,
	}
	resp := &types.ActionBaseResp{}

	err = a.client.Post(service, commandCreateContactGroup, req, resp)
	return
}

// DeleteContactGroup 删除会话分组数据
// App 管理员可以删除会话分组。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/81915
func (a *api) DeleteContactGroup(userId string, groupNames []string) (err error) {
	req := &deleteContactGroupReq{
		UserId:     userId,
		GroupNames: groupNames,
	}
	resp := &types.ActionBaseResp{}

	err = a.client.Post(service, commandDelContactGroup, req, resp)
	return
}

// UpdateContactGroup 更新会话分组数据
// App 管理员可以更新会话分组。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/81916
func (a *api) UpdateContactGroup(userId string, groupName string, newName string, contactIds []string) (err error) {
	req := &updateContactGroupReq{
		UserId:     userId,
		GroupName:  groupName,
		NewName:    newName,
		ContactIds: contactIds,
	}
	resp := &types.ActionBaseResp{}

	err = a.client.Post(service, commandUpdateContactGroup, req, resp)
	return
}

// SearchContactGroup 搜索会话分组标记数据
// App 管理员可以搜索会话分组标记数据。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/81917
func (a *api) SearchContactGroup(userId string, groupName string, markType int) (contacts []Contact, err error) {
	req := &searchContactGroupReq{
		UserId:    userId,
		GroupName: groupName,
		MarkType:  markType,
	}
	resp := &searchContactGroupResp{}

	if err = a.client.Post(service, commandSearchContactGroup, req, resp); err != nil {
		return
	}

	contacts = resp.Contacts
	return
}

// MarkContact 创建或更新会话标记数据
// App 管理员可以创建或更新会话标记。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/81918
func (a *api) MarkContact(userId string, contactId string, markType int) (err error) {
	req := &markContactReq{
		UserId:    userId,
		ContactId: contactId,
		MarkType:  markType,
	}
	resp := &types.ActionBaseResp{}

	err = a.client.Post(service, commandMarkContact, req, resp)
	return
}

// GetContactGroup 拉取会话分组标记数据
// App 管理员可以拉取会话分组标记数据。
// 点击查看详细文档:
// https://cloud.tencent.com/document/product/269/81919
func (a *api) GetContactGroup(userId string) (groups []ContactGroup, err error) {
	req := &getContactGroupReq{
		UserId: userId,
	}
	resp := &getContactGroupResp{}

	if err = a.client.Post(service, commandGetContactGroup, req, resp); err != nil {
		return
	}

	groups = resp.Groups
	return
}
