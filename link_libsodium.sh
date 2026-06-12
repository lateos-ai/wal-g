#!/bin/bash

set -e

readonly CWD=$PWD
readonly OS=$(uname)
readonly ARCH=$(uname -m)
readonly LIBSODIUM_VERSION=${LIBSODIUM_VERSION:-1.0.21}

# When system libsodium-dev is installed, copy headers/libs to tmp/libsodium/
# so the CGo directive in crypter.go resolves correctly.
echo "=== libsodium diagnostics ==="
echo "gcc: $(which gcc 2>/dev/null || echo 'NOT FOUND')"
echo "gcc version: $(gcc --version 2>&1 | head -1 || echo 'N/A')"
echo "go version: $(go version 2>&1 || echo 'N/A')"
echo "CGO_ENABLED: $(go env CGO_ENABLED 2>/dev/null || echo 'N/A')"
echo "CC: $(go env CC 2>/dev/null || echo 'N/A')"
echo "pkg-config: $(which pkg-config 2>/dev/null || echo 'NOT FOUND')"

if command -v pkg-config &>/dev/null && pkg-config --exists libsodium 2>/dev/null; then
  echo "info: system libsodium found via pkg-config"
  echo "pkg-config version: $(pkg-config --modversion libsodium 2>&1)"
  echo "pkg-config cflags: $(pkg-config --cflags libsodium 2>&1)"
  echo "pkg-config libs: $(pkg-config --libs libsodium 2>&1)"

  INCDIR=$(pkg-config --variable=includedir libsodium 2>/dev/null)
  LIBDIR=$(pkg-config --variable=libdir libsodium 2>/dev/null)
  echo "INCDIR: $INCDIR"
  echo "LIBDIR: $LIBDIR"

  if [ -n "$INCDIR" ]; then
    echo "  sodium.h exists: $(test -f "$INCDIR/sodium.h" && echo YES || echo NO)"
    echo "  sodium/ subdir: $(test -d "$INCDIR/sodium" && echo YES || echo NO)"
  fi

  echo "=== testing gcc preprocessor ==="
  echo '#include <sodium.h>' | gcc -E -xc - -o /dev/null 2>&1 && echo "gcc -E: OK" || echo "gcc -E: FAILED"

  echo "=== populating tmp/libsodium/ from system ==="
  test -d tmp/libsodium/include || mkdir -p tmp/libsodium/include
  test -d tmp/libsodium/lib || mkdir -p tmp/libsodium/lib
  if [ -n "$INCDIR" ] && [ -d "$INCDIR" ]; then
    [ -f "$INCDIR/sodium.h" ] && cp -v "$INCDIR/sodium.h" tmp/libsodium/include/ 2>&1
    [ -d "$INCDIR/sodium" ] && cp -v -r "$INCDIR/sodium" tmp/libsodium/include/ 2>&1
  fi
  # Fallback: search common include paths
  if [ ! -f tmp/libsodium/include/sodium.h ]; then
    for dir in /usr/include /usr/local/include; do
      [ -f "$dir/sodium.h" ] && cp -v "$dir/sodium.h" tmp/libsodium/include/ && cp -v -r "$dir/sodium" tmp/libsodium/include/ 2>&1 && break
    done
  fi
  if [ -n "$LIBDIR" ] && [ -d "$LIBDIR" ]; then
    echo "  libsodium files in $LIBDIR:"
    ls -la "$LIBDIR/libsodium"* 2>&1 || echo "    (none found)"
    cp -v -f "$LIBDIR/libsodium"* tmp/libsodium/lib/ 2>&1 || true
  fi
  # Fallback: search common lib paths
  if [ ! -f tmp/libsodium/lib/libsodium.a ]; then
    for dir in /usr/lib/x86_64-linux-gnu /usr/lib/aarch64-linux-gnu /usr/local/lib; do
      [ -f "$dir/libsodium.a" ] && cp -v -f "$dir/libsodium"* tmp/libsodium/lib/ 2>&1 && break
    done
  fi

  echo "=== tmp/libsodium content ==="
  find tmp/libsodium -type f -o -type d 2>/dev/null | sort
  echo "=== end diagnostics ==="
  cd ${CWD}
  exit 0
fi

echo "System libsodium not found, building from source"

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
