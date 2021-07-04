package main

import(  
  "crypto/tls"
  "crypto/x509"
  "io/ioutil"
  "net"
  "gopkg.in/mgo.v2"
  log "github.com/sirupsen/logrus"
)

func main(){
  // set log level
  log.SetLevel(log.DebugLevel)
  // Update Root ca pool
  roots := x509.NewCertPool()
  if ca, err := ioutil.ReadFile("./mongoinit/ca.crt"); err == nil { 
    roots.AppendCertsFromPEM(ca)
    log.Debugln("Updated Root ca pool")
  }else {
    log.Errorln("Failed to update Root ca pool: ", err) 
    return
  }
  tlsConfig := &tls.Config{}
  tlsConfig.RootCAs = roots
  // Load publi/private key pair
  if cer, err := tls.LoadX509KeyPair("./mongoinit/server.crt", "./mongoinit/server.key");err != nil {
     log.Errorln("Failed to load public/private keys: ",err)
     return
  }else{
     log.Debugln("Successfully loaded public/private keys")
     tlsConfig.Certificates = []tls.Certificate{cer}
  }
  //connection URL: "mongodb://<username>:<password>@<hostname>:<port>,<hostname>:<port>/<db-name>
  connectionURL := "mongodb://testuser:testuser@localhost:27017/mongodbssl"
  if dialInfo, err := mgo.ParseURL(connectionURL);err != nil{
     log.Errorln("Failed to parse Url: ",err)
  }else{
     dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
       conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
       return conn, err
     }
     // connect to mongodb
     if session, err := mgo.DialWithInfo(dialInfo);err != nil{
       log.Errorln("Error Connecting to mongodb : ", err)
     }else{
       log.Debugln("Successfully Connected : ", session) 
     }
  }
}