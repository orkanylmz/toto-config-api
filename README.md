# Toto Config API

Toto Config API helps mobile apps to get a randomly distributed subscription sku based
on some parameters.

# Docs
Documentation is created using our OpenAPI spec

[Documentation](https://orkanylmz.github.io/config-api-docs/)


# Demo
The demo is deployed in GCP (region `europe-west1`)\
Endpoint Url: https://skuconfig-http-cfby56oe6q-ew.a.run.app/api/config

The current architecture is deployed using the `smallest` machine configurations available in GCP for Redis, PostgreSQL and Cloud Run

Also, `authentication` is not implemented for simplicity.

## Benchmark Results
I have run a load test using https://github.com/vearutop/plt

Country Code extraction is disabled with a custom header `X-Custom-CC` for test purposes, because the 3rd party geo API is limiting.
(We can replace it with ours after.)

Command To Run Tests
```go
plt --live-ui --duration=120s --rate-limit=50 curl -X GET "https://skuconfig-http-cfby56oe6q-ew.a.run.app/api/config?package=com.x" -H "X-Custom-CC: TR"
```

![LatencyGraph](https://www.linkpicture.com/q/Screen-Shot-2022-06-22-at-22.53.36.png)

Latency Distributions

![LatencyDistribution](https://www.linkpicture.com/q/Screen-Shot-2022-06-22-at-22.53.49.png)

So, the test runs 120s with a 50 request per second, total request count was 6001 with 0 errors.

# Design
This API is designed with DDD (Domain Driven Design) and CleanArch approach to make it robust, well organized, easy and open to develop. It can be used as a micro service without any further development.
It's free of any third party technologies, and any solution that fits our need can be integrated with ease.

![CleanArch](https://miro.medium.com/max/772/1*B7LkQDyDqLN3rRSrNYkETA.jpeg)
# Caching
In order to serve our API low latency, system will use a caching mechanism. Currently supported integration is Redis.

### Redis Caching Idea
We will use a sorted set with a key name by combining the `Country Code` and `Package Name`, and storing all configurations for this pair with  their `Percentile Max` values being our score for sorted set.
### Example
```go
Country Code: US
Package Name: com.mytestapp.x
Possible SKUs for this pair:
	- 0,  25  => sku1
	- 25, 50  => sku2
	- 50, 75  => sku3
	- 75, 100 => sku4

Our key name for redis will be: us_com.mytestapp.x

We will create our sorted set
	ZADD us_com.mytestapp.x 25 sku1 50 sku2 75 sku3 100 sku4

We can obtain a value from cache with our randomly generate (e.g 33) number as following
	ZRANGEBYSCORE us_com.mytestapp.x 33 +inf LIMIT 0 1

This will return sku2 for the condition 25 < 33 < 50
```

# Request Lifecycle

When the request comes;
 - Our handler first validates the request for query parameter `package`.
 - The `country code` will be obtained by reading user IP address from request headers.
 - Endpoint will create a random number between `0 and 100` and check the cache first, If it finds the configuration, it returns the response.
 - If It can't find in the cache, It will check the DB, if the query satisfies, It will automatically sync the
cache and returns the response. (So further requests can be served directly from cache).

# Geolocation Extraction
Since the API deployed in Cloud Run, It handles the Country extraction from a custom API, since Cloud Run does not provide Country Code to us through headers.

# Deployment
### Running Locally
```go
docker-compose up
```

### Cloud Architecture
- CloudSQL Postgres
- Cloud Run
- Cloud Build
- Cloud Source Repositories
- Memorystore (Redis)
- Docker

Our CloudSQL and MemoryStore will not be publicly exposed so,
Cloud Run instance connects to SQL by using `cloudsqlproxy` and to MemoryStore with custom `vpc` configuration


When we commit and push our code to Cloud Source Repository, Cloud Build triggers and start building our docker image using `cloudbuild.yaml`, after success build, it will deploy it to Cloud Run.

![Architecture](https://i.postimg.cc/JnF1YLCT/toto-arch.png)



### Google Cloud Deployment (Terraform)
You can deploy the whole architecture with just a one command to Google Cloud.

#### Required
* Terraform
* gcloud CLI
* Docker

You need to authorize to google cloud cli first.
```
gcloud auth login
gcloud auth application-default login
```

After

```go
> cd terraform
> make

Fill all required parameters:
    project [current: toto-test]:                   # <----- put your Toto Config API Google Cloud project name here (it will be created)
    user [current: email@gmail.com]:                # <----- put your Google (Gmail, G-suite etc.) e-mail here
    billing_account [current: My billing account]:  # <----- your billing account name, can be found here https://console.cloud.google.com/billing
    region [current: europe-west1]:
    zone [current: europe-west1-b]: 
```

#### Destroy

If you want to tear down the project, run `make destroy`.

If you want to create it again, make sure to:
* Use different project name.
* Remove `terraform.tfstate` file.
