go install omo.msa.acm
mkdir _build
mkdir _build/bin

cp -rf /root/go/bin/omo.msa.acm _build/bin/
cp -rf conf _build/
cd _build
tar -zcf msa.acm.tar.gz ./*
mv msa.acm.tar.gz ../
cd ../
rm -rf _build
