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

type jsonTemplateHandler struct {
    jsonData map[string]interface{}
    templateFilename string
}

func (h *jsonTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    t, err := template.ParseFiles(h.templateFilename)

	if err == nil {
		t.Execute(w, h.jsonData)
	}
}


func main () {

	hostPtr := flag.String("host", "petstore.swagger.io", "Swagger Host")
	schemePtr := flag.String("scheme", "https", "Scheme")
	swaggerHtmlTemplatePtr := flag.String("template", "swaggerTemplate.html", "Swagger HTML template")
	verbosePtr := flag.Bool("verbose", false, "Verbose output")
	endpointPtr := flag.String("endpoint", "", "Endpoint override")
	prettyPrintPtr := flag.Bool("pretty", false, "Pretty print JSON")
	webServerPtr := flag.Bool("server", false, "Run Web Server")
	httpPortPtr := flag.String("httpPort", ":8080", "WebServer")

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
	}

	if *webServerPtr {
		fmt.Println("Running web server running on port", *httpPortPtr)
		http.Handle("/", &jsonTemplateHandler{jsonData: jsonData, templateFilename: *swaggerHtmlTemplatePtr})
    	http.ListenAndServe(*httpPortPtr, nil)
	} else if *prettyPrintPtr {
		fmt.Println(prettyJson(jsonData))
	} else {
		generateSwaggerHtml(jsonData, *swaggerHtmlTemplatePtr)
	}
}




