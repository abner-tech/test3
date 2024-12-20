package main

import (
	"flag"
	"log"
	"net/http"
)

//create a simple HTML page with JS added.
// However, in a professional settings,we would have the JS code in a script file

const html = `
<!DOCTYPE html>
<html lang="en">
<head>

	<meta charset="UTF-8">
</head>

<body>
	<h1>
		Appletree CORS
	</h1>
	<div id="output"></div>

	<script>
		document.addEventListener('DOMContentLoaded', function() {
			fetch("http://localhost:4000/api/v1/healthcheck")
				.then(function(response) {
					response.text().then(function (text) {
						document.getElementById("output").innerHTML = text;
					});
				}, 
				function(err) {
				document.getElementById("output").innerHTML = err ;
				}
			);
		});
	</script>
</body>
</html>
`

// a very simple HTTP server
func main() {
	addr := flag.String("addr", ":9000", "Server Address")
	flag.Parse()

	log.Printf("starting Server on %s", *addr)
	err := http.ListenAndServe(*addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(html))
	}))
	log.Fatal(err)
}
