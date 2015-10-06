ROOT=$(dirname "${BASH_SOURCE}")/../..
VERSION=${VERSION:-v2.2.0}
curl -L  https://github.com/coreos/etcd/releases/download/${VERSION}/etcd-${VERSION}-linux-amd64.tar.gz -o etcd-${VERSION}-linux-amd64.tar.gz
tar xzvf etcd-${VERSION}-linux-amd64.tar.gz
mv etcd-${VERSION}-linux-amd64 etcd