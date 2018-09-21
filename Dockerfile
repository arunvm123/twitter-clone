FROM scratch

COPY ./cmd/main .

ENTRYPOINT [ "./main" ]