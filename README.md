# Tailscale Talk

Running Server:

```
OAUTH_CLIENT_ID=$(op item get tailscaleoauth --fields label=id) OAUTH_CLIENT_SECRET=$(op item get tailscaleoauth --fields label=secret) TAILNET="tarunpothulapati@gmail.com" TSNET_FORCE_LOGIN=1 go run cmd/todo-list-server/main.go 
```
