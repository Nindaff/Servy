package main

import (
	"os"
	"path/filepath"

	"github.com/Nindaff/Servy/server"
	"github.com/codegangsta/cli"
)

const errMsg = `
<!DOCTYPE html>
<html>
<head><title>Servy Error</title>
<style>
body,html{height:100%;width:100%}body{background:#6495ed;font-family:sans-serif}h1{color:#f8f8ff}p{color:#2E2C2F}.big{font-size:160px}.middle{display:table;margin:40px auto}
</style></head>
<body><div class="middle"><h1 class="big">Oops!</h1><h1>Servy Could not find "index.html"</h1><p>You can specify the index html file with the "-i" argument</p></div>
</body></html>
`

func main() {

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	app := cli.NewApp()
	app.Name = "Servy"
	app.Usage = "Serve up static files for Quick Tests!"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "port, p",
			Value: 8080,
			Usage: "Set the port for the server",
		},
		cli.StringFlag{
			Name:  "index, i",
			Value: "index.html",
			Usage: "Set the default file to be served at `/`",
		},
		cli.StringFlag{
			Name:  "directory, dir",
			Value: cwd,
			Usage: "Set the base directory for the server",
		},
		cli.BoolFlag{
			Name:  "no-cache, nc",
			Usage: "Force the server to send all files with spoofed modification time",
		},
	}
	app.Action = func(c *cli.Context) {
		var dir string

		if "." == c.String("dir")[:1] {
			dir = filepath.Join(cwd, c.String("dir")[1:])
		} else {
			dir = c.String("dir")
		}

		staticServer := server.NewStaticServer(dir, c.String("i"), c.Int("p"), c.Bool("nc"), errMsg)
		staticServer.Serve()
	}
	app.Run(os.Args)
}
