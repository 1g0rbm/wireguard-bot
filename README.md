# wireguard bot

Бот для создания и управления аккаунтами VPN по протоколу wireguard.

### Requirements
1. go v1.23
2. docker
3. docker-compose

### Install
```sh
# install binary deps
make install-deps
```
```sh
# создать .env, заменить значения на свои
cp .env.example .env
``` 
```sh
# install modules
make update-modules
```

### Run
```sh
# запустить инфраструктуру(бд и пр.)
make run-infra
```
```sh
# накатить миграции
make migration-up
```
```sh
# запустить бот
make run-bot
```
