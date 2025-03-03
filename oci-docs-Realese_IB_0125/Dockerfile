FROM golang:1.18 as builder

RUN apt-get update && apt-get install -y libltdl-dev libpcsclite-dev openssl ca-certificates \
   fontconfig fontconfig-config fonts-dejavu-core libfontconfig1 libfontenc1 \
  libfreetype6 libjpeg62-turbo libpng16-16 libx11-6 libx11-data libxau6 \
  libxcb1 libxdmcp6 libxext6 libxrender1 x11-common xfonts-75dpi xfonts-base \
  xfonts-encodings xfonts-utils

RUN echo "pre install" && date
RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb
RUN apt install ./wkhtmltox_0.12.6-1.buster_amd64.deb


ENV CGO_ENABLED 1

COPY . /src

WORKDIR /src

RUN echo "pre build" && date
RUN go build -v -o ./oci-docs ./main.go
RUN echo "post build" && date


FROM debian:10-slim

RUN sed -i'.bak' 's/$/ contrib/' /etc/apt/sources.list
RUN apt-get update && apt-get install -y libltdl-dev libpcsclite-dev ca-certificates wget \
 fontconfig fontconfig-config fonts-dejavu-core libfontconfig1 libfontenc1 \
  libfreetype6 libjpeg62-turbo libpng16-16 libx11-6 libx11-data libxau6 \
  libxcb1 libxdmcp6 libxext6 libxrender1 x11-common xfonts-75dpi xfonts-base \
  xfonts-encodings xfonts-utils ttf-mscorefonts-installer ttf-mscorefonts-installer fontconfig

RUN wget https://raw.githubusercontent.com/mrbvrz/segoe-ui-linux/master/install.sh && chmod +x install.sh && ./install.sh

RUN wget https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb
RUN apt install ./wkhtmltox_0.12.6-1.buster_amd64.deb

RUN ln -sf /usr/share/zoneinfo/Asia/Aqtobe /etc/localtime

ENV CGO_ENABLED 1
#ENV LD_LIBRARY_PATH :/opt/kalkancrypt/:/opt/kalkancrypt/lib/engines

COPY ./data/sdk/usr/lib /usr/lib
COPY ./data/opt /opt
COPY ./data/ca /opt/ca

RUN mkdir -p /usr/local/share/ca-certificates/extra
COPY ./data/ca/production/extra /usr/local/share/ca-certificates/extra/
COPY ./data/ca/production/certs /etc/ssl/certs/
COPY ./data/ca/test/extra /usr/local/share/ca-certificates/extra/
COPY ./data/ca/test/certs /etc/ssl/certs/
RUN update-ca-certificates


WORKDIR /srv
COPY --from=builder /src/oci-docs /srv/app
COPY --from=builder /src/ca /srv/ca

EXPOSE 8091

ENTRYPOINT ["/srv/app"]
