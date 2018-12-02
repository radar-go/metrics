#!/bin/bash

set -o errexit
set -o nounset

export CGO_ENABLED=0

set -e

: ${REPORTS_DIR:?}

mkdir -p "${REPORTS_DIR}"

COVER_FILE="${REPORTS_DIR}/cover.out"
COVERAGE_REPORT="${REPORTS_DIR}/coverage.xml"
JUNIT_REPORT="${REPORTS_DIR}/junit-report.xml"

# List of tools that used to generate Quality Gate reports
tools=(
	github.com/axw/gocov/gocov
	github.com/AlekSi/gocov-xml
	github.com/jstemmer/go-junit-report
)

# Install missed tools
for tool in ${tools[@]}; do
	which $(basename ${tool}) > /dev/null || go get -u -v ${tool}
done

echo "Running unit tests."

PACKAGES=$(go list ./...)

# Generate tests report
go test -v ${PACKAGES} -coverprofile ${COVER_FILE} | tee /dev/tty | go-junit-report > "${JUNIT_REPORT}"; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

# Print code coverage details
go tool cover -func "${COVER_FILE}"

# Generate coverage report
echo "Generate coverage report."
gocov convert "${COVER_FILE}" | gocov-xml  > ${COVERAGE_REPORT}; test ${PIPESTATUS[0]} -eq 0 || status=${PIPESTATUS[0]}

if  [ ${status:-0} -ne 0 ]; then
    echo "FAIL"
    exit 1
fi
echo "PASS"

echo "Checking gofmt: "
ERRS=$(go fmt ${PACKAGES} || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL - the following files need to be gofmt'ed:"
    for e in ${ERRS}; do
        echo "    $e"
    done
    echo
    exit 1
fi
echo "PASS"

echo "Checking go vet: "
ERRS=$(go vet ${PACKAGES} 2>&1 || true)
if [ -n "${ERRS}" ]; then
    echo "FAIL"
    echo "${ERRS}"
    echo
    exit 1
fi
echo "PASS"

if [ ${GENERATE_REPORT} -eq 0 ]; then
	go test -json -v ${PACKAGES} > ${REPORTS_DIR}/test-report.json
fi

if [ ${GENERATE_REPORT} -eq 0 ]; then
	ERRS=$(go vet ${PACKAGES} 2>&1 || true)
else
	ERRS=$(go vet ${PACKAGES} 2> ${REPORTS_DIR}/govet-report.out || true)
fi
if [ -n "${ERRS}" ]; then
	echo "FAIL"
	echo "${ERRS}"
	echo
	exit 1
fi
echo "PASS"
