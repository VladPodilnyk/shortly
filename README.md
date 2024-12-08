### Shortly

Shortly is an open-sourced, free url-shortener service.
Currently, the service is as simple as it gets and doesn't
provide any _fancy_ features except shortening URLs :laughing:

To request a feature or to contribute one, please look at the [contribution guide](#contributing)

#### Internals
Shortly is leveraging Golang rich std library where possible and has a few dependencies.
Major dependencies are `viper` for service configurations and `mongo-driver` for MongoDB connector.
The app has the following structure (unnecessary files are omitted):
```
.
├── Makefile
├── cmd/shortly - contains the app entry point and public/static assets (for frontend)
├── compose.yml
├── deploy - contains a docker file definition
├── internal - that's where the implementation lies
│   ├── app - contains controllers, services and middleware
│   ├── config - contains config readers
│   ├── encoder - contains the core of the service that produces short URLs
│   ├── helpers - helper functions
│   ├── models - contains various data structures
│   ├── storage - defines storage interfaces
│   └── tests
└── scripts - contains scripts for routine tasks
```

Currently, Shortly supports only MongoDB as storage. The encoding algorithm is based on Base58 alphabet.
If you are curious about implementation details, please take a look at `./encoder/encoder.go`.

#### Running locally
There are a few ways to play with Shortly on your local machine.
The first one is a bit complicated and requires from you to have `Golang toolchain`,
`make` and `tailwindcss` to be installed. In case this works for you, use `make help`
to list available commands and actions in the repo.

The second option is to use `deploy/Dockerfile`. This one requires from you to have
`docker` on your local machine.
To build an image use the following command
```
docker build -f ./deploy/Dockerfile --tag shortly .
```
After this you'll be able to run a container.

__Service configuration__
Important note here is that to run Shortly you have to provide a service config.
Shortly supports JSON configs and environment variables.
When you run the service using `make run-dev` Shortly expects to have `.application.dev.json`
file in the root of the `shortly` folder. In case you start a binary yourself, it's possible to
supply a config file like that `./shortly --config <your-config-file>`.


In case you use a docker image, you must pass environment variables or `.env file`
```
docker run --env-file /path/to/env/file
docker run -d -e PORT='8080' -e PREFIX="awesome-domain.io"
```
Example `.env` config
```env
MONGO_DB_URI="YOUR-URI"
MONGO_DB_NAME="test"
MONGO_DB_COLLECTION="test"
PREFIX="awesome-domain.io" # this is essentially the domain you own, the final URL will look like http://cool/sdlfjkd
ENVIRONMENT="stage" # whatever string you want, as it's more for a health check purposes
PORT=4000 # a port to which a Golang server is bind.
ALIAS_MAX_SIZE=6 # the max number of character for alias strings
REQUEST_PER_MINUTE=10 # configures rage limiter
```

Example JSON config
```json
{
    "server": {
        "port": "8080"
    },
    "storage": {
        "uri": "your-uri",
        "name": "test",
        "collection": "test"
    },
    "env": "stage",
    "prefix": "awesome-domain.io",
    "alias_max_size": 6,
    "requests_per_minute": 50
}
```


#### Self-host 🚧
Now, frankly speaking, Shortly is not very customizable, but it's already possible
to self host it.
The best way to do that is use `deploy/Dockerfile` that can be published to whatever
container registry you want or deployed to a machine that you have. Please refer to
section [Running locally](#running-locally) for service configuration instructions.

#### Contributing
The process is very simple. If you have an idea or a feature that you would like to have
please request it using Github issues. PRs are welcome too 🙌
