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
docker buildx build --platform linux/amd64,linux/arm64 \
  -f app/Dockerfile \
  --target gateway \
  --push \
  -t 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/gateway:latest \
  .

# auth microservice
docker buildx build --platform linux/amd64,linux/arm64 \
  -f app/Dockerfile \
  --target auth \
  --push \
  -t 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/auth:latest \
  .

# pipeline microservice
docker buildx build --platform linux/amd64,linux/arm64 \
  -f app/Dockerfile \
  --target pipeline \
  --push \
  -t 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/pipeline:latest \
  .

# web frontend
docker buildx build --platform linux/amd64,linux/arm64 \
  -f web/Dockerfile \
  --target production \
  --push \
  -t 483411732137.dkr.ecr.us-east-1.amazonaws.com/nice-dvc-lab/web:latest \
  .
```

```text
[connectivity]
enable-quic-frontend=true
web-x-frame-options="ALLOW-FROM https://dev6166.d2j67odfn9a2fq.amplifyapp.com"
web-extra-http-headers=[("Content-Security-Policy", "frame-ancestors https://dev6166.d2j67odfn9a2fq.amplifyapp.com/")]
```