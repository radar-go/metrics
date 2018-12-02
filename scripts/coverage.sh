#!/bin/bash

set -e

go get -u github.com/AlekSi/gocov-xml
go get -u github.com/axw/gocov/gocov

# For Go versions < 1.9, 'go list' returns vendor packages too.
packages=$(go list ./...)

# Create cover profile
COVER_MODE="count"
echo "mode: $COVER_MODE" > profile.cov
for pkg in ${packages}; do
   go test -cover -covermode=${COVER_MODE} \
        -coverprofile=profile.cov.tmp ${pkg}
   cat profile.cov.tmp | tail -n +2 >> profile.cov
   rm profile.cov.tmp
done

# Print code coverage details
go tool cover -func profile.cov


# make reports directory (if it doesn't exist)
: ${REPORTS_DIR:?}
mkdir -p ${REPORTS_DIR}

GENERATE_REPORT="${GENERATE_REPORT:-0}"

# Generate coverage report
REPORT="${REPORTS_DIR}/coverage.xml"
if [ ${GENERATE_REPORT} -eq 0 ]; then
	gocov convert profile.cov | gocov-xml > ${REPORT}
else
	gocov convert profile.cov
fi
