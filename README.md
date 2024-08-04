# Rincon

<img align="right" width="159px" src="https://github.com/BK1031/Rincon/blob/main/assets/rincon-circle.png?raw=true" alt="rincon-logo">

[![Build Status](https://github.com/BK1031/Rincon/actions/workflows/test.yml/badge.svg)](https://github.com/BK1031/Rincon/actions/workflows/test.yml)
[![codecov](https://codecov.io/gh/BK1031/Rincon/graph/badge.svg?token=R4NMABYGOZ)](https://codecov.io/gh/BK1031/Rincon)
[![GoDoc](https://pkg.go.dev/badge/github.com/bk1031/rincon?status.svg)](https://pkg.go.dev/github.com/bk1031/rincon?tab=doc)
[![Docker Pulls](https://img.shields.io/docker/pulls/bk1031/rincon?style=flat-square)](https://hub.docker.com/repository/docker/bk1031/rincon)
[![Release](https://img.shields.io/github/release/bk1031/rincon.svg?style=flat-square)](https://github.com/bk1031/rincon/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


Rincon is a cloud-native service registry written in [Go](https://go.dev/). 
It is designed to be fast, lightweight, and highly scalable.
Rincon is also platform-agnostic, and can run in the cloud, a container, or even on bare-metal,
making it perfect for both local development and production environments.

Rincon makes it easy for services to register themselves and to discover other services.
Built-in support for health checking allows monitoring service status and prevents routing to unavailable services.
External services such as SaaS vendors can also be registered to create a unified discovery interface.

## Getting Started

The easiest way to get started with Rincon is to use the official Docker image.
You can pull it from [Docker Hub](https://hub.docker.com/r/bk1031/rincon).

```bash
$ docker run -d -p 10311:10311 bk1031/rincon:latest
```

Alternatively if you have an existing compose file, you can add Rincon as a service.
This way you can easily connect Rincon to your existing database.

```yml
rincon:
    image: bk1031/rincon:latest
    restart: unless-stopped
    environment:
      PORT: "10311"
      STORAGE_MODE: "sql"
      DB_DRIVER: "postgres"
      DB_HOST: "localhost"
      DB_PORT: "5432"
      DB_NAME: "rincon"
      DB_USER: "postgres"
      DB_PASSWORD: "password"
    ports:
      - "10311:10311"
```

By default Rincon will run on port `10311`. Once Rincon is running, we can connect to it from our application using the default username and password `admin`. Let's register a service called `Service A`.

```
curl -X "POST" "http://localhost:10311/rincon/services" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:admin' \
     -d $'{
  "endpoint": "localhost:8080",
  "name": "Service A",
  "health_check": "localhost:8080/health",
  "version": "1.0.0"
}'
```

The response body will contain an ID for our newly registered service. Including this ID in future requests will allow you to update the service's registration. Now we can register a route that `Service A` will handle.

```
curl -X "POST" "http://localhost:10311/rincon/routes" \
     -H 'Content-Type: application/json; charset=utf-8' \
     -u 'admin:admin' \
     -d $'{
  "route": "/service",
  "service_name": "Service A",
  "method": "*"
}'
```

Now we can verify that our service and route are properly registered by making a request to the route matching endpoint.

```
curl "http://localhost:10311/rincon/match?route=service&method=GET"
```

You should see the service registration for `Service A` returned in the response body. We have now successfully registered a service and route!

> [!TIP]
> If you're using Go, check out our client sdk [here](https://github.com/bk1031/rincon-go).

## Services

## Health Checking

## Routing

## Configuration

## API Endpoints

## Roadmap

## Contributing

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b gh-username/my-amazing-feature`)
3. Commit your Changes (`git commit -m 'Add my amazing feature'`)
4. Push to the Branch (`git push origin gh-username/my-amazing-feature`)
5. Open a Pull Request

## License

Distributed under the MIT License. See `LICENSE.txt` for more information.