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

// Client will do something
type Client struct {
	Hostname string
	Port     int
	ID       string
	Username string
	Password string
}

// ReturnURL will do sometihng
func (mqtt_client Client) ReturnURL() string {
	return "mqtt://internal:internal@localhost:1883/testing"
}

// Connect will do sometihng
func (mqtt_client Client) Connect(postfixID string) mqtt.Client {
	// opts := createClientOptions(clientId, uri)
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", mqtt_client.Hostname, mqtt_client.Port))
	opts.SetUsername(mqtt_client.Username)
	opts.SetPassword(mqtt_client.Password)
	opts.SetClientID(mqtt_client.ID + "/" + postfixID)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	for !token.WaitTimeout(3 * time.Second) {
	}
	if err := token.Error(); err != nil {
		log.Fatal(err)
	}
	return client
}

// Listen will do something
func (mqtt_client Client) Listen(topic string) {
	fmt.Println("Listen")
	client := mqtt_client.Connect("sub")
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

// EndConnection will disconnect client from broker
func (mqtt_client Client) EndConnection(client mqtt.Client) {
	client.Disconnect(250)
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
func (mqtt_client Client) Publish(client mqtt.Client, topic string, message string) {
	hostname, _ := os.Hostname()
	ipData, _ := net.LookupHost(hostname)
	ip := "unknown"
	if len(ipData) > 0 {
		ip = fmt.Sprintf("%v", ipData[0])
	}
	res2D := &deviceInfo{
		Address:  ip,
		DeviceID: mqtt_client.ID,
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
	token := client.Publish(topic, 2, false, strings.Replace(string(helpers.GetJSON(res1D)), "\\\"", "\"", -1))
	token.Wait()
	if token.Error() != nil {
		log.Fatal(token.Error())
	}
}

// PrepareData will do something
// func (mqtt_client Client) PrepareData(messageID string, inputData map[string]string) string {

// 	hostname, _ := os.Hostname()
// 	ip := "unknown"
// 	ipData, _ := net.LookupHost(hostname)
// 	if len(ipData) > 0 {
// 		ip = fmt.Sprintf("%v", ipData[0])
// 	}
// 	res2D := &deviceInfo{
// 		Address:  ip,
// 		DeviceID: mqtt_client.ID,
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
