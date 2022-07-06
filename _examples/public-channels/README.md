Send Message to the Public Channels
===

This example demonstrates how you can use the Channelize package to publish messages
to the different public channels.

How Does It Work?
===

The purpose of this example is creating a websocket server that is able to publish different
types of messages into the related channels. A channel is a stream of events with the same type.

This application publishes three different event type:
* News
* Alert
* Notification

Each event should be published in the channel that it belongs to. We will have the following channels:
* news: Is a stream of `News` events.
* alerts: Is a stream of `Alert` events.
* notifications: Is a stream of `Notification` events.

When you start the server, it registers the channels and calls `publish` function to send events 
to a specific channel in a separate goroutine. The `publish` simulate a message broker consumer 
function, and it accepts a `messageProducerFunc` function type as input which simulates a message
broker producer.

How to Run It?
===

Use the following command to run the `pubic-channels` app:
```shell
go run ./...
```

After executing above command, the application will run an HTTP server on port `8080`. 
The server has only one rout, `/ws` which creates a websocket connection.  

How to see the result?
===

After running the application, you can use postman or any browser plugin to connect to the
websocket server and send message to it.

First thing you should do is connecting to the websocket:
![connect to the websocket](images/connect.png?raw=true "Connect")

Then you should send a message to the open connection to subscribe to some existing channels:
![subscribe channels](images/subscribe.png?raw=true "Subscribe")

Also, you can unsubscribe channels to stop receiving stream:
![unsubscribe channels](images/unsubscribe.png?raw=true "Unsubscribe")

In the following gif, you can see all of these steps with three different users:
![all steps](https://user-images.githubusercontent.com/11541936/177541494-3b163571-2364-4a65-93e9-6f07797044a6.gif?raw=true "All")
