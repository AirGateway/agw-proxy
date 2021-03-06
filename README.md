# agw-proxy
Reverse proxy to use our NDC services with a frontend application without
sharing secrets

## Usage

The proxy works in a way that it receives requests in `/:version/:method` and
proxies them to `:upstream/:version/:method`, adding the `Authorization` header.

Public Docker image is available as `airgateway/agw-proxy` on Docker Hub.

It could be tested directly just by doing:

```
$ docker run -it -p 8000:8000 -e API_KEY=<YOUR_API_KEY> airgateway/agw-proxy
```

or point to production environemnt by doing:

```
$ docker run -it -p 8000:8000 -e API_KEY=<YOUR_API_KEY> API_URL=<API_PRODUCTION_URL> airgateway/agw-proxy
```

### Why is this important?

Because if you are using our API in your frontend application, you **must** use
a proxy backend to avoid, at least, having your API key publicly available at
your frontend's source code.

## Configuration

All the configuration is made using ENV vars.

- **API_KEY** (required): the API key that identifies your company when calling
AirGateway API.

- **API_URL**: the upstream URL to which the calls are redirected. By default
it's `https://proxy.airgateway.net`, our sandbox environment.

- **PORT**: the listening port. By default is 8000.

- **LISTEN**: the listen URL of the server. By default is 0.0.0.0 (from
everywhere).

- **CORS_ORIGIN**: the url (in format `schema://url:port`, like
`http://ndc.myapp.com` or `https://ndc.myapp.com:8081`) from which the requests
are enabled using CORS.

**Note about `LISTEN`**: take care that in most docker environments (Kubernetes,
Docker Compose, etc.) the IP from which the communications are received is not
the real remote IP, as there are network layers between the container and
internet.

## Security concerns

More layers of security can be added. But at least we encourage anyone using
project to:

- **Set up CORS headers**: to avoid receiving calls from other domains.

- **Set up CSP headers**: to avoid calling this service using iframes and similar
attacks.

- **Set up CSRF tokens**: not included here, because CSRF tokens are very tied to
your frontend implementation, but is highly recommended to implement it to
avoid your proxy server being used by third parties.

## Contributions

Contributions are welcome, and the Issues and Pull Requests of this project
remains open to view, evaluate and accept any kind of contributions.

## License

This project, it;s source code and the docker image built from it is licensed
under MIT License by AirGateway GmbH., as can be seen on `LICENSE` file.
