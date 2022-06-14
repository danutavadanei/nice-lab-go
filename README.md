### Deploy stack using CF

awslocal cloudformation deploy --stack-name kali-nice-dcv --template-file "./cf/kali-nice-dcv.yaml"

### Migrate db
```shell
docker compose -f docker-compose.migrate.yml up -d
docker compose -f docker-compose.migrate.yml exec migration bash
/usr/go/bin/sql-migrate up -config="./migrations/dbconfig.yml"
```

### Build and push images to AWS ECR
```shell
# gateway microservice
docker build -f app/Dockerfile --target gateway -t nice-dvc-lab/gateway .
docker tag nice-dvc-lab/gateway:latest 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/gateway:latest
docker push 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/gateway:latest

# auth microservice
docker build -f app/Dockerfile --target auth -t nice-dvc-lab/auth .
docker tag nice-dvc-lab/auth:latest 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/auth:latest
docker push 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/auth:latest

# pipeline microservice
docker build -f app/Dockerfile --target pipeline -t nice-dvc-lab/pipeline .
docker tag nice-dvc-lab/pipeline:latest 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/pipeline:latest
docker push 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/pipeline:latest
```