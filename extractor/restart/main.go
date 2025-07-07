package main

import (
	"bufio"
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"path/filepath"
	"restart/internal/repository"
	"strings"
)

type PdfFile struct {
	path string
	name string
}

func getExamFiles(dir string) (map[string]PdfFile, error) {
	// string: name
	pdfFiles := make(map[string]PdfFile)
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == "~BROMIUM" {
				return filepath.SkipDir
			}

			return nil
		}

		currentPdfFile := PdfFile{
			path: path,
			name: info.Name(),
		}

		pdfFiles[currentPdfFile.name] = currentPdfFile

		return nil
	})
	if err != nil {
		return nil, err
	}

	return pdfFiles, nil
}

type FileContent struct {
	SubjectName string
	year        int
	HeaderExam
	File PdfFile
}

type HeaderExam struct {
	ExamType     string `json:"exam_type"`
	MaterialType string `json:"material_type"`
	TaskLabel    string `json:"task_label"`
	Difficulty   string `json:"difficulty"`
	WorkTime     int    `json:"work_time"`
}

func ParseHeaderExam(headerString string, client *genai.Client) (HeaderExam, error) {
	/* Example layout but the layout can change from file to file
	| Zentralabitur 2017 | Biologie | Schülermaterial           |
	|--------------------|----------|---------------------------|
	| Aufgabe I          | gA       | Bearbeitungszeit: 220 min |
	*/

	responseSchema := genai.Schema{
		Type:        genai.TypeObject,
		Description: "The response schema",
		Properties: map[string]*genai.Schema{
			"exam_type": {
				Type:        genai.TypeString,
				Description: "The exam type without the year eg: Zentralabitur 2017 -> exam_type=Zentralabitur under all circumstances DO NOT return the year. If no exam_type is found return an string that is 'none'. an exam type is not for example 'Natron und Soda – Vom Hausmittel zum Schlüsselprodukt'",
			},
			"material_type": {
				Type:        genai.TypeString,
				Enum:        []string{"exam", "answer", "other"},
				Description: "Say weather the file is an exam or the answers to one, or something other like deckblätter etc. Here are some examples: Schülermaterial=exam, Lehrermaterial=answer, Material für Schülerinnen und Schüler=exam, Erwartungshorizont=answer,  Material für Prüflinge=exam etc.",
			},
			"task_label": {
				Type:        genai.TypeString,
				Description: "The task label eg: Aufgabe I or Wahlteil Rechnertyp: GTR, Prüfungsteil B Rechnertyp: CAS Analysis but not for example Betriebswirtschaft mit Rechnungswesen-Controlling.",
			},
			"difficulty": {
				Type:        genai.TypeString,
				Enum:        []string{"gA", "eA", "none"},
				Description: "The difficulty eg: gA or eA, if not decipherable return an none string",
			},
			"work_time": {
				Type:        genai.TypeInteger,
				Description: "The work time in minutes if no time is found return -1 eg: Bearbeitungszeit: 220 min -> work_time=220",
			},
		},
		Required: []string{"exam_type", "material_type", "task_label", "difficulty", "work_time"},
	}
	var thinkingBudget int32 = 0
	thinkingConfig := genai.ThinkingConfig{ThinkingBudget: &thinkingBudget}

	ctx := context.Background()
	prompt := fmt.Sprintf("you are given a header of an file. Decipher the fields from it, DO NOT RETRUN new line or tabs etc only THE TEXT:" + headerString)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite-preview-06-17",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ThinkingConfig:   &thinkingConfig,
			ResponseMIMEType: "application/json",
			ResponseSchema:   &responseSchema,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(result.Text())

	var headerExam HeaderExam
	err = json.Unmarshal([]byte(result.Text()), &headerExam)
	if err != nil {
		return HeaderExam{}, err
	}

	fmt.Println("Used tokens:", result.UsageMetadata.TotalTokenCount)

	return headerExam, nil

}

func worker(id int, jobs <-chan Job, results chan<- FileContent) {
	for j := range jobs {
		fmt.Println("worker", id, "started Job", j.file)

		file, err := os.Open(j.file.path)
		if err != nil {
			log.Fatal(err)
		}
		scanner := bufio.NewScanner(file)

		currentFileContent := FileContent{}

		searching := false

		// only scans files and uses llm if data not already generated
		if _, err := os.Stat("./file_header_mapping.json"); os.IsNotExist(err) {
			searching = true
		}

		for scanner.Scan() && searching {
			line := scanner.Text()

			// header detected
			if strings.Contains(line, "Zentralabitur") || strings.Contains(line, "Erwartungshorizont") {
				//fmt.Println("Doc:", j.file.path)

				headerString := fmt.Sprintln(line)
				// gets next line // i < 2 would be enough, but it is not always clear by only the header if a file is really an answer or exam, so we take something from the main part as well
				for i := 0; i < 7; i++ {
					scanner.Scan()
					headerString += fmt.Sprintln(scanner.Text())
				}
				// removes <br> for easier parsing for the llm and for reducing token use
				headerString = strings.ReplaceAll(headerString, "<br>", " ")
				headerString = strings.ReplaceAll(headerString, "\t", " ")
				headerString = strings.ReplaceAll(headerString, "\n", " ")
				headerString = strings.ReplaceAll(headerString, "  ", " ")
				headerString = strings.ReplaceAll(headerString, "**", "")

				header, err := ParseHeaderExam(headerString, j.llm)
				if err != nil {
					log.Fatal(err)
				}

				currentFileContent.HeaderExam = header

				// stops the scanner from scanning more lines after this
				searching = false
			}
		}
		// if not, any string checked for is found in the entire file
		if searching {
			// TODO: check with the filename for the header infos as a last resort
			currentFileContent.HeaderExam = HeaderExam{
				ExamType:     "",
				MaterialType: "other",
				TaskLabel:    "",
				Difficulty:   "",
				WorkTime:     -1,
			}
		}

		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("worker", id, "finished Job", j.file)

		currentFileContent.File = j.file

		results <- currentFileContent
	}
}

type Job struct {
	file PdfFile
	llm  *genai.Client
}

type ParseSubjectYearResponse struct {
	SubjectName string `json:"subject_name"`
	Year        int    `json:"year"`
}

var enumSubjects = []string{
	"Deutsch",
	"Englisch",
	"Französisch",
	"Spanisch",
	"Latein",
	"Griechisch",
	"Kunst",
	"Musik",
	"Erdkunde",
	"Geschichte",
	"Politik-Wirtschaft",
	"Evangelische Religion",
	"Katholische Religion",
	"Werte und Normen",
	"Mathematik",
	"Mathematik (berufliches Gymnasium)",
	"Mathematik (zweiter Bildungsweg)",
	"Mechatronik", // This subject exists two times on the official website. Why is that?
	"Biologie",
	"Chemie",
	"Physik",
	"Informatik",
	"Sport",
	"Ernährung",
	"Betriebswirtschaft mit Rechnungswesen-Controlling",
	"Pädagogik-Psychologie",
	"Betriebs- und Volkswirtschaft",
	"Volkswirtschaft",
	"Gesundheit-Pflege",
}

func ParseSubjectAndYearFromDirectoryName(directoryName string, client *genai.Client) (ParseSubjectYearResponse, error) {
	responseSchema := genai.Schema{
		Type:        genai.TypeObject,
		Description: "The response schema",
		Properties: map[string]*genai.Schema{
			"subject_name": {
				Type:        genai.TypeString,
				Enum:        enumSubjects,
				Description: "The subject name eg: 2016VW -> Volkswirtschaft",
			},
			"year": {
				Type:        genai.TypeInteger,
				Description: "The year of the exam eg: 2016VW  -> 2016",
			},
		},
		Required: []string{"subject_name", "year"},
	}

	var thinkingBudget int32 = 0
	thinkingConfig := genai.ThinkingConfig{ThinkingBudget: &thinkingBudget}

	ctx := context.Background()
	prompt := fmt.Sprintf("you are given a directory name. Decipher the fields from it:" + directoryName)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-lite-preview-06-17",
		genai.Text(prompt),
		&genai.GenerateContentConfig{
			ThinkingConfig:   &thinkingConfig,
			ResponseMIMEType: "application/json",
			ResponseSchema:   &responseSchema,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result.Text())

	fmt.Println("Used tokens:", result.UsageMetadata.TotalTokenCount)

	var resp ParseSubjectYearResponse
	err = json.Unmarshal([]byte(result.Text()), &resp)
	if err != nil {
		return ParseSubjectYearResponse{}, err
	}

	return resp, nil
}

type directorySubjectMapping map[string]ParseSubjectYearResponse

func getSubjectDirectoryOfFile(file PdfFile) string {
	// transforms the path so that only the subject directory remains
	leftRemoved := strings.TrimPrefix(file.path, "../../exams/markdown/")
	subjectDirectoryName := strings.SplitN(leftRemoved, "/", 2)[0]

	return subjectDirectoryName
}

//go:embed schema.sql
var ddl string

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	files, err := getExamFiles("../../exams/markdown")
	if err != nil {
		log.Fatal(err)
	}

	// creates the AI client
	llmCtx := context.Background()
	client, err := genai.NewClient(llmCtx, nil)
	if err != nil {
		log.Fatal(err)
	}

	numJobs := len(files)
	jobs := make(chan Job, numJobs)
	results := make(chan FileContent, numJobs)

	// creates workers
	// scary number because the workers call a llm, so it can be costly when mistakes occur
	workerAmount := 500
	for w := 1; w <= workerAmount; w++ {
		go worker(w, jobs, results)
	}

	// push work to workers
	for _, file := range files {
		jobs <- Job{
			file: file,
			llm:  client,
		}
	}
	close(jobs)

	// gets results
	collectedResults := make([]FileContent, numJobs)
	for a := 1; a <= numJobs; a++ {
		collectedResults[a-1] = <-results
	}
	fmt.Println("Done with header parsing.")
	// saves file content struct to disk so that it does not need to regenerate it with the llm
	if _, err := os.Stat("./file_header_mapping.json"); os.IsNotExist(err) {
		fmt.Println("Saving headers to disk.")

		jsonCollectedResults, err := json.Marshal(collectedResults)
		if err != nil {
			log.Fatal(err)
		}
		// permission: an owner can edit everyone else can read
		err = os.WriteFile("./file_header_mapping.json", jsonCollectedResults, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		jsonCollectedResults, err := os.ReadFile("./file_header_mapping.json")
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonCollectedResults, &collectedResults)
	}

	mapping := directorySubjectMapping{}
	// Check if mapping is already done. Saved to disk for faster development by not needing to call the llm api every time
	if _, err := os.Stat("./year_subject_mapping.json"); os.IsNotExist(err) {

		// maps the subject name and year to the path

		for i, result := range collectedResults {
			subjectDirectoryName := getSubjectDirectoryOfFile(result.File)

			if _, subjectMapped := mapping[subjectDirectoryName]; !subjectMapped {
				fmt.Println("Done:", i)
				resp, err := ParseSubjectAndYearFromDirectoryName(subjectDirectoryName, client)

				if err != nil {
					log.Fatal(err)
				}

				mapping[subjectDirectoryName] = resp
			}
		}
		jsonMapping, err := json.Marshal(mapping)
		if err != nil {
			log.Fatal(err)
		}
		// permission: an owner can edit everyone else can read
		err = os.WriteFile("./year_subject_mapping.json", jsonMapping, 0644)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		mappingFileContent, err := os.ReadFile("./year_subject_mapping.json")
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(mappingFileContent, &mapping)
		if err != nil {
			log.Fatal(err)
		}
	}

	sqlCtx := context.Background()
	conn, err := sql.Open("sqlite", "db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer func(conn *sql.DB) {
		err = conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(conn)

	_, err = conn.ExecContext(sqlCtx, ddl)
	if err != nil {
		log.Fatal(err)
	}

	querries := repository.New(conn)

	// Begin transaction
	tx, err := conn.BeginTx(sqlCtx, nil)
	if err != nil {
		log.Fatal("Couldn't start transaction:", err)
	}
	queryTx := querries.WithTx(tx)

	// string is the dir subject name
	subjectMapping := map[string]repository.Subject{}
	// insert subject
	for _, subject := range enumSubjects {
		insertSubject, err := queryTx.InsertSubject(sqlCtx, subject)
		if err != nil {
			log.Fatal("Insert error:", err)
		}
		subjectMapping[subject] = insertSubject
	}

	fmt.Println(subjectMapping)

	// string is a filepath
	fileMapping := map[string]repository.File{}
	// insert file
	for _, result := range collectedResults {
		yearAndSubject := mapping[getSubjectDirectoryOfFile(result.File)]

		savedFile, err := queryTx.InsertFile(sqlCtx, repository.InsertFileParams{
			FilePath: result.File.path,
			Year:     int64(yearAndSubject.Year),
		})
		if err != nil {
			log.Fatal("Insert error:", err)
		}

		fileMapping[result.File.path] = savedFile
	}

	// insert exam, answer and other
	for _, result := range collectedResults {
		yearAndSubject := mapping[getSubjectDirectoryOfFile(result.File)]
		subject := subjectMapping[yearAndSubject.SubjectName]
		file := fileMapping[result.File.path]

		switch result.MaterialType {
		case "exam":
			_, err = queryTx.InsertExam(sqlCtx, repository.InsertExamParams{
				SubjectID:         subject.ID,
				FileID:            file.FileID,
				EmbeddingID:       -1,
				ExamType:          result.HeaderExam.ExamType,
				Difficulty:        result.HeaderExam.Difficulty,
				TaskLabel:         result.HeaderExam.TaskLabel,
				WorkTimeInMinutes: int64(result.WorkTime),
			})
			if err != nil {
				log.Fatal("Insert error:", err)
			}
		case "answer":
			_, err = queryTx.InsertAnswer(sqlCtx, repository.InsertAnswerParams{
				SubjectID:   subject.ID,
				FileID:      file.FileID,
				EmbeddingID: -1,
			})
			if err != nil {
				log.Fatal("Insert error:", err)
			}
		case "other":
			_, err = queryTx.InsertOther(sqlCtx, repository.InsertOtherParams{
				FileID:      file.FileID,
				SubjectID:   subject.ID,
				EmbeddingID: -1,
			})
		}

	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		log.Fatal("commit error:", err)
	}

}
