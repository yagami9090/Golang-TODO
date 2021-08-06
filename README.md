# 2021/08/05/ - 2021/08/06 Golang Programming with Clean Architecture

## 1. Create Network
### 1.1 for create network
```
docker network create --gateway="172.18.0.1" --subnet=172.18.0.0/16 golang
```

## 2. RUN MYSQL AND SETUP DATABASE
### 2.1 for pull image mySQL

```
docker pull mysql
```

### 2.2 run container mySQL (for test version)
```
docker run -d --name mysql --rm -p 8888:3306 -e MYSQL_ROOT_PASSWORD=pass --network=golang mysql
```

### 2.3 follow this command
```
docker exec -it mysql sh
mysql -u root -p
create database test;
```
## 3. BUILD IMAGE AND RUN CONTAINER
### 3.1 for build docker image
```
docker build -t golang:dev ./
```
### 3.2 run container
```
docker run -it --rm --name golang --network=golang -p 9090:9090 golang:dev
```


# etc.

```
go run main.go
go test ./ -cover -v
```
#### for show database detail
```
show databases;
show tables;
```
