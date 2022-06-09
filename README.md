### Deploy stack using CF

awslocal cloudformation deploy --stack-name kali-nice-dcv --template-file "./cf/kali-nice-dcv.yaml"

### Migrate db
```shell
docker compose -f docker-compose.migrate.yml up -d
docker compose -f docker-compose.migrate.yml exec migration bash
/usr/go/bin/sql-migrate up -config="./migrations/dbconfig.yml"
```