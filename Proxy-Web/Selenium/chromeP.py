from selenium import webdriver
from selenium.webdriver.chrome.service import Service
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
import time

# Set up Chrome options
chrome_options = Options()
chrome_options.add_argument("--headless")  
chrome_options.add_argument("--disable-gpu") 
chrome_options.add_argument("--no-sandbox") 
chrome_options.add_argument("--autoplay-policy=no-user-gesture-required")  
chrome_options.add_argument("--user-data-dir=/home/vagrant/data")  
chrome_options.add_argument("--mute-audio")  
chrome_options.add_argument("--ignore-certificate-errors")  

chrome_options.add_argument("--disable-backgrounding-occluded-windows")
chrome_options.add_argument("--disable-background-timer-throttling")
chrome_options.add_argument("--disable-renderer-backgrounding")

chrome_options.set_capability("goog:loggingPrefs", {"browser": "ALL"})

chromedriver_path = "/usr/local/bin/chromedriver"

service = Service(chromedriver_path)
driver = webdriver.Chrome(service=service, options=chrome_options)

try:
    url = "http://localhost:3000/video"
    driver.get(url)

    # logs = driver.get_log("browser")
    # if logs:
    #     print("Browser console logs:")
    #     for entry in logs:
    #         print(f"[{entry['level']}] {entry['message']}")
    # else:
    #     print("No console logs found.")
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

        # input_element = driver.find_element(By.ID, 'inputField')
        # entered_value = input_element.get_attribute('value')
        # print(entered_value)
        # a = driver.find_element(By.ID, 'bits')
        # a = a.get_attribute('value')
        # print(a)# Keep the script alive to maintain the browser session
        
        # input_elementa = driver.find_element(By.ID, 'bitsa')
        # entered_valuea = input_elementa.get_attribute('value')
        # print(entered_valuea)
        
        # input_elements = driver.find_element(By.ID, 'inputFielda')
        # entered_values = input_elements.get_attribute('value')
        # print(entered_values)

finally:
    driver.quit()
    print("Browser closed.")
