FROM ubuntu
COPY se_project_apiserver /bin/apiserver
COPY env /env
EXPOSE 8080
RUN /bin/bash -c "source /env"
ENTRYPOINT ["/bin/apiserver"]
