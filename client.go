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
	Cert           string
	CertKey        string
	StorageType    int
	Store          string
	RequestTimeout time.Duration
}

const (
	Memory int = 0
	File   int = 1
)

type Client struct {
	TlsConfig tls.Config
	MQTT      mqtt.Client
	ClientID  uuid.UUID
	timeout   time.Duration
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
	client.timeout = options.RequestTimeout
	opts := mqtt.NewClientOptions()
	opts.AddBroker("ssl://gateway.unitgrid.in:8883")
	opts.SetUsername("device")
	opts.SetClientID(client.ClientID.String())
	opts.SetTLSConfig(&client.TlsConfig)
	//opts.SetCleanSession(true)
	switch options.StorageType {
	case Memory:
		opts.SetStore(mqtt.NewMemoryStore())
	case File:
		opts.SetStore(mqtt.NewFileStore(options.Store + "mqttstore"))
	}
	opts.SetConnectRetryInterval(10 * time.Second)
	opts.SetConnectRetry(true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(ugconnectHandler)
	//opts.SetConnectionLostHandler(ugconnectionLostHandler)
	opts.SetCleanSession(true)
	client.MQTT = mqtt.NewClient(opts)
	clientID = client.GetClientID()
	return client, nil
}

func (c *Client) Connect() Token {
	token := c.MQTT.Connect()
	c.safeSubcscribe("device/"+c.GetClientID()+"/api/response/cbor", 1, defaultResponseHandler)
	responseSubscribed = true
	return token
}

func (c *Client) Close() {
	c.MQTT.Disconnect(1000)
}

func (c *Client) IsConnected() bool {
	return c.MQTT.IsConnectionOpen()
}

func (c *Client) GetClientID() string {
	return c.ClientID.String()
}

func (c *Client) safeSubcscribe(topic string, qos byte, callaback mqtt.MessageHandler) error {
	token := c.MQTT.Subscribe(topic, qos, callaback)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	temp := Subscripiton{
		Topic:    topic,
		Qos:      qos,
		Callback: callaback,
	}
	subscriptionPool = append(subscriptionPool, temp)
	return nil
}

var ugconnectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connection Established")
	fmt.Println("Subscribing the topics")
	for _, v := range subscriptionPool {
		token := client.Subscribe(v.Topic, v.Qos, v.Callback)
		if token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
		fmt.Printf("Subscribed: %s\n", v.Topic)
	}
}

func (c *Client) safeSubAppend(topic string, qos byte, callaback mqtt.MessageHandler) {
	temp := Subscripiton{
		Topic:    topic,
		Qos:      qos,
		Callback: callaback,
	}
	subscriptionPool = append(subscriptionPool, temp)
}
