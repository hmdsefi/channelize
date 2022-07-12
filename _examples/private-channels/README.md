Send Message to the Private Channels
===

This example demonstrates how you can use the Channelize package to publish messages
to a private channel.

How Does It Work?
===

The purpose of this example is creating a websocket server that is able to publish messages
on a private channel. The private channel name that publishes user related data is `notifications`.
The client should expect receiving events from `notifications` channel that is related to its
userID.

In `auth.go` file we implemented a minimal version of an auth service to create a `Token` struct. The `Token`
object includes unique `token`, `userID`, and an `expiresAt` field. The default Time To Live(TTL) of a token
is 30 minutes.

To create a token, you should call the following request:

```shell
curl --location --request POST 'localhost:8080/ws/token'
```

The response should be something like this:

```json
{
  "token": "a77b49689caa58d209958420a881f5b81a31ae2a9964ba4479ed87eccfe9"
}
```

How to Run It?
===

Use the following command to run the `private-channels` app:

```shell
go run ./...
```

After executing above command, the application will run an HTTP server on port `8080`.
The server has two APIs:

* `/ws`: To create a new websocket connection and connect to the websocket server.
* `/ws/token`: To generate a new websocket auth token.

How to see the result?
===

After running the application , you should call the `/ws/token` API to generate an auth token.

```shell
curl --location --request POST 'localhost:8080/ws/token'
```

response:

```json
{
  "token": "2c5d8f8ddcb1a1c3bbb0761cc757783cfb6428e1dad8cb4a2478f493430d"
}
```

Then, use the `/ws` and connect to the websocket server. After connecting to the websocket server
you can subscribe to the `notifications` channel by using the generated token from previous step by
sending the following message:

```json
{
  "type": "subscribe",
  "params": {
    "channels": [
      "notifications"
    ],
    "token": "2c5d8f8ddcb1a1c3bbb0761cc757783cfb6428e1dad8cb4a2478f493430d"
  }
}
```

As you can see in the following image, after subscribing to the `notifications` channel with an auth token,
we are getting only notifications with userID `81f37eec-76ce-486b-9e40-0d2f0eafddcc`.

![subscribe channels](images/subscribe.png?raw=true "Subscribe")

To unsubscribe a private channel you don't need the token. You can send the following message to the websocket
server to unsubscribe `notifications` channel:

```json
{
  "type": "unsubscribe",
  "params": {
    "channels": [
      "notifications"
    ]
  }
}
```