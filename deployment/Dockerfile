FROM centos:7.2.1511

RUN mkdir /opt/tulc
RUN touch /opt/tulc/serverlist

RUN mkdir -p /opt/tulc
ADD application.yml /opt/tulc/application.yml
ADD serverlist /opt/tulc/serverlist
ADD serverlistforwebpay /opt/tulc/serverlistforwebpay

RUN mkdir /app
WORKDIR /app
ADD tulc-api ./tulc-api-server

RUN rm -f etc/localtime; ln -sf /usr/share/zoneinfo/Asia/Ho_Chi_Minh /etc/localtime
EXPOSE 8088/tcp
# Run
CMD ["/app/tulc-api-server", "-m", "http"]
