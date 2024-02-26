Run `make run`

Server on port 8080


Set up postgres server on linux
Run the following in order, might need to copy the command in the make file and append sudo before idk why can't bypass permission error
`make postgresinit`
`make postgres`
`make createdb` 
Postgres server is in detached mode

Golang-migrate
Create postgres tables, kinda like git for db tables
`make createmigration`
`make migrateup`


Check docker db tables
`make postgres`
`\l`
Connect to go-chat db: `\c go-chat`
TO get items `\d`