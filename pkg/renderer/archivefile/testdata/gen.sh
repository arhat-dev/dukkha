#!/bin/sh

set -eux


# usage: go to the dir where this script live, run `sh gen.sh`

cd _archive_content

tar -cf ../001.tar .
tar -czf ../002.tar.gz .
tar -cjf ../003.tar.bz2 .
tar --lzma -cf ../004.tar.lzma .
tar -cJf ../005.tar.xz .
zip -r ../101.zip .

cd -

# externally compressed zip files
rm -f 102.zip.gz || true
gzip -k 101.zip
mv 101.zip.gz 102.zip.gz

rm -f 103.zip.bz2 || true
bzip2 -k 101.zip
mv 101.zip.bz2 103.zip.bz2

rm -f 104.zip.lzma || true
lzma -k 101.zip
mv 101.zip.lzma 104.zip.lzma

rm -f 105.zip.xz || true
xz -k 101.zip
mv 101.zip.xz 105.zip.xz
