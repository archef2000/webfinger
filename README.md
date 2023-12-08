# webfinger

[![Docker WebFinger build](https://github.com/Archef2000/webfinger/actions/workflows/main.yml/badge.svg)](https://github.com/Archef2000/webfinger/actions/workflows/main.yml)

[Github](https://github.com/Archef2000/webfinger/)

Explanation: [webfinger.net](https://webfinger.net/)

## docker-compose

```yaml
version: "3"
services:
  webfinger:
    image: archef2000/webfinger:latest
    container_name: webfinger
    ports:
      - 8080:8080/tcp
    environment:
      - PORT=8080
      - LINKS_1_REL=http://openid.net/specs/connect/1.0/issuer
      - LINKS_1_HREF=${issuer URL}
      - LINKS_2_REL=${rel URL}
      - LINKS_2_HREF=${issuer URL}
    restart: unless-stopped
```
