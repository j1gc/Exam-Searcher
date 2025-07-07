 CREATE TABLE IF NOT EXISTS subject (
    id INTEGER PRIMARY KEY AUTOINCREMENT ,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS file (
    file_id INTEGER PRIMARY KEY AUTOINCREMENT,
    file_path TEXT NOT NULL UNIQUE,
    year INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS exam (
    exam_id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    exam_type TEXT NOT NULL,
    difficulty TEXT NOT NULL,
    task_label TEXT NOT NULL,
    work_time_in_minutes INTEGER NOT NULL,
    embedding_id INTEGER NOT NULL,
    foreign key (embedding_id) references document_embedding(embedding_id),
    FOREIGN KEY (subject_id) REFERENCES subject(id),
    FOREIGN KEY (file_id) REFERENCES file(file_id)
);

CREATE TABLE IF NOT EXISTS answer (
    answer_id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY (embedding_id) references document_embedding(embedding_id),
    FOREIGN KEY (subject_id) REFERENCES subject(id),
    foreign key (file_id) references file(file_id)
);

CREATE TABLE IF NOT EXISTS other (
    other_id INTEGER PRIMARY KEY AUTOINCREMENT,
    subject_id INTEGER NOT NULL,
    file_id INTEGER NOT NULL,
    embedding_id INTEGER NOT NULL,
    FOREIGN KEY (embedding_id) references document_embedding(embedding_id),
    FOREIGN KEY (file_id) REFERENCES file(file_id)
);

CREATE TABLE IF NOT EXISTS document_embedding (
    embedding_id INTEGER PRIMARY KEY AUTOINCREMENT ,
    embedding BLOB NOT NULL
);