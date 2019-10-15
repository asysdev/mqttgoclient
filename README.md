# mqttgoclient
MQTT Client developed in golang enhancing capability of MQTT Client. It uses paho client underneath for MQTT communication further facilitating developers to focus on application specific code. 

# Usage

1. Create mqttgoclient object as below:

        var client mqtt.Client

2. Configure broker information:

 	   client.Broker = "tcp://test.mosquitto.org:1883"     
  
3. Configure session details, it's recommended to start with clean session:

        client.CleanSession = true     
     
4. Assign client ID, this is optional field. mqttgoclient will generate random and unique client id if this is not provided. 

        client.ClientId = "goclienttest"

5. Register subscription handler to receive subscribed messages
        
        client.MessageArrived = MessageArrived
        
6. Init client :

	    err := client.Init()

	    if err != nil {
		     log.Panic(err)
	    }

7. Connect client to broker:

        client.Connect()
        
8. Subscription 

        err := client.Subscribe("goclient_sub", 0)
      
9. Publish

        client.Publish("goclient_pub", 0, false, "Hello")
