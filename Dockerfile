FROM golang:alpine as builder
COPY . /go/src/github.com/SWCE/keyval-resource
ENV CGO_ENABLED 0
ENV GOPATH /go/src/github.com/SWCE/keyval-resource/Godeps/_workspace:${GOPATH}
ENV PATH /go/src/github.com/SWCE/keyval-resource/Godeps/_workspace/bin:${PATH}
RUN go build -o /assets/out github.com/SWCE/keyval-resource/out
RUN go build -o /assets/in github.com/SWCE/keyval-resource/in
RUN go build -o /assets/check github.com/SWCE/keyval-resource/check
RUN set -e; for pkg in $(go list ./...); do \
		go test -o "/tests/$(basename $pkg).test" -c $pkg; \
	done

FROM alpine:edge AS resource
RUN apk add --update bash tzdata
COPY --from=builder /assets /opt/resource

FROM resource AS tests
COPY --from=builder /tests /tests
RUN set -e; for test in /tests/*.test; do \
		$test; \
	done

FROM resource
