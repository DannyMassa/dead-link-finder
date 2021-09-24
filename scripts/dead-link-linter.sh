set -xe
: "${LINTER_VERSION:="0.0.6"}"

URL="https://github.com/DannyMassa/dead-link-finder/releases/download/${LINTER_VERSION}/dead-link-linter-${LINTER_VERSION}-linux-amd64.tar.gz"
rm -rf /tmp/dead-link-linter
mkdir /tmp/dead-link-linter
sudo -E curl -sSLo "/tmp/dead-link-linter/dead-link-linter-${LINTER_VERSION}-linux-amd64.tar.gz" ${URL}
tar xvf /tmp/dead-link-linter/dead-link-linter-${LINTER_VERSION}-linux-amd64.tar.gz /tmp/dead-link-linter/
sudo install -m 755 -o root /tmp/dead-link-linter/dead-link-linter /usr/local/bin
dead-link-linter
status=$?
[ $status -eq 0 ] && echo "Dead Link Linter Passed" || echo "Dead Link Linter Failed" && exit 1
