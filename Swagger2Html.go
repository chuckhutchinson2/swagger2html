package main

import "flag"
import "fmt"
import "encoding/json"
import "io/ioutil"
import "net/http"
import "html/template"
import "os"

func buildEndpoint(scheme string, host string) string {
	return fmt.Sprintf("%s://%s/v2/swagger.json", scheme, host)
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

func prettyJson(jsonData map[string]interface{}) string {

	bytes, err := json.MarshalIndent(jsonData, "", "   ")
	if err == nil {
		return string(bytes)
	}

	return fmt.Sprintf("Error %s", err)
}

func generateSwaggerHtml(jsonData map[string]interface{}, swaggerHtmlTemplate string) {
	t, err := template.ParseFiles(swaggerHtmlTemplate)

	if err == nil {
		t.Execute(os.Stdout, jsonData)
	}
}

func main () {

	hostPtr := flag.String("host", "petstore.swagger.io", "Swagger Host")
	schemePtr := flag.String("scheme", "https", "Scheme")
	swaggerHtmlTemplatePtr := flag.String("template", "swaggerTemplate.html", "Swagger HTML template")
	verbosePtr := flag.Bool("verbose", false, "Verbose output")
	endpointPtr := flag.String("endpoint", "", "Endpoint override")

	flag.Parse();

	endpoint := buildEndpoint(*schemePtr, *hostPtr)

	if len(*endpointPtr) > 0 {
		endpoint = *endpointPtr
	}

	jsonData := callEndpoint(endpoint)

	if (*verbosePtr) {
		fmt.Println("Host:", *hostPtr)
		fmt.Println("Scheme:", *schemePtr)

		fmt.Println("Endpoint: ", endpoint)

		fmt.Println("Swagger HTML Template", *swaggerHtmlTemplatePtr)

		fmt.Println("Data: ", prettyJson(jsonData))
	}

	generateSwaggerHtml(jsonData, *swaggerHtmlTemplatePtr)
}




