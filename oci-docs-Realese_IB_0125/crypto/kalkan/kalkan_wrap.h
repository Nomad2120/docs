
#ifndef __KALKAN_WRAP___
#define __KALKAN_WRAP___

int init_lib();

int free_lib();

void get_last_error(char *err_str, int errLen);

int load_key_store(int storage, char* container, char *password, char *alias) ;

int sign_data(char *alias, int flags, char *in_data, int in_data_len, char *in_sign, int in_sign_len, unsigned char *out_sign, int *out_sign_length);

int extract_cert_cms(char *cms, int cms_len, int sign_id, int flags, char *out_cert, int *out_cert_len);

int sign_wsse(char *alias, int flags, char *in_data, int in_data_len, unsigned char *out_sign, int *out_sign_length, char *sign_node_id);

#endif
