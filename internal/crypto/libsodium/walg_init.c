#include <sodium.h>
#include <walg_config.h>

int walg_sodium_init(void)
{
    return sodium_init();
}

int walg_secretstream_init_push(
    walg_secretstream_state *state,
    unsigned char *header,
    const unsigned char *key)
{
    return crypto_secretstream_xchacha20poly1305_init_push(
        (crypto_secretstream_xchacha20poly1305_state *)state,
        header, key);
}

int walg_secretstream_init_pull(
    walg_secretstream_state *state,
    const unsigned char *header,
    const unsigned char *key)
{
    return crypto_secretstream_xchacha20poly1305_init_pull(
        (crypto_secretstream_xchacha20poly1305_state *)state,
        header, key);
}

int walg_secretstream_push(
    walg_secretstream_state *state,
    unsigned char *out, unsigned long long *out_len,
    const unsigned char *in, unsigned long long in_len,
    const unsigned char *ad, unsigned long long ad_len,
    unsigned char tag)
{
    return crypto_secretstream_xchacha20poly1305_push(
        (crypto_secretstream_xchacha20poly1305_state *)state,
        out, out_len, in, in_len, ad, ad_len, tag);
}

int walg_secretstream_pull(
    walg_secretstream_state *state,
    unsigned char *out, unsigned long long *out_len,
    unsigned char *tag,
    const unsigned char *in, unsigned long long in_len,
    const unsigned char *ad, unsigned long long ad_len)
{
    return crypto_secretstream_xchacha20poly1305_pull(
        (crypto_secretstream_xchacha20poly1305_state *)state,
        out, out_len, tag, in, in_len, ad, ad_len);
}
