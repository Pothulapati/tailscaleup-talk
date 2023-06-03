# Tailscale Up Talk

This repository contains the code for my talk at [Tailscale Up 2023](https://tailscale.dev/up).
It was about Embedding Tailscale into your applications (i.e Android, Web and Server) to solve
the networking problems when you are self-hosting cross-device applications.

Slides: <https://docs.google.com/presentation/d/1I9pI93UroS-iEHjO_OGDP8Aa0_qq2z0QMcdsgXYNHI4/edit?usp=sharing>

## Demo

In this Demo, We will run a TODO App which includes server, and client Golang packages
generated from the [OpenAPI Spec](./swagger.yml) using [Go Swagger](https://github.com/go-swagger/go-swagger).

The end applications would involve a single server component and multiple client components (i.e Android, Web). None of these components expect Tailscale to be installed on the host device.
(It's the whole point of the talk).

### Tailscale and ACL

First, We want to make sure that the server and clients can only talk to each other and not to any other device on the network. This is where Tailscale ACLs come in handy.

Create a tag called `tailtodo` and make sure you add the following ACL to it:

```json
  {
   "action": "accept",
   "src":    ["tag:tailtodo"],
   "dst":    ["tag:tailtodo:*"],
  },
```

### Server

The server is a simple HTTP server which serves the OpenAPI Spec and the Swagger UI. The interesting part is that it uses [tsnet](https://pkg.go.dev/tailscale.com/tsnet) to serve the HTTP server on a Tailscale IP address in a specific Tailnet. This is awesome as it means you can run
multiple services on the same host machine with no Host Name conflicts, etc.

For this to work, You need [OAUTH credentials](https://tailscale.com/kb/1215/oauth-clients/) on your tailnet and you need to set the required environment variables below:

```bash
OAUTH_CLIENT_ID=<> OAUTH_CLIENT_SECRET=<> TAILNET=<> TSNET_FORCE_LOGIN=1 go run cmd/todo-list-server/main.go
```

### Web App

For the Webapp, We need a way to embed tailscale into the website. This is a hard problem obviously
but the folks at tailscale have made it easy by providing a [Web Assembly](https://tailscale.com/kb/1216/embedded/) version of the tailscale client. This is awesome as it means you can embed tailscale into your website and use it to connect to your tailnet.

But by default it only allows `fetch` requests on a Tailscale client. My [fork of the same](./cmd/todo-web/wasm/wasm_js.go) adds support for other POST, PUT that are required for this Todo App.

You can run the web app using the following command:

```bash
cd cmd/todo-web
go run . build && go run . serve
```

Once the webapp is running, You can head over to `localhost:9090` which asks you to login using your Tailscale credentials. Once you login, You can see the TODOs that are fetched from the server.
You can add and delete TODOs as well. For this, The webapp client talks to the previously mentioned server component using the Tailscale IP address. Magical, Right?

### Android App

Me being Me, I'm not an Android developer especially not a one who can write an App that can install, manipulate VPN profiles, TUN interfaces as there are updates in your tailnet. So, I decided to *fork the official tailscale android app* and make it a TODO App. You heard it right.

Follow the instructions in the tailscale-android README to install dependencies but after that, You can run the following command to install the app on your device:

```bash
make rundebug
```

Once the app Opens, It'll ask you to login using your Tailscale credentials. Once you login, You can see the TODOs that are fetched from the server. As you update the TODOs on the webapp, You can see them being updated on the Android app as well. Magical, Right?

## Conclusion

My goal with this talk was to show how Tailscale can now move one level up (from device/host) and
be embedded into your applications. This is a huge deal as it means that your application users can now just login using their Tailscale credentials and start using your application without having to install Tailscale on their devices. Pretty Cool.

I hope this demo was useful to you. If you have any questions, Please feel free to reach out to me on [Twitter](https://twitter.com/tarrooon).
