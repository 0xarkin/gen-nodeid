package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"github.com/aherve/gopool"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/staking"
)

const (
    certsDir = "certs"
)

var WHITELIST = []string {
	"AVAX",
	"ARKIN",
}

func saveBytesToFile(bytes []byte, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("couldn't create file: %w", err)
	}
	if _, err := file.Write(bytes); err != nil {
		return fmt.Errorf("couldn't write file: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("couldn't close file: %w", err)
	}
	return nil
}

func generateCertificate(i int, pool *gopool.GoPool) {
	defer pool.Done()

	if i%10000 == 0 {
		fmt.Println(i, "certificates generated")
	}

	certBytes, keyBytes, err := staking.NewCertAndKeyBytes()
	if err != nil {
		return
	}
	cert, err := tls.X509KeyPair(certBytes, keyBytes)
	if err != nil {
		return
	}
	cert.Leaf, err = x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}
	nodeIDFromCert := ids.NodeIDFromCert(cert.Leaf)
	nodeID := nodeIDFromCert.String()
	for _, symbol := range WHITELIST {
		if nodeID == symbol {
			fmt.Println("OK", nodeID, i)
			os.MkdirAll(fmt.Sprintf("certs/%s", nodeID), os.ModePerm)
			saveBytesToFile(certBytes, fmt.Sprintf("%s/%s/file.cert", certsDir, nodeID))
			saveBytesToFile(keyBytes, fmt.Sprintf("%s/%s/file.key", certsDir, nodeID))
		}
	}
}

func main() {
	fmt.Println("running programm...")
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