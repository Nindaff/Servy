Servy
=============

Make your cwd into your server root in one command!

## Installation
  Servy is not compiled for windows yet, but if you have Go setup just use `go get`. 
  ```sh
  $ go get github.com/Nindaff/Servy
  ```
  # Mac, Linux
  Download one of the <a href="https://github.com/Nindaff/Servy/releases">binaries</a>
  ```sh
  $ sudo chmod +x servy && sudo mv servy /usr/local/bin
  ```
  Note if you download the `servy_mac` binary you will want to run this command instead
  ```sh
  $ sudo chmod +x servy_mac && sudo mv servy_mac /usr/local/bin/servy
  ```

## Usage
### Set the server root.
  ```sh
  $ servy -dir /your/directory
  ```
  Default will always be your current working directory.
### Set the port
  ```sh
  $ servy -p 8000
  ```
  Default port is 8080.
### Set the index html file
  ```sh
  $ servy -i default.html
  ```
### Force no caching
  ```sh
  $ servy -nc 
  ```
  Servy doesn't actually cache any content in memory, but you can force Servy to send all files with a spoofed modification time so that the browser will not respond with a previously cached resource.
  By default all files are sent with a spoofed modtime the first time being served on a session.

## Note
  static server only accepts GET requests, all other methods will recieve a 405 response status
