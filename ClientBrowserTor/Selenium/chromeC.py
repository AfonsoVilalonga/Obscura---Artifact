from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
import time

chrome_options = Options()
chrome_options.add_argument("--headless")  
chrome_options.add_argument("--disable-gpu")
chrome_options.add_argument("--no-sandbox")
chrome_options.add_argument("--autoplay-policy=no-user-gesture-required")
chrome_options.add_argument("--user-data-dir=/home/vagrant/data")
chrome_options.add_argument("--mute-audio")
chrome_options.add_argument("--ignore-certificate-errors")
chrome_options.add_argument("--disable-features=WebRtcHideLocalIpsWithMdns")
chrome_options.add_argument("--force-webrtc-ip-handling-policy=default")

chrome_options.add_argument("--disable-backgrounding-occluded-windows")
chrome_options.add_argument("--disable-background-timer-throttling")
chrome_options.add_argument("--disable-renderer-backgrounding")

#Capture all logs from the browser
chrome_options.set_capability("goog:loggingPrefs", {"browser": "ALL"})  

# Path to the ChromeDriver CHANGE
chromedriver_path = "/usr/local/bin/chromedriver"


service = Service(chromedriver_path)
driver = webdriver.Chrome(service=service, options=chrome_options)

try:
    url = "http://localhost:3010/video"
    driver.get(url)


    print("Browser is running. Press Ctrl+C to stop.")
    while True:
        # logs = driver.get_log("browser")
        # if logs:
        #     print("Browser console logs:")
        #     for entry in logs:
        #         print(f"[{entry['level']}] {entry['message']}")
        # else:
        #     print("No console logs found.")

        time.sleep(1)

finally:
    driver.quit()
    print("Browser closed.")