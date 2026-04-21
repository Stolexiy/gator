# gator

## Init database
```
cd sql/schema
goose postgres "postgres://postgres:postgres@localhost:5432/gator" up
```

## Generate db connector
`sqlc generate`
