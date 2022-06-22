# Toto Config API

Toto Config API helps the applications to get a randomly distributed subscription skus based
on some parameters.

# Design
This API is designed with DDD (Domain Driven Design) and CleanArch approach to make it robust, well organized, easy and open to develop. It can be used as a micro service without any further development. 
Main idea in the design as following;
It's free of any third party technologies, and any solution that fits our need can be integrated with ease.

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
 - Our handler first validates the request for query parameter `package`
 - The `country code` will be obtained from headers
 - Endpoint will create a random number between `0 and 100` and check the cache first, If it finds the configuration, it returns the response
 - If It can't find in the cache, It will check the DB, if the query satisfies, It will automatically sync the
cache and returns the response. (So further requests can be served directly from cache.)


# Cloud Architecture
- CloudSQL Postgres
- Cloud Run
- Cloud Build
- Cloud Source Repositories
- Memorystore (Redis)
- Docker

Our CloudSQL will not be publicly exposed, our Cloud Run instance connects to it using `cloudsqlproxy`

When we commit and push our code to Cloud Source Repository, Cloud Build triggers and start building our code using `cloudbuild.yaml`, after success build, it will deploy it to Cloud Run.

![Architecture](https://i.postimg.cc/JnF1YLCT/toto-arch.png)

### Running Locally
```go
docker-compose up
```

### Google Cloud Deployment (Terraform)
You can deploy the API and all required parts with just a one command to Google Cloud.

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
