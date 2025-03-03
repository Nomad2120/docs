#ifdef _WIN32
#include <windows.h>
#else
#include <dlfcn.h>
#include <dirent.h>
#endif

#include <stdio.h>
#include <string.h>

#ifndef _WIN32
#define _DLL
#undef KALKAN_EXPORTS
#endif

#include "KalkanCrypt.h"
#include "kalkan_wrap.h"

#define LENGHT 2048

typedef int (*KC_GetFunctionList1)(stKCFunctionsType **KCfunc);
KC_GetFunctionList1 lib_funcList = NULL;
stKCFunctionsType *kc_funcs;

#ifdef _WIN32
HMODULE handle = NULL;
#else
void    *handle;
#endif
const char* tsaurl = "http://tsp.pki.gov.kz:80";


int init_lib() {
#ifdef _WIN32
    handle = LoadLibrary("KalkanCrypt.dll");//("KalkanCrypt_x64.dll");
	if (!handle) {
		int err = GetLastError();
        return err;
	}

	lib_funcList = (int(*)(stKCFunctionsType**))GetProcAddress(handle, "KC_GetFunctionList");
	lib_funcList(&kc_funcs);
	int err = kc_funcs->KC_Init();
    if (err) return err;
	kc_funcs->KC_TSASetUrl((char*)tsaurl);
	return 0;
#else	
  handle = dlopen("libkalkancryptwr-64.so.2.0.9",  RTLD_LAZY);
	if (!handle) {
			fputs (dlerror(), stderr);
			return 1;
	}

  lib_funcList = (KC_GetFunctionList1)dlsym(handle, "KC_GetFunctionList");
	lib_funcList(&kc_funcs);
	int rv = kc_funcs->KC_Init();
	if (rv != 0) {
		return rv;
	}
	kc_funcs->KC_TSASetUrl((char*)tsaurl);
  return rv;
#endif  
}

int load_key_store(int storage, char* container, char *password, char *alias)  {
	return kc_funcs->KC_LoadKeyStore(storage, password, strlen(password), container, strlen(container), alias);
}

int free_lib() {
#ifdef _WIN32
	BOOL res = FreeLibrary(handle);
   return res ? 0 : 1;
#else
  return dlclose(handle);
 #endif
}


void get_last_error(char *err_str, int errLen) {
	int rv = kc_funcs->KC_GetLastErrorString(err_str, &errLen);
}

int sign_data(char *alias, int flags, char *in_data, int in_data_len, char *in_sign, int in_sign_len, unsigned char *out_sign, int *out_sign_length) {
  int rv = kc_funcs->SignData(alias, flags, in_data, in_data_len, (unsigned char*) in_sign, in_sign_len, out_sign, out_sign_length);
  return rv;
}

int extract_cert_cms(char *cms, int cms_len, int sign_id, int flags, char *out_cert, int *out_cert_len) {
	int rv = kc_funcs->KC_GetCertFromCMS(cms, cms_len, sign_id, flags, out_cert, out_cert_len);
    return rv;
}

int sign_wsse(char *alias, int flags, char *in_data, int in_data_len, unsigned char *out_sign, int *out_sign_length, char *sign_node_id) {
	int rv = kc_funcs->SignWSSE(alias, flags, in_data, in_data_len, out_sign, out_sign_length, sign_node_id);
	return rv;
}

