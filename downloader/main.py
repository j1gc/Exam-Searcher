import os
import uuid
import zipfile

from playwright.sync_api import sync_playwright

def download_exams_on_page(page):
    # separate columns
    columns = page.get_by_label("#download > ul").all()

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

def download_exams_on_url( urls, output_path):
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


def main():
    zipfile_output_path = "../exams/zips/entire/"
    output_path_subject_zips = "../exams/zips/subjects/"
    output_path_subject = "../exams/subjects/"
    # urls = ["https://za-aufgaben.nibis.de/index.php?jahr=1", "https://za-aufgaben.nibis.de/index.php?jahr=1"]
    # download_exams_on_url(urls, zipfile_output_path)

    unzip_zips_in_dir(zipfile_output_path, output_path_subject_zips)

    unzip_zips_in_dir(output_path_subject_zips, output_path_subject)



if __name__ == "__main__":
    main()
