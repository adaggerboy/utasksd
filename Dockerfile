FROM golang:1.21-alpine AS build
LABEL stage=gobuild
RUN apk update --no-cache
WORKDIR /src
COPY . .
WORKDIR /src/cmd/utasksd
RUN go build -o /bin/utasksd

FROM scratch
COPY --from=build /bin/utasksd /bin/utasksd
CMD ["/bin/utasksd"]