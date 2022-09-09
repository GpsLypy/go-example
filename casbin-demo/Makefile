cfssl gencert \
         -ca=ca.pem \
         -ca-key=ca-key.pem \
         -config=test/ca-config.json \
         -profile=client \
         -cn="root" \
         test/client-csr.json | cfssljson -bare root-client
 
 cfssl gencert \
         -ca=ca.pem \
         -ca-key=ca-key.pem \
         -config=test/ca-config.json \
         -profile=client \
         -cn="nobody" \
         test/client-csr.json | cfssljson -bare nobody-client