package main

import "flag"
import "fmt"
import "encoding/json"
import "io/ioutil"
import "net/http"
import "html/template"
import "os"


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

func loadHtmlTemplate(htmlTemplate string) *template.Template {
	t, err := template.ParseFiles(htmlTemplate)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	return t
}

func generateSwaggerHtml(jsonData map[string]interface{}, htmlTemplate string) {
	t := loadHtmlTemplate(htmlTemplate)

	t.Execute(os.Stdout, jsonData)
}

type jsonTemplateHandler struct {
    endpoint string
    template *template.Template
}

func (h *jsonTemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.template.Execute(w, callEndpoint(h.endpoint))
}


func main () {

	htmlTemplatePtr := flag.String("template", "swaggerTemplate.html", "Swagger HTML template")
	verbosePtr := flag.Bool("verbose", false, "Verbose output")
	endpointPtr := flag.String("endpoint", "https://petstore.swagger.io/v2/swagger.json", "Endpoint")
	prettyPrintPtr := flag.Bool("pretty", false, "Pretty print JSON")
	webServerPtr := flag.Bool("server", false, "Run Web Server")
	httpPortPtr := flag.String("httpPort", ":8080", "WebServer")

	flag.Parse();

	endpoint := *endpointPtr

	jsonData := callEndpoint(endpoint)

	if (*verbosePtr) {
		fmt.Println("Endpoint: ", endpoint)

		fmt.Println("HTML Template", *htmlTemplatePtr)
	}

	if *webServerPtr {
		fmt.Println("Running web server running on port", *httpPortPtr)
		http.Handle("/", &jsonTemplateHandler{endpoint: endpoint, template: loadHtmlTemplate(*htmlTemplatePtr)})
    	http.ListenAndServe(*httpPortPtr, nil)
	} else if *prettyPrintPtr {
		fmt.Println(prettyJson(jsonData))
	} else {
		generateSwaggerHtml(jsonData, *htmlTemplatePtr)
	}
}




