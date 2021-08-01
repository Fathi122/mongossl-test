package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

func main() {
	if len(os.Args) != 2 {
		panic("ERROR: Expecting one argument")
	}
	deployment := os.Args[1]
	mongoCACrt := "./mongoinit/ca.crt"
	mongoPubKey := "./mongoinit/server.crt"
	mongoPrivKey := "./mongoinit/server.key"
	//connection URL: "mongodb://<username>:<password>@<hostname>:<port>,<hostname>:<port>/<db-name>
	connectionURL := "mongodb://testuser:testuser@localhost:27017/mongodbssl"
	if deployment == "k8s" {
		mongoCACrt = "/opt/certs/ca.crt"
		mongoPubKey = "/opt/certs/server.crt"
		mongoPrivKey = "/opt/certs/server.key"
		connectionURL = "mongodb://testuser:testuser@mongo.default.svc.cluster.local:27017/mongodbssl"
	}
	// set log level
	log.SetLevel(log.DebugLevel)
	// TLS config structure
	tlsConfig := &tls.Config{}
	// Update Root ca pool
	roots := x509.NewCertPool()
	if ca, err := ioutil.ReadFile(mongoCACrt); err != nil {
		log.Errorln("Failed to update Root ca pool: ", err.Error())
		return
	} else {
		roots.AppendCertsFromPEM(ca)
		log.Debugln("Updated Root ca pool")
		tlsConfig.RootCAs = roots
	}
	// Load public/private key pair
	if cer, err := tls.LoadX509KeyPair(mongoPubKey, mongoPrivKey); err != nil {
		log.Errorln("Failed to load public/private keys: ", err.Error())
		return
	} else {
		log.Debugln("Successfully loaded public/private keys")
		tlsConfig.Certificates = []tls.Certificate{cer}
	}
	if dialInfo, err := mgo.ParseURL(connectionURL); err != nil {
		log.Errorln("Failed to parse Url: ", err.Error())
	} else {
		dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
		dialInfo.Timeout = 60 * time.Minute
		// connect to mongodb
		session, err := mgo.DialWithInfo(dialInfo)
		defer func() {
			log.Debugln("Closing mongo session")
			if session != nil {
				session.Close()
			}
		}()
		if err != nil {
			log.Errorln("Error Connecting to mongodb : ", err.Error())
		} else {
			if session.Ping() == nil {
				log.Debugln("Successfully Connected")
			}
		}
	}
}
