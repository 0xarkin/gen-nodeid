package main

import (
	"crypto/tls"
	"crypto/x509"
	"strings"
	"fmt"
	"os"

	"github.com/aherve/gopool"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
	log "github.com/sirupsen/logrus"
)

const (
    certsDir = "./certs"
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

func containsI(a string, b string) bool {
	return strings.Contains(
		strings.ToLower(a),
		strings.ToLower(b),
	)
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
		if containsI(nodeID, symbol) {
			log.Info(i, " ", nodeID)
			dir := fmt.Sprintf("%s/%s", certsDir, nodeID)
			os.MkdirAll(dir, 0755)
			saveBytesToFile(certBytes, fmt.Sprintf("%s/cert.crt", dir))
			saveBytesToFile(keyBytes, fmt.Sprintf("%s/cert.key", dir))
		}
	}
}

func main() {
	logLevel := os.Getenv("LOG_LEVEL")
	level, err := log.ParseLevel(logLevel)
	if err != nil {
		log.Warn("couldn't parse log level: %w", err)
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(level)
	}
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.Debug("debug mode enabled")
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
