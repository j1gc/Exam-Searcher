package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Subject struct {
	SubjectName string  `json:"SubjectName"`
	Exams       []Exams `json:"Exams"`
}

type Exams []struct {
	AnswerPart string `json:"AnswerPart"`
	ExamPart   string `json:"ExamPart"`
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

func main() {
	content, err := os.ReadFile("../exam_mapping.json")
	if err != nil {
		log.Fatal(err)
	}

	var subjects []Subject
	err = json.Unmarshal(content, &subjects)
	if err != nil {
		log.Fatal(err)
	}

	pdfFiles, err := getExamFiles("../../exams/subjects")
	if err != nil {
		log.Fatal(err)
	}

	rootDirectoryName := "../../exams/subjects_mapped"
	for _, subject := range subjects {
		for examIdx, exam := range subject.Exams {
			for examPartIdx, examPart := range exam {
				pdfFileExam, ok := pdfFiles[examPart.ExamPart]
				if !ok {
					log.Println("File not found:", examPart.ExamPart)
				}

				pdfFileAnswer, ok := pdfFiles[examPart.AnswerPart]
				if !ok {
					log.Println("File not found:", examPart.AnswerPart)
				}

				pathString := fmt.Sprintf("%s/%s/exam_%d/exam_part_%d/", rootDirectoryName, subject.SubjectName, examIdx, examPartIdx)

				// exam
				err = os.MkdirAll(pathString+"/exam/", os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}

				err = os.Rename(pdfFileExam.path, pathString+"/exam/"+pdfFileExam.name)
				if err != nil {
					log.Fatal(err)
				}

				// answers
				err = os.MkdirAll(pathString+"/answer/", os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}

				err = os.Rename(pdfFileAnswer.path, pathString+"/answer/"+pdfFileAnswer.name)
				if err != nil {
					log.Fatal(err)
				}

				fmt.Println(pdfFileExam.path)
				fmt.Println(pdfFileAnswer.path)
			}
		}
	}
}
