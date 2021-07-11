FROM mongo:4.0

RUN mkdir -p /opt/mongossl

COPY mongod.conf /etc/mongod.conf
COPY ca.crt /opt/mongossl/ca.crt
COPY mongodb.pem /opt/mongossl/mongodb.pem
COPY entrypoint.sh /bin/

RUN chmod +x /bin/entrypoint.sh
