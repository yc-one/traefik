FROM node:6.16.0

ENV WEBUI_DIR /src/webui
RUN mkdir -p $WEBUI_DIR

RUN rm -fr /var/lib/apt/lists/

COPY package.json $WEBUI_DIR/
COPY yarn.lock $WEBUI_DIR/

WORKDIR $WEBUI_DIR
RUN npm set progress=false
RUN yarn install

COPY . $WEBUI_DIR/

EXPOSE 8080
