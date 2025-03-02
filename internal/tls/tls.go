package tls

import (
	"crypto/tls"
	"crypto/x509"
	"os"
)

// AWS IoT Core接続用のTLS設定を作成します
func NewTLSConfig(rootCAPath, certPath, keyPath string) (*tls.Config, error) {
	// ルート証明書を読み込む
	rootCA, err := os.ReadFile(rootCAPath)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(rootCA)

	// クライアント証明書を読み込む
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	// TLS設定を返す
	return &tls.Config{
		RootCAs:      pool,
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}, nil
}
