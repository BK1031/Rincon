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
Built-in support for health checking and load balancing allows monitoring service status and prevents routing to unavailable services.
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

Services are the core of Rincon. They represent instances of a specific application that you want to enable discovery of. Each service has the following properties:

- `id`: The unique identifier for the service.
- `name`: The name of the service.
- `version`: The version of the service.
- `endpoint`: The endpoint of the service.
- `health_check`: The health check endpoint of the service.

When you register a service, you must provide the name, version, endpoint, and health check endpoint. Upon registration, Rincon will return the following Service object.

```json
{
  "id": 820522,
  "name": "rincon",
  "version": "2.0.0",
  "endpoint": "http://localhost:10311",
  "health_check": "http://localhost:10311/rincon/ping",
  "updated_at": "2024-08-04T19:32:40.109239344-07:00",
  "created_at": "2024-08-04T19:32:40.109239386-07:00"
}
```

The `id` field is an auto-generated integer with the length specified by the `SERVICE_ID_LENGTH` configuration option. It is a unique identifier for a specific instance of a service. Instances are tied to a service by the `name` field. Note that the service name will be converted to lower-snakecase. The `endpoint` must also be unique across all service instances. When updating an existing service, the `endpoint` will be used as the primary identifier.

### Example

Say we register a service `New York` with the following definition:

`[POST] http://localhost:10311/rincon/services`

```json
{
  "name": "New York",
  "version": "1.0.0",
  "endpoint": "http://localhost:3000",
  "health_check": "http://localhost:3000/ping",
}
```

Rincon will return the following service object to us:
```json
{
  "id": 820522,
  "name": "new_york",
  "version": "1.0.0",
  "endpoint": "http://localhost:3000",
  "health_check": "http://localhost:3000/ping",
  "updated_at": "2024-08-04T19:32:40.109239344-07:00",
  "created_at": "2024-08-04T19:32:40.109239386-07:00"
}
```

Any more services that we register using `New York` (will be converted to `new_york`) as the `name` will be considered another instance of the `New York` service.

## Health Checking

Rincon supports health-checking to ensure that registered services are available through heartbeats. Any services determined to be unhealthy will be removed from the registry and no longer discoverable by other services.

### Server Heartbeats

This is the default heartbeat mode, where Rincon will ping each service's `health_check` endpoint at an interval defined by `HEARTBEAT_INTERVAL`.
Any `2xx` status code will mark the heartbeat as successful. Any other response will mark the heartbeat as failed, resulting in the service being removed from the registry.

### Client Heartbeats

When operating in client heartbeat mode, Rincon expects pings from the registered services at least once every heartbeat interval. This ping should be made to the service registration endpoint, with all the service fields correctly populated. The ping can be cofirmed by ensuring the `updated_at` field in the response has increased. At each heartbeat interval, Rincon will remove all services with an `updated_at` timestamp older than the interval.

## Routing

Services register routes to tell Rincon what requests they can handle. Each route has the following properties:

- `id`: The unique identifier for the route.
- `route`: The actual route path.
- `method`: The http methods associated with the route.
- `service_name`: The service associated with the route.

When registering a route, you must provide the route, method, and service associated with that route. Upon registration, Rincon will return the following Route object.

```json
{
  "id": "/rincon/ping-[*]",
  "route": "/rincon/ping",
  "method": "*",
  "service_name": "rincon",
  "created_at": "2024-08-04T01:05:37.262645179-07:00"
}
```

The `id` is a generated field in the format `route-[method]`. Note that routes are tied to services, not any specific instance of a service. So if we register 10 routes to the `new_york` service and then spin up 2 more instances of the `new_york` service, all instances of the `new_york` service are considered able to handle those 10 routes.

### Supported Methods

Rincon currently supports the following HTTP Methods:
- `GET`
- `POST`
- `PUT`
- `DELETE`
- `PATCH`
- `OPTIONS`
- `HEAD`
- `*` (wilcard, route can handle all methods)

### Wildcard Routes

Dynamic routing is supported through the use of wildcards, enabling services to handle variable path segments efficiently.

#### `/*` - Any Wildcard

This wildcard can be used to allow any string to match to a certain path segment.

```
/users/*
---
/users/123
/users/abc
```

You can also use this wildcard in the middle of your route.

```
/users/*/profile
---
/users/123/profile
/users/abc/profile
```

#### `/**` - All Wildcard

This wildcard can be used to allow any string to match to all remaining path segments.

```
/users/**
---
/users/123
/users/abc/edit
/users/a1b2c3/account/settings
```

You can even use both wildcards for more specific routing.

```
/users/*/profile/**
---
/users/123/profile/edit
/users/abc/profile/settings/notifications
```

> [!WARNING]
> While you can have the all wildcard (`**`) in the middle of a route path, when the route graph is computed all proceeding segments are ignored.
> ```
> /users/**/profile
> ---
> /users/123
> /users/abc/profile
> /users/a1b2c3/profile/edit
> ```

### Stacking Routes

Routes with the same path and service name will automatically be "stacked". This just means that their methods will be combined into one registration in Rincon.

```c
New York: /users [GET]
New York: /users [POST]
---
New York: /users [GET,POST]
```

If the existing or new route method is `*`, then the stacked method will simply be the wildcard method.

```c
New York: /users [*]
New York: /users [POST]
---
New York: /users [*]
```

### Conflicting Route Registrations

By default, `OVERWRITE_ROUTES` is set to `false`. This means that if you attempt to register a route that has a conflict with an existing route, it will be rejected. Usually these conflicts arise from two routes attached to different services having an overlap in their methods.

```c
New York: /users [GET]
San Francisco: /users [GET] // cannot be registered
```

Even if the newer route has a higher method coverage than the existing route, the registration will be rejected as long as `OVERWRITE_ROUTES` is set to `false`.

```c
New York: /users [GET]
San Francisco: /users [GET,POST] // cannot be registered
```

To ensure that your routes are registered successfully, make sure there are no method overlaps.

```c
New York: /users [GET]
San Francisco: /users [POST,PUT] // no conflict, will be registered successfully!
```

### Overwriting Routes

When `OVERWRITE_ROUTES` is set to `true`, any conflicting registrations will not be rejected. Instead the new registration will replace the existing one.

```c
New York: /users [GET]
San Francisco: /users [GET]
---
San Francisco: /users [GET]
```

If there are multiple conflicting routes, they will all be replaced.

```c
New York: /users [GET]
Boston: /users [POST]
Los Angelos /users [DELETE]
San Francisco: /users [GET,POST,DELETE]
---
San Francisco: /users [GET,POST,DELETE]
```

> [!CAUTION]
> Existing routes will be replaced even if they have a higher route coverage than the new route. Be careful when overwriting routes!
> ```c
> New York: /users [GET,POST]
> San Francisco: /users [GET]
> ---
> San Francisco: /users [GET] // route for New York was completely removed!
> ```

### Route Matching



### Example

Using our `New York` service from the previous example, let's register the following route.

`[POST] http://localhost:10311/rincon/routes`

```json
{
  "route": "/users",
  "method": "*",
  "service_name": "New York",
  }
```

Rincon will return the following route object to us.

```json
{
  "id": "/users-[*]",
  "route": "/users",
  "method": "*",
  "service_name": "new_york",
  "created_at": "2024-08-27T14:04:43.688527-07:00"
}
```

Now we can confirm our route was correctly registered by making a request to the route matching endpoint.

`[GET] http://localhost:10311/rincon/match?route=users&method=GET`

```json
{
  "id": 416156,
  "name": "new_york",
  "version": "1.0.0",
  "endpoint": "http://localhost:3000",
  "health_check": "http://localhost:3000/health",
  "updated_at": "2024-08-27T14:04:23.172214-07:00",
  "created_at": "2024-08-27T14:04:23.172214-07:00"
}
```

As expected, Rincon returned our `New York` service definition. Now let's try to register a route for a different service (assume that `San Francisco` is a service that has already been registered with Rincon).

`[POST] http://localhost:10311/rincon/routes`

```json
{
  "route": "/users",
  "method": "POST",
  "service_name": "San Francisco",
}
```

This time, we get the following error from Rincon.

```json
{
  "message": "route with id /users-[POST] overlaps with existing routes [[*] /users (new_york)]"
}
```

This is because `New York` was already registered to handle all methods (based on the `*` wildcard method definition) on the `/users` route. By default a new service cannot register a route with a conflicting method. This can be changed by setting `OVERWRITE_ROUTES`.

## Load Balancing

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