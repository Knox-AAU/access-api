# access-api

## Connect to server via ssh

```bash
ssh <student_mail>@knox-proxy01.srv.aau.dk -L <your_port>:localhost:80
```

## Deploy new version

```bash
sudo docker run -p 0.0.0.0:80:8080 -d -e INTERNAL_KEY=<internal_key> ghcr.io/knox-aau/access-api:main
```

## Add new service

### Find your service's ip

Connect to the server your service is running on and run:

```bash
hostname -I
```

Output should be something like:

```bash
<your_ip> <something_else>
```

### Update the service list

Create a pull request with the following changes:

- Add your service to the list in `services.json`

```json
[
    // other services here ...
    {
      "name": "<your_service_name>",
      "base_url": "http://<your_ip>"
    }
]
```

### Use the service

[Connect](#connect-to-server-via-ssh) to the server and send a request to `http://localhost:<your_port>/<your_service_name>/<your_service_endpoint>`.

header `access-authorization` must be set to `<internal_key>` to access the service. Contact the authors of this repository or the current Product Owner of KNOX to get the key.

## Authors

- Casper Bruun Christensen <caschr21@student.aau.dk>
- Emily Treadwell Pedersen <emiped21@student.aau.dk>
- Malthe Reipurth <mreipu21@student.aau.dk>
- Matthias Munch Jakobsen <mattja21@student.aau.dk>
- Moritz Marcus Hönscheidt <mhoens21@student.aau.dk>
- Rasmus Louie Jensen <rjen20@student.aau.dk>
