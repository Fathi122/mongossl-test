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
  // TLS config structure
  tlsConfig := &tls.Config{}
  // Update Root ca pool
  roots := x509.NewCertPool()
  if ca, err := ioutil.ReadFile("./mongoinit/ca.crt"); err != nil {
    log.Errorln("Failed to update Root ca pool: ", err.Error())
    return
  }else {
    roots.AppendCertsFromPEM(ca)
    log.Debugln("Updated Root ca pool")
    tlsConfig.RootCAs = roots
  }
  // Load public/private key pair
  if cer, err := tls.LoadX509KeyPair("./mongoinit/server.crt", "./mongoinit/server.key");err != nil {
     log.Errorln("Failed to load public/private keys: ",err.Error())
     return
  }else{
     log.Debugln("Successfully loaded public/private keys")
     tlsConfig.Certificates = []tls.Certificate{cer}
  }
  //connection URL: "mongodb://<username>:<password>@<hostname>:<port>,<hostname>:<port>/<db-name>
  connectionURL := "mongodb://testuser:testuser@localhost:27017/mongodbssl"
  if dialInfo, err := mgo.ParseURL(connectionURL);err != nil{
     log.Errorln("Failed to parse Url: ",err.Error())
  }else{
     dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
       conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
       return conn, err
     }
     // connect to mongodb
    session, err := mgo.DialWithInfo(dialInfo);
    defer func(){
      log.Debugln("Closing mongo session")
      session.Close()
    }()
    if err != nil{
       log.Errorln("Error Connecting to mongodb : ", err.Error())
     }else{
       if session.Ping() == nil{
          log.Debugln("Successfully Connected")
       }
     }
  }
}