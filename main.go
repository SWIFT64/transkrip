package main

import (
	"context"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

type Data struct {
	NamaMhs        string
	NomorIndukMhs  int
	TtlMhs         string
	TtlMhsEng      string
	TahunMasukMhs  string
	FakultasMhs    string
	FakultasMhsEng string
	ProdiMhs       string
	NoTranskrip    string // Add this field
	Subjects       []Subject
	CreditsTotal   int // Capitalized first letter to make it accessible outside the package
	GradeTotal     float64
	GraduationDate string
	PredikatMhs    string
	JudulSkriptsi  []JudulSkriptsi
	TempatTerbit   string
	TanggalTerbit  string
	NamaDekan      string
	NikDekan       int
}

type Subject struct {
	Index       int // 1-based index
	Subjname    string
	Subjnameeng string
	Credits     int
	Grade       string
}

type JudulSkriptsi struct {
	JudulIndonesia string
	JudulInggris   string
}

func main() {
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

	// Sample data representing student information and grades
	data := Data{
		NamaMhs:        "John Doe",
		NomorIndukMhs:  813619637,
		TtlMhs:         "Jakarta, 23 Januari 2000",
		TtlMhsEng:      "Jakarta, 23 January 2000",
		TahunMasukMhs:  "2019/2020 Ganjil",
		FakultasMhs:    "Sekolah Vokasi",
		FakultasMhsEng: "Vocation School",
		ProdiMhs:       "D4 Teknik Informatika",
		NoTranskrip:    "12345",
		GradeTotal:     3.9,
		GraduationDate: "Januari 2, 2023",
		PredikatMhs:    "Dengan Pujian",
		TempatTerbit:   "Bandung",
		TanggalTerbit:  "21 April 2021",
		NamaDekan:      "Aliffathur M. R.",
		NikDekan:       763784563,
		Subjects: []Subject{
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
			{0, "Matematika", "Math", 4, "A"},
			{0, "Sains", "Science", 3, "B+"},
			{0, "Sejarah", "History", 2, "A-"},
		},
		JudulSkriptsi: []JudulSkriptsi{
			{"AIUEO", "AIUEO"},
			{"BBBBB", "BBBBB"},
			{"BBBBB", "BBBBB"},
			{"BBBBB", "BBBBB"},
		},
	}

	// Calculate total credits
	var creditsTotal int
	for _, subject := range data.Subjects {
		creditsTotal += subject.Credits
	}
	data.CreditsTotal = creditsTotal

	// Increment the index for each subject to make it 1-based
	for i := range data.Subjects {
		data.Subjects[i].Index = i + 1
	}

	// Create a new file to write the output
	outputFile, err := os.Create("output.html")
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}

	println("HTML file generated successfully.")

	// Create a new file to write the HTML output
	htmlOutputFile, err := os.Create("output.html")
	if err != nil {
		panic(err)
	}
	defer htmlOutputFile.Close()

	// Execute the template and write the HTML content to the file
	err = tmpl.Execute(htmlOutputFile, data)
	if err != nil {
		log.Println("Error executing template:", err)
		return
	}

	// Get the absolute path of the HTML file
	absPath, err := filepath.Abs(htmlOutputFile.Name())
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
	err = ioutil.WriteFile("output.pdf", pdfData, 0644)
	if err != nil {
		log.Println("Error writing PDF file:", err)
		return
	}

	println("PDF file generated successfully.")
}
