# Dotnet Core Sample App

## Building

`pack build dotnet-runtime-sample --buildpack gcr.io/paketo-buildpacks/dotnet-core`

## Running

`docker run --interactive --tty --env PORT=8080 --publish 8080:8080 dotnet-runtime-sample`

## Viewing

`curl http://localhost:8080`
