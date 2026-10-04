# Quick Start

## Docker Run

Great for testing out the application, but not recommended for stable use. Checkout the docker-compose for the recommended deployment.

For each image there are two tags, respectively the regular tag and $TAG-rootless, which uses a non-root image.

```sh
# If using the rootless image, ensure data
# folder has correct permissions
$ mkdir -p /path/to/data/folder
$ chown 65532:65532 -R /path/to/data/folder
# ---------------------------------------
# Run the image
$ docker run -d \
  --name catalog \
  --restart unless-stopped \
  --publish 3100:7745 \
  --env TZ=Europe/Bucharest \
  --volume /path/to/data/folder/:/data \
  ghcr.io/sdmad737/catalog:latest
# ghcr.io/sdmad737/catalog:latest-rootless

```

## Docker-Compose

```yaml
services:
  catalog:
    image: ghcr.io/sdmad737/catalog:latest
#   image: ghcr.io/sdmad737/catalog:latest-rootless
    container_name: catalog
    restart: always
    environment:
    - CATALOG_LOG_LEVEL=info
    - CATALOG_LOG_FORMAT=text
    - CATALOG_WEB_MAX_UPLOAD_SIZE=10
    volumes:
      - catalog-data:/data/
    ports:
      - 3100:7745

volumes:
   catalog-data:
     driver: local
```

!!! note
    If you use the `rootless` image, and instead of using named volumes you would prefer using a hostMount directly (e.g., `volumes: [ /path/to/data/folder:/data ]`) you need to `chown` the chosen directory in advance to the `65532` user (as shown in the Docker example above).

## Env Variables & Configuration

| Variable                             | Default                | Description                                                                        |
| ------------------------------------ | ---------------------- | ---------------------------------------------------------------------------------- |
| CATALOG_MODE                            | production             | application mode used for runtime behavior  can be one of: development, production |
| CATALOG_WEB_PORT                        | 7745                   | port to run the web server on, if you're using docker do not change this           |
| CATALOG_WEB_HOST                        |                        | host to run the web server on, if you're using docker do not change this           |
| CATALOG_OPTIONS_ALLOW_REGISTRATION      | true                   | allow users to register themselves                                                 |
| CATALOG_OPTIONS_AUTO_INCREMENT_ASSET_ID | true                   | auto increments the asset_id field for new items                                   |
| CATALOG_OPTIONS_CURRENCY_CONFIG         |                        | json configuration file containing additional currencies                           |
| CATALOG_WEB_MAX_UPLOAD_SIZE             | 10                     | maximum file upload size supported in MB                                           |
| CATALOG_WEB_READ_TIMEOUT                | 10                     | read timeout of HTTP server                                                         |
| CATALOG_WEB_WRITE_TIMEOUT               | 10                     | write timeout of HTTP server                                                        |
| CATALOG_WEB_IDLE_TIMEOUT                | 30                     | idle timeout of HTTP server                                                         |
| CATALOG_STORAGE_DATA                    | /data/                 | path to the data directory, do not change this if you're using docker              |
| CATALOG_STORAGE_SQLITE_URL              | /data/catalog.db?_fk=1 | sqlite database URL, if you're using docker do not change this                     |
| CATALOG_LOG_LEVEL                       | info                   | log level to use, can be one of: trace, debug, info, warn, error, critical         |
| CATALOG_LOG_FORMAT                      | text                   | log format to use, can be one of: text, json                                       |
| CATALOG_MAILER_HOST                     |                        | email host to use, if not set no email provider will be used                       |
| CATALOG_MAILER_PORT                     | 587                    | email port to use                                                                  |
| CATALOG_MAILER_USERNAME                 |                        | email user to use                                                                  |
| CATALOG_MAILER_PASSWORD                 |                        | email password to use                                                              |
| CATALOG_MAILER_FROM                     |                        | email from address to use                                                          |
| CATALOG_SWAGGER_HOST                    | 7745                   | swagger host to use, if not set swagger will be disabled                           |
| CATALOG_SWAGGER_SCHEMA                  | http                   | swagger schema to use, can be one of: http, https                                  |

!!! tip "CLI Arguments"
      If you're deploying without Docker, use command-line arguments to configure the application. Run the compiled API binary with `--help` for more information.

      ```sh
      Usage: api [options] [arguments]

      OPTIONS
        --mode/$CATALOG_MODE                                                        <string>  (default: development)
        --web-port/$CATALOG_WEB_PORT                                                <string>  (default: 7745)
        --web-host/$CATALOG_WEB_HOST                                                <string>
        --web-max-upload-size/$CATALOG_WEB_MAX_UPLOAD_SIZE                          <int>     (default: 10)
        --storage-data/$CATALOG_STORAGE_DATA                                        <string>  (default: ./.data)
        --storage-sqlite-url/$CATALOG_STORAGE_SQLITE_URL                            <string>  (default: ./.data/catalog.db?_fk=1)
        --log-level/$CATALOG_LOG_LEVEL                                              <string>  (default: info)
        --log-format/$CATALOG_LOG_FORMAT                                            <string>  (default: text)
        --mailer-host/$CATALOG_MAILER_HOST                                          <string>
        --mailer-port/$CATALOG_MAILER_PORT                                          <int>
        --mailer-username/$CATALOG_MAILER_USERNAME                                  <string>
        --mailer-password/$CATALOG_MAILER_PASSWORD                                  <string>
        --mailer-from/$CATALOG_MAILER_FROM                                          <string>
        --swagger-host/$CATALOG_SWAGGER_HOST                                        <string>  (default: localhost:7745)
        --swagger-scheme/$CATALOG_SWAGGER_SCHEME                                    <string>  (default: http)
        --demo/$CATALOG_DEMO                                                        <bool>
        --debug-enabled/$CATALOG_DEBUG_ENABLED                                      <bool>    (default: false)
        --debug-port/$CATALOG_DEBUG_PORT                                            <string>  (default: 4000)
        --options-allow-registration/$CATALOG_OPTIONS_ALLOW_REGISTRATION            <bool>    (default: true)
        --options-auto-increment-asset-id/$CATALOG_OPTIONS_AUTO_INCREMENT_ASSET_ID  <bool>    (default: true)
        --options-currency-config/$CATALOG_OPTIONS_CURRENCY_CONFIG                  <string>
        --help/-h
        display this help message
      ```
