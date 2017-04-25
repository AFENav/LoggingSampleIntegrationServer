package main

import (
  "github.com/antchfx/xquery/xml"
  "github.com/jclement/multiwritercloser"
  "io"
  "log"
  "net/http"
  "os"
  "path"
  "time"
)

const listen = "127.0.0.1:7878"
const eventsPath = "events"
const validationsPath = "validations"
const numberingPath = "numbering"
const validationResponseFilename = "validation_response.xml"

func writeEmptySoapBody(w io.Writer) {
  io.WriteString(w, "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\">")
  io.WriteString(w, "<soap:Body>")
  io.WriteString(w, "</soap:Body>")
  io.WriteString(w, "</soap:Envelope>")
}

func processAndLogRequest(logPath string, w http.ResponseWriter, r *http.Request) (*xmlquery.Node, io.WriteCloser, string, error) {

  // generate a timestamp based filename
  filenameBase := path.Join(logPath, time.Now().Format("2006-01-02T15-04-05Z07-00"))
  filenameRequest := filenameBase + "_request.xml"
  filenameResponse := filenameBase + "_response.xml"

  // create a file to log XML to
  f, err := os.Create(filenameRequest)
  if err != nil {
    return nil, nil, "", err
  }

  defer r.Body.Close()
  defer f.Close()

  // stream XML to file and XML reader
  tmp := io.TeeReader(r.Body, f)

  // parse XML request
  root, err := xmlquery.Parse(tmp)
  if err != nil {
    return nil, nil, "", err
  }

  responseFile, err := os.Create(filenameResponse)
  if err != nil {
    return nil, nil, "", err
  }

  output := multiwritercloser.MultiWriterCloser(responseFile, w)

  return root, output, filenameBase, nil
}

// processEvent handles inbound event calls and logs them to a file
func processEvent(w http.ResponseWriter, r *http.Request) {

  root, output, filename, err := processAndLogRequest(eventsPath, w, r)
  defer output.Close()

  if err != nil {
    log.Panicln(err)
    return
  }

  // read the event type from the XML body // /Envelope/soap:Body/AFEEvent/event:Event
  eventType := ""
  if node := xmlquery.FindOne(root, "//event:Event"); node != nil {
    eventType = node.InnerText()
  }

  // pull out AFE's DocumentID
  docID := ""
  if node := xmlquery.FindOne(root, "//event:AFE/afe:DocumentID"); node != nil {
    docID = node.InnerText()
  }

  log.Printf("Event[%v on %v] - (%v %v // Action=%v) >> %v\n", eventType, docID, r.Method, r.URL.EscapedPath(), r.Header["Soapaction"][0], filename)

  writeEmptySoapBody(output)
}

func processValidation(w http.ResponseWriter, r *http.Request) {

  root, output, filename, err := processAndLogRequest(validationsPath, w, r)
  defer output.Close()

  if err != nil {
    log.Panicln(err)
    return
  }

  // pull out AFE's DocumentID
  docID := ""
  if node := xmlquery.FindOne(root, "//validate:AFE/afe:DocumentID"); node != nil {
    docID = node.InnerText()
  }

  log.Printf("Validation[%v] - (%v %v // Action=%v) >> %v\n", docID, r.Method, r.URL.EscapedPath(), r.Header["Soapaction"][0], filename)

  f, err := os.Open(validationResponseFilename)
  if err != nil {
    log.Panicln(err)
    return
  }

  defer f.Close()
  io.Copy(output, f)
}

func processNumbering(w http.ResponseWriter, r *http.Request) {

  root, output, filename, err := processAndLogRequest(numberingPath, w, r)
  defer output.Close()

  if err != nil {
    log.Panicln(err)
    return
  }

  afenumber := time.Now().Format("2006-01-02T15-04-05Z07-00")

  // pull out AFE's DocumentID
  docID := ""
  if node := xmlquery.FindOne(root, "//AFE/afe:DocumentID"); node != nil {
    docID = node.InnerText()
  }

  log.Printf("Numbering[%v on %v] - (%v %v // Action=%v) >> %v\n", afenumber, docID, r.Method, r.URL.EscapedPath(), r.Header["Soapaction"][0], filename)

  io.WriteString(output, "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\">")
  io.WriteString(output, "<soap:Body>")
  io.WriteString(output, "<AFENumberResult xmlns=\"http://energynavigator.com/xml/afe/afenumber-result/2\" xmlns:xsi2001=\"http://www.w3.org/2001/XMLSchema-instance\" xsi2001:schemaLocation=\"http://energynavigator.com/xml/afe/afenumber-result/2 afe-afenumber-result.xsd\">")
  io.WriteString(output, "<Success><AFENumber>"+afenumber+"</AFENumber></Success></AFENumberResult>")
  io.WriteString(output, "</soap:Body>")
  io.WriteString(output, "</soap:Envelope>")
}

func main() {
  log.SetOutput(os.Stdout)

  log.Println("AFE Nav Integration Test Server")
  log.Println("======================================")
  log.Println("Listening: ", listen)

  // create output folders if they don't exist
  for _, outputPath := range []string{eventsPath, validationsPath, numberingPath} {
    if _, err := os.Stat(outputPath); err != nil {
      if os.IsNotExist(err) {
        log.Printf("Creating '%s' folder\n", outputPath)
        if err = os.Mkdir(outputPath, os.ModeDir); err != nil {
          panic("Unable to create folder: " + outputPath)
        }
      }
    }
  }

  // verify that validation response file exists
  if _, err := os.Stat(validationResponseFilename); err != nil {
    if os.IsNotExist(err) {
      panic("validation response file does not exist: " + validationResponseFilename)
    }
  }

  http.HandleFunc("/event/", processEvent)
  http.HandleFunc("/validate/", processValidation)
  http.HandleFunc("/numbering/", processNumbering)

  log.Fatal(http.ListenAndServe(listen, nil))
}
