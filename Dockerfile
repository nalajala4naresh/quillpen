FROM scratch
WORKDIR /code
COPY quillpen .

ENTRYPOINT ["./quillpen"]