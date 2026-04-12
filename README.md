# wireguard bot

Бот для создания и управления аккаунтами VPN по протоколу wireguard.

### Requirements
1. [mise](https://mise.jdx.dev/)
2. docker
3. docker-compose

### Install
```sh
# install tools from mise.toml (Go and others)
mise install
```
```sh
# создать .env, заменить значения на свои
cp .env.example .env
``` 
```sh
# install modules
mise run update-modules
```

### Run
```sh
# запустить инфраструктуру(бд и пр.)
mise run run-infra
```
```sh
# накатить миграции
mise run migration-up
```
```sh
# запустить бот
mise run run-bot
```
