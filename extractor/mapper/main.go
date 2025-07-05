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

func MapExamsInDirectoryToAnswers(ctx context.Context, client *genai.Client, directoryTree string) []Exams {
	// creates the config
	responseSchema := genai.Schema{
		Type: genai.TypeArray,
		Items: &genai.Schema{
			Title: "Exams",
			Type:  genai.TypeArray,
			Items: &genai.Schema{
				Type: genai.TypeObject,
				Properties: map[string]*genai.Schema{
					"ExamPart": {
						Type: genai.TypeString,
					},
					"Answer part": {
						Type: genai.TypeString,
					},
				},
			},
		},
	}
	thinkingBudget := int32(0) // Disables thinking
	thinkingConfig := genai.ThinkingConfig{ThinkingBudget: &thinkingBudget}

	// queries the AI for the structure
	prompt := fmt.Sprintf("you are given an file tree, "+
		"organize the filenames according to that. Under all circumstances only give back the filename not the entire path for the fields and do not change the file name under all circumstances in adition to that all fields are rquierd and every field can only take one file. "+
		"This is an example of how you can give it back: [\n  [\n    {\n      \"Answer part\": \"2017DeutschEAA1L.md\",\n      \"ExamPart\": \"2017DeutschEAAufg1.md\"\n    },\n    {\n      \"Answer part\": \"2017DeutschEAA2L.md\",\n      \"ExamPart\": \"2017DeutschEAAufg2.md\"\n    },\n    {\n      \"Answer part\": \"2017DeutschEAA3L.md\",\n      \"ExamPart\": \"2017DeutschEAAufg3.md\"\n    }\n  ],\n  [\n    {\n      \"Answer part\": \"2017DeutschGAA1L.md\",\n      \"ExamPart\": \"2017DeutschGAAufg1.md\"\n    },\n    {\n      \"Answer part\": \"2017DeutschGAA2L.md\",\n      \"ExamPart\": \"2017DeutschGAAufg2.md\"\n    },\n    {\n      \"Answer part\": \"2017DeutschGAA3L.md\",\n      \"ExamPart\": \"2017DeutschGAAufg3.md\"\n    }\n  ]\n]"+
		"but please do not change the name to suit that example the name should be still the same as the input: %s", directoryTree)
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash-preview-04-17",
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

	// unmarshals the JSON returned by the AI into a struct
	var exams []Exams
	err = json.Unmarshal([]byte(result.Text()), &exams)
	if err != nil {
		log.Fatal(err)
	}

	return exams
}

type PdfFile struct {
	path        string
	name        string
	alreadyUsed bool
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
			path:        path,
			name:        info.Name(),
			alreadyUsed: false,
		}

		pdfFiles[currentPdfFile.name] = currentPdfFile

		return nil
	})
	if err != nil {
		return nil, err
	}

	fmt.Println(pdfFiles)

	return pdfFiles, nil
}

type Subject struct {
	SubjectName string  `json:"SubjectName"`
	Exams       []Exams `json:"Exams"`
}

type Exams []struct {
	AnswerPart string `json:"AnswerPart"`
	ExamPart   string `json:"ExamPart"`
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

	// That all could be run in go routines. But cloud scars me, so slower is better
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

			log.Println(string(directoryStructure))

			exams := MapExamsInDirectoryToAnswers(ctx, client, string(directoryStructure))

			for _, exam := range exams {
				for _, file := range exam {
					answerFileUsed, answerFileExists := fileMap[file.AnswerPart]
					examFileUsed, examFileExists := fileMap[file.ExamPart]

					if answerFileExists && examFileExists && !answerFileUsed.alreadyUsed && !examFileUsed.alreadyUsed {
						dataIsCorrect = true

						answerFileUsed.alreadyUsed = true
						examFileUsed.alreadyUsed = true

						fileMap[file.AnswerPart] = answerFileUsed
						fileMap[file.ExamPart] = examFileUsed
					} else {
						fmt.Println("File not found:", file.AnswerPart, file.ExamPart)
					}
				}
			}
			if dataIsCorrect {
				// ads the current subject to the subject array if the AI does not hallucinate the data
				subjects = append(subjects, Subject{
					SubjectName: SubjectName,
					Exams:       exams,
				})
			}
		}
	}

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
