#!/bin/bash

set -e

readonly CWD=$PWD
readonly OS=$(uname)
readonly ARCH=$(uname -m)
readonly LIBSODIUM_VERSION=${LIBSODIUM_VERSION:-1.0.21}

# When a system libsodium-dev (or equivalent) package is installed,
# populate tmp/libsodium/ from it. This lets builds using the libsodium
# build tag find the headers and static library even under -mod=vendor.
# (Go >= 1.21 sanitizes #cgo CFLAGS/LDFLAGS directives; we avoid path
# directives in source and rely on CGO_ENABLED + CGO_CFLAGS/CGO_LDFLAGS
# from the Makefile + this tree.)
if command -v pkg-config >/dev/null 2>&1 && pkg-config --exists libsodium 2>/dev/null; then
	echo "info: system libsodium found via pkg-config"
	INCDIR=$(pkg-config --variable=includedir libsodium 2>/dev/null)
	LIBDIR=$(pkg-config --variable=libdir libsodium 2>/dev/null)

	mkdir -p tmp/libsodium/include tmp/libsodium/lib

	if [ -n "$INCDIR" ] && [ -d "$INCDIR" ]; then
		[ -f "$INCDIR/sodium.h" ] && cp -f "$INCDIR/sodium.h" tmp/libsodium/include/
		[ -d "$INCDIR/sodium" ] && cp -rf "$INCDIR/sodium" tmp/libsodium/include/
	fi

	# Fallback search for common system locations
	if [ ! -f tmp/libsodium/include/sodium.h ]; then
		for d in /usr/include /usr/local/include; do
			if [ -f "$d/sodium.h" ]; then
				cp -f "$d/sodium.h" tmp/libsodium/include/
				[ -d "$d/sodium" ] && cp -rf "$d/sodium" tmp/libsodium/include/
				break
			fi
		done
	fi

	# Extra broad search (covers cases where sodium.h lives in unexpected subdirs)
	if [ ! -f tmp/libsodium/include/sodium.h ]; then
		for candidate in $(find /usr -name sodium.h -type f 2>/dev/null | head -20); do
			# Prefer a top-level sodium.h (not the one inside a sodium/ subdir) for #include <sodium.h>
			if [ -f "$candidate" ] && [[ "$candidate" != */sodium/sodium.h ]]; then
				cp -f "$candidate" tmp/libsodium/include/sodium.h
				break
			fi
		done
	fi

	# Always ensure the sodium/ subdir with detailed headers is present (internal includes use "sodium/...")
	for sdir in "$INCDIR/sodium" /usr/include/sodium /usr/local/include/sodium; do
		if [ -d "$sdir" ]; then
			mkdir -p tmp/libsodium/include
			cp -rf "$sdir" tmp/libsodium/include/ 2>/dev/null || true
			break
		fi
	done

	# Last-resort redirect wrapper so #include <sodium.h> works even if only subdir was found
	if [ ! -f tmp/libsodium/include/sodium.h ] && [ -d tmp/libsodium/include/sodium ]; then
		cat > tmp/libsodium/include/sodium.h << 'EOT'
/* Auto-generated redirect for wal-g libsodium builds (system package layout) */
#include "sodium/sodium.h"
EOT
	fi

	if [ -n "$LIBDIR" ] && [ -d "$LIBDIR" ]; then
		cp -f "$LIBDIR"/libsodium* tmp/libsodium/lib/ 2>/dev/null || true
	fi
	if [ ! -f tmp/libsodium/lib/libsodium.a ]; then
		for d in /usr/lib/x86_64-linux-gnu /usr/lib/aarch64-linux-gnu /usr/local/lib; do
			if [ -f "$d/libsodium.a" ]; then
				cp -f "$d"/libsodium* tmp/libsodium/lib/ 2>/dev/null || true
				break
			fi
		done
	fi

	# Diagnostics + hard verification (prevents silent "info printed, later cgo failure")
	echo "info: pkg-config INCDIR=$INCDIR LIBDIR=$LIBDIR"
	echo "info: libsodium tree after population:"
	ls -lR tmp/libsodium || true
	if [ ! -f tmp/libsodium/include/sodium.h ]; then
		echo "ERROR: failed to populate tmp/libsodium/include/sodium.h"
		echo "Searched candidates:"
		find /usr -name 'sodium.h' 2>/dev/null | head -10 || true
		exit 1
	fi
	if [ ! -f tmp/libsodium/lib/libsodium.a ] && [ ! -f tmp/libsodium/lib/libsodium.so ] && [ ! -f tmp/libsodium/lib/libsodium.dylib ]; then
		echo "WARNING: no libsodium static/shared lib found in tmp/libsodium/lib; link step may fail"
	fi

	cd "$CWD"
	exit 0
fi

echo "info: no system libsodium, building from source"

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

# Verification for the source-build case (should have correct layout from configure --prefix)
echo "info: libsodium tree after source build:"
ls -lR tmp/libsodium || true
if [ ! -f tmp/libsodium/include/sodium.h ]; then
	echo "ERROR: source build did not produce tmp/libsodium/include/sodium.h"
	exit 1
fi
if [ ! -f tmp/libsodium/lib/libsodium.a ]; then
	echo "WARNING: source build did not produce tmp/libsodium/lib/libsodium.a"
fi
