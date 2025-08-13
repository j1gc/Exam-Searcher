use std::{fs, io, };
use std::collections::HashMap;
use std::fmt::format;
use std::fs::DirEntry;
use std::path::Path;
use std::sync::Arc;
use axum::http::{Method, StatusCode};
use axum::{Json, Router};
use axum::extract::State;
use axum::routing::post;
use markdown::mdast::Node;
use serde::{Serialize, Deserialize};
use sqlx::{ FromRow, Pool,  Sqlite};
use tower_http::cors::{Any, CorsLayer};

fn get_text_from_markdown_tree(node: &Node, text_buffer: &mut String) {
    match node {
        Node::Text(text) => {
            text_buffer.push_str(&text.value);
        }
        Node::InlineCode(code) => {
            text_buffer.push_str(&code.value);
        }
        Node::Code(code) => {
            text_buffer.push_str(&code.value);
        }
        _ => {
            if let Some(children) = node.children() {
                for child in children {
                    get_text_from_markdown_tree(child, text_buffer);
                }
            }
        }
    }
}

// https://doc.rust-lang.org/nightly/std/fs/fn.read_dir.html#examples
fn visit_dirs<F>(dir: &Path, cb: &mut F) -> io::Result<()>
where
    F: FnMut(&DirEntry),
{
    if dir.is_dir() {
        for entry in fs::read_dir(dir)? {
            let entry = entry?;
            let path = entry.path();
            if path.is_dir() {
                visit_dirs(&path, cb)?;
            } else {
                cb(&entry);
            }
        }
    }
    Ok(())
}

fn collect_text_from_markdown_file(file: &DirEntry) -> io::Result<String> {
    if !file.path().is_file() {
        return Err(io::Error::new(io::ErrorKind::Other, "Not a file"));
    }
    if file.path().extension().unwrap() != "md" {
        return Err(io::Error::new(io::ErrorKind::Other, "Not a markdown file"));
    }

    let path = file.path();
    let mut contents = fs::read_to_string(&path)?;
    contents = contents.replace("<br>", " ");

    Ok(contents)
}

fn collect_words_in_document(document: &String) -> Vec<String> {
    let words: Vec<&str> = document.split_whitespace().collect();

    let filtered_words: Vec<String> = words
        .iter()
        .filter(|word| !word.is_empty())
        .map(|word| {
            word.chars()
                .filter(|c| {
                    // checks if the character is alphanumeric or math symbols
                    c.is_alphanumeric() || "+-*/=<>".contains(*c)
                })
                .flat_map(|c| c.to_lowercase())
                .collect()
            // removes words that are empty or only contain math symbols like "--------"
        }).filter(|s: &String| !s.is_empty() && s.chars().any(|c| {c.is_alphanumeric()}))
        .collect();

    filtered_words
}

fn count_term_occurrences(words: &Vec<String>) ->  HashMap<String, u32> {
    let mut term_occurrences : HashMap<String, u32> = HashMap::new();

    for word in words.iter() {
        *term_occurrences.entry(word.to_string()).or_insert(0) += 1;
    }

    term_occurrences
}

fn compute_term_frequency(word_occurrences: &HashMap<String, u32>, word_amount: u32) -> HashMap<String, f32> {
    let mut term_frequency : HashMap<String, f32> = HashMap::new();
    for (term, count) in word_occurrences.iter() {
        // normalizes the term frequency by dividing it by the word amount of the document
        term_frequency.insert(term.to_string(), (*count as f32) / (word_amount as f32));
    }
    term_frequency
}


fn cosine_similarity(
    vec1: &HashMap<String, f32>,
    vec2: &HashMap<String, f32>
) -> f32 {
    let mut dot_product = 0.0;
    let mut norm1 = 0.0;
    let mut norm2 = 0.0;

    for (key, val) in vec1.iter() {
        norm1 += val * val;
        if let Some(val2) = vec2.get(key) {
            dot_product += val * val2;
        }
    }

    for val in vec2.values() {
        norm2 += val * val;
    }

    if norm1 == 0.0 || norm2 == 0.0 {
        return 0.0;
    }

    dot_product / (norm1.sqrt() * norm2.sqrt())
}
#[derive(Serialize, Deserialize, Debug)]
struct Document {
    id: i64,
    path: String,
    words: Vec<String>,
    word_occurrences: HashMap<String, u32>,
    tf: HashMap<String, f32>,
    tf_idf: HashMap<String, f32>,
}

impl Document {
    fn count_word_occurrences(&mut self) {
        let mut term_occurrences : HashMap<String, u32> = HashMap::new();

        for word in self.words.iter() {
            *term_occurrences.entry(word.to_string()).or_insert(0) += 1;
        }

        self.word_occurrences = term_occurrences;
    }

    fn get_word_occurrences(&self) -> HashMap<String, u32> {
        self.word_occurrences.clone()
    }

    fn compute_tf(&mut self) {
        let term_occurrences = count_term_occurrences(&self.words);
        let word_amount = self.words.len();
        self.tf = compute_term_frequency(&term_occurrences, word_amount as u32);
    }

    fn compute_tf_idf(&mut self, idf: &HashMap<String, f32>) {
        let mut tf_idf : HashMap<String, f32> = HashMap::new();
        for (term, tf_val) in self.tf.iter() {
            if let Some(idf_val) = idf.get(term) {
                tf_idf.insert(term.clone(), tf_val * idf_val);
            }
        }
        self.tf_idf = tf_idf;
    }
}

#[derive(Debug)]
#[derive(Clone)]
#[derive(serde::Serialize)]
struct QueryReturn {
    embedding_id: i64,
    similarity: f32,
}

#[derive(Serialize, Deserialize, Debug)]
struct Searcher {
    documents: Vec<Document>,
    idf: HashMap<String, f32>,
}

impl Searcher {
    fn compute_idf(&mut self) -> HashMap<String, f32> {
        let mut idf : HashMap<String, f32> = HashMap::new();
        let document_amount = self.documents.len();

        // gets the number of documents that contain a term
        let documents_word_occurrences: Vec<_> = self.documents.iter().map(|doc| doc.get_word_occurrences()).collect();
        for document_word_occurrences in documents_word_occurrences.iter() {
            for (term, _) in document_word_occurrences.iter() {
                *idf.entry(term.to_string()).or_insert(0.0) += 1.0;

            }
        }

        // calculates the inverse document frequency
        for (_, val) in idf.iter_mut() {
            *val = (document_amount as f32 / *val).log10()
        }

        idf
    }

    fn load_embeddings(&mut self, file_path: String) {
        let file_contents = fs::read(file_path).unwrap();
        let deserialized_embeddings: Searcher = postcard::from_bytes(&file_contents).unwrap();
        self.documents = deserialized_embeddings.documents;
        self.idf = deserialized_embeddings.idf;
    }

    async fn load_documents(&mut self, path: String, db_connection: &Pool<Sqlite>) {
        let mut collected_documents = Vec::new();

        // load words
        visit_dirs(path.as_ref(), &mut |entry| {
            match collect_text_from_markdown_file(entry) {
                Ok(contents) => {
                    let document_path = entry.path().to_str().unwrap().to_string();

                    let mut current_document = Document {
                        id: -1,
                        path: document_path,
                        words: Vec::new(),
                        word_occurrences: HashMap::new(),
                        tf: HashMap::new(),
                        tf_idf: HashMap::new(),
                    };

                    let markdown_ast = markdown::to_mdast(&*contents, &markdown::ParseOptions::default()).unwrap();

                    let mut current_words = String::new();
                    get_text_from_markdown_tree(&markdown_ast, &mut current_words);

                    current_document.words = collect_words_in_document(&current_words);
                    collected_documents.push(current_document);
                }
                Err(e) => {
                    println!("{:?}", e);
                }
            }
        }).unwrap();
        // assigns the word list to the documents
        self.documents = collected_documents;

        // computes the tf and occurrences of words for the documents
        self.documents.iter_mut().for_each(|doc| {
            doc.count_word_occurrences();
            doc.compute_tf();
        });

        self.idf = self.compute_idf();

        self.documents.iter_mut().for_each(|doc| {
            println!("Computing tf_idf for doc:{:?}", doc.path);
            doc.compute_tf_idf(&self.idf);
        });

        for doc in self.documents.iter_mut() {
            let db_file_path = format!("../{}", doc.path);
            println!("Saving embedding in sqlite: {:?}", doc.path);

            let db_embedding: DBEmbedding = sqlx::query_as( r#"INSERT INTO document_embedding (embedding) VALUES (?) RETURNING *;"#)
                .bind(&db_file_path)
                .fetch_one(db_connection)
                .await
                .unwrap();

            doc.id = db_embedding.embedding_id.unwrap();

            sqlx::query("UPDATE file SET embedding_id = ? WHERE file_path = ?").bind(db_embedding.embedding_id).bind(db_file_path).execute(db_connection).await.unwrap();
        }
    }

    fn search_documents(&self, query: String) -> HashMap<i64, QueryReturn> {
        let query_words = collect_words_in_document(&query);
        let mut query_doc = Document{
            id: 0,
            path: "query".to_string(),
            words: query_words,
            word_occurrences: HashMap::new(),
            tf: HashMap::new(),
            tf_idf: HashMap::new(),
        };
        query_doc.compute_tf();
        query_doc.compute_tf_idf(&self.idf);
        //let query_doc_as_vector: Vec<_> = query_doc.tf_idf.values().collect();

        let mut similarities= HashMap::new();

        for doc in self.documents.iter() {
            //let doc_as_vector: Vec<_> = doc.tf_idf.values().collect();

            let similarity = cosine_similarity(&query_doc.tf_idf, &doc.tf_idf );
            similarities.insert(doc.id,QueryReturn {
                embedding_id: doc.id,
                similarity,
            });
        }


        similarities
    }

     fn save_embeddings(&self, file_path: String) {
        let embedding_file_content: Vec<u8> = postcard::to_stdvec(&self).unwrap();
        fs::write(file_path, embedding_file_content).unwrap();
    }
}

#[derive(Debug, FromRow, Serialize)]
#[derive(Clone)]
struct ExamFile {
    // exam
    exam_id: Option<i64>,
    subject_id: Option<i64>,
    exam_type: Option<String>,
    difficulty: Option<String>,
    task_label: Option<String>,
    work_time_in_minutes: Option<i64>,
    // file
    file_id: Option<i64>,
    year: Option<i64>,
    file_path: Option<String>,
    embedding_id: Option<i64>,
}

#[derive(Debug, FromRow, Serialize)]
#[derive(Clone)]
struct AnswerFile {
    // answer
    answer_id: Option<i64>,
    subject_id: Option<i64>,
    // file
    file_id: Option<i64>,
    year: Option<i64>,
    file_path: Option<String>,
    embedding_id: Option<i64>,
}

#[derive(Debug, FromRow, Serialize)]
#[derive(Clone)]
struct OtherFile {
    // other
    other_id: Option<i64>,
    subject_id: Option<i64>,
    // file
    file_id: Option<i64>,
    year: Option<i64>,
    file_path: Option<String>,
    embedding_id: Option<i64>,
}

#[derive(serde::Deserialize)]
#[derive(Debug)]
struct SearchRequest {
    query: String,
    file_types: Option<Vec<String>>,
    subject_ids: Option<Vec<i64>>,
    years: Option<Vec<i64>>,
}

struct AppState {
    searcher: Searcher,
    db_connection: Pool<Sqlite>,
}

#[derive(Debug, Serialize)]
#[derive(Clone)]
enum FileTypes {
    Exam(ExamFile),
    Answer(AnswerFile),
    Other(OtherFile),
}

#[derive(Debug, Serialize, Clone)]
struct FileResponsePart {
    file_type: String,
    file: FileTypes,
    similarity: f32
}

impl FileResponsePart {
    pub fn get_file_path(&self) -> Option<&String> {
        match &self.file {
            FileTypes::Exam(exam) => exam.file_path.as_ref(),
            FileTypes::Answer(answer) => answer.file_path.as_ref(),
            FileTypes::Other(other) => other.file_path.as_ref(),
        }
    }
}

struct SearchResponse {
    files: Vec<FileResponsePart>,
}

fn combine_vec_to_string<T: std::fmt::Display>(vec: &Option<Vec<T>>) -> Option<String> {

    vec.as_ref().and_then(|vec| {
        if vec.is_empty() {
            None
        } else {
            if vec.len() == 1 {
                return Some(format!(",{},", vec[0].to_string()));
            }

            let joined_vec = vec.iter()
                .map(|t| t.to_string())
                .collect::<Vec<String>>()
                .join(",");
            Some(format!(",{},", joined_vec))
        }
    })
}

#[derive(Deserialize)]
struct Pagination {
    index: usize,
    page_size: usize,
}

#[axum::debug_handler]
async fn search(State(state): State<Arc<AppState>>, pagination: axum::extract::Query<Pagination>, Json(req): Json<SearchRequest>) -> (StatusCode, Json<Vec<FileResponsePart>>) {

    let subject_id_string = combine_vec_to_string(&req.subject_ids);
    println!("subject id string: {:?}", subject_id_string);
    let year_string = combine_vec_to_string(&req.years);
    println!("year string: {:?}", year_string);


    let doc_similarities = state.searcher.search_documents(req.query.clone());

    let mut files: Vec<FileResponsePart> = Vec::new();
    
    if let Some(mut file_types) = req.file_types.clone() {
        // for deduplication sorting is needed first
        file_types.sort();
        file_types.dedup();

        for file_type in &file_types {
            match file_type.as_str() {
                "exam" => {
                    let exams: Vec<ExamFile> = sqlx::query_as!(ExamFile, r#"
                            SELECT
                                file.file_id,
                                file.file_path,
                                file.year,
                                file.embedding_id,
                                e.exam_id,
                                e.subject_id,
                                e.exam_type,
                                e.difficulty,
                                e.task_label,
                                e.work_time_in_minutes
                            FROM
                                file
                                    INNER JOIN main.exam e ON file.file_id = e.file_id
                                    AND (?1 IS NULL OR INSTR(?1, ',' || e.subject_id || ','))
                            WHERE
                               (?2 IS NULL OR INSTR(?2, ',' || file.year || ','))
                        "#, subject_id_string, year_string).fetch_all(&state.db_connection).await.unwrap();

                    for exam in exams.iter() {
                        files.push(FileResponsePart {
                            file_type: "exam".to_string(),
                            file: FileTypes::Exam(exam.clone()),
                            similarity: doc_similarities.get(&exam.embedding_id.unwrap()).unwrap().similarity,
                        });
                    }
                }
                "answer" => {
                    let answers: Vec<AnswerFile> = sqlx::query_as!(AnswerFile, r#"
                            SELECT
                                file.file_id,
                                file.file_path,
                                file.year,
                                file.embedding_id,
                                a.answer_id,
                                a.subject_id
                            FROM
                                file
                                    INNER JOIN main.answer a ON file.file_id = a.file_id
                                    AND (?1 IS NULL OR INSTR(?1, ',' || a.subject_id || ','))
                            WHERE
                               (?2 IS NULL OR INSTR(?2, ',' || file.year || ','))
                        "#, subject_id_string, year_string).fetch_all(&state.db_connection).await.unwrap();

                    for answer in answers.iter() {
                        files.push(FileResponsePart {
                            file_type: "answer".to_string(),
                            file: FileTypes::Answer(answer.clone()),
                            similarity: doc_similarities.get(&answer.embedding_id.unwrap()).unwrap().similarity,
                        });
                    }
                }
                "other" => {
                    let other: Vec<OtherFile> = sqlx::query_as!(OtherFile, r#"
                            SELECT
                                file.file_id,
                                file.file_path,
                                file.year,
                                file.embedding_id,
                                o.other_id,
                                o.subject_id
                            FROM
                                file
                                    INNER JOIN main.other o ON file.file_id = o.file_id
                                    AND (?1 IS NULL OR INSTR(?1, ',' || o.subject_id || ','))
                            WHERE
                               (?2 IS NULL OR INSTR(?2, ',' || file.year || ','))
                        "#, subject_id_string, year_string).fetch_all(&state.db_connection).await.unwrap();

                    for other in other.iter() {
                        files.push(FileResponsePart {
                            file_type: "other".to_string(),
                            file: FileTypes::Other(other.clone()),
                            similarity: doc_similarities.get(&other.embedding_id.unwrap()).unwrap().similarity,
                        });
                    }
                }
                _ => {
                }
            }
        }
    }


    files.sort_by(|a, b| b.similarity.partial_cmp(&a.similarity).unwrap());

    println!("{:?}", req);
    let paginated_files = files.chunks(pagination.page_size).nth(pagination.index).unwrap_or(&[]).to_vec();

    return (StatusCode::OK, Json(paginated_files));
}

#[derive(Debug, sqlx::FromRow)]
struct DBEmbedding {
    embedding_id: Option<i64>,
    embedding: Option<String>,
}


#[tokio::main]
async fn main() {
    let pool = sqlx::sqlite::SqlitePool::connect("sqlite://./db.sqlite3").await.unwrap();

    let mut s = Searcher {
        documents: Vec::new(),
        idf: HashMap::new(),
    };

    let embedding_path = "./embeddings.postcard";

    if !Path::new(embedding_path).exists() {
        s.load_documents("../exams/markdown/".to_string(), &pool).await;
        s.save_embeddings(embedding_path.to_string());
        println!("Saved embeddings to file");
    } else {
        s.load_embeddings(embedding_path.to_string());
        println!("Loaded embeddings from file");
    }

    let app_state = Arc::new(AppState {
        searcher: s,
        db_connection: pool,
    });



    let mut app = Router::new().route("/search", post(search)).with_state(app_state);

    if (std::env::var("DEBUG").unwrap_or("".to_string()).to_lowercase() == "true") {
        app = app.layer(CorsLayer::permissive());
        println!("Running in debug mode");
    } else {
        let cors = CorsLayer::new()
            // allow any headers
            .allow_headers(Any)
            .allow_methods([Method::POST])
            // allow requests from origins
            .allow_origin(["http://exam.jonas-floerke.de".parse().unwrap(), "https://exam.jonas-floerke.de".parse().unwrap()]);


        println!("Running in production mode");
        app = app.layer(cors);
    }


    match std::env::var("PORT") {
        Ok(val) => {
            println!("Listening on http://0.0.0.0:{}", val);
            let listener = tokio::net::TcpListener::bind(format!("0.0.0.0:{}", val)).await.unwrap();
            axum::serve(listener, app).await.unwrap();
        }
        Err(_) => {
            println!("Listening on http://0.0.0.0:3000");
            let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
            axum::serve(listener, app).await.unwrap();
        }
    }


}
