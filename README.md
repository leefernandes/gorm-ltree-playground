create `.env` from `env.example` 

run it
```sh
export $(cat .env | xargs) && go run .
```