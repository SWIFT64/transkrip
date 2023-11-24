package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/aiteung/athelper/fiber"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Config struct to hold configuration values
type Config struct {
	APIURL      string `json:"apiURL"`
	LoginToken  string `json:"loginToken"`
}

type Data struct {
	NamaMhs           string        `json:"nama_mhs"`
	NomorIndukMhs     int           `json:"nomor_induk_mhs"`
	TtlMhs            string        `json:"ttl_mhs"`
	TtlMhsEng         string        `json:"ttl_mhs_eng"`
	TahunMasukMhs     string        `json:"tahun_masuk_mhs"`
	FakultasMhs       string        `json:"fakultas_mhs"`
	FakultasMhsEng    string        `json:"fakultas_mhs_eng"`
	ProdiMhs          string        `json:"prodi_mhs"`
	ProdiMhsEng       string        `json:"prodi_mhs_eng"`
	NoTranskrip       string        `json:"no_transkrip"`
	Subjects          []Subject     `json:"subjects"`
	CreditsTotal      int           `json:"credits_total"`
	GradeTotal        float64       `json:"grade_total"`
	GraduationDate    string        `json:"graduation_date"`
	GraduationDateEng string        `json:"graduation_date_eng"`
	PredikatMhs       string        `json:"predikat_mhs"`
	PredikatMhsEng    string        `json:"predikat_mhs_eng"`
	JudulSkriptsi     JudulSkriptsi `json:"judul_skriptsi"`
	TempatTerbit      string        `json:"tempat_terbit"`
	TanggalTerbit     string        `json:"tanggal_terbit"`
	TanggalTerbitEng  string        `json:"tanggal_terbit_eng"`
	NamaDekan         string        `json:"nama_dekan"`
	NikDekan          int           `json:"nik_dekan"`
}

type Subject struct {
	Index       int    `json:"index"`
	Subjname    string `json:"subjname"`
	Subjnameeng string `json:"subjnameeng"`
	Credits     int    `json:"credits"`
	Grade       string `json:"grade"`
}

type JudulSkriptsi struct {
	JudulIndonesia string `json:"judul_indonesia"`
	JudulInggris   string `json:"judul_inggris"`
}

func main() {
	// Read configuration from file
	config, err := readConfig("config.json")
	if err != nil {
		log.Println("Error reading configuration:", err)
		return
	}
		// Replace the mock API URL with the Ulbi API URL

	apiURL := config.APIURL

	// Include the token in the request headers
	request, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		log.Println("Error creating API request:", err)
		return
	}
	// Set the authorization header with the provided token
	request.Header.Set("LOGIN", config.LoginToken)

	// Send the request and get the response
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error fetching data from API:", err)
		return
	}
	defer response.Body.Close()

	// Read the JSON response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading API response:", err)
		return
	}

	// Parse JSON data into a single Data struct
	fmt.Printf("%s\n\n", body)
	var raw fiber.ReturnData[Data]
	err = json.Unmarshal(body, &raw)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return
	}
	student := raw.Data
	// Calculate total credits
	var creditsTotal int
	for _, subject := range student.Subjects {
		creditsTotal += subject.Credits
	}
	student.CreditsTotal = creditsTotal

	// Increment the index for each subject to make it 1-based
	for i := range student.Subjects {
		student.Subjects[i].Index = i + 1
	}

	// Read the HTML template from the file
	templateFile := "template/transcript.html"
	htmlTemplate, err := ioutil.ReadFile(templateFile)
	if err != nil {
		panic(err)
	}

	// Parse the HTML template
	tmpl, err := template.New("transcript").Parse(string(htmlTemplate))
	if err != nil {
		panic(err)
	}

	// Create a new file to write the output
	outputFileName := "output_" + strconv.Itoa(student.NomorIndukMhs) + ".html"
	outputFile, err := os.Create(outputFileName)
	if err != nil {
		log.Println("Error creating output file:", err)
		return
	}
	defer outputFile.Close()

	// Execute the template and write the HTML content to the file
	err = tmpl.Execute(outputFile, student)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}

	// Get the absolute path of the HTML file
	absPath, err := filepath.Abs(outputFile.Name())
	if err != nil {
		log.Println("Error getting absolute path:", err)
		return
	}

	// Print the file path for debugging
	log.Println("HTML file path:", absPath)

	// Create a context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// Open the HTML file using chromedp
	err = chromedp.Run(ctx,
		chromedp.Navigate("file://"+absPath),
		chromedp.WaitReady("body"),
	)
	if err != nil {
		log.Println("Error navigating to HTML file:", err)
		return
	}

	// Print the page to PDF
	var pdfData []byte
	err = chromedp.Run(ctx,
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = page.PrintToPDF().Do(ctx)
			return err
		}),
	)
	if err != nil {
		log.Println("Error printing to PDF:", err)
		return
	}

	// Write the PDF data to a file
	pdfFileName := "output_" + strconv.Itoa(student.NomorIndukMhs) + ".pdf"
	err = ioutil.WriteFile(pdfFileName, pdfData, 0644)
	if err != nil {
		log.Println("Error writing PDF file:", err)
		return
	}

	log.Printf("PDF file generated successfully for student %d.\n", student.NomorIndukMhs)
}

// Function to read configuration from a file
func readConfig(filename string) (Config, error) {
	var config Config

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}