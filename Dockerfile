# Default to Go 1.11
ARG GO_VERSION=1.23

# Start from golang v1.11 base image
FROM golang:${GO_VERSION}-alpine AS builder

# Add Maintainer Infos
LABEL maintainer="Mohammad Kaium <mohammadkaiom79@gmail.com>"

# Create the user and group files that will be used in the running container to
# run the process as an unprivileged user.
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group

# Install the Certificate-Authority certificates for the app to be able to make
# calls to HTTPS endpoints.
RUN apk add --no-cache ca-certificates tzdata

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /src

# Import the code from the context.
COPY ./ ./

# Build the Go app
RUN CGO_ENABLED=0 GOFLAGS=-mod=vendor GOOS=linux go build -a -installsuffix 'static' -o /app cmd/*.go

######## Start a new stage from scratch #######
# Final stage: the running container.
FROM scratch AS final

# Import the user and group files from the first stage.
COPY --from=builder /user/group /user/passwd /etc/

# Import the compiled executable from the first stage.
COPY --from=builder /app /app

# As we're going to run the executable as an unprivileged user, we can't bind
# to ports below 1024.
EXPOSE 8601

# Perform any further action as an unprivileged user.
USER nobody:nobody

# Run the compiled binary.
ENTRYPOINT ["/app"]