FROM scratch

LABEL maintainer="Andrei Varabyeu <andrei_varabyeu@epam.com>"
LABEL version=3.1.16

ENV APP_DOWNLOAD_URL https://dl.bintray.com/epam/reportportal/3.1.16/service-index_linux_amd64

ADD ${APP_DOWNLOAD_URL} /service-index

EXPOSE 8080
ENTRYPOINT ["/service-index"]
