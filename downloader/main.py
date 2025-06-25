import time
import uuid

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

def main():
    # urls = ["https://za-aufgaben.nibis.de/index.php?jahr=1", "https://za-aufgaben.nibis.de/index.php?jahr=1"]
    # download_exams_on_url(urls, "../exams/zips/")

    print("Downloading exams...")

if __name__ == "__main__":
    main()
