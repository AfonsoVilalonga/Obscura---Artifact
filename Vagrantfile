Vagrant.configure("2") do |config|
  # Common resources for all VMs
  config.vm.provider "virtualbox" do |vb|
    vb.memory = "8192"  # 4GB of RAM
    vb.cpus = 4         # 4 vCPUs
    vb.customize ["modifyvm", :id, "--vram", "128"]  # Set VRAM to 128MB
    vb.customize ["modifyvm", :id, "--nictype1", "Am79C973"]
    vb.customize ["modifyvm", :id, "--nictype2", "Am79C973"]
    vb.customize ["modifyvm", :id, "--nictype3", "Am79C973"]
    vb.customize ["modifyvm", :id, "--nictype4", "Am79C973"]

    vb.customize ["modifyvm", :id, "--cableconnected1", "on"]  
    vb.customize ["modifyvm", :id, "--cableconnected2", "on"]  
    vb.customize ["modifyvm", :id, "--cableconnected3", "on"]  
    vb.customize ["modifyvm", :id, "--cableconnected4", "on"] 
  end

  # ---- BROKER ----
  config.vm.define "broker" do |broker|
    broker.vm.box = "ubuntu/jammy64"
    broker.vm.network "private_network", ip: "192.168.10.10", virtualbox__intnet: "broker_client"  #enp0s8
    broker.vm.network "private_network", ip: "192.168.20.10", virtualbox__intnet: "proxy_broker"   #enp0s9
    broker.vm.synced_folder "./Broker", "/home/vagrant/Broker"  
    
    broker.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/Broker" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Broker/entry.sh
    SHELL
  end

  # ---- CLIENT PION ----
  config.vm.define "client" do |client|
    client.vm.box = "ubuntu/jammy64"
    client.vm.network "private_network", ip: "192.168.10.20", virtualbox__intnet: "broker_client" #enp0s8
    client.vm.network "private_network", ip: "192.168.30.20", virtualbox__intnet: "proxy_client"  #enp0s9

    #client.vm.provider "virtualbox" do |vb|
    #  vb.customize ["modifyvm", :id, "--nic3", "natnetwork"]
    #  vb.customize ["modifyvm", :id, "--nat-network3", "test"]
    #end

    client.vm.synced_folder "./Client", "/home/vagrant/Client" 

    client.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools python3 python3-pip libnetfilter-queue-dev libnfnetlink-dev
      
      
      pip3 install scapy
      pip3 install NetfilterQueue
      

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/Client" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Client/entry.sh
      chmod +x /home/vagrant/Client/t.sh
    SHELL

  end

   # ---- CLIENT PION ANIMATION ----
  config.vm.define "client_a" do |client_a|
    client_a.vm.box = "ubuntu/jammy64"
    client_a.vm.network "private_network", ip: "192.168.10.20", virtualbox__intnet: "broker_client" #enp0s8
    client_a.vm.network "private_network", ip: "192.168.30.20", virtualbox__intnet: "proxy_client"  #enp0s9
    client_a.vm.synced_folder "./ClientAnimation", "/home/vagrant/Client" 

    client_a.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools python3 python3-pip libnetfilter-queue-dev libnfnetlink-dev
      sudo apt install -y ffmpeg
      
      
      pip3 install scapy
      pip3 install NetfilterQueue
      

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/Client" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Client/entry.sh
      chmod +x /home/vagrant/Client/t.sh
    SHELL

  end

  # ---- CLIENT PION TOR ----
  config.vm.define "client_tor" do |client_tor|
    client_tor.vm.box = "ubuntu/jammy64"
    client_tor.vm.network "private_network", ip: "192.168.10.20", virtualbox__intnet: "broker_client" #enp0s8
    client_tor.vm.network "private_network", ip: "192.168.30.20", virtualbox__intnet: "proxy_client"  #enp0s9
    client_tor.vm.synced_folder "./ClientTor", "/home/vagrant/Client" 

    client_tor.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

    
      sudo apt install -y apt-transport-https
      echo -e "deb [signed-by=/usr/share/keyrings/deb.torproject.org-keyring.gpg] https://deb.torproject.org/torproject.org jammy main\n\
deb-src [signed-by=/usr/share/keyrings/deb.torproject.org-keyring.gpg] https://deb.torproject.org/torproject.org jammy main" \
        | sudo tee /etc/apt/sources.list.d/tor.list > /dev/null

      sudo apt install -y gnupg
      wget -qO- https://deb.torproject.org/torproject.org/A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89.asc | sudo gpg --dearmor | sudo tee /usr/share/keyrings/deb.torproject.org-keyring.gpg >/dev/null
      sudo apt update
      sudo apt install -y tor deb.torproject.org-keyring
      cp /home/vagrant/Client/torrc /etc/tor/torrc
      

      echo "cd /home/vagrant/Client" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Client/entry.sh
      chmod +x /home/vagrant/Client/t.sh
    SHELL

  end

  # ---- CLIENT WEB ----
  config.vm.define "client_web" do |client_web|
    client_web.vm.box = "ubuntu/jammy64"
    client_web.vm.network "private_network", ip: "192.168.10.20", virtualbox__intnet: "broker_client"  #enp0s8
    #client_web.vm.network "public_network", bridge: "enp66s0f0", ip: "192.168.70.20"          
    client_web.vm.provider "virtualbox" do |vb|
      vb.customize ["modifyvm", :id, "--nic3", "natnetwork"]
      vb.customize ["modifyvm", :id, "--nat-network3", "test"]
    end
    
    client_web.vm.synced_folder "./ClientBrowser", "/home/vagrant/Client" 
    
    client_web.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update 
      sudo apt-get install -y wget curl gnupg unzip jq python3 python3-pip software-properties-common apt-transport-https net-tools libnetfilter-queue-dev libnfnetlink-dev

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

       wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
      sudo dpkg -i google-chrome-stable_current_amd64.deb || sudo apt-get install -f -y
      rm google-chrome-stable_current_amd64.deb

      CHROME_VERSION=$(google-chrome-stable --version | awk '{ print $3 }')
      CHROME_MAJOR_VERSION=$(echo $CHROME_VERSION | cut -d '.' -f 1)
      DRIVER_VERSION=$(curl -sSL https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions.json | jq -r '.channels.Stable.version')
      wget https://storage.googleapis.com/chrome-for-testing-public/$DRIVER_VERSION/linux64/chromedriver-linux64.zip
      unzip chromedriver-linux64.zip
      sudo mv chromedriver-linux64/chromedriver /usr/local/bin/
      sudo chmod +x /usr/local/bin/chromedriver
      rm -rf chromedriver-linux64 chromedriver-linux64.zip

      pip3 install selenium
      pip3 install scapy
      pip3 install NetfilterQueue

      curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
      sudo apt-get install -y nodejs

      wget -O firefox.tar.xz "https://download.mozilla.org/?product=firefox-latest&os=linux64&lang=en-US"
      tar -xJf firefox.tar.xz
      sudo mv firefox /opt/firefox
      sudo ln -s /opt/firefox/firefox /usr/bin/firefox
      rm firefox.tar.xz

      GECKO_LATEST=$(curl -s https://api.github.com/repos/mozilla/geckodriver/releases/latest | jq -r '.tag_name')
      wget https://github.com/mozilla/geckodriver/releases/download/${GECKO_LATEST}/geckodriver-${GECKO_LATEST}-linux64.tar.gz
      tar -xzf geckodriver-${GECKO_LATEST}-linux64.tar.gz
      sudo mv geckodriver /usr/local/bin/
      sudo chmod +x /usr/local/bin/geckodriver
      rm geckodriver-${GECKO_LATEST}-linux64.tar.gz

      echo "cd /home/vagrant/Client" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Client/entry.sh
      chmod +x /home/vagrant/Client/t.sh

      #sudo route add default gw 192.168.70.30
    SHELL

  end


  # ---- CLIENT WEB TOR ----
  config.vm.define "client_web_tor" do |client_web_tor|
    client_web_tor.vm.box = "ubuntu/jammy64"
    client_web_tor.vm.network "private_network", ip: "192.168.10.20", virtualbox__intnet: "broker_client"  #enp0s8
    
    #enp0s9
    client_web_tor.vm.provider "virtualbox" do |vb|
      vb.customize ["modifyvm", :id, "--nic3", "natnetwork"]
      vb.customize ["modifyvm", :id, "--nat-network3", "test"]
    end
    
    client_web_tor.vm.synced_folder "./ClientBrowserTor", "/home/vagrant/Client" 
    
    client_web_tor.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update 
      sudo apt-get install -y wget curl gnupg unzip jq python3 python3-pip software-properties-common apt-transport-https net-tools

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
      sudo dpkg -i google-chrome-stable_current_amd64.deb || sudo apt-get install -f -y
      rm google-chrome-stable_current_amd64.deb

      CHROME_VERSION=$(google-chrome-stable --version | awk '{ print $3 }')
      CHROME_MAJOR_VERSION=$(echo $CHROME_VERSION | cut -d '.' -f 1)
      DRIVER_VERSION=$(curl -sSL https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions.json | jq -r '.channels.Stable.version')
      wget https://storage.googleapis.com/chrome-for-testing-public/$DRIVER_VERSION/linux64/chromedriver-linux64.zip
      unzip chromedriver-linux64.zip
      sudo mv chromedriver-linux64/chromedriver /usr/local/bin/
      sudo chmod +x /usr/local/bin/chromedriver
      rm -rf chromedriver-linux64 chromedriver-linux64.zip

      pip3 install selenium

      curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
      sudo apt-get install -y nodejs

      wget -O firefox.tar.xz "https://download.mozilla.org/?product=firefox-latest&os=linux64&lang=en-US"
      tar -xJf firefox.tar.xz
      sudo mv firefox /opt/firefox
      sudo ln -s /opt/firefox/firefox /usr/bin/firefox
      rm firefox.tar.xz

      GECKO_LATEST=$(curl -s https://api.github.com/repos/mozilla/geckodriver/releases/latest | jq -r '.tag_name')
      wget https://github.com/mozilla/geckodriver/releases/download/${GECKO_LATEST}/geckodriver-${GECKO_LATEST}-linux64.tar.gz
      tar -xzf geckodriver-${GECKO_LATEST}-linux64.tar.gz
      sudo mv geckodriver /usr/local/bin/
      sudo chmod +x /usr/local/bin/geckodriver
      rm geckodriver-${GECKO_LATEST}-linux64.tar.gz

      sudo apt install -y apt-transport-https
      echo -e "deb [signed-by=/usr/share/keyrings/deb.torproject.org-keyring.gpg] https://deb.torproject.org/torproject.org jammy main\n\
deb-src [signed-by=/usr/share/keyrings/deb.torproject.org-keyring.gpg] https://deb.torproject.org/torproject.org jammy main" \
      | sudo tee /etc/apt/sources.list.d/tor.list > /dev/null

      sudo apt install -y gnupg
      wget -qO- https://deb.torproject.org/torproject.org/A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89.asc | sudo gpg --dearmor | sudo tee /usr/share/keyrings/deb.torproject.org-keyring.gpg >/dev/null
      sudo apt update
      sudo apt install -y tor deb.torproject.org-keyring
      cp /home/vagrant/Client/torrc /etc/tor/torrc

      echo "cd /home/vagrant/Client" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Client/entry.sh
      chmod +x /home/vagrant/Client/t.sh
    SHELL

  end


  # ---- PROXY PION ----
  config.vm.define "proxy" do |proxy|
    proxy.vm.box = "ubuntu/jammy64"
    proxy.vm.network "private_network", ip: "192.168.20.30", virtualbox__intnet: "proxy_broker"  #enp0s8
    proxy.vm.network "private_network", ip: "192.168.30.30", virtualbox__intnet: "proxy_client"  #enp0s9

    #proxy.vm.provider "virtualbox" do |vb|
    #  vb.customize ["modifyvm", :id, "--nic4", "natnetwork"]
    #  vb.customize ["modifyvm", :id, "--nat-network4", "test"]
    #end

    proxy.vm.network "private_network", ip: "192.168.40.30", virtualbox__intnet: "proxy_bridge"  #enp0s10
    proxy.vm.synced_folder "./Proxy", "/home/vagrant/Proxy" 

    proxy.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools python3 python3-pip libnetfilter-queue-dev libnfnetlink-dev

      pip3 install scapy
      pip3 install NetfilterQueue

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/Proxy" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Proxy/entry.sh
    SHELL
  end

   # ---- PROXY PION ANIMATION ----
  config.vm.define "proxy_a" do |proxy_a|
    proxy_a.vm.box = "ubuntu/jammy64"
    proxy_a.vm.network "private_network", ip: "192.168.20.30", virtualbox__intnet: "proxy_broker"  #enp0s8
    proxy_a.vm.network "private_network", ip: "192.168.30.30", virtualbox__intnet: "proxy_client"  #enp0s9
    proxy_a.vm.network "private_network", ip: "192.168.40.30", virtualbox__intnet: "proxy_bridge"  #enp0s10
    proxy_a.vm.synced_folder "./ProxyAnimation", "/home/vagrant/Proxy" 

    proxy_a.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools python3 python3-pip libnetfilter-queue-dev libnfnetlink-dev
      sudo apt install -y ffmpeg

      pip3 install scapy
      pip3 install NetfilterQueue

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/Proxy" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Proxy/entry.sh
    SHELL
  end

  # ---- PROXY WEB ----
  config.vm.define "proxy_web" do |proxy_web|
    proxy_web.vm.box = "ubuntu/jammy64"
    
    proxy_web.vm.network "private_network", ip: "192.168.20.30", virtualbox__intnet: "proxy_broker"  #enp0s8
    #proxy_web.vm.network "public_network", bridge: "enp66s0f0", ip: "192.168.80.20"  
    
    #enp0s10
    proxy_web.vm.provider "virtualbox" do |vb|
      vb.customize ["modifyvm", :id, "--nic4", "natnetwork"]
      vb.customize ["modifyvm", :id, "--nat-network4", "test"]
    end
    
    proxy_web.vm.network "private_network", ip: "192.168.40.30", virtualbox__intnet: "proxy_bridge"  #enp0s9
    proxy_web.vm.synced_folder "./Proxy-Web", "/home/vagrant/Proxy" 

    proxy_web.vm.provision "shell", inline: <<-SHELL
      set -e
      
      sudo apt update 
      sudo apt-get install -y wget curl gnupg unzip jq python3 python3-pip software-properties-common apt-transport-https net-tools libnetfilter-queue-dev libnfnetlink-dev

      wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
      sudo dpkg -i google-chrome-stable_current_amd64.deb || sudo apt-get install -f -y
      rm google-chrome-stable_current_amd64.deb

      CHROME_VERSION=$(google-chrome-stable --version | awk '{ print $3 }')
      CHROME_MAJOR_VERSION=$(echo $CHROME_VERSION | cut -d '.' -f 1)
      DRIVER_VERSION=$(curl -sSL https://googlechromelabs.github.io/chrome-for-testing/last-known-good-versions.json | jq -r '.channels.Stable.version')
      wget https://storage.googleapis.com/chrome-for-testing-public/$DRIVER_VERSION/linux64/chromedriver-linux64.zip
      unzip chromedriver-linux64.zip
      sudo mv chromedriver-linux64/chromedriver /usr/local/bin/
      sudo chmod +x /usr/local/bin/chromedriver
      rm -rf chromedriver-linux64 chromedriver-linux64.zip

      pip3 install selenium
      pip3 install scapy
      pip3 install NetfilterQueue

      curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
      sudo apt-get install -y nodejs

      wget -O firefox.tar.xz "https://download.mozilla.org/?product=firefox-latest&os=linux64&lang=en-US"
      tar -xJf firefox.tar.xz
      sudo mv firefox /opt/firefox
      sudo ln -s /opt/firefox/firefox /usr/bin/firefox
      rm firefox.tar.xz

      GECKO_LATEST=$(curl -s https://api.github.com/repos/mozilla/geckodriver/releases/latest | jq -r '.tag_name')
      wget https://github.com/mozilla/geckodriver/releases/download/${GECKO_LATEST}/geckodriver-${GECKO_LATEST}-linux64.tar.gz
      tar -xzf geckodriver-${GECKO_LATEST}-linux64.tar.gz
      sudo mv geckodriver /usr/local/bin/
      sudo chmod +x /usr/local/bin/geckodriver
      rm geckodriver-${GECKO_LATEST}-linux64.tar.gz

      echo "cd /home/vagrant/Proxy" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Proxy/entry.sh

      #sudo route add default gw 192.168.80.30
    SHELL
  end

  # ---- BRIDGE ----
  config.vm.define "bridge" do |bridge|
    bridge.vm.box = "ubuntu/jammy64"
    bridge.vm.network "private_network", ip: "192.168.40.40", virtualbox__intnet: "proxy_bridge" #enp0s8
    bridge.vm.network "private_network", ip: "192.168.50.40", virtualbox__intnet: "bridge_http"  #enp0s9
    bridge.vm.synced_folder "./Bridge", "/home/vagrant/Bridge"  

    bridge.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/Bridge" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Bridge/entry.sh
    SHELL
  end

  # ---- BRIDGE TOR ----
  config.vm.define "bridge_tor" do |bridge_tor|
    bridge_tor.vm.box = "ubuntu/jammy64"
    bridge_tor.vm.network "private_network", ip: "192.168.40.40", virtualbox__intnet: "proxy_bridge" #enp0s8
    bridge_tor.vm.network "private_network", ip: "192.168.50.40", virtualbox__intnet: "bridge_http"  #enp0s9
    bridge_tor.vm.synced_folder "./BridgeTor", "/home/vagrant/Bridge"  

    bridge_tor.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      sudo apt install -y apt-transport-https
      echo -e "deb [signed-by=/usr/share/keyrings/deb.torproject.org-keyring.gpg] https://deb.torproject.org/torproject.org jammy main\n\
deb-src [signed-by=/usr/share/keyrings/deb.torproject.org-keyring.gpg] https://deb.torproject.org/torproject.org jammy main" \
        | sudo tee /etc/apt/sources.list.d/tor.list > /dev/null

      sudo apt install -y gnupg
      wget -qO- https://deb.torproject.org/torproject.org/A3C4F0F979CAA22CDBA8F512EE8CBC9E886DDD89.asc | sudo gpg --dearmor | sudo tee /usr/share/keyrings/deb.torproject.org-keyring.gpg >/dev/null
      sudo apt update
      sudo apt install -y tor deb.torproject.org-keyring
      cp /home/vagrant/Bridge/torrc /etc/tor/torrc

      echo "cd /home/vagrant/Bridge" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/Bridge/entry.sh
    SHELL
  end

  # ---- HTTP SERVER ----
  config.vm.define "http" do |http|
    http.vm.box = "ubuntu/jammy64"
    http.vm.network "private_network", ip: "192.168.50.50", virtualbox__intnet: "bridge_http" #enp0s8 
    http.vm.synced_folder "./HTTP", "/home/vagrant/HTTP"  

    http.vm.provision "shell", inline: <<-SHELL
      set -e
      GO_VERSION=1.21.2

      sudo apt update && sudo apt install -y net-tools

      wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz

      sudo rm -rf /usr/local/go
      sudo tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz

      rm go${GO_VERSION}.linux-amd64.tar.gz

      sudo su

      echo PATH=\"$PATH:/usr/local/go/bin\" >> /etc/environment
      echo GOPATH=\"/home/vagrant/go\"

      echo "cd /home/vagrant/HTTP" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/HTTP/entry.sh
    SHELL
  end

  # ---- MIDDLEBOX ----
  config.vm.define "m" do |m|
    m.vm.box = "ubuntu/jammy64"
    m.vm.network "public_network", bridge: "enp66s0f0", ip: "192.168.80.30"
    m.vm.network "public_network", bridge: "enp66s0f0", ip: "192.168.70.30"   

    m.vm.synced_folder "./Middlebox", "/home/vagrant/M" 

    m.vm.provision "shell", inline: <<-SHELL
      set -e

      sudo apt update 
      sudo apt-get install -y wget curl gnupg unzip jq python3 python3-pip software-properties-common apt-transport-https net-tools libnetfilter-queue-dev libnfnetlink-dev iptables

      pip3 install selenium
      pip3 install scapy
      pip3 install NetfilterQueue

      # Enable IP forwarding
      echo "net.ipv4.ip_forward=1" | sudo tee -a /etc/sysctl.conf
      sudo sysctl -p

      # NAT between client_middle (enp0s8) → proxy_middle (enp0s9)
      sudo iptables -t nat -A POSTROUTING -o enp0s9 -j MASQUERADE

      echo "cd /home/vagrant/M" >> /home/vagrant/.bashrc

      chmod +x /home/vagrant/M/entry.sh
    SHELL
  end

end
