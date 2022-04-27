FROM golang as builder

COPY . /buildfolder
WORKDIR /buildfolder
RUN go get -v -t -d ./... && \
  CGO_ENABLED=0 GOOS=linux go build -o pagerduty-slack-ic-sync

FROM scratch

COPY --from=builder /buildfolder/pagerduty-slack-ic-sync .

CMD ["./pagerduty-slack-ic-sync"]
