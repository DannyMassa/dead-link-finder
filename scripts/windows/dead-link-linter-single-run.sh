set -xe

: "${LINTER_VERSION:="1.0.2"}"
: "${OS_TYPE:="windows"}"
: "${ARCHITECTURE:="amd64"}"

rm dead-link-linter-${LINTER_VERSION}-${OS_TYPE}-${ARCHITECTURE}.zip || true
rm dead-link-linter.exe || true

URL="https://github.com/DannyMassa/dead-link-linter/releases/download/${LINTER_VERSION}/dead-link-linter-${LINTER_VERSION}-${OS_TYPE}-${ARCHITECTURE}.zip"
wget ${URL}
unzip dead-link-linter-${LINTER_VERSION}-${OS_TYPE}-${ARCHITECTURE}.zip

./dead-link-linter.exe

status=$?
[ $status -eq 0 ] && \
  echo "Dead Link Linter Passed" && \
  rm dead-link-linter-${LINTER_VERSION}-${OS_TYPE}-${ARCHITECTURE}.zip && \
  rm dead-link-linter.exe && \
  exit 0 \
  || \
  echo "Dead Link Linter Failed" && \
  rm dead-link-linter-${LINTER_VERSION}-${OS_TYPE}-${ARCHITECTURE}.zip && \
  rm dead-link-linter.exe && \
  exit 1
