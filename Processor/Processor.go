// 处理收到的信息事件
package Processor

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hoshinonyaruko/gensokyo/callapi"
	"github.com/hoshinonyaruko/gensokyo/config"
	"github.com/hoshinonyaruko/gensokyo/mylog"
	"github.com/hoshinonyaruko/gensokyo/wsclient"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

// Processor 结构体用于处理消息
type Processors struct {
	Api             openapi.OpenAPI                   // API 类型
	Apiv2           openapi.OpenAPI                   //群的API
	Settings        *config.Settings                  // 使用指针
	Wsclient        []*wsclient.WebSocketClient       // 指针的切片
	WsServerClients []callapi.WebSocketServerClienter //ws server被连接的客户端
}

type Sender struct {
	Nickname string `json:"nickname"`
	TinyID   string `json:"tiny_id"`
	UserID   int64  `json:"user_id"`
	Role     string `json:"role,omitempty"`
}

// 频道信息事件
type OnebotChannelMessage struct {
	ChannelID   string      `json:"channel_id"`
	GuildID     string      `json:"guild_id"`
	Message     interface{} `json:"message"`
	MessageID   string      `json:"message_id"`
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	SelfID      int64       `json:"self_id"`
	SelfTinyID  string      `json:"self_tiny_id"`
	Sender      Sender      `json:"sender"`
	SubType     string      `json:"sub_type"`
	Time        int64       `json:"time"`
	Avatar      string      `json:"avatar"`
	UserID      int64       `json:"user_id"`
	RawMessage  string      `json:"raw_message"`
	Echo        string      `json:"echo,omitempty"`
}

// 群信息事件
type OnebotGroupMessage struct {
	RawMessage  string      `json:"raw_message"`
	MessageID   int         `json:"message_id"`
	GroupID     int64       `json:"group_id"` // Can be either string or int depending on p.Settings.CompleteFields
	MessageType string      `json:"message_type"`
	PostType    string      `json:"post_type"`
	SelfID      int64       `json:"self_id"` // Can be either string or int
	Sender      Sender      `json:"sender"`
	SubType     string      `json:"sub_type"`
	Time        int64       `json:"time"`
	Avatar      string      `json:"avatar"`
	Echo        string      `json:"echo,omitempty"`
	Message     interface{} `json:"message"` // For array format
	MessageSeq  int         `json:"message_seq"`
	Font        int         `json:"font"`
	UserID      int64       `json:"user_id"`
}

// 私聊信息事件
type OnebotPrivateMessage struct {
	RawMessage  string        `json:"raw_message"`
	MessageID   int           `json:"message_id"` // Can be either string or int depending on logic
	MessageType string        `json:"message_type"`
	PostType    string        `json:"post_type"`
	SelfID      int64         `json:"self_id"` // Can be either string or int depending on logic
	Sender      PrivateSender `json:"sender"`
	SubType     string        `json:"sub_type"`
	Time        int64         `json:"time"`
	Avatar      string        `json:"avatar"`
	Echo        string        `json:"echo,omitempty"`
	Message     interface{}   `json:"message"`     // For array format
	MessageSeq  int           `json:"message_seq"` // Optional field
	Font        int           `json:"font"`        // Optional field
	UserID      int64         `json:"user_id"`     // Can be either string or int depending on logic
}

type PrivateSender struct {
	Nickname string `json:"nickname"`
	UserID   int64  `json:"user_id"` // Can be either string or int depending on logic
}

func FoxTimestamp() int64 {
	return time.Now().Unix()
}

// ProcessInlineSearch 处理内联查询
func (p *Processors) ProcessInlineSearch(data *dto.WSInteractionData) error {
	//ctx := context.Background() // 或从更高级别传递一个上下文

	// 在这里处理内联查询
	// 这可能涉及解析查询、调用某些API、获取结果并格式化为响应
	// ...

	// 示例：发送响应
	// response := "Received your interaction!"            // 创建响应消息
	// err := p.api.PostInteractionResponse(ctx, response) // 替换为您的OpenAPI方法
	// if err != nil {
	// 	return err
	// }

	return nil
}

//return nil

//下面是测试时候固定代码
//发私信给机器人4条机器人不回,就不能继续发了

// timestamp := time.Now().Unix() // 获取当前时间的int64类型的Unix时间戳
// timestampStr := fmt.Sprintf("%d", timestamp)

// dm := &dto.DirectMessage{
// 	GuildID:    GuildID,
// 	ChannelID:  ChannelID,
// 	CreateTime: timestampStr,
// }

// PrintStructWithFieldNames(dm)

// // 发送默认回复
// toCreate := &dto.MessageToCreate{
// 	Content: "默认私信回复",
// 	MsgID:   data.ID,
// }
// _, err = p.Api.PostDirectMessage(
// 	context.Background(), dm, toCreate,
// )
// if err != nil {
// 	mylog.Println("Error sending default reply:", err)
// 	return nil
// }

// 打印结构体的函数
func PrintStructWithFieldNames(v interface{}) {
	val := reflect.ValueOf(v)

	// 如果是指针，获取其指向的元素
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()

	// 确保我们传入的是一个结构体
	if typ.Kind() != reflect.Struct {
		mylog.Println("Input is not a struct")
		return
	}

	// 迭代所有的字段并打印字段名和值
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)
		mylog.Printf("%s: %v\n", field.Name, value.Interface())
	}
}

// 将结构体转换为 map[string]interface{}
func structToMap(obj interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	j, _ := json.Marshal(obj)
	json.Unmarshal(j, &out)
	return out
}

// 修改函数的返回类型为 *Processor
func NewProcessor(api openapi.OpenAPI, apiv2 openapi.OpenAPI, settings *config.Settings, wsclient []*wsclient.WebSocketClient) *Processors {
	return &Processors{
		Api:      api,
		Apiv2:    apiv2,
		Settings: settings,
		Wsclient: wsclient,
	}
}

// 修改函数的返回类型为 *Processor
func NewProcessorV2(api openapi.OpenAPI, apiv2 openapi.OpenAPI, settings *config.Settings) *Processors {
	return &Processors{
		Api:      api,
		Apiv2:    apiv2,
		Settings: settings,
	}
}

// 发信息给所有连接正向ws的客户端
func (p *Processors) SendMessageToAllClients(message map[string]interface{}) error {
	var result *multierror.Error

	for _, client := range p.WsServerClients {
		// 使用接口的方法
		err := client.SendMessage(message)
		if err != nil {
			// Append the error to our result
			result = multierror.Append(result, fmt.Errorf("failed to send to client: %w", err))
		}
	}

	// This will return nil if no errors were added
	return result.ErrorOrNil()
}

// 方便快捷的发信息函数
func (p *Processors) BroadcastMessageToAll(message map[string]interface{}) error {
	var errors []string

	// 发送到我们作为客户端的Wsclient
	for _, client := range p.Wsclient {
		err := client.SendMessage(message)
		if err != nil {
			errors = append(errors, fmt.Sprintf("error sending private message via wsclient: %v", err))
		}
	}

	// 发送到我们作为服务器连接到我们的WsServerClients
	for _, serverClient := range p.WsServerClients {
		err := serverClient.SendMessage(message)
		if err != nil {
			errors = append(errors, fmt.Sprintf("error sending private message via WsServerClient: %v", err))
		}
	}

	// 在循环结束后处理记录的错误
	if len(errors) > 0 {
		return fmt.Errorf(strings.Join(errors, "; "))
	}

	return nil
}
