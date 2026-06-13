// This file provides Go 1.25 cgo DWARF type information for preamble
// declarations under -mod=vendor. It mirrors the declarations from
// walg_config.h without needing include paths or sodium headers.
typedef struct { unsigned char d[52]; } walg_secretstream_state;
int walg_sodium_init();
int walg_secretstream_init_push(walg_secretstream_state *, unsigned char *, const unsigned char *);
int walg_secretstream_init_pull(walg_secretstream_state *, const unsigned char *, const unsigned char *);
int walg_secretstream_push(walg_secretstream_state *, unsigned char *, unsigned long long *, const unsigned char *, unsigned long long, const unsigned char *, unsigned long long, unsigned char);
int walg_secretstream_pull(walg_secretstream_state *, unsigned char *, unsigned long long *, unsigned char *, const unsigned char *, unsigned long long, const unsigned char *, unsigned long long);
