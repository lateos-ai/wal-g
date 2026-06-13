// Weak stub definitions for Go 1.25 cgo DWARF analysis under -mod=vendor.
#include "gen/walg_config.h"

int walg_secretstream_init_push(walg_secretstream_state *state, unsigned char *header, const unsigned char *key) __attribute__((weak));
int walg_secretstream_init_push(walg_secretstream_state *state, unsigned char *header, const unsigned char *key) { return 0; }

int walg_secretstream_init_pull(walg_secretstream_state *state, const unsigned char *header, const unsigned char *key) __attribute__((weak));
int walg_secretstream_init_pull(walg_secretstream_state *state, const unsigned char *header, const unsigned char *key) { return 0; }

int walg_secretstream_push(walg_secretstream_state *state, unsigned char *out, unsigned long long *out_len, const unsigned char *in, unsigned long long in_len, const unsigned char *ad, unsigned long long ad_len, unsigned char tag) __attribute__((weak));
int walg_secretstream_push(walg_secretstream_state *state, unsigned char *out, unsigned long long *out_len, const unsigned char *in, unsigned long long in_len, const unsigned char *ad, unsigned long long ad_len, unsigned char tag) { return 0; }

int walg_secretstream_pull(walg_secretstream_state *state, unsigned char *out, unsigned long long *out_len, unsigned char *tag, const unsigned char *in, unsigned long long in_len, const unsigned char *ad, unsigned long long ad_len) __attribute__((weak));
int walg_secretstream_pull(walg_secretstream_state *state, unsigned char *out, unsigned long long *out_len, unsigned char *tag, const unsigned char *in, unsigned long long in_len, const unsigned char *ad, unsigned long long ad_len) { return 0; }
