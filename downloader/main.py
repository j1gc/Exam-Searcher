import time

from playwright.sync_api import sync_playwright

def main():
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=False)
        page = browser.new_page()
        page.goto("https://za-aufgaben.nibis.de/index.php?jahr=0")
        time.sleep(10)
        print(page.title())



if __name__ == "__main__":
    main()
