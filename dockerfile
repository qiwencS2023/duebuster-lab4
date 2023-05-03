FROM ubuntu

# update/upgrade
RUN apt update -y; apt upgrade -y
RUN apt install -y make bc build-essential libssl-dev libffi-dev

# Install Java
RUN apt install -y openjdk-19-jdk openjdk-19-jre

RUN apt install -y golang-go

WORKDIR /project

CMD ["/bin/bash", "echo", "14736"]
