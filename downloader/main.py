import os
import uuid
import zipfile

from marker.converters.pdf import PdfConverter
from marker.models import create_model_dict
from marker.output import text_from_rendered
from playwright.sync_api import sync_playwright


def download_exams_on_page(page):
    # check boxes of content to download
    download_checkboxes = page.locator('input[type=checkbox]').all()
    for checkbox in download_checkboxes:
        # downloading entire years does produce a bug on the website where nothing is being exams
        if checkbox.is_checked() or checkbox.get_attribute("data-fach") == "alle":
            continue
        checkbox.click()

    # execute download
    with page.expect_download() as download_info:
        page.get_by_role("button", name="Auswahl herunterladen").click()

    download = download_info.value

    return download


def download_exams_on_url(urls, output_path):
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        page = browser.new_page()

        for url in urls:
            # navigate to page
            page.goto(url)
            print(page.title())

            download = download_exams_on_page(page)
            # save download at
            download.save_as(output_path + uuid.uuid4().hex + ".zip")

        page.close()


def unzip_zips_in_dir(directory, output_dir):
    files = os.listdir(directory)
    print(files)
    # unpack zips
    for file in files:
        # if file is not a zip, skip
        if not file.endswith(".zip"):
            print("Skipping " + file)
            continue

        # unzips downloaded files
        with zipfile.ZipFile(directory + file, "r") as exam_zip:
            exam_zip.extractall(output_dir)


def walk_files(src_dir) -> list[str]:
    filepath_list = []

    for root, dirs, files in os.walk(src_dir):
        # get files in current walk
        for file in files:

            # root handling
            if root == ".":
                root_path = os.getcwd() + "/"
            else:
                root_path = root

            # checks if an extra / is needed
            if (root_path != src_dir) and (root != "."):
                filepath = root_path + "/" + file
            else:
                filepath = root_path + file

            # appends filepath to filelist if not already in it
            if filepath not in filepath_list:
                filepath_list.append(filepath)

    return filepath_list


def convert_pdf_to_markdown(pdf_file, markdown_file_path: str, converter):
    rendered = converter(pdf_file)
    text, _, _ = text_from_rendered(rendered)

    # creates the dir that the markdown file will live in
    os.makedirs(os.path.dirname(markdown_file_path), exist_ok=True)

    # handles save logic of md file
    with open(markdown_file_path, "w") as markdown_file:
        # writes to the created markdown file
        markdown_file.write(text)



def handle_markdown_conversion(pdf_dir, markdown_dir):
    # creates pdf converter with specified config
    converter = PdfConverter(
        artifact_dict=create_model_dict()
    )

    # gets paths pdf files
    files = walk_files(pdf_dir)

    index = 0
    num_pdf_files = len(files)
    for file in files:
        if not file.endswith(".pdf"):
            print("NONE PDF FILE: " + file)
            continue

        print("File:", index, "of", num_pdf_files, "~", round(index/num_pdf_files*100, 2), "%")

        markdown_file_path = markdown_dir + file.removeprefix(pdf_dir).removesuffix(".pdf") + ".md"
        print("converting file:", file, "to:", markdown_file_path)

        convert_pdf_to_markdown(file, markdown_file_path , converter)

        index += 1


def main():
    zipfile_output_path = "../exams/zips/entire/"
    output_path_subject_zips = "../exams/zips/subjects/"
    output_path_subject = "../exams/subjects/"
    output_path_markdown = "../exams/markdown/"

    # urls = ["https://za-aufgaben.nibis.de/index.php?jahr=1", "https://za-aufgaben.nibis.de/index.php?jahr=1"]
    # download_exams_on_url(urls, zipfile_output_path)

    # unzip_zips_in_dir(zipfile_output_path, output_path_subject_zips)

    # unzip_zips_in_dir(output_path_subject_zips, output_path_subject)

    handle_markdown_conversion(output_path_subject, output_path_markdown)


if __name__ == "__main__":
    main()
