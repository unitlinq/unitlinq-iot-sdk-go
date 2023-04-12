package unitlinq

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
)

type ClientOptions struct {
	Cert    string
	CertKey string
}

type Client struct {
	TlsConfig tls.Config
	MQTT      mqtt.Client
	ClientID  uuid.UUID
}

type Token = mqtt.Token

func NewInstance(options ClientOptions) (Client, error) {
	var client Client
	certPool := x509.NewCertPool()
	cert, err := tls.LoadX509KeyPair(options.Cert, options.CertKey)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return Client{}, err
	}
	clientcert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return Client{}, err
	}
	client.ClientID, err = uuid.FromString(clientcert.Subject.CommonName)
	if err != nil {
		return Client{}, err
	}

	serverCert := `-----BEGIN CERTIFICATE-----
	MIIClTCCAhugAwIBAgIIQE4JchcUCD0wCgYIKoZIzj0EAwIwdzELMAkGA1UEBhMC
	SU4xEDAOBgNVBAgTB0d1amFyYXQxDzANBgNVBAcTBkluZGlhbjEVMBMGA1UEChMM
	RmxpbnQgRW5lcmd5MREwDwYDVQQLEwhVbml0Z3JpZDEbMBkGA1UEAxMSVW5pdGdy
	aWQgUm9vdCBDQSAxMCAXDTIyMDkwOTAwMDAwMFoYDzIwNTIwOTA4MjM1OTU5WjB3
	MQswCQYDVQQGEwJJTjEQMA4GA1UECBMHR3VqYXJhdDEPMA0GA1UEBxMGSW5kaWFu
	MRUwEwYDVQQKEwxGbGludCBFbmVyZ3kxETAPBgNVBAsTCFVuaXRncmlkMRswGQYD
	VQQDExJVbml0Z3JpZCBSb290IENBIDEwdjAQBgcqhkjOPQIBBgUrgQQAIgNiAARm
	49Ttjevys9q21Ue3bN2LdjCbgSME7T2rONjBCH6onja6U2PfS4hw9ALyjqprBcx0
	/qlmGpnAuiVEm4qJWn1CyNQfYEfZPpvuV0NkBON4Nx3GEzEgthvuMtij/NROxYaj
	cjBwMA8GA1UdEwEB/wQFMAMBAf8wHQYDVR0OBBYEFHD+22jmlpoT7can4bXC1Oic
	c8BLMAsGA1UdDwQEAwIBBjARBglghkgBhvhCAQEEBAMCAAcwHgYJYIZIAYb4QgEN
	BBEWD3hjYSBjZXJ0aWZpY2F0ZTAKBggqhkjOPQQDAgNoADBlAjBl5LeydupNqGzs
	YHObNZqZuHEz1XRH0siSs70ZGLHe3jbJxzDZF9lJeS8pNio3S6wCMQDYkWMf2wuR
	wbd/MYG7mURoAkzwlIJXHN94USIQ9KQ9hWCKTl7UyzGChPKpunDAJe0=
	-----END CERTIFICATE-----`

	certPool.AppendCertsFromPEM([]byte(serverCert))

	client.TlsConfig = tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{cert},
	}

	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://gateway.unitgrid.in:8883")
	opts.SetUsername("device")
	opts.SetClientID(client.ClientID.String())
	opts.SetTLSConfig(&client.TlsConfig)
	//opts.SetCleanSession(true)
	//opts.SetStore(mqtt.NewFileStore(config.GlobalConfig.GetString("store") + "mqstore"))
	opts.SetConnectRetryInterval(10 * time.Second)
	opts.SetConnectRetry(true)
	opts.SetAutoReconnect(true)
	//opts.SetConnectionLostHandler(ugconnectionLostHandler)
	opts.SetCleanSession(false)
	client.MQTT = mqtt.NewClient(opts)
	return client, nil
}

func (c *Client) Connect() Token {
	return c.MQTT.Connect()
}

func (c *Client) Close() {
	c.MQTT.Disconnect(1000)
}

func (c *Client) IsConnected() bool {
	return c.MQTT.IsConnectionOpen()
}
