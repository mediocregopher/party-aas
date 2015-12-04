package main

import (
	"fmt"
	"net/http"
)

var index = `<html>
	<head><title>PARTYPARTYPARTYPARTYPARTY</title></head>
	<body>
		<form action="/partyfy" method="post" enctype="multipart/form-data">
			<p>
    			<input type="file" name="file" />
			</p>
			<p>
				Total frames, the number of frames in the gif (aka one single spin):</br>
				<input type="text" name="totalFrames" value="20" />
			</p>
			<p>
				FPS, number of frames displayed every second.<br/>
				(Total frames) / FPS = gif duration:</br>
				<input type="text" name="fps" value="20" />
			</p>
			<p>
				Max width of the output gif:</br>
				<input type="text" name="maxWidth" value="128" />
			</p>
			<p>
				Max height of the output gif:</br>
				<input type="text" name="maxHeight" value="128" />
			</p>
    		<input type="submit" value="PARTYPARTYPARTY" />
		</form>
	</body>
</html>`

func indexHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, index)
}
