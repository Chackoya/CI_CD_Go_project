/*
- reqhttp.go
Functions in this file take care of making HTTP GET requests and error handling. This isolates HTTP request logic from other application logic.

*/
package main

/*
// makeGETRequest performs a GET request to the given URL and returns the response body as a byte slice.
func makeGETRequest(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	return body
}
*/
