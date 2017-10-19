FROM mongo:3.4.9-jessie

#RUN apt-get update \
#	&& apt-get install -y --no-install-recommends \
#		wget \
#		git \
#	&& rm -rf /var/lib/apt/lists/*

#RUN wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz \
#	&& tar -C /usr/local -zxf go1.9.1.linux-amd64.tar.gz \
#	&& rm go1.9.1.linux-amd64.tar.gz \
#	&& echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile

ENV PATH="${PATH}:/usr/local/go/bin"

COPY docker/init /srv/init

COPY dump /srv/dump

COPY webapp /srv/app/src/eaciit/apichain/webapp

RUN export GOPATH=/srv/app/ \
	&& cd /srv/app/src/eaciit/apichain/webapp \
	&& go get ... \
	&& go build .

COPY docker/config/app.json /srv/app/src/eaciit/apichain/config/app.json

EXPOSE 9150

CMD ["/srv/init/init.sh"]
