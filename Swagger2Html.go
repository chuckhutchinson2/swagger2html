package main

import "flag"
import "fmt"
// import "bytes"
import "encoding/json"
import "io/ioutil"
import "net/http"


func buildEndpoint(protocol string, host string) string {
	return fmt.Sprintf("%s://%s/v2/swagger.json", protocol, host)
}

func callEndpoint(endpoint string) map[string]interface{} {

	response, err := http.Get(endpoint)

	if err != nil {
		fmt.Printf("Error calling %s Error: %s\n", endpoint, err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		var jsonData map[string]interface{}
		json.Unmarshal(data, &jsonData)

		return jsonData
	}

	return nil
}


func main () {

	hostPtr := flag.String("host", "petstore.swagger.io", "Swagger Host")
	protocolPtr := flag.String("protocol", "https", "Protocol")

	flag.Parse();

	fmt.Println("Host:", *hostPtr)
	fmt.Println("Protocol:", *protocolPtr)

	endPoint := buildEndpoint(*protocolPtr, *hostPtr)

	fmt.Println("Endpoint: ", endPoint)

	fmt.Println("Data: ", callEndpoint(endPoint))
}




