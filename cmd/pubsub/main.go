package main

import (
	"fmt"
	"time"

	"github.com/dev-satoshi/aws-mqtt-pubsub/internal/config"
	"github.com/dev-satoshi/aws-mqtt-pubsub/internal/mqtt"
)

func main() {
	cfg, err := config.LoadConfig("configs/config.json")
	if err != nil {
		fmt.Println("設定ファイルの読み込みエラー:", err)
		return
	}

	client, err := mqtt.NewClient(cfg)
	if err != nil {
		fmt.Println("MQTTクライアント作成エラー:", err)
		return
	}

	if err := client.Connect(); err != nil {
		fmt.Println("接続エラー:", err)
		return
	}
	fmt.Println("AWS IoT Coreに接続しました")

	if err := client.Subscribe(cfg.SubTopic); err != nil {
		fmt.Println("サブスクライブエラー:", err)
		return
	}
	fmt.Printf("トピック '%s' をサブスクライブしました\n", cfg.SubTopic)

	// メッセージを定期的に送信
	counter := 0
	for {
		counter++
		text := fmt.Sprintf("こんにちは #%d", counter)
		if err := client.Publish(cfg.PubTopic, text); err != nil {
			fmt.Printf("送信エラー: %v\n", err)
			continue
		}
		fmt.Printf("送信: トピック: %s, メッセージ: %s\n", cfg.PubTopic, text)
		time.Sleep(5 * time.Second)
	}
}
