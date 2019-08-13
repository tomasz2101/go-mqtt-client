package mqttclient

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/tomasz2101/go-helpers"
)

// MQTTClient will do something
type MQTTClient struct {
	Client   mqtt.Client
	Hostname string
	Port     int
	ID       string
	Username string
	Password string
}

// ReturnURL will do sometihng
func (m *MQTTClient) ReturnURL() string {
	return "mqtt://internal:internal@localhost:1883/testing"
}

// Connect will do sometihng
func (m *MQTTClient) Connect(postfixID string) (mqtt.Client, error) {
	// opts := createClientOptions(clientId, uri)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", m.Hostname, m.Port))
	opts.SetUsername(m.Username)
	opts.SetPassword(m.Password)
	opts.SetClientID(m.ID + "/" + postfixID)
	m.Client = mqtt.NewClient(opts)
	token := m.Client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return m.Client, nil
}

// Listen will do something
// func (m *MQTTClient) Listen(topic string) {
// 	fmt.Println("Listen")
// 	client := m.Client.Connect("sub")
// 	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
// 		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
// 	})
// }

// EndConnection will disconnect client from broker
func (m *MQTTClient) EndConnection() error {
	if m.Client.IsConnected() {
		m.Client.Disconnect(20)
		fmt.Println("client disconnected")
	}
	return nil
}

type deviceInfo struct {
	Address  string `json:"address"`
	DeviceID string `json:"deviceid"`
	Hostname string `json:"hostname"`
	Mac      string `json:"mac"`
}

type mqttMessage struct {
	Time       string `json:"time"`
	ID         string `json:"id"`
	Message    string `json:"message"`
	DeviceInfo string `json:"device_info"`
}

// Publish will do something
func (m *MQTTClient) Publish(topic string, message string) {
	hostname, _ := os.Hostname()
	ipData, _ := net.LookupHost(hostname)
	ip := "unknown"
	if len(ipData) > 0 {
		ip = fmt.Sprintf("%v", ipData[0])
	}
	res2D := &deviceInfo{
		Address:  ip,
		DeviceID: m.ID,
		Mac:      helpers.GetMacAddr(),
		Hostname: hostname}
	res2B, err := json.Marshal(res2D)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	res1D := &mqttMessage{
		Time:       helpers.GetDate(),
		Message:    message,
		DeviceInfo: string(res2B)}
	token := m.Client.Publish(topic, 2, false, strings.Replace(string(helpers.GetJSON(res1D)), "\\\"", "\"", -1))
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// PrepareData will do something
// func (m MQTTClient) PrepareData(messageID string, inputData map[string]string) string {

// 	hostname, _ := os.Hostname()
// 	ip := "unknown"
// 	ipData, _ := net.LookupHost(hostname)
// 	if len(ipData) > 0 {
// 		ip = fmt.Sprintf("%v", ipData[0])
// 	}
// 	res2D := &deviceInfo{
// 		Address:  ip,
// 		DeviceID: m.ID,
// 		Mac:      helpers.GetMacAddr(),
// 		Hostname: hostname}
// 	res2B, err := json.Marshal(res2D)
// 	jsonInput, err := json.Marshal(inputData)
// 	if err != nil {
// 		fmt.Printf("%s\n", err)
// 	}

// 	res1D := &mqttMessage{
// 		Time:       helpers.GetDate(),
// 		ID:         messageID,
// 		Message:       string(jsonInput),
// 		DeviceInfo: string(res2B)}
// 	return strings.Replace(string(helpers.GetJSON(res1D)), "\\\"", "\"", -1)
// }
