package kalkan_test

import (
	"crypto/x509"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/common"
	"gitlab.enterprise.qazafn.kz/oci/oci-docs/crypto/kalkan"
)

func TestExtractCertFromCMS(t *testing.T) {
	pem := `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`
	kal, err := kalkan.NewKalkanCrypto()
	if err != nil {
		t.Error(err)
	}
	defer kal.Close()
	cert, err := kal.ExtractCertFromCMS(pem, 1, kalkan.KC_IN_BASE64|kalkan.KC_OUT_BASE64)
	if err != nil {
		t.Error(err)
	}
	b, err := base64.StdEncoding.DecodeString(string(cert))
	if err != nil {
		t.Error(err)
	}
	crt, err := x509.ParseCertificate(b)
	if err != nil {
		t.Error(err)
	}
	iin := common.ExtractIIN(crt.Subject.String())
	assert.Equal(t, "741012301763", iin)
}

func TestSignWSSE(t *testing.T) {
	data := `<SOAP-ENV:Envelope xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
	<SOAP-ENV:Header>
	 <wsse:Security xmlns:ds="http://www.w3.org/2000/09/xmldsig#" 
					xmlns:wsse="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-secext-1.0.xsd" 
					xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" 
					xmlns:xenc="http://www.w3.org/2001/04/xmlenc#" SOAP-ENV:mustUnderstand="1">
	 <wsse:BinarySecurityToken  
				   EncodingType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-soap-message-security-1.0#Base64Binary" 
				   ValueType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509" 
				   wsu:Id="x509cert00">MIIChDCCAe2gAwIBAgIBADANBgkqhkiG9w0BAQUFADAwMQswCQYDVQQGEwJHQjEMMAoGA1UEChMD
									   SUJNMRMwEQYDVQQDEwpXaWxsIFlhdGVzMB4XDTA2MDEzMTAwMDAwMFoXDTA3MDEzMTIzNTk1OVow
									   MDELMAkGA1UEBhMCR0IxDDAKBgNVBAoTA0lCTTETMBEGA1UEAxMKV2lsbCBZYXRlczCBnzANBgkq
									   hkiG9w0BAQEFAAOBjQAwgYkCgYEArsRj/n+3RN75+jaxuOMBWSHvZCB0egv8qu2UwLWEeiogePsR
									   6Ku4SuHbBwJtWNr0xBTAAS9lEa70yhVdppxOnJBOCiERg7S0HUdP7a8JXPFzA+BqV63JqRgJyxN6
									   msfTAvEMR07LIXmZAte62nwcFrvCKNPCFIJ5mkaJ9v1p7jkCAwEAAaOBrTCBqjA/BglghkgBhvhC
									   AQ0EMhMwR2VuZXJhdGVkIGJ5IHRoZSBTZWN1cml0eSBTZXJ2ZXIgZm9yIHovT1MgKFJBQ0YpMDgG
									   ZQVRFU0BVSy5JQk0uQ09ggdJQk0uQ09NhgtXV1cuSUJNLkNPTYcECRRlBjAO
	 </wsse:BinarySecurityToken>
	 <ds:Signature xmlns:ds="http://www.w3.org/2000/09/xmldsig#">
	  <ds:SignedInfo xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/" 
					 xmlns:ds="http://www.w3.org/2000/09/xmldsig#" 
					 xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" 
					 xmlns:xenc="http://www.w3.org/2001/04/xmlenc#">
	   <ds:CanonicalizationMethod Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#">
		<c14n:InclusiveNamespaces xmlns:c14n="http://www.w3.org/2001/10/xml-exc-c14n#" PrefixList="ds wsu xenc SOAP-ENV "/>
	   </ds:CanonicalizationMethod>
	   <ds:SignatureMethod Algorithm="http://www.w3.org/2000/09/xmldsig#rsa-sha1"/>
	   <ds:Reference URI="#TheBody">
		<ds:Transforms>
		 <ds:Transform Algorithm="http://www.w3.org/2001/10/xml-exc-c14n#">
		   <c14n:InclusiveNamespaces xmlns:c14n="http://www.w3.org/2001/10/xml-exc-c14n#" PrefixList="wsu SOAP-ENV "/>
		 </ds:Transform>
		</ds:Transforms>
		<ds:DigestMethod Algorithm="http://www.w3.org/2000/09/xmldsig#sha1"/>
		<ds:DigestValue>QORZEA+gpafluShspHxhrjaFlXE=</ds:DigestValue>
	   </ds:Reference>
	  </ds:SignedInfo>
	  <ds:SignatureValue>drDH0XESiyN6YJm27mfK1ZMG4Q4IsZqQ9N9V6kEnw2lk7aM3if77XNFnyKS4deglbC3ga11kkaFJ 
						 p4jLOmYRqqycDPpqPm+UEu7mzfHRQGe7H0EnFqZpikNqZK5FF6fvYlv2JgTDPwrOSYXmhzwegUDT
						 lTVjOvuUgXYrFyaO3pw=</ds:SignatureValue>
	   <ds:KeyInfo>
		<wsse:SecurityTokenReference>
		  <wsse:Reference URI="#x509cert00" 
						  ValueType="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-x509-token-profile-1.0#X509"/> 
		</wsse:SecurityTokenReference>
	   </ds:KeyInfo>
	  </ds:Signature>
	 </wsse:Security>
	</SOAP-ENV:Header>
	<SOAP-ENV:Body xmlns:wsu="http://docs.oasis-open.org/wss/2004/01/oasis-200401-wss-wssecurity-utility-1.0.xsd" wsu:Id="TheBody">
	 <getVersion xmlns="http://msgsec.wssecfvt.ws.ibm.com"/>
	</SOAP-ENV:Body>
	</SOAP-ENV:Envelope>`
	kal, err := kalkan.NewKalkanCrypto()
	require.NoError(t, err)

	err = kal.LoadKeyStore(kalkan.KCST_PKCS12, "../../testdata/GOSTKNCA.p12", "Qwerty12", "test")
	require.NoError(t, err)

	//dataBase64 := base64.StdEncoding.EncodeToString([]byte(data))

	b, err := kal.SignWSSE("test",
		0,
		//  kalkan.KC_SIGN_CMS|
		// kalkan.KC_WITH_CERT|
		//   kalkan.KC_IN_BASE64|
		//   kalkan.KC_OUT_BASE64,
		data,
		"x509cert00")
	require.NoError(t, err)
	sign := string(b)
	n := 0
	if strings.Contains(sign, "<?xml") {
		n = strings.Index(sign, ">")
		if n == -1 {
			n = 0
		}
		if n != 0 && len(sign) > n+2 {
			if sign[n+1] == '\n' || sign[n+1] == '\r' {
				n++
			}
			if sign[n+2] == '\n' || sign[n+2] == '\r' {
				n++
			}
		}
	}
	println(sign[n+1:])
}
