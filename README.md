# access-api

The Access API is accessible from AAU's network on `knox-proxy01.srv.aau.dk`. You can use that as url, without any ssh connection established.

## Connect to server via ssh

```bash
ssh <student_mail>@knox-proxy01.srv.aau.dk -L <your_port>:localhost:80
```

## Deploy new version manually

Deployment is normally handled by watchtower on push to main. However, in case of the need of manual deployment, run

```bash
docker run -p 0.0.0.0:80:8080 -d -e INTERNAL_KEY=*** -e KNOX_DATABASE_AUTHORIZATION=*** -e API_SECRET=*** ghcr.io/knox-aau/access-api:main
```

## Add new service

Add your service to the list in [services.json](https://github.com/Knox-AAU/access-api/blob/main/services.json) like this:

```json
[
    // other services here ...
    {
      "name": "<your_service_name>",
      "base_url": "http://<your_server_domain>:<your_port>",
      "authorization_key_identifier":"<your_env_key>"
    }
]
```

For examples, check out the actual [services.json](https://github.com/Knox-AAU/access-api/blob/main/services.json) file.

After the changes are merged into main, view the Github Action tab for the status of the deployment. After it succeeds, there can go up to a minute for Watchtower to update the container. After that, your service should be available. See [Use the service](#use-the-service) for more information.

### Authorization of your service

The `authorization_key_identifier` can be ignored if your service does not require authentication. 

This is the name of the key, not the value.

Your service *must* accept the header `Authorization` with the value of the key. If your service does not accept this header, the access api will not be able to access your service.

For the environment variable to be set in the Access API, `docker run` has to be executed again on it's server, with the new environment variable injected into the container. For deployment, see [Deploy new version manually](#deploy-new-version-manually), or contact the authors.

### Use the service

Send a request to `http://knox-proxy01.srv.aau.dk/<your_service_name>`. All parameters and headers will be forwarded to the service. The response will be forwarded to the client. The response will be the same as if you would have sent the request directly to the service.

Header `access-authorization` must be set to an internal key to access the service. Contact the authors of this repository or the current Product Owner of KNOX to get the key.

## Authors

- Casper Bruun Christensen <caschr21@student.aau.dk>
- Emily Treadwell Pedersen <emiped21@student.aau.dk>
- Malthe Reipurth <mreipu21@student.aau.dk>
- Matthias Munch Jakobsen <mattja21@student.aau.dk>
- Moritz Marcus Hönscheidt <mhoens21@student.aau.dk>
- Rasmus Louie Jensen <rjen20@student.aau.dk>
