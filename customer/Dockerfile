FROM openjdk:8 AS build-java

# Install maven
RUN echo "deb [check-valid-until=no] http://cdn-fastly.deb.debian.org/debian bullseye main" > /etc/apt/sources.list.d/bullseye.list
RUN sed -i '/deb http:\/\/deb.debian.org\/debian bullseye-updates main/d' /etc/apt/sources.list
RUN apt-get -o Acquire::Check-Valid-Until=false update && apt-get install -y maven

WORKDIR /code

# # Prepare by downloading dependencies
ADD pom.xml /code/pom.xml
RUN ["mvn", "dependency:resolve"]
# RUN ["mvn", "verify"]

# Adding source, compile and package into a fat jar
# This assumes you've configured such a goal in pom.xml
ADD src /code/src
RUN ["mvn", "package"]

FROM openjdk:8-jdk-alpine

COPY --from=build-java /code/target/*.jar app.jar
ENTRYPOINT ["java","-jar","/app.jar"]