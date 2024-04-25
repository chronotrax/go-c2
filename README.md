# go-c2
A Go c2 server and agent.

Goose commands to migrate database:
- Upgrade schema: 
```BASH
goose -dir sql/schema sqlite3 ./db.sqlite up 
```
- Downgrade schema:
```BASH
goose -dir sql/schema sqlite3 ./db.sqlite down
```