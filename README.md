# access-api

The Access API is accessible from AAU's network on `knox-proxy01.srv.aau.dk`. You can use that as url, without any ssh connection established. 

## Connect to server via ssh

```bash
ssh <student_mail>@knox-proxy01.srv.aau.dk -L <your_port>:localhost:80
```

## Deploy new version manually

Deployment is normally handled by watchtower on push to main. However, in case of the need of manual deployment, run

```bash
sudo docker run -p 0.0.0.0:80:8080 -d -e INTERNAL_KEY=*** -e API_SECRET=*** ghcr.io/knox-aau/access-api:main
```

## Add new service

Before adding your service, you have to make sure that you deploy your service from port `80` on it's server. This means that when you deploy your service, you have to be connected to your server via `<your_port>:localhost:80`, and in your deployment command, you need to specify the ports like `-p 0.0.0.0:80:<your_port>`.

Example:

```bash
ssh <STUDENT_MAIL>@knox-kb01.srv.aau.dk -L <your_port>:localhost:80
docker run -p 0.0.0.0:80:<your_port> --add-host=host.docker.internal:host-gateway -d ghcr.io/knox-aau/databaselayer_server:main
```

### Find your service's ip

Connect to the server your service is running on via ssh and run:

```bash
ifconfig
```

Output includes a list of things. The right one is usually prefixed with `ens160: ` and is the only one that has a normal ip address und `inet`, eg, inet 192.38.54.90.

### Update the service list

Create a pull request to the access api github [repository](https://github.com/Knox-AAU/access-api) with the following changes:

- Add your service to the list in `services.json`

```json
[
    // other services here ...
    {
      "name": "<your_service_name>",
      "base_url": "http://<your_ip>",
      "authorization_key_identifier":"<your_env_key>"
    }
]
```

The `authorization_key_identifier` can be ignored if your service does not require authentication. This is the key, not the value. The value has to be injected into the Access API as environment variable using the `docker run` command. To work, the run command has to be executed again on the server the access api is running on.

### Use the service

Send a request to `http://130.225.57.13/<your_service_name>/<your_service_endpoint>`.

header `access-authorization` must be set to `<internal_key>` to access the service. Contact the authors of this repository or the current Product Owner of KNOX to get the key.

## Authors

- Casper Bruun Christensen <caschr21@student.aau.dk>
- Emily Treadwell Pedersen <emiped21@student.aau.dk>
- Malthe Reipurth <mreipu21@student.aau.dk>
- Matthias Munch Jakobsen <mattja21@student.aau.dk>
- Moritz Marcus Hönscheidt <mhoens21@student.aau.dk>
- Rasmus Louie Jensen <rjen20@student.aau.dk>
