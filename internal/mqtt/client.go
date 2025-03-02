package mqtt

import (
	"fmt"
	"time"

	"github.com/dev-satoshi/aws-mqtt-pubsub/internal/config"
	"github.com/dev-satoshi/aws-mqtt-pubsub/internal/tls"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTクライアントのラッパー
type Client struct {
	mqttClient mqtt.Client
	config     *config.Config
}

// メッセージ受信ハンドラー
var messageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("受信: トピック: %s, メッセージ: %s\n", msg.Topic(), msg.Payload())
}

// 新しいMQTTクライアントを作成
func NewClient(cfg *config.Config) (*Client, error) {
	// TLS設定を作成
	tlsConfig, err := tls.NewTLSConfig(cfg.RootCAPath, cfg.CertPath, cfg.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("TLS設定エラー: %w", err)
	}

	// MQTT接続オプションを設定
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcps://%s:8883", cfg.Endpoint))
	opts.SetClientID(cfg.ClientID)
	opts.SetTLSConfig(tlsConfig)
	opts.SetDefaultPublishHandler(messageHandler)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)

	return &Client{
		mqttClient: client,
		config:     cfg,
	}, nil
}

// AWS IoT Coreに接続
func (c *Client) Connect() error {
	if token := c.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// トピックをサブスクライブ
func (c *Client) Subscribe(topic string) error {
	if token := c.mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// トピックにメッセージを発行
func (c *Client) Publish(topic, message string) error {
	token := c.mqttClient.Publish(topic, 0, false, message)
	token.Wait()
	return token.Error()
}
