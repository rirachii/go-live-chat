Run `make`

starts server on port 8080


Set up postgres server on linux
1. Download docker first then run:

`sudo docker run --name postgres15 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:15-alpine`

Postgres server is in detached mode

Golang-migrate
Create postgres tables, kinda like git for db tables

Check db tables
`make postgre`
TO list
`\l`
Connect to go-chat db: `\c go-chat`
Then `\d`