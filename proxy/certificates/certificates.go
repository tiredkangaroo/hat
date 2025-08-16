package certificates

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"main/proxy/config"
	"math/big"
	"net"
	"os"
	"sync"
	"time"
)

type Service struct {
	cache sync.Map // cache for certificates

	cert    *x509.Certificate
	key     any
	Ready   bool
	Enabled bool
}

func (c *Service) Init() error {
	if config.DefaultConfig.MITM.CertificateFile == "" || config.DefaultConfig.MITM.KeyFile == "" {
		return fmt.Errorf("certificate file or key file not provided")
	}

	// read the files
	caCertRaw, err := os.ReadFile(config.DefaultConfig.MITM.CertificateFile)
	if err != nil {
		return fmt.Errorf("read cert file: %w", err)
	}
	caKeyRaw, err := os.ReadFile(config.DefaultConfig.MITM.KeyFile)
	if err != nil {
		return fmt.Errorf("read key file: %w", err)
	}

	// pem decode
	certBlock, _ := pem.Decode(caCertRaw)
	keyBlock, _ := pem.Decode(caKeyRaw)
	if certBlock == nil || keyBlock == nil {
		return fmt.Errorf("decode cert or key file: pem decode failed")
	}

	// parse
	c.cert, err = x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse cert file: %w", err)
	}
	c.key, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return fmt.Errorf("parse key file: %w", err)
	}

	c.Ready = true
	c.Enabled = true // defaults to enabled if cert and key are loaded

	return nil
}

func (c *Service) getTLSCert(host string) (tls.Certificate, error) {
	// check if the certificate is already in the cache
	cachedCertificate, ok := c.cache.Load(host)
	if ok {
		if time.Until(cachedCertificate.(tls.Certificate).Leaf.NotAfter) < time.Minute {
			c.cache.Delete(host) // delete the certificate if it is expired or about to expire
		} else {
			return cachedCertificate.(tls.Certificate), nil
		}
	}

	// generate a new private key for the new certificate
	pk, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("an error occured while attempting to generate an ecdsa key: %s", err.Error())
	}

	// create a serial number for certificate
	sn, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("an error occured while attempting to generate a certificate serial number: %s", err.Error())
	}

	// create cert config
	config := &x509.Certificate{
		SerialNumber: sn,
		Subject: pkix.Name{
			Country:      []string{"US"},
			Organization: []string{"N/A"},
		},
		DNSNames:              []string{host},
		NotBefore:             time.Now().Add(-(time.Hour * 7200)),
		NotAfter:              time.Now().Add(time.Hour * time.Duration(config.DefaultConfig.MITM.CertificateLifetime)),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
	}

	// create certificate
	cert, err := x509.CreateCertificate(rand.Reader, config, c.cert, &pk.PublicKey, c.key)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("creating the x509 certificate: %w", err)
	}

	// encode certificate
	pemCert := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert})
	if pemCert == nil {
		return tls.Certificate{}, fmt.Errorf("encode the cert with pe (unknown error)")
	}

	// encode the private key
	privBytes, err := x509.MarshalPKCS8PrivateKey(pk)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("marshal private key: %w", err)
	}
	pemKey := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: privBytes})
	if pemKey == nil {
		return tls.Certificate{}, fmt.Errorf("encode the private key with pem (unknown error)")
	}

	// create the certificate
	tlscert, err := tls.X509KeyPair(pemCert, pemKey)
	if err != nil {
		return tls.Certificate{}, fmt.Errorf("create x509 key pair: %w", err)
	}

	// store the certificate in the cache
	c.cache.Store(host, tlscert)
	return tlscert, nil

}

func (c *Service) TLSConn(conn net.Conn, host string) (*tls.Conn, error) {
	cert, err := c.getTLSCert(host)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		PreferServerCipherSuites: true,
		CurvePreferences:         []tls.CurveID{tls.X25519, tls.CurveP256},
		MinVersion:               tls.VersionTLS10,
		MaxVersion:               tls.VersionTLS13,
		Certificates:             []tls.Certificate{cert},
	}
	return tls.Server(conn, tlsConfig), nil
}

// GetServices creates and initializes a new certificate service instance. It returns the instance or
// an error if initialization fails.
func GetService() (*Service, error) {
	certService := &Service{}
	if err := certService.Init(); err != nil {
		return nil, fmt.Errorf("initialize certificate service: %w", err)
	}
	return certService, nil
}
