package unitlinq

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

// ClientOptions - A struct describing the client options
type ClientOptions struct {
	Cert           string
	CertKey        string
	StorageType    int
	Store          string
	RequestTimeout time.Duration
	OnConnect      OnConnectCallback
	OnConnLost     OnConnLostCallback
}

const (
	Memory int = 0
	File   int = 1
	SQLite int = 2
)

type Client struct {
	TlsConfig tls.Config
	MQTT      mqtt.Client
	ClientID  uuid.UUID
	timeout   time.Duration
	options   ClientOptions
	db        *gorm.DB
	destroy   chan bool
	dbLock    sync.Mutex
}

var opts *mqtt.ClientOptions

type Token = mqtt.Token

// Create a new instance of unitlinq client. You can create as many instance you want. The only restriction is that you can not use same client
// certificate more than once.
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
	fmt.Println("Parsed:", client.ClientID.String())

	ok := certPool.AppendCertsFromPEM([]byte(serverCertEmbedded))
	if !ok {
		fmt.Println("Server cert not parsed")
	}

	client.TlsConfig = tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{cert},
	}
	client.timeout = options.RequestTimeout
	client.options = options
	opts = mqtt.NewClientOptions()
	opts.AddBroker("ssl://gateway.unitgrid.in:8883")
	opts.SetUsername("device")
	opts.SetClientID(client.ClientID.String())
	opts.SetTLSConfig(&client.TlsConfig)
	opts.SetCleanSession(true)
	switch options.StorageType {
	case Memory:
		opts.SetStore(mqtt.NewMemoryStore())
	case File:
		opts.SetStore(mqtt.NewFileStore(options.Store + "mqttstore"))
	case SQLite:
		opts.SetStore(mqtt.NewMemoryStore())
		// Initialize SQLite DB
		client.initStorage(options.Store)
	}
	opts.SetConnectRetryInterval(2 * time.Second)
	opts.SetConnectRetry(true)
	opts.SetAutoReconnect(true)
	opts.SetOnConnectHandler(client.onConnectActions)
	opts.SetReconnectingHandler(client.onRetryAction)
	opts.SetConnectionLostHandler(client.onConnLostAction)
	//opts.SetConnectionLostHandler(ugconnectionLostHandler)
	opts.SetCleanSession(true)
	opts.SetWriteTimeout(10 * time.Second)
	client.MQTT = mqtt.NewClient(opts)
	clientID = client.GetClientID()
	return client, nil
}

// Establish a connection with unitlinq server
func (c *Client) Connect() Token {
	token := c.MQTT.Connect()
	c.safeSubcscribe("device/"+c.GetClientID()+"/api/response/cbor", 1, defaultResponseHandler)
	responseSubscribed = true
	go c.SubmitMessagesTask()
	return token
}

// Close the connection
func (c *Client) Close() {
	c.MQTT.Disconnect(1000)
	DB, err := c.db.DB()
	if err != nil {
		fmt.Println(err)
	}
	DB.Close()
}

// Check whether connection is active or not. Returns true if connection is active
func (c *Client) IsConnected() bool {
	return c.MQTT.IsConnectionOpen()
}

// Get client ID string parsed from client certificate
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

func (c *Client) onConnectActions(client mqtt.Client) {
	//fmt.Println("Connection Established")
	//fmt.Println("Subscribing the topics")
	fmt.Println("On connect handler for Client with CLientID:", c.GetClientID())
	for _, v := range subscriptionPool {
		token := client.Subscribe(v.Topic, v.Qos, v.Callback)
		if token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
		}
		//fmt.Printf("Subscribed: %s\n", v.Topic)
	}
	// Issue callback
	c.options.OnConnect(*c)
}

func (c *Client) onConnLostAction(client mqtt.Client, err error) {
	c.options.OnConnLost(*c, err)
}

func (c *Client) onRetryAction(cl mqtt.Client, opts *mqtt.ClientOptions) {
	fmt.Println("Retrying to connect with broker!")
}

func (c *Client) safeSubAppend(topic string, qos byte, callaback mqtt.MessageHandler) {
	temp := Subscripiton{
		Topic:    topic,
		Qos:      qos,
		Callback: callaback,
	}
	subscriptionPool = append(subscriptionPool, temp)
}
