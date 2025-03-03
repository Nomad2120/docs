package kalkan

// #cgo CFLAGS: -I .
// #cgo linux LDFLAGS:  -L . -ldl
// #include "stdlib.h"
// #include "KalkanCrypt.h"
// #include "kalkan_wrap.h"
import "C"
import (
	"errors"
	"fmt"
	"sync"
	"unsafe"

	"gitlab.enterprise.qazafn.kz/oci/oci-docs/crypto"
)

const (
	KCST_PKCS12   = C.KCST_PKCS12
	KCST_KZIDCARD = C.KCST_KZIDCARD

	KC_SIGN_DRAFT      = C.KC_SIGN_DRAFT
	KC_SIGN_CMS        = C.KC_SIGN_CMS
	KC_IN_PEM          = C.KC_IN_PEM
	KC_IN_DER          = C.KC_IN_DER
	KC_IN_BASE64       = C.KC_IN_BASE64
	KC_IN2_BASE64      = C.KC_IN2_BASE64
	KC_DETACHED_DATA   = C.KC_DETACHED_DATA
	KC_WITH_CERT       = C.KC_WITH_CERT
	KC_WITH_TIMESTAMP  = C.KC_WITH_TIMESTAMP
	KC_OUT_PEM         = C.KC_OUT_PEM
	KC_OUT_DER         = C.KC_OUT_DER
	KC_OUT_BASE64      = C.KC_OUT_BASE64
	KC_PROXY_OFF       = C.KC_PROXY_OFF
	KC_PROXY_ON        = C.KC_PROXY_ON
	KC_PROXY_AUTH      = C.KC_PROXY_AUTH
	KC_IN_FILE         = C.KC_IN_FILE
	KC_NOCHECKCERTTIME = C.KC_NOCHECKCERTTIME
	KC_HASH_SHA256     = C.KC_HASH_SHA256
	KC_HASH_GOST95     = C.KC_HASH_GOST95

	KC_CERT_CA           = C.KC_CERT_CA
	KC_CERT_INTERMEDIATE = C.KC_CERT_INTERMEDIATE
	KC_CERT_USER         = C.KC_CERT_USER
)

type kalkanCrypt struct {
}

func NewKalkanCrypto() (crypto.CryptoAPI, error) {
	var (
		once sync.Once
		err  error
	)
	onceInit := func() {
		res := int(C.init_lib())
		if res != 0 {
			err = fmt.Errorf("Error init kalkan library: %d", res)
		}
	}
	once.Do(onceInit)
	if err != nil {
		return nil, err
	}
	return &kalkanCrypt{}, nil
}

func (k *kalkanCrypt) LoadKeyStore(storage int, container, password, alias string) error {
	ccont := C.CString(container)
	cpass := C.CString(password)
	calias := C.CString(alias)
	defer func() {
		C.free(unsafe.Pointer(ccont))
		C.free(unsafe.Pointer(cpass))
		C.free(unsafe.Pointer(calias))
	}()

	res := int(C.load_key_store(C.int(storage), ccont, cpass, calias))

	if res != 0 {
		return k.lastError()
	}
	return nil
}

// func (k *kalkanSigner) LoadKeyStoreP12(container, password, alias string) error {
// 	return k.LoadKeyStore(KCST_PKCS12, container, password, alias)
// }

func (k *kalkanCrypt) Close() error {
	var (
		once sync.Once
		err  error
	)
	onceInit := func() {
		res := int(C.free_lib())
		if res != 0 {
			err = fmt.Errorf("Error close kalkan library: %d", res)
		}
	}
	once.Do(onceInit)

	return err
}

func (k *kalkanCrypt) SignData(alias string, flags int, inData string, inSign string) ([]byte, error) {
	outSignLen := 50000 + 2*len(inData)
	cAlias := C.CString(alias)
	cInData := C.CString(inData)
	cInSign := C.CString(inSign) // (*C.uchar)(unsafe.Pointer(&inSign[0]))
	ptr := C.malloc(C.sizeof_char * C.size_t(outSignLen))
	sigLen := C.int(outSignLen)

	defer func() {
		C.free(unsafe.Pointer(cAlias))
		C.free(unsafe.Pointer(cInData))
		C.free(unsafe.Pointer(cInSign))
		C.free(unsafe.Pointer(ptr))
	}()

	res := int(C.sign_data(cAlias, C.int(flags), cInData, C.int(len(inData)), cInSign, C.int(len(inSign)), (*C.uchar)(ptr), &sigLen))
	if res != 0 {
		return nil, k.lastError()
	}
	sign := C.GoBytes(unsafe.Pointer(ptr), sigLen)
	return sign, nil
}

func (k *kalkanCrypt) SignCMSBase64(alias string, inData string, inSign string) (string, error) {
	res, err := k.SignData(alias, KC_SIGN_CMS|KC_IN_BASE64|KC_OUT_BASE64|KC_WITH_TIMESTAMP|KC_WITH_CERT, inData, inSign)
	//res, err := k.SignData(KC_SIGN_CMS|KC_IN_BASE64|KC_OUT_BASE64|KC_WITH_CERT, inData, inSign)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func (k *kalkanCrypt) ExtractCertFromCMSBase64(cms string, signID int) ([]byte, error) {
	return k.ExtractCertFromCMS(cms, signID, KC_IN_BASE64|KC_OUT_BASE64)
}

func (k *kalkanCrypt) ExtractCertFromCMS(cms string, signID int, flags int) ([]byte, error) {
	outCertLen := 32768
	cCMS := C.CString(cms)
	ptr := C.malloc(C.sizeof_char * C.size_t(outCertLen))
	certLen := C.int(outCertLen)

	defer func() {
		C.free(unsafe.Pointer(cCMS))
		C.free(unsafe.Pointer(ptr))
	}()

	res := int(C.extract_cert_cms(cCMS, C.int(len(cms)), C.int(signID), C.int(flags), (*C.char)(ptr), &certLen))
	if res != 0 {
		return nil, k.lastError()
	}
	cert := C.GoBytes(unsafe.Pointer(ptr), certLen)
	return cert, nil
}

func (k *kalkanCrypt) SignWSSE(alias string, flags int, inData, signNodeID string) ([]byte, error) {
	outSignLen := 50000 + 2*len(inData)
	cAlias := C.CString(alias)
	cSignNodeID := C.CString(signNodeID)
	cInData := C.CString(inData)
	ptr := C.malloc(C.sizeof_char * C.size_t(outSignLen))
	sigLen := C.int(outSignLen)

	defer func() {
		C.free(unsafe.Pointer(cAlias))
		C.free(unsafe.Pointer(cSignNodeID))
		C.free(unsafe.Pointer(cInData))
		C.free(unsafe.Pointer(ptr))
	}()

	res := int(C.sign_wsse(cAlias, C.int(flags), cInData, C.int(len(inData)), (*C.uchar)(ptr), &sigLen, cSignNodeID))
	if res != 0 {
		return nil, k.lastError()
	}
	sign := C.GoBytes(unsafe.Pointer(ptr), sigLen-1)
	return sign, nil
}

func (k *kalkanCrypt) lastError() error {
	size := 2048
	ptr := C.malloc(C.sizeof_char * C.size_t(size))
	defer C.free(unsafe.Pointer(ptr))
	C.get_last_error((*C.char)(ptr), C.int(size))
	str := C.GoString((*C.char)(ptr))
	return errors.New(str)
}
