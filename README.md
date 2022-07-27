# channelize

![build](https://github.com/hmdsefi/channelize/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/hmdsefi/channelize)](https://goreportcard.com/report/github.com/hmdsefi/channelize)
[![codecov](https://codecov.io/gh/hmdsefi/channelize/branch/master/graph/badge.svg?token=6IUFW3MADN)](https://codecov.io/gh/hmdsefi/channelize)
[![Go Reference](https://pkg.go.dev/badge/github.com/hmdsefi/channelize.svg)](https://pkg.go.dev/github.com/hmdsefi/channelize)

Channelize is a websocket framework based on gorilla websocket that manages inbound and outbound streams. It gives the client the option of
receiving events in different channels. Channelize is useful when your websocket server is reading different types
of events from a message broker. e.g., you might have multiple kafka topics or NATs subjects.

Channelize gives you this option to categorize the outbound messages. You can choose a name for each event type
and register them as a channel in the Channelize and send messages to those channels. Then client can subscribe
to one or more channels to receive the streams.

The main usage of Channelize is in cryptocurrency exchanges, when you have multiple channels like user's balances, 
user's orders, user's trades, market trades, mini-ticker data, etc. The front-end should be able to subscribe to 
the related channels in each page and unsubscribe those that doesn't need any more.

![channelize](images/channelize.png?raw=true "channelize")

## Table of Contents

* [Install](#Install)
* [How to use](#How-to-use)
    * [Public channels](#Public-channels)
    * [Private channels](#Private-channels)
* [Metrics](#Metrics)
* [Examples](https://github.com/hmdsefi/channelize/tree/master/_examples)

### Install

To use Channelize, first you should use `go get` command to get the latest version of this library:

```shell
go get -u github.com/hmdsefi/channelize
```

Next step is importing channelize to your project:

```go
import "github.com/hmdsefi/channelize"
```

### How to use

To use Channelize, first you should know about the channel concept. A channel is a unique name for a stream of events
that have the same type. A channel can be public or private. A public channel sends the outbound messages to the all
available connections that subscribed to that channel. A private channel needs authentication, it sends the outbound
messages to a specific connection that already subscribed to that private channel with a valid token.

#### Public channels

To use public channels, first you should register your public channels with one of the following functions:

```go
channel := channelize.RegisterPublicChannel("my-public-channel")
```

```go
channels := channelize.RegisterPublicChannels("my-public-channel1", "my-public-channel2")
```

Registering same channel more than once won't break anything. It will override the previous one.
To send messages to the channels, you should create an instance of Channelize struct to be able to use the library
APIs.

```go
chlz := channelize.NewChannelize()
```

Then you can call the following function in your consumer function to send the messages to the proper channel:

```go
err := chlz.SendPublicMessage(ctx, channel, message)
if err != nil {
return err
}
```

Channelize struct gives you this option to create your handler or using the handler that channelize creates for you.
If you want to implement the handler by yourself, then you can use the following method:

```go
connection := chlz.CreateConnection(ctx context.Context, wsConn *websocket.Conn, options ...conn.Option)
```

Or you can create the handler by calling the following method:

```go
handler := chlz.MakeHTTPHandler(appCtx context.Context, upgrader websocket.Upgrader, options ...conn.Option)
```

To subscribe to one or more public channels, client should send the following message to the server:

```json
{
  "type": "subscribe",
  "params": {
    "channels": [
      "my-public-channel1",
      "my-public-channel2"
    ]
  }
}
```

To unsubscribe one or more channels, client should send the following message to the server:

```json
{
  "type": "unsubscribe",
  "params": {
    "channels": [
      "my-public-channel1"
    ]
  }
}
```

#### Private channels

To use private channels, first you should register your private channels with one of the following functions:

```go
channel := channelize.RegisterPrivateChannel("my-private-channel")
```

```go
channels := channelize.RegisterPrivateChannels("my-private-channel1", "my-private-channel2")
```

Private channels need authentication. To provide authentication you should implement the function type that is defined
in `auth` package:

```go
type AuthenticateFunc func (token string) (*Token, error)
```

Then you should pass as an option to the Channelize constructor:

```go
chlz := channelize.NewChannelize(channelize.WithAuthFunc(MyAuthFunc)))
```

You can use `CreateConnection` or `MakeHTTPHandler` to create the connection for the client just like public channels.
To send the message to the client you should use the following function:

```go
err := chlz.SendPrivateMessage(ctx, channel, userID, message)
if err != nil {
return err
}
```

To subscribe to private channels, client should fill the token field with the proper value:

```json
{
  "type": "subscribe",
  "params": {
    "channels": [
      "my-public-channel1",
      "my-private-channel1"
    ],
    "token": "618bb5b00161cbd68bc744b2ea84c96601d6705f31cc7d32e01c3371f6e7"
  }
}
```

To unsubscribe, client can use the following message:

```json
{
  "type": "unsubscribe",
  "params": {
    "channels": [
      "my-private-channel1"
    ]
  }
}
```

### Metrics

You can find the following prometheus metrics in Channelize:

| METRIC                             | TYPE  | DESCRIPTION                                    |
|------------------------------------|-------|------------------------------------------------|
| open_connections                   | gauge | Number of open connections.                    |
| private_connections                | gauge | Number of private connections.                 |
| private_connections_storage_length | gauge | Number of stored private connections.          |
| open_connections_storage_length    | gauge | Number of stored open connections.             |
| subscribed_channels_storage_length | gauge | Number of subscribed channels that are stored. |

## License

MIT License, please see [LICENSE](https://github.com/hmdsefi/channelize/blob/master/LICENSE) for details.
