#!/usr/bin/env bash

PACKAGE_NAME="chirpstack-pg-to-sqlite"
PACKAGE_VERSION=$1
REV="r1"

DIR=`dirname $0`
PACKAGE_DIR="${DIR}/package"

# Cleanup
rm -rf $PACKAGE_DIR

# CONTROL
mkdir -p $PACKAGE_DIR/CONTROL
cat > $PACKAGE_DIR/CONTROL/control << EOF
Package: $PACKAGE_NAME
Version: $PACKAGE_VERSION-$REV
Architecture: all
Maintainer: Orne Brocaar <info@brocaar.com>
Priority: optional
Section: network
Source: N/A
Description: ChirpStack PostgreSQL to SQLite
EOF

cat > $PACKAGE_DIR/CONTROL/postinst << EOF
mkdir -p /srv/chirpstack

cp /opt/$PACKAGE_NAME/chirpstack.empty.sqlite /tmp/chirpstack.sqlite
/opt/$PACKAGE_NAME/chirpstack-pg-to-sqlite -sqlite-path /tmp/chirpstack.sqlite

cp /tmp/chirpstack.sqlite /srv/chirpstack/chirpstack.sqlite
rm /tmp/chirpstack.sqlite
EOF
chmod 755 $PACKAGE_DIR/CONTROL/postinst

# Files
mkdir -p $PACKAGE_DIR/opt/$PACKAGE_NAME
cp ../../chirpstack-pg-to-sqlite $PACKAGE_DIR/opt/$PACKAGE_NAME
cp ../../chirpstack.empty.sqlite $PACKAGE_DIR/opt/$PACKAGE_NAME

# Package
opkg-build -c -o root -g root $PACKAGE_DIR

# Cleanup
rm -rf $PACKAGE_DIR
