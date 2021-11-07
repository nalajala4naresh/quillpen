FROM scratch
WORKDIR /code
COPY main .

ENTRYPOINT ["./main"]