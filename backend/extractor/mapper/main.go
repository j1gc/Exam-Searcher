package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"google.golang.org/genai"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func printModelList(ctx context.Context, client *genai.Client) {
	modelList, err := client.Models.List(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, model := range modelList.Items {
		fmt.Println(model.Name)
	}
}

func MapExamsInDirectoryToAnswers(ctx context.Context, client *genai.Client, directoryTree string) (Subject, int32) {
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

	// creates the config
	responseSchema := genai.Schema{
		Type:        genai.TypeObject,
		Description: "The response schema",
		Properties: map[string]*genai.Schema{
			"friendlySubjectName": {
				Type:        genai.TypeString,
				Enum:        enumSubjects,
				Description: "The friendly subject name from the enum provided eg: 2016VW -> Volkswirtschaft",
			},
			"year": {
				Type:        genai.TypeString,
				Description: "The year of the exam eg: 2016VW  -> 2016",
			},
			"exams": {
				Type: genai.TypeArray,
				Items: &genai.Schema{
					Title: "exams",
					Type:  genai.TypeArray,
					Items: &genai.Schema{
						Type: genai.TypeObject,
						Properties: map[string]*genai.Schema{
							"examPart": {
								Type:        genai.TypeString,
								Description: "The exam part eg: 2017DeutschEAAufg1.md",
							},
							"answerPart": {
								Type:        genai.TypeString,
								Description: "The answer part eg: 2017DeutschEAA1L.md. Do not use a path use only the file name",
							},
							"additionalParts": {
								Type: genai.TypeArray,
								Items: &genai.Schema{
									Type: genai.TypeString,
								},
								Description: "The additional parts eg: Deckblatt, Hinweise, Material: 2016ChemieEAmitExpA1LHinweise.md",
							},
						},
						Required: []string{"examPart", "answerPart"},
					},
				},
			},
		},
		Required: []string{"friendlySubjectName", "year", "exams"},
	}
	thinkingBudget := int32(512) // Disables thinking if set to 0
	thinkingConfig := genai.ThinkingConfig{ThinkingBudget: &thinkingBudget}

	// queries the AI for the structure
	prompt := fmt.Sprintf("you are given an file tree, "+
		"organize the filenames according to that. Under all circumstances only give back the filename not the entire path for the fields and do not change the file name under all circumstances in adition to that all fields are rquierd and every field can only take one file. Every file from the input should be included in the output. Also give back an friendly subject name and a year eg:. 2016VW -> FriendlyName=Volkswirtschaft, Year=2016., the field adtional parts can be emtpy but if a file is a Deckblatt or a Hinsweis etc. it should be in that array. eg: 2016ChemieEAmitExpA1LHinweise.md and 2016ChemieEAmitExpA1LMaterial.md"+
		""+
		"but please do not change the name to suit that example the name should be still the same as the input: %s", directoryTree)
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

	// unmarshals the JSON returned by the AI into a struct
	var subject Subject
	err = json.Unmarshal([]byte(result.Text()), &subject)
	if err != nil {
		log.Fatal(err)
	}

	return subject, result.UsageMetadata.TotalTokenCount
}

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

type Subject struct {
	FriendlySubjectName string  `json:"friendlySubjectName"`
	SubjectName         string  `json:"subjectName"`
	Year                string  `json:"year"`
	Exams               []Exams `json:"exams"`
}

type Exams []struct {
	AnswerPart      string   `json:"answerPart"`
	ExamPart        string   `json:"examPart"`
	AdditionalParts []string `json:"additionalParts"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	entries, err := os.ReadDir("../../exams/markdown")
	var SubjectNames []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		SubjectNames = append(SubjectNames, entry.Name())
	}
	fmt.Println(SubjectNames)

	// creates the AI client
	ctx := context.Background()
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fileMap, err := getExamFiles("../../exams/markdown")
	if err != nil {
		log.Fatal(err)
	}

	var subjects []Subject

	totalTokensUsed := int32(0)

	// That all could be run in go routines. But cloud scars me, so slower is better
	// SubjectNames[9:11]
	for idx, SubjectName := range SubjectNames {
		dataIsCorrect := false
		// very risky for loop, let's hope that I do not end up with a crazy cloud bill
		for !dataIsCorrect {
			fmt.Println("Subject:", SubjectName, "Index:", idx, "Total:", len(SubjectNames))

			// get the directory tree of the subject
			cmdCommands := fmt.Sprintf("cd ../../exams/markdown/%s && tree", SubjectName)
			cmd := exec.Command("sh", "-c", cmdCommands)
			directoryStructure, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatal(err)
			}

			subject, tokensUsed := MapExamsInDirectoryToAnswers(ctx, client, string(directoryStructure))
			totalTokensUsed += tokensUsed

			//filesAlreadyUsed := make(map[string]bool)
			examFilesCorrect := 0
			examFilesTotal := 0
			for _, exam := range subject.Exams {
				for _, filesForExam := range exam {
					examFilesTotal += 1

					_, answerFileExist := fileMap[filesForExam.AnswerPart]
					_, examFileExists := fileMap[filesForExam.ExamPart]

					additionalPartsExist := true
					if len(filesForExam.AdditionalParts) > 0 {
						for _, additionalPart := range filesForExam.AdditionalParts {
							_, additionalPartsExist = fileMap[additionalPart]
							if !additionalPartsExist {
								break
							}
						}
					}

					if answerFileExist && examFileExists && additionalPartsExist {
						//filesAlreadyUsed[filesForExam.ExamPart] = true

						examFilesCorrect += 1

					} else {
						fmt.Println("File not found:", filesForExam.AnswerPart, filesForExam.ExamPart, "Because:", filesForExam.AdditionalParts,
							"Answer:", answerFileExist, "Exam:", examFileExists, "AdditionalParts:", additionalPartsExist,
							"UsedAnswer:" /*, usedAnswer "UsedExam:", usedExam,*/, "UsedAdditionalPart:",
						)
						//fmt.Println("Files already used:", filesAlreadyUsed)
					}
				}
			}
			if examFilesCorrect == examFilesTotal {
				dataIsCorrect = true

				// ads the current subject to the subject array if the AI does not hallucinate the data
				subject.SubjectName = SubjectName
				subjects = append(subjects, subject)
			}
		}
	}

	fmt.Println("Total tokens used:", totalTokensUsed)

	// Saves the subjects as a JSON file
	bytes, err := json.Marshal(subjects)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create("../exam_mapping.json")
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)
	_, err = file.Write(bytes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(subjects)
}
