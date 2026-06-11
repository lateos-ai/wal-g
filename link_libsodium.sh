#!/bin/bash

set -e

readonly CWD=$PWD
readonly OS=$(uname)
readonly ARCH=$(uname -m)
readonly LIBSODIUM_VERSION=${LIBSODIUM_VERSION:-1.0.21}

# Prefer system-installed libsodium via pkg-config when available.
# This avoids CGo linking issues with source-built static libraries.
if command -v pkg-config &>/dev/null && pkg-config --exists libsodium 2>/dev/null; then
  echo "info: system libsodium found via pkg-config, using it"
  test -d tmp/libsodium/include || mkdir -p tmp/libsodium/include
  test -d tmp/libsodium/lib || mkdir -p tmp/libsodium/lib
  INCDIR=$(pkg-config --variable=includedir libsodium 2>/dev/null || pkg-config --cflags-only-I libsodium 2>/dev/null | sed 's/-I//')
  LIBDIR=$(pkg-config --variable=libdir libsodium 2>/dev/null || pkg-config --libs-only-L libsodium 2>/dev/null | sed 's/-L//')
  [ -n "$INCDIR" ] && [ -d "$INCDIR" ] && cp -r "$INCDIR"/sodium.h "$INCDIR"/sodium tmp/libsodium/include/ 2>/dev/null || true
  [ -n "$LIBDIR" ] && [ -d "$LIBDIR" ] && cp -f "$LIBDIR"/libsodium* tmp/libsodium/lib/ 2>/dev/null || true
  cd ${CWD}
  exit 0
fi

test -d tmp/libsodium || mkdir -p tmp/libsodium

cd tmp/libsodium

curl --retry 5 --retry-delay 0 -sL https://github.com/jedisct1/libsodium/releases/download/$LIBSODIUM_VERSION-RELEASE/libsodium-$LIBSODIUM_VERSION.tar.gz -o libsodium-$LIBSODIUM_VERSION.tar.gz
tar xfz libsodium-$LIBSODIUM_VERSION.tar.gz --strip-components=1

CONFIGURE_ARGS="--prefix ${PWD} --disable-debug --disable-dependency-tracking --enable-static --disable-shared"
if [[ "${OS}" == "SunOS" ]]; then
  # On Illumos / Solaris libssp causes linking issues when building wal-g.
  CONFIGURE_ARGS="${CONFIGURE_ARGS} --disable-ssp"
fi   

LOCAL_CFLAGS="-O2"
if [[ "${OS}" == "Linux" ]] && [[ "${ARCH}" == *arm* || "${ARCH}" == "aarch64" ]]; then
  LOCAL_CFLAGS="${LOCAL_CFLAGS} -flax-vector-conversions"
fi

CFLAGS="${LOCAL_CFLAGS}" ./configure ${CONFIGURE_ARGS}
make && make check && make install

# Remove shared libraries for using static
rm -f lib/*.so lib/*.so.* lib/*.dylib

cd ${CWD}
