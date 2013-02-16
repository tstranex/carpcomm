#!/bin/bash

# 1. Update version:
VERSION=0.19
# 2. Update version in setup.py.
# 3. Update version in client.py.
# 4. Update debian/changelog with new version.
# 5. Commit to git.
# 6. Run this script.

export DEBFULLNAME="Carpcomm GmbH (Linux Packaging)"
export DEBEMAIL="linux-packaging@carpcomm.com"

#dh_make -s -f ../carpsd-0.1.tar.gz

python setup.py sdist
cd dist
cp carpsd-$VERSION.tar.gz carpsd_$VERSION.orig.tar.gz
tar -xzf carpsd-$VERSION.tar.gz
cd carpsd-$VERSION/
cp -R ../../debian .
debuild -us -uc

# 7. Add to apt repository:
# cd ~/carp/apt
# reprepro --ask-passphrase -Vb . includedeb lucid carpsd_0.1-1_all.deb
# Add to gs1 repository:
# cd ~/tmp/gs1/apt
# reprepro --ask-passphrase -Vb . includedeb gs1 carpsd_0.1-1_all.deb

# 8. Update date at the bottom of docs/index.html.
# 9. Update files on live http server.