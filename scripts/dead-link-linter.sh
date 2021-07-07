: "${LINTER_VERSION:="0.0.6"}"

wget https://github.com/DannyMassa/dead-link-finder/releases/download/${LINTER_VERSION}/dead-link-linter-${LINTER_VERSION}-linux-amd64.tar.gz
tar -xvf dead-link-linter-${LINTER_VERSION}-linux-amd64.tar.gz
./dead-link-linter
status=$?
rm dead-link-linter-${LINTER_VERSION}-linux-amd64.tar.gz
[ $status -eq 0 ] && echo "Dead Link Linter Passed" || echo "Dead Link Linter Failed" && exit 1
