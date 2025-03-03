package signer

import (
	"crypto/x509"
	"encoding/base64"
	"errors"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	crypto "gitlab.enterprise.qazafn.kz/oci/oci-docs/crypto"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/crypto/kalkan"
	"go.uber.org/zap"
)

type Signer interface {
	SignCMSBase64(alias, data string) (string, error)
	GetIIN(cms string, signID int) (string, *x509.Certificate, error)
	SignWSSE(alias, data string, signNodeID string) (string, error)
}

type CAContainer struct {
	Container string
	Password  string
	Alias     string
}

type signer struct {
	crypt      crypto.CryptoAPI
	log        *zap.SugaredLogger
	containers []CAContainer
}

func NewSigner(log *zap.SugaredLogger, crypt crypto.CryptoAPI, containers []CAContainer) Signer {
	for _, cont := range containers {
		if err := crypt.LoadKeyStore(kalkan.KCST_PKCS12, cont.Container, cont.Password, cont.Alias); err != nil {
			log.DPanicf("LoadKeyStore %v", err)
		}
	}

	return &signer{
		crypt:      crypt,
		log:        log,
		containers: containers,
	}
}

func (s *signer) getContainer(alias string) *CAContainer {
	for _, cont := range s.containers {
		if cont.Alias == alias {
			return &cont
		}
	}
	return nil
}

func (s *signer) loadCert(alias string) error {
	cont := s.getContainer(alias)
	if cont == nil {
		return errors.New("not found container by alias: " + alias)
	}

	return s.crypt.LoadKeyStore(kalkan.KCST_PKCS12, cont.Container, cont.Password, cont.Alias)
}

func (s *signer) SignCMSBase64(alias, data string) (string, error) {
	if err := s.loadCert(alias); err != nil {
		return "", err
	}

	res, err := s.crypt.SignData(alias, kalkan.KC_SIGN_CMS|kalkan.KC_IN_BASE64|kalkan.KC_OUT_BASE64|kalkan.KC_WITH_TIMESTAMP|kalkan.KC_WITH_CERT, data, data)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (s *signer) GetIIN(cms string, signID int) (string, *x509.Certificate, error) {
	s.log.Debugf("GetIIN cms %s", cms)
	cert, err := s.crypt.ExtractCertFromCMS(cms, signID, kalkan.KC_IN_BASE64|kalkan.KC_OUT_BASE64)
	if err != nil {
		return "", nil, err
	}
	s.log.Debugf("GetIIN cert %s", string(cert))
	b, err := base64.StdEncoding.DecodeString(string(cert))
	if err != nil {
		return "", nil, err
	}
	s.log.Debug("GetIIN  pre ParseCertificate")
	crt, err := x509.ParseCertificate(b)
	if err != nil {
		return "", nil, err
	}
	s.log.Infof("Cert subject %s", crt.Subject.String())
	return common.ExtractIIN(crt.Subject.String()), crt, nil
}

func (s *signer) SignWSSE(alias, data string, signNodeID string) (string, error) {
	if err := s.loadCert(alias); err != nil {
		return "", err
	}

	res, err := s.crypt.SignWSSE(alias, 0, data, signNodeID)
	if err != nil {
		return "", err
	}
	return string(res), nil
}
