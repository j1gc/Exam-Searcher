-- name: GetSubjects :many
SELECT * FROM subject;

-- name: GetAllExams :many
SELECT * FROM exam;

-- name: GetAllAnswers :many
SELECT * FROM answer;

-- name: InsertSubject :one
INSERT INTO subject (name)
VALUES  (?)
RETURNING *;

-- name: InsertFile :one
INSERT INTO file (file_path, year)
VALUES (?, ?)
RETURNING *;

-- name: InsertExam :one
INSERT INTO exam (subject_id, file_id, embedding_id,exam_type, difficulty, work_time_in_minutes, task_label)
VALUES (?, ?, ?, ?, ?, ?, ?)
RETURNING *;

-- name: InsertAnswer :one
INSERT INTO answer (subject_id, file_id, embedding_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: InsertOther :one
INSERT INTO other (file_id, subject_id, embedding_id)
VALUES (?, ?, ?)
RETURNING *;

-- name: InsertDocumentEmbedding :one
INSERT INTO document_embedding (embedding)
VALUES (?)
RETURNING *;