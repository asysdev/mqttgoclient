/*
*    Copyright (C) 2019  asysdev
*    All Rights Reserved.
*
*    This program is free software: you can redistribute it and/or modify
*    it under the terms of the GNU General Public License as published by
*    the Free Software Foundation, version 3 of the License.
*
*    This program is distributed in the hope that it will be useful,
*    but WITHOUT ANY WARRANTY; without even the implied warranty of
*    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*    GNU General Public License for more details.
*
*    You should have received a copy of the GNU General Public License
*    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*
 */

package mqtt

import (
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"log"
	"math/rand"
)

// Receiver type allows user to receive subscribed messages in application
// layer. User can define multiple handles of Receiver for various subscriptions.
type Receiver func(*Client, string, []byte)

// Client struct contains MQTT client configuration and handle of client connection
// for further access of publish and subscriptions. There is not any restriction on
// number of client instances. User can define multiple instances for same broker or
// different brokers in single application.
type Client struct {
	Broker             string      // Broker URL to connect
	User               string      // Broker User name for authentication
	Password           string      // Broker password for authentication
	ClientId           string      // Client ID
	CleanSession       bool        // Clear session at Broker on connect, default false
	Qos                int         // Quality of Service, (i.e, 0: QOS0, 1 QOS1, 2 QOS2)
	DisconnectInterval uint        // Client will disconnect after wait of this period, default 0
	connObj            MQTT.Client // MQTT Client Connection object
	MessageArrived     Receiver    // Subscribe receiver handler for this client object
}

// messageReceiver is internal handler to mqttgoclient for managing subscription handlers
func (client *Client) messageReceiver(mqttClient MQTT.Client, message MQTT.Message) {

	// call client registered message handler so client can perform application specific operations
	client.MessageArrived(client, message.Topic(), message.Payload())

}

// Init is an initialization function for defined Client object. Each client object has Init function
// which must be called before accessing any other member functions of struct Client.
func (client *Client) Init() (err error) {

	log.Println("MQTT Client Init Started")

	// Generate error if client doesn't contain any host URL
	if client.Broker == "" {
		err = errors.New("Empty Broker URL")
		return
	}

	// Generate random client id if user has not defined client id
	if client.ClientId == "" {

		log.Println("Empty Client ID, Generating Random One")

		client.ClientId = fmt.Sprintf("mqttgoclient_id", rand.Int())

		log.Printf("Generated RandomID: '%s'", client.ClientId)
	}

	// Fill options instance with mqtt client configurations to connect to broker
	opts := MQTT.NewClientOptions()
	opts.AddBroker(client.Broker)
	opts.SetClientID(client.ClientId)

	// If user and password is empty, it means broker does not require password authentication
	// If broker is expecting authentication then these fields must be filled with appropriate credentials.
	if client.User != "" && client.Password != "" {
		opts.SetUsername(client.User)
		opts.SetPassword(client.Password)
	}
	opts.SetCleanSession(client.CleanSession)

	client.connObj = MQTT.NewClient(opts)

	log.Println("MQTT Client Init Finished, Ready to Connect.")

	return
}

// Connect initiates MQTT Connect sequence as per MQTT protocol. On receiving CONN ACK from
// broker, it changes state to connected. On failure returns error.
func (client *Client) Connect() (err error) {

	log.Println("Connecting...")

	// Connect to broker and wait until acknowledgement.
	if token := client.connObj.Connect(); token.Wait() && token.Error() != nil {
		err = token.Error()
		log.Println("Connect Error: ", err)
		return
	}

	log.Println("MQTT Client Connected")

	return
}

// Disconnect drops connection to MQTT broker.
func (client *Client) Disconnect() {

	client.connObj.Disconnect(client.DisconnectInterval)

	return
}

// Publish sends given message to specific topic on MQTT broker.
func (client *Client) Publish(topic string, qos byte, retain bool, payload interface{}) (err error) {

	token := client.connObj.Publish(topic, qos, retain, payload)
	token.Wait()

	return
}

// Subscribe performs subscription on MQTT broker for specific topic with QOS
// Client library will call message handler on reception of any message on
// subscribed topic.
func (client *Client) Subscribe(topic string, qos byte) (err error) {

	if topic == "" {
		err = errors.New("Empty Subscribe Topic")
		return
	}

	if token := client.connObj.Subscribe(topic, qos, client.messageReceiver); token.Wait() && token.Error() != nil {
		err = token.Error()
		return
	}

	log.Println("Subscribed Successfully")

	return
}

// Unsubscribe is used to remove subscription of specific topic
// from MQTT broker. After this operation, client will not receive
// any message published on unsubscribed topic.
func (client *Client) Unsubscribe(topic string) {
	client.connObj.Unsubscribe(topic)
}

// Connected returns true if client is connected to broker actively.
// It returns false if client is in reconnect or disconnected mode
//
// This can be used to get absolute connection state to broker.
func (client *Client) Connected() bool {
	return client.connObj.IsConnectionOpen()
}

// Alive returns true if client is connected to broker. It returns true
// even if client is in reconnecting state.
//
// This can be used to check whether client is connected to broker or not
// regardless of various MQTT connection stages.
func (client *Client) Alive() bool {
	return client.connObj.IsConnected()
}
