package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/aherve/gopool"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	log "github.com/sirupsen/logrus"
)

const (
    certsDir = "certs"
)

var WHITELIST = []string {
	"AVAX",
	"ARKIN",
}

func saveBytesToFile(bytes []byte, path string) {
	file, err := os.Create(path)
	if err != nil {
		log.Warn("couldn't create file: %w", err)
	}
	if _, err := file.Write(bytes); err != nil {
		log.Warn("couldn't write file: %w", err)
	}
	if err := file.Close(); err != nil {
		log.Warn("couldn't close file: %w", err)
	}
}

func generateCertificate(i int, pool *gopool.GoPool) {
	defer pool.Done()

	if i%10000 == 0 {
		log.Info(i, " certificates generated")
	}

	certBytes, keyBytes, err := staking.NewCertAndKeyBytes()
	if err != nil {
		log.Warn("couldn't generate certificate: %w", err)
		return
	}
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		log.Warn("couldn't parse certificate: %w", err)
		return
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		log.Warn("couldn't parse certificate: %w", err)
		return
	}
	nodeIDFromCert := ids.NodeIDFromCert(cert.Leaf)
	nodeID := nodeIDFromCert.String()
	log.Debug(i, " ", nodeID)
	for _, symbol := range WHITELIST {
		if nodeID == symbol {
			log.Info(nodeID)
			os.MkdirAll(fmt.Sprintf("certs/%s", nodeID), os.ModePerm)
			saveBytesToFile(certBytes, fmt.Sprintf("%s/%s/file.cert", certsDir, nodeID))
			saveBytesToFile(keyBytes, fmt.Sprintf("%s/%s/file.key", certsDir, nodeID))
		}
	}
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
		
	})
	log.Info("starting programm...")

	i := 0
	pool := gopool.NewPool(8)
	for {
		pool.Add(1)
		go generateCertificate(i, pool)
		i++
		// better with a stop condition
	}
	// pool.Wait() // never reached (infinite loop)
}