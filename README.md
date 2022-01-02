
``` sh
docker run -d \
  --restart=unless-stopped \
  --name=kuboard-spray \
  -p 80:80/tcp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v /root/kuboard-spray-data:/data \
  eipwork/kuboard-spray:v1.0.0-beta.1
```