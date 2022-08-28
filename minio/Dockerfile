
# FROM  minio/minio:latest
FROM minio/minio:RELEASE.2022-04-12T06-55-35Z.fips

# WORKDIR /app

RUN curl -L -o certgen-linux-amd64 https://github.com/minio/certgen/releases/download/v0.0.2/certgen-linux-amd64

RUN chmod +x certgen-linux-amd64
# RUN ls ~/


# replace with your IP or hostname
RUN  ./certgen-linux-amd64 -ca -host "minio,localhost,127.0.0.1,172.77.0.45"  
#if you get below error, replace last ip/ 172.77.0.45--> with you minio ip from docker logs
#x509: cannot validate certificate for <ip> because it doesn't contain any IP SANs

RUN  mkdir -p ~/.minio/certs/
RUN  mv private.key public.crt ~/.minio/certs/
# RUN sudo systemctl restart minio