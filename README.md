# `raker`
[![Go Reference](https://pkg.go.dev/badge/github.com/AppleGamer22/raker.svg)](https://pkg.go.dev/github.com/AppleGamer22/raker) [![Release](https://github.com/AppleGamer22/raker/actions/workflows/release.yml/badge.svg)](https://github.com/AppleGamer22/raker/actions/workflows/release.yml)
<!-- [![Update Documentation](https://github.com/AppleGamer22/raker/actions/workflows/tag.yml/badge.svg)](https://github.com/AppleGamer22/raker/actions/workflows/tag.yml) -->

<!-- [![Test](https://github.com/AppleGamer22/raker/actions/workflows/test.yml/badge.svg)](https://github.com/AppleGamer22/raker/actions/workflows/test.yml) [![CodeQL](https://github.com/AppleGamer22/raker/actions/workflows/codeql.yml/badge.svg)](https://github.com/AppleGamer22/raker/actions/workflows/codeql.yml) -->

![](assets/icons/raker.svg)

## `raker`'s Logo References
* The Raker [Gopher](https://go.dev/blog/gopher) is based on an original work by [Renee French](https://reneefrench.blogspot.com/).
	* Ueda, T. (2016, July 27). Gopher Vector Graphics. GitHub. <https://github.com/golang-samples/gopher-vector>
* Google. (2022, February 4). Google Chrome Icon. Google. <https://www.google.com/chrome/whats-new/m100>

## Description
`raker` is full-stack and command-line interface for a social media scraper for Instagram, TikTok and VSCO. Both the server and CLI are written in Go, and the web interface is server-side rendered. Both Instagram and TikTok scraping require authentication cookies, which are stored locally wither on the server's MongoDB instance after provided, or on a file system accesible by the CLI.

## Usage Responsibilities
* You should use this software with responsibility and with accordance to [Instagram's terms of use](https://help.instagram.com/581066165581870):
> * **You can't attempt to create accounts or access or collect information in unauthorized ways.**
> This includes creating accounts or collecting information in an automated way without our express permission.
* You should use this software with responsibility and with accordance to [TikTok's terms of use](https://www.tiktok.com/legal/terms-of-use):
> You may not:
> * use automated scripts to collect information from or otherwise interact with the Services;
* You should use this software with responsibility and with accordance to [VSCO's terms of use](https://vsco.co/about/terms_of_use):
> **C Service Rules**  
> You agree not to engage in any of the following prohibited activities:
> * **(I)** copying, distributing, or disclosing any part of the Service in any medium, including without limitation by any automated or non-automated “scraping”,
> * **(II)** using any automated system, including without limitation “robots,” “spiders,” “offline readers,” etc., to access the Service in a manner that sends more request messages to the VSCO servers than a human can reasonably produce in the same period of time by using a conventional on-line web browser (except that VSCO grants the operators of public search engines revocable permission to use spiders to copy materials from vsco.co for the sole purpose of and solely to the extent necessary for creating publicly available searchable indices of the materials but not caches or archives of such materials),
> * **(XI)** accessing any content on the Service through any technology or means other than those provided or authorized by the Service,
> * **(XII)** bypassing the measures we may use to prevent or restrict access to the Service, including without limitation features that prevent or restrict use or copying of any content or enforce limitations on use of the Service or the content therein.

## Installation
### Docker Compose
The `URI` environment variable shown below is suitable for when the [database is also managed by `docker-compose`](https://github.com/AppleGamer22/raker/wiki/Database#docker-compose-yaml). For any other scenario, the URI should be changed accordingly.

```yaml
version: "3"
services:
  raker:
    container_name: raker
    build: .
    environment:
      SECRET: a secret
      URI: mongodb://database:27017
      DATABASE: raker
    ports:
      - 4100:4100
    volumes:
      - /run/media/applegamer22/RPI4HDD/hdd1/.raker/storage:/raker/storage
    depends_on:
      - database
  database:
      container_name: database
      image: mongo:5.0.8
      environment:
        - PUID=1000
        - PGID=1000
      volumes:
        - /run/media/applegamer22/RPI4HDD/hdd1/.raker/database/:/data/db
      ports:
        - 27017:27017
  # migration is work in progress
  # docker exec -it mariadb mariadb -u root -p
  # mariadb:
  #   container_name: mariadb
  #   image: mariadb:11.7.2-ubi9
  #   user: "1000"
  #   ports:
  #     - 3306:3306
  #   volumes:
  #     - ./mariadb:/var/lib/mysql:Z
  #   environment:
  #     TZ: Australia/Melbourne
  #     MARIADB_ROOT_PASSWORD: mariadb
  #     MARIADB_DATABASE: raker
  #     MARIADB_USER: ${USER}
  # docker exec -it postgres psql -U $USER raker
  postgres:
    container_name: postgres
    image: postgres:17.5-alpine3.22
    user: "1000"
    ports:
      - 5432:5432
    volumes:
      - ./postgres:/var/lib/postgresql/data
      - /etc/passwd:/etc/passwd:ro
    environment:
      TZ: Australia/Melbourne
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: raker
      POSTGRES_USER: ${USER}
      POSTGRES_INITDB_ARGS: --data-checksums
```