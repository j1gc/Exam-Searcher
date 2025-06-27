use std::{fs, io, };
use std::fs::DirEntry;
use std::path::Path;
use std::sync::{Arc, Mutex};
use markdown::mdast::Node;

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
fn visit_dirs(dir: &Path, cb: &dyn Fn(&DirEntry)) -> io::Result<()> {
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
    let contents = fs::read_to_string(&path)?;

    Ok(contents)
}

fn main() {
    let pages = Arc::new(Mutex::new(Vec::new() as Vec<String>));
    
    visit_dirs("../exams/markdown/".as_ref(), &|entry| {
        match collect_text_from_markdown_file(entry) {
            Ok(contents) => {
                let markdown_ast = markdown::to_mdast(&*contents, &markdown::ParseOptions::default()).unwrap();

                let mut current_page = String::new();
                get_text_from_markdown_tree(&markdown_ast, &mut current_page);

                pages.lock().unwrap().push(current_page);
            }
            Err(e) => {
                println!("{:?}", e);
            }
        }
    }).unwrap();
    
    
}
