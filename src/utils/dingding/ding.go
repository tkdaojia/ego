package dingding

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DingDing struct{}

type DingTalkClient struct {
	AccessToken string
	AgentID     int64
	HTTPClient  *http.Client
}

type FormComponentValue struct {
	Name          string `json:"name"`
	ComponentType string `json:"componentType"`
	Value         string `json:"value"`
	ID            string `json:"id"`
}

type OperationRecord struct {
	UserID        string `json:"userId"`
	ShowName      string `json:"showName"`
	Result        string `json:"result"`
	Remark        string `json:"remark"`
	CreateTime    string `json:"createTime"`
	OperationType string `json:"operationType"`
}

type TaskRecord struct {
	TaskID     int64  `json:"taskId"` // 注意这里是数字/长整型
	UserID     string `json:"userId"`
	Status     string `json:"status"` // RUNNING, COMPLETED
	Result     string `json:"result"` // NONE, AGREE
	CreateTime string `json:"createTime"`
	MobileUrl  string `json:"mobileUrl"`
	PcUrl      string `json:"pcUrl"`
	ActivityID string `json:"activityId"`
}

type DingTalkProcessInstanceRoot struct {
	Success bool                    `json:"success"`
	Result  ProcessInstanceDetailV2 `json:"result"`
}

type ProcessInstanceDetailV2 struct {
	Title               string               `json:"title"`
	BusinessID          string               `json:"businessId"`
	Status              string               `json:"status"` // RUNNING, COMPLETED, TERMINATED
	Result              string               `json:"result"` // NONE, agree, refuse
	OriginatorUserID    string               `json:"originatorUserId"`
	OriginatorDeptID    string               `json:"originatorDeptId"`
	OriginatorDeptName  string               `json:"originatorDeptName"`
	CreateTime          string               `json:"createTime"`
	FinishTime          string               `json:"finishTime"`
	FormComponentValues []FormComponentValue `json:"formComponentValues"` // 对应 json 中的 formComponentValues
	OperationRecords    []OperationRecord    `json:"operationRecords"`    // 对应 json 中的 operationRecords
	Tasks               []TaskRecord         `json:"tasks"`
}

// 新建客户端并获取 access_token
func (dingding *DingDing) NewDingTalkClient(appKey, appSecret string, AgentID int64) (*DingTalkClient, error) {
	client := &DingTalkClient{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
	}

	// 获取 access_token
	tokenURL := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s",
		appKey, appSecret)

	resp, err := client.HTTPClient.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("获取 access_token 失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if result["errcode"] != nil && int(result["errcode"].(float64)) != 0 {
		return nil, fmt.Errorf("获取 access_token 失败: %v", result["errmsg"])
	}

	client.AccessToken = result["access_token"].(string)
	fmt.Printf("✅ 获取 access_token 成功: %s\n", client.AccessToken)

	client.AgentID = AgentID
	return client, nil
}

// 发送文本消息给指定用户
func (c *DingTalkClient) SendTextMessage(userIDs []string, content string) error {
	type TextMessage struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}

	// 发送消息响应
	type MessageResponse struct {
		Errcode   int    `json:"errcode"`
		Errmsg    string `json:"errmsg"`
		RequestID string `json:"request_id"`
		TaskID    int64  `json:"task_id"` // 任务ID，用于查询发送状态
	}

	apiURL := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2"

	// 构造消息对象
	msg := TextMessage{
		MsgType: "text",
		Text: struct {
			Content string `json:"content"`
		}{
			Content: content,
		},
	}

	// 序列化消息对象为JSON字符串
	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 构造表单数据
	formData := url.Values{}
	formData.Set("access_token", c.AccessToken)
	formData.Set("agent_id", fmt.Sprintf("%d", c.AgentID))
	formData.Set("msg", string(msgJSON))
	formData.Set("userid_list", strings.Join(userIDs, ","))
	formData.Set("to_all_user", "false")

	fmt.Printf("🔍 发送表单数据: %s\n", formData.Encode())

	// 发送POST请求，使用application/x-www-form-urlencoded格式
	resp, err := c.HTTPClient.PostForm(apiURL, formData)
	if err != nil {
		return fmt.Errorf("发送消息请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var result MessageResponse
	json.Unmarshal(body, &result)

	if result.Errcode != 0 {
		return fmt.Errorf("发送消息失败: %s", result.Errmsg)
	}

	fmt.Printf("✅ 消息已提交发送，请求ID: %s, 任务ID: %d\n", result.RequestID, result.TaskID)
	return nil
}

func (c *DingTalkClient) GetProcessInstanceDetailV2(ctx context.Context, processInstanceID string) (*ProcessInstanceDetailV2, error) {
	accessToken := c.AccessToken

	apiURL := fmt.Sprintf("https://api.dingtalk.com/v1.0/workflow/processInstances?processInstanceId=%s", processInstanceID)

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("x-acs-dingtalk-access-token", accessToken)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("接口调用失败，HTTP 状态码: %d", resp.StatusCode)
	}

	fmt.Println(string(body))
	// 实例化最外层根结构
	var root DingTalkProcessInstanceRoot
	if err := json.Unmarshal(body, &root); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if !root.Success {
		return nil, fmt.Errorf("钉钉接口返回 success 为 false")
	}

	// 直接返回核心的 result 节点数据
	return &root.Result, nil
}
