# ODoH Proxy

Modified version of the ODoH server from Cloudflare, with only the proxying functionality retained.

# Local development

To deploy the server locally, first acquire a TLS certificate using [mkcert](https://github.com/FiloSottile/mkcert) as follows:

~~~
$ mkcert -key-file key.pem -cert-file cert.pem 127.0.0.1 localhost
~~~

Then build and run the server as follows:

~~~
$ make all
$ CERT=cert.pem KEY=key.pem PORT=4567 ./odoh-server
~~~

By default, the proxy listens on `/proxy` and the target listens on `/dns-query`.

### Reverse proxy

You need to deploy a reverse proxy with a valid TLS server certificate
for clients to be able to authenticate the target or proxy.

The simplest option for this is using [Caddy](https://caddyserver.com).
Caddy will automatically provision a TLS certificate using ACME from [Let's Encrypt](https://letsencrypt.org).

For instance:

```
caddy reverse-proxy --from https://odoh.example.net:443 --to 127.0.0.1:8080
```

Alternatively, use a Caddyfile similar to:

```
odoh.example.net

reverse_proxy localhost:8080
```
and run `caddy start`.
