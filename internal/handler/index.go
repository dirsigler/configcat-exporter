package handler

import "net/http"

// IndexHandler provides an index page with reference links
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	msg := `
	<html>
		<head>
    		<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>ConfigCat Exporter</title>
		</head>
		<body>
			<h1>ConfigCat Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			<p><a href="/health">Healthcheck</a></p>
		</body>
	</html>
			`

	w.WriteHeader(http.StatusOK)
	//nolint:errcheck
	w.Write([]byte(msg))
}
