# errors

[![Build Status](https://drone.owncloud.com/api/badges/owncloud-ops/errors/status.svg)](https://drone.owncloud.com/owncloud-ops/errors)

This service provides handlers for a default backend which can be used within Kubernetes Ingress controllers like Nginx. It displayes nice error pages and it is customizable on all aspects.

## Environment Variables

```Shell
# Path to optional config file
ERRORS_CONFIG_FILE=
# Set logging level
ERRORS_LOG_LEVEL=info
#Enable colored logging
ERRORS_LOG_COLOR=true
# Enable pretty logging
ERRORS_LOG_PRETTY=true

# Address to bind the metrics
ERRORS_METRICS_ADDR=0.0.0.0:8081
# Token to make metrics secure
ERRORS_METRICS_TOKEN=

# Address to bind the server
ERRORS_SERVER_ADDR=0.0.0.0:8080
# Enable pprof debugging
ERRORS_SERVER_PPROF=false
# Root path of the server
ERRORS_SERVER_ROOT=/
# External access to server
ERRORS_SERVER_HOST=http://localhost:8080
# Path to cert for SSL encryption
ERRORS_SERVER_CERT=
# Path to key for SSL encryption
ERRORS_SERVER_KEY=
# Use strict SSL curves
ERRORS_SERVER_STRICT_CURVES=false
# Use strict SSL ciphers
ERRORS_SERVER_STRICT_CIPHERS=false
# Folder for custom templates
ERRORS_SERVER_TEMPLATES=
# Path for overriding errors
ERRORS_SERVER_ERRORS=
```

## Ports

- 8080
- 8081

## Build

Make sure you have a working Go environment, for further reference or a guide take a look at the [install instructions](https://golang.org/doc/install.html). This project requires Go >= v1.18.

```Shell
git clone https://github.owncloud.com/owncloud-ops/errors.git
cd errors

make generate build
./bin/errors --help
```

To build the container use:

```Shell
docker build -f Dockerfile -t errors:latest .
```

## License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.
