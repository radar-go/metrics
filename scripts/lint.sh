#!/bin/bash
set -e

# install gometalinter
export PATH=${PWD}/bin:${PATH}
which gometalinter || curl -L https://git.io/vp6lP | sh

# make reports directory (if it doesn't exist)
: ${REPORTS_DIR:?}
mkdir -p ${REPORTS_DIR}

GENERATE_REPORT="${GENERATE_REPORT:-0}"

# run gometalinter and generate report
REPORT="${REPORTS_DIR}/checkstyle-report.xml"
set +e
if [ ${GENERATE_REPORT} -eq 0 ]; then
	# show on screen
    gometalinter --config ${GOMETALINTER_CONFIG} --sort=path ./... > ${REPORT}
	cat ${REPORT}
	# show number of issues
	ISSUES=$(wc -l < ${REPORT})
	echo "gometalinter found ${ISSUES} issues."
else
	# generate report
	gometalinter --config ${GOMETALINTER_CONFIG} --sort=path --checkstyle ./... > ${REPORT}
	golint ./... | grep -v vendor > ${REPORTS_DIR}/linter.out
fi
status=$?
set -e

# We need to catch error codes that are bigger then 2, they signal that gometalinter
# exited because of underlying error. In this case we should fail travis build to maintain
# consistent QG reports. Otherwise different error reports will be produced based on
# gometalinter execution speed.
# more information here: https://github.com/alecthomas/gometalinter#exit-status
# If gometalinter exited because of deadline hit - increase deadline in reasonable manner.
if [ ${status} -ge 2 ]; then
    echo "gometalinter exited with code ${status}, check gometalinter errors"
    exit ${status}
fi

exit 0
