package xhs

import (
	"encoding/json"
	"fmt"
	"github.com/Esword618/xhs-go/consts"
	"github.com/Esword618/xhs-go/utils"
	"github.com/bitly/go-simplejson"
)

// GetNoteById 获取笔记通过 noteId
func (x *xhsClient) GetNoteById(noteId string) (*simplejson.Json, error) {
	res, err := x.client.R().
		SetBody(fmt.Sprintf(`{"source_note_id":"%s"}`, noteId)).
		Post("/api/sns/web/v1/feed")
	if err != nil {
		return simplejson.New(), err
	}
	return simplejson.NewFromReader(res.Body)
}

// GetUserInfo 获取用户信息通过 userId
func (x *xhsClient) GetUserInfo(userId string) (*simplejson.Json, error) {
	res, err := x.client.R().
		SetQueryParams(map[string]string{
			"target_user_id": userId,
		}).
		Get("/api/sns/web/v1/user/otherinfo")
	if err != nil {
		return simplejson.New(), err
	}
	return simplejson.NewFromReader(res.Body)
}

// GetUserNotes 获取用户笔记通过 userId，cursor（可选参数）
func (x *xhsClient) GetUserNotes(userId string, cursor ...string) (*simplejson.Json, error) {
	cursorValue := ""
	if len(cursor) > 0 {
		cursorValue = cursor[0]
	}
	res, err := x.client.R().
		SetQueryParams(map[string]string{
			"num":     "30",
			"cursor":  cursorValue,
			"user_id": userId,
		}).
		Get("/api/sns/web/v1/user_posted")
	if err != nil {
		return simplejson.New(), err
	}
	return simplejson.NewFromReader(res.Body)
}

// GetNoteByKeyword 获取笔记通过关键词 keyword，page默认1,page_size默认20，sort默认排序，noteType默认全部笔记类型
func (x *xhsClient) GetNoteByKeyword(keyword string, args ...interface{}) (*simplejson.Json, error) {
	page := 1
	pageSize := 20
	sort := consts.General
	noteType := consts.All

	if len(args) >= 1 {
		if arg, ok := args[0].(int); ok {
			page = arg
		}
	}

	if len(args) >= 2 {
		if arg, ok := args[1].(int); ok {
			pageSize = arg
		}
	}

	if len(args) >= 3 {
		if arg, ok := args[2].(consts.SearchSortType); ok {
			sort = arg
		}
	}

	if len(args) >= 4 {
		if arg, ok := args[3].(consts.SearchNoteType); ok {
			noteType = arg
		}
	}
	data := map[string]any{
		"keyword":   keyword,
		"page":      page,
		"page_size": pageSize,
		"search_id": utils.GetSearchID(),
		"sort":      sort,
		"note_type": noteType,
	}
	dataStr, err := json.Marshal(data)
	if err != nil {
		return simplejson.New(), err
	}
	res, err := x.client.R().
		SetBody(dataStr).
		Post("/api/sns/web/v1/search/notes")
	if err != nil {
		return simplejson.New(), err
	}
	return simplejson.NewFromReader(res.Body)
}

// GetHomeFeed ....
func (x *xhsClient) GetHomeFeed(feedType consts.FeedType) (*simplejson.Json, error) {
	data := map[string]any{
		"cursor_score":         "",
		"num":                  40,
		"refresh_type":         1,
		"note_index":           0,
		"unread_begin_note_id": "",
		"unread_end_note_id":   "",
		"unread_note_count":    0,
		"category":             feedType,
	}
	dataStr, err := json.Marshal(data)
	if err != nil {
		return simplejson.New(), err
	}
	res, err := x.client.R().
		SetBody(dataStr).
		Post("/api/sns/web/v1/homefeed")
	if err != nil {
		return simplejson.New(), err
	}
	return simplejson.NewFromReader(res.Body)
}
