# Exam-Searcher

If you are doing your A Level in the State of Lower Saxony, Germany, you have likely experienced the bad user experience of using the official [website](https://za-aufgaben.nibis.de/) of the Ministry of Education here. The site only provides the cental A Level examinations as ZIPs containing all exams/answers for a selected year and subject. This requires you to download large ZIPs, extract their contents, and manually search through the big file structure of these ZIPs to get to the specific exam youre looking for. This process is inefficient and time consuming, especially when youre in class and need to find the exact exam that the teacher wants you to work on.

To make your search faster, this project offers an online search engine designed for the A Level exams in Lower Saxony by implementing TF-IDF. In addition to that, it also provides filtering options by year, subject and file type, so that you can quickly find, instead of searching through large ZIPs, the exam or answer that youre currently working on.

## The Current state of the project:

- [x] Scraped the exam and answer files
- [x] Converted the PDF files to markdown for easier data extraction
- [x] Created the TF-IDF for the files and used the cosine similarity between the query and the TF-IDF of the documents as a relevancy metric
- [x] Extracted data for filtering, like the difficulty, whether a file is an exam or an answer, etc., with the help of an LLM API from the markdown documents
- [x] Saved the filtering data, etc., in a SQLite database for querying
- [x] Expanded the backend API to make filtering by the frontend possible
- [x] Created a frontend
- [x] Deployed the front- and backend code to Google Cloud Platforms Cloud Run service

Please be advised that the code quality is currently at a hackathon level

## Implementation:

The backend is built with Rust using axum, SQLx with SQLite for data storage, and the Postcard format for storing the TF-IDF of the documents. Also located in the backend is a data enrichment pipeline, which consists of Python with Playwright for document scraping, Golang for the classification of the documents with the help of Gemini-flash-lite, and the creation of the SQLite DB with SQLC.

The frontend is a SvelteKit app with tailwind for styling, shadcn-svelte, Zod for parsing JSON, and TanStack Query for infinite scrolling.
