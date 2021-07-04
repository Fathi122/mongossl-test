# fluentd/Dockerfile
FROM mongo:4.0

#COPY createDBUser.js /docker-entrypoint-initdb.d/
COPY mongod.conf /etc/mongod.conf
RUN mkdir -p /opt/mongossl
COPY ca.crt /opt/mongossl/ca.crt
COPY mongodb.pem /opt/mongossl/mongodb.pem
#COPY entrypoint.sh /etc/entrypoint.sh

COPY entrypoint.sh /bin/
RUN chmod +x /bin/entrypoint.sh
