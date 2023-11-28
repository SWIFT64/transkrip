package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/cdp"

	"github.com/aiteung/athelper/fiber"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// Config struct to hold configuration values
type Config struct {
	APIURL     string `json:"apiURL"`
	LoginToken string `json:"loginToken"`
}

type Data struct {
	NamaMhs           string        `json:"nama_mhs"`
	NomorIndukMhs     string          `json:"nomor_induk_mhs"`
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
	NikDekan          string        `json:"nik_dekan"`
}

type PageOnePDF struct {
	NamaMhs        string    `json:"nama_mhs"`
	NomorIndukMhs  string       `json:"nomor_induk_mhs"`
	TtlMhs         string    `json:"ttl_mhs"`
	TtlMhsEng      string    `json:"ttl_mhs_eng"`
	TahunMasukMhs  string    `json:"tahun_masuk_mhs"`
	FakultasMhs    string    `json:"fakultas_mhs"`
	FakultasMhsEng string    `json:"fakultas_mhs_eng"`
	ProdiMhs       string    `json:"prodi_mhs"`
	ProdiMhsEng    string    `json:"prodi_mhs_eng"`
	NoTranskrip    string    `json:"no_transkrip"`
	Subjects       []Subject `json:"subjects"`
}
type PageTwoPDF struct {
	CreditsTotal      int           `json:"credits_total"`
	GradeTotal        float64       `json:"grade_total"`
	GraduationDate    string        `json:"graduation_date"`
	GraduationDateEng string        `json:"graduation_date_eng"`
	Subjects          []Subject     `json:"subjects"`
	PredikatMhs       string        `json:"predikat_mhs"`
	PredikatMhsEng    string        `json:"predikat_mhs_eng"`
	JudulSkriptsi     JudulSkriptsi `json:"judul_skriptsi"`
	TempatTerbit      string        `json:"tempat_terbit"`
	TanggalTerbit     string        `json:"tanggal_terbit"`
	FakultasMhs       string        `json:"fakultas_mhs"`
	FakultasMhsEng string    `json:"fakultas_mhs_eng"`

	TanggalTerbitEng string `json:"tanggal_terbit_eng"`
	NamaDekan        string `json:"nama_dekan"`
	NoTranskrip      string `json:"no_transkrip"`
	NikDekan         string    `json:"nik_dekan"`
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
	//fmt.Printf("%s\n\n", body)
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

	pdfPageOne := PageOnePDF{
		NamaMhs:        student.NamaMhs,
		NomorIndukMhs:  student.NomorIndukMhs,
		TtlMhs:         student.TtlMhs,
		TtlMhsEng:      student.TtlMhsEng,
		TahunMasukMhs:  student.TahunMasukMhs,
		FakultasMhs:    student.FakultasMhs,
		FakultasMhsEng: student.FakultasMhsEng,
		ProdiMhs:       student.ProdiMhs,
		ProdiMhsEng:    student.ProdiMhsEng,
		NoTranskrip:    student.NoTranskrip,
	}

	pdfPageTwo := PageTwoPDF{
		CreditsTotal:      student.CreditsTotal,
		GradeTotal:        student.GradeTotal,
		GraduationDate:    student.GraduationDate,
		GraduationDateEng: student.GraduationDateEng,
		PredikatMhs:       student.PredikatMhs,
		PredikatMhsEng:    student.PredikatMhsEng,
		JudulSkriptsi:     student.JudulSkriptsi,
		TempatTerbit:      student.TempatTerbit,
		TanggalTerbit:     student.TanggalTerbit,
		TanggalTerbitEng:  student.TanggalTerbitEng,
		NamaDekan:         student.NamaDekan,
		NoTranskrip:       student.NoTranskrip,
		FakultasMhs:       student.FakultasMhs,
		FakultasMhsEng:    student.FakultasMhsEng,
		NikDekan:          student.NikDekan,
	}

	switch len(student.Subjects) > 43 {
	case true:
		pdfPageOne.Subjects = student.Subjects[:45]
		pdfPageTwo.Subjects = student.Subjects[44:]
	case false:
		pdfPageOne.Subjects = student.Subjects
	}

	// Read the HTML template from the file
	templateFile1 := "template/page_1.html"
	templateFile2 := "template/page_2.html"
	htmlTemplate1, err := os.ReadFile(templateFile1)
	if err != nil {
		panic(err)
	}

	htmlTemplate2, err := os.ReadFile(templateFile2)
	if err != nil {
		panic(err)
	}

	// Parse the HTML template
	tmpl1, err := template.New("page_1").Parse(string(htmlTemplate1))
	if err != nil {
		panic(err)
	}
	// Parse the HTML template
	tmpl2, err := template.New("page_2").Parse(string(htmlTemplate2))
	if err != nil {
		panic(err)
	}

	// Create a new file to write the output
	//outputFileName := "output_" + strconv.Itoa(student.NomorIndukMhs) + ".html"
	//outputFile, err := os.Create(outputFileName)
	//if err != nil {
	//	log.Println("Error creating output file:", err)
	//	return
	//}
	//defer outputFile.Close()
	//
	//// Create a new file to write the output
	//outputFileName2 := "output2_" + strconv.Itoa(student.NomorIndukMhs) + ".html"
	//outputFile2, err := os.Create(outputFileName2)
	//if err != nil {
	//	log.Println("Error creating output file:", err)
	//	return
	//}
	//defer outputFile.Close()

	pdf1 := new(bytes.Buffer)
	pdf2 := new(bytes.Buffer)

	// Execute the template and write the HTML content to the file
	err = tmpl1.Execute(pdf1, pdfPageOne)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}
	// Execute the template and write the HTML content to the file
	err = tmpl2.Execute(pdf2, pdfPageTwo)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}

	// Get the absolute path of the HTML file
	//absPath, err := filepath.Abs(outputFile.Name())
	//if err != nil {
	//	log.Println("Error getting absolute path:", err)
	//	return
	//}
	//
	//// Print the file path for debugging
	//log.Println("HTML file path:", absPath)

	// Create a context
	opt := append(chromedp.DefaultExecAllocatorOptions[:], chromedp.ExecPath("/usr/bin/chromium"))

	allocatorCtx, allocatorCancel := chromedp.NewExecAllocator(
		context.Background(),
		opt...,
	)
	defer allocatorCancel()

	ctx, cancel := chromedp.NewContext(allocatorCtx)
	defer cancel()

	// Open the HTML file using chromedp
	//err = chromedp.Run(ctx,
	//	chromedp.Navigate("file://"+absPath),
	//	chromedp.WaitReady("body"),
	//)
	//if err != nil {
	//	log.Println("Error navigating to HTML file:", err)
	//	return
	//}

	// Print the page to PDF
	var pdfData []byte
	var pdfData2 []byte
	//err = chromedp.Run(ctx,
	//	chromedp.ActionFunc(func(ctx context.Context) error {
	//		var err error
	//		pdfData, _, err = page.PrintToPDF().Do(ctx)
	//		return err
	//	}),
	//)

	// Generate Page 1
	err = chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(c context.Context) (err error) {
			chromepCtx := chromedp.FromContext(c)
			executor := cdp.WithExecutor(c, chromepCtx.Target)
			frame, err := page.GetFrameTree().Do(executor)
			if err != nil {
				fmt.Printf("error frame tree %+v \n ", err)
				return
			}

			err = page.SetDocumentContent(frame.Frame.ID, pdf1.String()).Do(executor)
			return
		}),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(c context.Context) (err error) {
			chromepCtx := chromedp.FromContext(c)
			pdfData, _, err = page.PrintToPDF().Do(cdp.WithExecutor(c, chromepCtx.Target))
			return
		}),
	)
	if err != nil {
		log.Println("Error printing to PDF:", err)
		return
	}

	// Generate Page 2
	err = chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(c context.Context) (err error) {
			chromepCtx := chromedp.FromContext(c)
			executor := cdp.WithExecutor(c, chromepCtx.Target)
			frame, err := page.GetFrameTree().Do(executor)
			if err != nil {
				fmt.Printf("error frame tree %+v \n ", err)
				return
			}

			err = page.SetDocumentContent(frame.Frame.ID, pdf2.String()).Do(executor)
			return
		}),
		chromedp.WaitReady("body"),
		chromedp.ActionFunc(func(c context.Context) (err error) {
			chromepCtx := chromedp.FromContext(c)
			pdfData2, _, err = page.PrintToPDF().Do(cdp.WithExecutor(c, chromepCtx.Target))
			return
		}),
	)
	if err != nil {
		log.Println("Error printing to PDF:", err)
		return
	}

	// Write the PDF data to a file
	pdfFileName := "output_" + student.NomorIndukMhs + ".pdf"
	err = os.WriteFile(pdfFileName, pdfData, 0644)
	if err != nil {
		log.Println("Error writing PDF file:", err)
		return
	}
	pdfFileName2 := "output2_" + student.NomorIndukMhs + ".pdf"
	err = os.WriteFile(pdfFileName2, pdfData2, 0644)
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
