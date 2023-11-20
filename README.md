# access-api

## Connect to server via ssh

```bash
ssh <student_mail>@knox-proxy01.srv.aau.dk -L <your_port>:localhost:80
```

## Deploy new version
```bash
sudo docker run -p 0.0.0.0:80:8080 -d -e INTERNAL_KEY=*** ghcr.io/knox-aau/access-api:main
```
