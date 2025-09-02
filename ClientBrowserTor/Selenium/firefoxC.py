from selenium import webdriver
from selenium.webdriver.firefox.service import Service
from selenium.webdriver.firefox.options import Options
from selenium.webdriver.common.by import By
import time

firefox_options = Options()

firefox_options.add_argument("--headless")  
firefox_options.set_preference("network.stricttransportsecurity.preloadlist", False)
firefox_options.set_preference("security.enterprise_roots.enabled", True)
firefox_options.set_preference("media.autoplay.default", 0)  
firefox_options.set_preference("media.autoplay.allow-muted", True)  
firefox_options.set_preference("media.autoplay.block-webaudio", False)  
firefox_options.set_preference("media.peerconnection.ice.loopback", True)

firefox_options.set_preference("dom.min_background_timeout_value", 0)
firefox_options.set_preference("dom.timeout.background_throttling_max_budget", -1)
firefox_options.set_preference("browser.tabs.remote.separateFileUriProcess", False)

geckodriver_path = "/usr/local/bin/geckodriver"

service = Service(geckodriver_path)
driver = webdriver.Firefox(service=service, options=firefox_options)

try:
    url = "http://localhost:3010/video"
    driver.get(url)
    print(f"Opened {url} successfully.")

    print("Firefox is running in headless mode. Press Ctrl+C to stop.")
        
    while True:
        time.sleep(1)

finally:
    driver.quit()
    print("Firefox browser closed.")
