# tap-image-factory
Image Factory is a microservice developed to be a part of TAP platform.
The Image Factory component is responsible for building Docker images from tar.gz package containing application binaries or sources (for interpreted languages) provided by the end-user.

## REQUIREMENTS

### Binary
Image Factory requires an access to tap-catalog, tap-blob-store, and rabbitMQ components.

### Compilation
* git (for pulling repository)
* go >= 1.6

## Compilation
To build project:
```
  git clone https://github.com/intel-data/tap-image-factory
  cd tap-image-factory
  make build_anywhere
```
Binaries are available in ./application directory.

## USAGE

To provide IP and port for the application, you have to setup system environment variables:
```
export BIND_ADDRESS=127.0.0.1
export PORT=80
```

As Image Factory depends on external components, you have to set following system environment variables:
* BLOB_STORE_PORT
* BLOB_STORE_HOST
* BLOB_STORE_USER
* BLOB_STORE_PASS
* CATALOG_HOST
* CATALOG_PORT
* CATALOG_USER
* CATALOG_PASS
* QUEUE_PORT
* QUEUE_HOST
* QUEUE_USER
* QUEUE_PASS
* QUEUE_NAME
* HUB_ADDRESS

Image Factory listens to rabbitMQ queue to check if some image needs to be built.

Image Factory endpoints are documented in swagger.yaml file.
Below you can find Image Factory usage.

#### Building Image
Assuming an image entry in Catalog component and its binary in Blob Store components exist, you can build image with:
```
curl -H "Content-Type: application/json" -X POST -d '{"id":"4fcee3a1-201a-4db2-782e-7eae3e654535"}' http://127.0.0.1/api/v1/image --user admin:password
```
After this operation a docker image is built and pushed to docker registry indicated by "HUB_USER" environment variable.
