############################
# STEP 1 build executable binary
############################

FROM golang AS builder

WORKDIR /app

COPY .. .

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -o main main.go

############################
# STEP 2 build a small image
############################
FROM public.ecr.aws/lambda/provided:al2023

COPY --from=builder /app/main ./main

ENTRYPOINT [ "./main" ]
