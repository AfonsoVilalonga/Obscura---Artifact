# Artifact Appendix

Paper title: **Obscura: Enabling Ephemeral Proxies for Traffic Encapsulation in WebRTC Media Streams Against Cost-Effective Censors**

Requested Badge(s):
  - [X] **Available**
  - [X] **Functional**
  - [] **Reproduced**


## Description
This repository provides the artifact associated with the following paper:

```
@article{Vilalonga2026Obscura,
  title = {Obscura: Enabling Ephemeral Proxies for Traffic Encapsulation in WebRTC Media Streams Against Cost-Effective Censors},
  author = {Afonso Vilalonga and Kevin Gallagher and João S. Resende and Henrique Domingos},
  journal = {Proceedings on Privacy Enhancing Technologies},
  volume = {2026},
  issue = {1},
  year = {2026},
  series = {PoPETS '26},
  doi = {<to add>},
  url = {<to add>}
}
```

Obscura is a censorship evasion system that exchanges covert traffic, encapsulated within video frames of a WebRTC-based video streaming web application, between a client and a proxy. Since proxies may be ephemeral and go offline mid-session, each proxy with an active connection to a client also establishes a connection with a bridge through a WebSocket. Bridges and clients maintain the session state, allowing proxies to be switched during a live session. In addition, the bridge acts either as the entry node to the Tor network if Obscura is configured to be used as a Tor pluggable transport, or as the final gateway to the destination the user is trying to reach if Obscura is configured as a standalone system, similar to a VPN. Obscura allows the WebRTC connection to be established between different client and proxy implementations, specifically a web-based proxy or client and a Pion-based proxy or client, with [Pion](https://github.com/pion/webrtc) being an open-source WebRTC implementation in Go. The web-based instances use a WebRTC video streaming application that we developed, which is served locally via a Node.js server. To access it, both the client and the proxy automatically start the Node.js server and launch a browser to open the local web page containing the video streaming application. The application then establishes the connection between the client and the proxy. All combinations of the different instance types of proxies and clients are possible as endpoints in a client-proxy connection.


### Security/Privacy Issues and Ethical Concerns
Obscura does not introduce security or privacy risks to the evaluator’s machine, as our artifact does not require disabling any security features, and it was developed to operate within virtual machines on isolated networks. As a censorship circumvention tool, using the system in the wild may carry inherent security risks for users when accessing censored domains. However, we did not test the system in real-world scenarios or deployments, nor did we use real users for its development and testing, as we only used a machine rented from a cloud provider during testing and accessed only an HTTP server under our control. Additionally, we disclaim that the Obscura artifact should be viewed as a proof of concept for the research community to study and experiment with, rather than as a production-ready system immediately suitable for real-world scenarios and deployments.

 
## Basic Requirements
The artifact does not require any special hardware or proprietary software to run.


### Hardware Requirements
The setup used in our testing consists of five components (for a detailed explanation of each component, see the "Environment" section in this file):
- Broker
- Client
- Proxy
- Bridge
- HTTP Server

During our tests and the development of our paper, each of the following components ran on a virtual machine with the following hardware specifications: 8 GB RAM, 4 vCPUs, 128 MB VRAM, and 40 GB disk space (maximum capacity). Our artifact assumes that all five virtual machines, with each one running a single component, are executed on the same physical machine. To this end, the machine running this artifact should have at least the following resources:
- 56 GB of RAM (to account for the five virtual machines as well as the host machine);
- 12 CPU cores;
- 300 GB of disk space (40 GB is the maximum capacity for each virtual machine, although the actual usage is lower (~10-15 GB));

The machine used during testing was an OVH Cloud RISE-4 dedicated server (which is still available for rental at the time of writing). We rented it for three months and configured it with the following hardware specifications:
- 128 GB DDR4 RAM 
- 1Gbps of public and private bandwidth
- 2 * 960 GB SSD NVMe Soft RAID
- AMD Epyc 7313 - 16c/32t - 3GHz/3.7GHz

**Note:** Running all five machines simultaneously places a heavy load on the host system, so machines with lower resources will likely affect the results of the throughput tests.


### Software Requirements
The software requirements vary widely depending on the component. For example, the client requires Tor to run in pluggable transport mode, while the proxy does not require Tor even when the system is running in pluggable transport mode. Some components also rely on different languages, interpreters, or frameworks that other components do not require. For instance, the web-based proxy and client require Node.js, whereas the bridge does not. To simplify setup, each component has a specific virtual machine configured via a Vagrant script that is easy to run. Additionally, we provide different virtual machines for the same components, each configured with a specific configuration. For example, there are separate VMs for the web-based proxy and client and the Pion-based proxy and client, as well as separate VMs for clients and bridges running as a Tor pluggable transport versus those running in standalone mode. Each VM is equipped with all the necessary software and does not require any additional configuration on the evaluator’s machine. However, the following software must be installed on the host machine for the artifact to be configured and run:

1. **Operating System**
   - The OS used within each virtual machine is Ubuntu 22.04 LTS.
   - The OS used on the host machine was Ubuntu 22.04 LTS, but other versions or distributions should also work, as long as the remaining software requirements in this section can be installed.


2. **VirtualBox** (the most current version should work)
   - To install VirtualBox, which will act as the virtualization provider for Vagrant, run the following commands (assuming Ubuntu 22.04 LTS as the operating system):

   1. **Install the GPG keys for the official VirtualBox repository.**
      ```bash
      wget -q https://www.virtualbox.org/download/oracle_vbox_2016.asc -O- | sudo apt-key add -
      ```
   2. **Add the VirtualBox repository to Ubuntu.**
      ```bash
      echo "deb [arch=amd64] https://download.virtualbox.org/virtualbox/debian $(lsb_release -cs) contrib" | sudo tee /etc/apt/sources.list.d/virtualbox.list
      ```

   3. **Install the latest version of VirtualBox** (at the time of writing, version 7.2.4 is the current version; however, for the installation, replace [version-number] with 7.2).
      ```bash
      sudo apt update
      sudo apt install virtualbox-[version-number]
      ```

   4. **Verify the VirtualBox installation** (check that VirtualBox was successfully installed by printing the version currently installed on the system).
      ```bash
      VBoxManage --version
      ```

   5. **Install the VirtualBox Extension Pack**  (at the time of writing, replace [version-number] with 7.2.4).
      ```bash
      wget https://download.virtualbox.org/virtualbox/[version-number]/Oracle_VirtualBox_Extension_Pack-[version-number].vbox-extpack
      ```

      ```bash
      sudo VBoxManage extpack install Oracle_VirtualBox_Extension_Pack-[version-number].vbox-extpack
      ```

   6. **Verify that the extension pack was installed correctly by printing the list of installed extension packs.**
      ```bash
      VBoxManage list extpacks
      ```

3. **Vagrant** (the most current version should work)
   - To install Vagrant, the tool used to manage and automate virtual machines, run the following commands:

   1. **Download and install the GPG key for HashiCorp’s APT repository.**
      ```bash
      wget -O - https://apt.releases.hashicorp.com/gpg | sudo gpg --dearmor -o /usr/share/keyrings/hashicorp-archive-keyring.gpg
      ```
   
   2. **Add the HashiCorp APT repository to Ubuntu.**
      ```bash
      echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] https://apt.releases.hashicorp.com $(grep -oP '(?<=UBUNTU_CODENAME=).*' /etc/os-release || lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
      ```

   3. **Install Vagrant.**
      ```bash
      sudo apt update && sudo apt install vagrant
      ```

   4. **Verify that Vagrant was installed correctly by printing the currently installed version.**
      ```bash
      vagrant --version
      ```

The following software is also required for testing the system and reproducing our results, as defined in the "Artifact Evaluation" section:

1. **FFmpeg**
- To install FFmpeg, the tool used to encode videos with different resolutions and bitrates, run the following command:

   1. **Install FFmpeg.**
      ```bash
      sudo apt install ffmpeg
      ```

   2. **Check if FFmpeg was installed correctly by printing the currently installed version.**
      ```bash
      ffmpeg -version
      ```

2. **Big Buck Bunny Open Video Project:**
- Big Buck Bunny is the video file we use for most of our tests and, for reproducibility, should also be used. To download the video, run the following commands:
   
   1. **Download the video file from the official webpage.**
      ```bash
      wget https://download.blender.org/peach/bigbuckbunny_movies/big_buck_bunny_720p_stereo.avi
      ```


### Estimated Time and Storage Consumption
To run the artifact and check whether everything is installed and functioning correctly, we estimate 2 to 3 hours of overall human and computer time. For disk space, although the maximum capacity of each VM is 40 GB and running the artifact requires five virtual machines, resulting in a maximum of 200 GB, on average, we expect the evaluator to need around 10-15 GB per virtual machine, or 50-75 GB in total.


## Environment
As previously described, the environment required to run the Obscura artifact is composed of five components: the client, which is the client-side software that runs on the user’s device; the proxy, which is the ephemeral server-side component of the WebRTC connection and runs on volunteer devices; the bridge, which serves as the fixed component maintaining the server-side session state between the user and the bridge; the broker, a necessary component for any WebRTC connection that enables the establishment of the peer-to-peer connection between the client and the proxy; and finally the HTTP server, which acts as an HTTP endpoint under our control and is used for testing purposes as the destination server for the user. Thus, the workflow is as follows: the client and the proxy connect to the broker and use it to establish a WebRTC connection between them. During the establishment process, the proxy is instructed to connect to the bridge, and once the WebRTC connection between the client and the proxy is active, the client can send traffic to the proxy. The proxy then forwards this traffic to the bridge, which, in our testing setup, forwards it to the HTTP server. Responses follow the reverse path.

The client and proxy can run either as Pion instances (Go programs) or as web instances. For Pion instances, both the proxy and the client are Go programs that contain all the necessary logic to function. For web instances, a browser and a Node.js server are also required. The Node.js server serves a webpage containing the code for the WebRTC-based video streaming application used by web-based proxies and browsers to establish a connection, and it is automatically started upon executing either the client or proxy software. The browsers are also launched automatically and are used to access the locally hosted webpage that runs the WebRTC video streaming application. The artifact supports running two different browsers, Firefox and Chrome, selectable through a configuration setting available in both the web-based client and proxy configuration files.
   
  
### Accessibility
The artifact can be accessed using the following link:
```
[https://github.com/AfonsoVilalonga/Obscura---Artifact](https://github.com/AfonsoVilalonga/Obscura---Artifact)
```

### Set up the environment
To set up the artifact, follow the commands below (this section assumes that the software indicated in the "Software Requirements" section is already installed):

1. **Clone the repository of the artifact.**   
   ```bash
   git clone https://github.com/AfonsoVilalonga/Obscura---Artifact.git
   ```
   
2. **Encode the video downloaded in the "Software Requirements" section (the Big Buck Bunny) for both the Pion and web-based proxy and client instances.** (the video will be encoded with a 2 Mbps bitrate, 30 fps, and 1280x720 resolution as a WebM video for the web-based clients and proxies, and as an IVF file for video and an OGG file for audio for the Pion instances. `INPUT_FILE` should be replaced with the name of the Big Buck Bunny downloaded file).
   ```bash
   ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 -b:v 2M output.ivf
   ffmpeg -i INPUT_FILE -c:a libopus -page_duration 20000 -vn output.ogg
   ffmpeg -i INPUT_FILE -vf scale=1280:720 -c:v libvpx -b:v 2M -g 30 output.webm
   ```

3. **Move the output video and audio files to the correct locations for the artifact to use.** ([PATH] should be the path to the Obscura root directory).

   ```bash
   cp output.ivf "[PATH]/Obscura---Artifact/Client/Media
   cp output.ogg "[PATH]/Obscura---Artifact/Client/Media
   cp output.ivf "[PATH]/Obscura---Artifact/Proxy/Media
   cp output.ogg "[PATH]/Obscura---Artifact/Proxy/Media
   cp output.webm "[PATH]/Obscura---Artifact/ClientBrowser/Node-Server/videos
   cp output.webm "[PATH]/Obscura---Artifact/Proxy-Web/videos
   ```

4. **For a peer-to-peer connection between the web-based client and proxy, a NAT network must be created within VirtualBox using the following commands.**
   ```bash
   VBoxManage natnetwork add --netname test --network "192.168.100.0/24" --dhcp on
   VBoxManage natnetwork start --netname test
   ```
   
5. **Test that the network was created by listing all VirtualBox NAT networks.**
   ```bash
   VBoxManage natnetwork list
   ```

After following all of these commands, the artifact should be ready to run for any type of connection between proxy and clients (i.e., Pion-to-web, web-to-web, web-to-Pion, and Pion-to-Pion) and in any mode (i.e., Tor pluggable transport or standalone system). However, the results presented in the main paper (excluding the appendices) are based only on connections between endpoints of the same type and in standalone mode. Therefore, in this document, we primarily discuss running the system with proxies and clients of the same type and only in the standalone version. For documentation purposes, to test the system in Tor pluggable transport mode, the same workflow as in standalone mode is used, and the video and audio files must be placed in the appropriate directories: inside the Media folder for the Pion-based client (ClientTor folder) or inside the NodeServer/videos folder for the web-based client (ClientBrowserTor folder).


### Testing the Environment
To test that the environment is properly configured and the artifact is functional, follow the commands below:

1. **To test that the Pion-to-Pion connection is working properly, run the following commands inside the Obscura---Artifact folder, at the root of the directory containing the Vagrantfile:**
   
   1. **Initiate the virtual machines for each component (run each command in a separate console) (the working directory of the console should be the root directory of the Obscura---Artifact folder, where the Vagrantfile is located) (wait for each VM to load before starting the next one).**
      ```bash
      vagrant up bridge
      vagrant up broker
      vagrant up client
      vagrant up proxy
      vagrant up http
      ```
   
   2. **On each specific console, connect to the corresponding virtual machine.**
      ```bash
      vagrant ssh bridge
      vagrant ssh broker
      vagrant ssh client
      vagrant ssh proxy
      vagrant ssh http
      ```

   3. **Start up the bridge and broker on their respective consoles and wait for the "READY" output in each console after running the following command on each.**
      ```bash
      ./entry.sh
      ```

   4. **Start up the HTTP server on its console and wait for the "Server is running" output after running the following command.**
      ```bash
      ./entry.sh
      ```

   5. **Start up the proxy on its console.**
      ```bash
      ./entry.sh
      ```
   
   6. **Start up the client on its console.**
      ```bash
      ./entry.sh
      ```

   7. **Once the client and proxy are connected, a "Connected" message should be printed on both the client and proxy consoles.**
   

2. **To test that the Web-to-Web connection is working properly, run the following commands inside the Obscura---Artifact folder, at the root of the directory containing the Vagrantfile (run each command in a separate console):**
   
   1. **Initiate the virtual machines for each component (run each command in a separate console) (the working directory of the console should be the root directory of the Obscura---Artifact folder, where the Vagrantfile is located) (wait for each VM to load before starting the next one).**
      ```bash
      vagrant up bridge
      vagrant up broker
      vagrant up client_web
      vagrant up proxy_web
      vagrant up http
      ```

   2. **On each specific console, connect to the corresponding virtual machine.**
      ```bash
      vagrant ssh bridge
      vagrant ssh broker
      vagrant ssh client_web
      vagrant ssh proxy_web
      vagrant ssh http
      ```

   3. **Start up the bridge and broker on their respective consoles and wait for the "READY" output in each console after running the following command on each.**
      ```bash
      ./entry.sh
      ```

   4. **Start up the HTTP server on its console and wait for the "Server is running" output after running the following command.**
      ```bash
      ./entry.sh
      ```

   5. **It is possible to choose whether the browser used is Firefox or Chrome. To select the browser for the client, do the following:**
      - Edit the `entry.sh` file located inside the ClientBrowser folder.  
      - Inside the file, there is a line containing the following code:
      ```bash
      BROWSER_N=${BROWSER_NAME:-c}
      ```

      - Change the `c` to `f` or vice versa to switch between Chrome (`c`) and Firefox (`f`).

   6. **Start up the client_web on its console.**
      ```bash
      ./entry.sh
      ```

   7. **It is possible to choose whether the browser used is Firefox or Chrome. To select the browser for the proxy, do the following:**
      - Edit the `entry.sh` file located inside the Proxy-Web folder.  
      - Inside the file, there is a line containing the following code:
      ```bash
      BROWSER_N=${BROWSER_NAME:-c}
      ```

      - Change the `c` to `f` or vice versa to switch between Chrome (`c`) and Firefox (`f`). For now, use the same browser that was defined in step 5 for the client.

   8. **Start up the proxy_web on its console.**
      ```bash
      ./entry.sh
      ```

   9. **Once the client_web and proxy_web are connected, a "READY" message should be printed on the client_web console.**

3. **After the connection is established, the system should be ready to use and fully functional. To test simulating the download of data (1 MB) from the HTTP server using Obscura, do the following:**

   1. **In a new console, connect to the virtual machine running the client (the working directory of the console should be the root directory of the Obscura---Artifact folder, where the Vagrantfile is located).**
      ```bash
      vagrant ssh client_web
      ```

      or 

      ```bash
      vagrant ssh client
      ```
   
   2. **Run the `t.sh` script, which executes 5 curl requests to the HTTP server, each downloading 1 MB of data. If a download is successful, the system is verified to be operating correctly (it is not necessary to wait for all 5 downloads at this stage).**
      ```bash
      ./t.sh
      ```

      **NOTE:** If `./t.sh` outputs a "bad interpreter" error, run the following commands:
         ```bash
         sed -i 's/\r$//' t.sh
         chmod +x t.sh
         ```

   3. **Repeat this process for all connection types (i.e., Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox).**

To shut down all machines, run the following command in the root directory (where the Vagrantfile is, inside the Obscura---Artifact folder):
   ```bash
   vagrant halt
   ```

To delete all machines, run the following command in the root directory (where the Vagrantfile is, inside the Obscura---Artifact folder):
   ```bash
   vagrant destroy -f
   ```  

All machines can be started again by following the steps described in this section.


## Artifact Evaluation
In this section, we present the main results and claims from the paper, excluding the appendices, as we do not consider them part of the core findings. The paper identifies six main results: three related to performance (presented first) and three related to the system's resistance to specific attacks defined within our threat model (presented later in this section). 

All performance tests are based on setting up the artifact and measuring system throughput by downloading a 1 MB file from the HTTP server, while varying either the network conditions or the video used as the carrier to observe their effects on throughput. In the paper, each throughput value corresponds to the average of 100 downloads. However, because reproducing all experiments would take approximately to 3 weeks to 5 weeks, we consider that, for reproducibility purposes, performing five measurements should be sufficient, even though the resulting values may differ slightly from those reported in the paper. Thus, all time estimates refer to this less time-consuming version of our experiments.

### Main Results and Claims
#### Main Result 1: Obscura is able to maintain reasonable throughput under different network conditions for all types of connections between a client and a proxy (i.e., Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox).
In this result, we evaluate the throughput of the different types of connections under varying network conditions. Specifically, we test the impact of different round-trip times across the entire system (from the client to the HTTP server), different bandwidth constraints on the client to proxy connection, and different packet loss values between the client and the proxy. All these conditions are simulated using the Linux tc command. As expected, and as shown in Figure 7 of the paper, more restrictive network conditions lead to a deterioration in throughput.  


#### Main Result 2: Different video encoding parameterizations result in different throughput values. However, even with low bitrate and resolution, Obscura achieves reasonable throughput for all types of connections between a client and a proxy.
In this result, we test the throughput of the different types of connections (i.e., Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox) under different encoding parameterizations for the same video. Specifically, we measure throughput using the same video re-encoded at different bitrates and resolutions. As expected, and as shown in Figure 8 of the paper, higher bitrates and higher resolutions result in better throughput values.


#### Main Result 3: Different video profiles (e.g., gaming, chatting, sports) achieve different throughput values due to their inherent visual content characteristics. 
In this result, we test the throughput of different types of connections (i.e., Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox) against different video profiles, specifically chatting, coding, gaming, and sports. As expected, more dynamic and visually complex videos produce higher throughput. The results are presented in Figure 9 of the paper.


#### Main Result 4: Web-to-web connections between a client and a proxy do not have DTLS fingerprints that differ from the DTLS fingerprint of the browser in use. 
In this result, we show that web-to-web connections have the same DTLS fingerprints as the browser in use, making them identical to those of other WebRTC-based applications running in the same browser. DTLS is the protocol used within WebRTC to exchange the encryption keys for the connection. Some real-world WebRTC-based censorship evasion systems (i.e., Snowflake) have been targeted by attacks that exploit DTLS fingerprints. As expected, because Obscura uses the browser for both the client and the proxy in web-to-web connections, the DTLS fingerprint is the same as that of the browser in use and is similar or identical to other WebRTC-based applications running in the browser. Results are presented in Section 5.2 of the paper.


#### Main Result 5: Pion-to-Pion connections between a client and a proxy resist the media-based variant of differential degradation attacks. 
In this result, we show that Pion-to-Pion connections are resistant to an attack against WebRTC-based censorship evasion systems known as [differential degradation attacks](https://arxiv.org/pdf/2409.06247). For this result, we focus on the media-based variant of the attack. In this variant, specific packets from frames are dropped across multiple frames, causing the targeted frames to be lost while keeping overall packet loss relatively low. This severely degrades system throughput while still maintaining the video stream as watchable, making the attack feasible in real-world conditions with low collateral damage.
      
Pion-to-Pion connections are expected to resist this attack because, unlike other WebRTC media-based censorship evasion systems that rely on browsers (similar to our web-to-web connections), they perform decapsulation at the packet level rather than at the frame level. This means that if a packet from a frame is lost, Pion-to-Pion connections do not need to wait for its retransmission to obtain the entire frame, nor do they lose the entire frame if a packet fails to arrive. Since multiple covert packets may be contained within a single frame, avoiding the loss of the entire frame when some packets are dropped helps mitigate the attack. Additionally, Pion-to-Pion provides a reliability layer and does not suffer from packet-loss-induced video quality degradation, where video quality decreases in response to network packet loss. Results are shown in Figures 10 and 11 of the paper. 


#### Main Result 6: Proxies can be ephemeral.
This is not a result presented in the paper, but rather a feature of the system that acts as a countermeasure against some attacks defined in the threat model. We state that the proxies are ephemeral and can go offline, and in this section, we provide an experiment to showcase this feature.
   

### Experiments
#### Experiment 1: Varying network conditions for Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox connections.   
To perform this experiment related to Main Result 1, follow the instructions below. It should take approximately 5-6 hours to complete (about 30 mins to 1 hour per network configuration for all connection types):

1. **The network conditions must be defined. In our paper, we defined 9 different conditions, all specified within the `entry.sh` files for each component folder. However, only the `entry.sh` files need to be modified for the client and proxy components when changing the network conditions.**

   **NOTE**: For functionality testing, the baseline network conditions (number 2) can be used to test all three connection types. However, by fully following this experiment as we describe it, it is possible to reproduce our results, either with the suggested 5 downloads for a less time-consuming version or the 100 downloads we performed in the paper.

   - For the Pion-to-Pion connection, `entry.sh` is located inside the Client folder for the client and inside the Proxy folder for the proxy.  
   - For the Chrome-to-Chrome and Firefox-to-Firefox connections, `entry.sh` is located inside the ClientBrowser folder for the client and inside the Proxy-Web folder for the proxy.  
         
   The 9 defined network conditions are as follows:    
      
   1. RTT of 80 ms for the system  
   2. RTT of 115 ms for the system  
   3. RTT of 165 ms for the system  
   4. Packet loss of 2% and RTT of 115 ms  
   5. Packet loss of 5% and RTT of 115 ms  
   6. Packet loss of 10% and RTT of 115 ms  
   7. Bandwidth constraint of 250 Kbps and RTT of 115 ms  
   8. Bandwidth constraint of 750 Kbps and RTT of 115 ms  
   9. Bandwidth constraint of 1500 Kbps and RTT of 115 ms  
   
   The network configuration is changed by modifying the `TEST_N` variable in the `entry.sh` file and assigning it a specific number. The client and proxy network condition numbers should always be the same to reproduce our results (the same number must be set in both the `entry.sh` file of the client and the `entry.sh` file of the proxy). The available numbers are 1 to 9, with each number corresponding to one of the network conditions defined above. Configure the network condition to be tested in the same order as listed above.

   **NOTE**: When testing bandwidth constraints, some connection types may not work under these conditions. The ones that fail are not shown in Figure 7. Therefore, if a connection type has no results in Figure 7 for any of these network configurations, it can be skipped. Additionally, for 10% packet loss, the browsers might exhibit unusual behavior, with high variability in the results, especially over extended periods (this is particularly true for F-F connections).

   
2. **To run the Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox connections, the evaluator must follow the "Testing the Environment" section of this document to run the artifact with all these connection types. Completely follow steps 1) and 2) to make the artifact run, but do not follow step 3) for the download test.** (we assume that the "Set up the environment" section of this document has already been completed. From this point onward, we also assume that the artifact is running with any of these three connection types, but only one at a time).
  

3. **In a separate console, connect to the client virtual machine (the working directory of the console must be the same as the root directory of the artifact, where the Vagrantfile is located).**
   - For web-based instances:
   ```bash
   vagrant ssh client_web
   ```
   - For Pion-based instances:
   ```bash
   vagrant ssh client
   ```
   
4. **Inside the client virtual machine, run the following command to download a file 5 times from the HTTP server.**
   ```bash
   ./t.sh > r.txt
   ```

5. **The output will show, per line (5 lines), the download speed for each attempt in B/s inside the r.txt file (once it is finished). In our results, we report speeds in Mb/s or Kb/s. The evaluator can convert the results, calculate the average, and compare them to the results shown in Figure 7, which should be similar. We provide a script that performs this conversion (Kb/s) and calculates the average. It takes a text file as input and prints the result. The script can be run using the following command inside the client virtual machine (in a new console, after the five downloads).**
   ```bash
   python3 helper.py r.txt
   ```


6. **Repeat the process until all network conditions have been tested for all connection types.**

7. **After finishing this experiment, halt or destroy the VMs for a clean start, as described in the "Testing the Environment" section.**


#### Experiment 2: Varying video parameterizations for Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox connections.
To perform this experiment, related to Main Result 2, follow the instructions below. It should take about 2–4 hours to perform all encodings and approximately 6-7 hours to complete the testing. These time estimates are highly pessimistic for some videos, as different encodings result in varying throughputs; higher throughputs can complete much faster, taking only 10–15 minutes. All encoded videos will require roughly 10 GB of storage.

1. **For this test, the same video must be encoded with different configurations. We have a total of 4 different bitrate configurations and 8 different resolution/bitrate configurations that must be encoded before performing this experiment. Execute the following commands to encode the videos that will be used in this experiment. Note that for web-to-web connections, even though we specify the bitrate, the browser codec will encode the video according to its own settings. Therefore, the FFmpeg encoding is primarily important for the Pion-to-Pion connections.**

   **NOTE**: For functionality testing, a single video can be used to test all three connection types. However, by fully following this experiment as we describe it, it is possible to reproduce our results, either with the suggested 5 downloads for a less time-consuming version or the 100 downloads we performed in the paper.

   - **Bitrate:** To test different bitrates, perform the following commands using the Big Buck Bunny video downloaded in the "Set up the Environment" section. `INPUT_FILE` should be the name of the downloaded video from the "Testing the Environment" section, and the commands should be executed in the same directory where the video is located.

   1. 500 Kbps
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 -b:v 500k 1280_720_500k.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -c:v libvpx -b:v 500k -g 30 1280_720_500k.webm
      ```

   2. 1 Mbps
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 -b:v 1M 1280_720_1M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -c:v libvpx -b:v 1M -g 30 1280_720_1M.webm
      ```

   3. 2 Mbps
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 -b:v 2M 1280_720_2M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -c:v libvpx -b:v 2M -g 30 1280_720_2M.webm
      ```

   4. 5 Mbps
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 -b:v 5M 1280_720_5M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -c:v libvpx -b:v 5M -g 30 1280_720_5M.webm
      ```

   - **Resolution:** To test different resolutions, perform the following commands using the Big Buck Bunny video downloaded in the "Set up the Environment" section. `INPUT_FILE` should be the name of the downloaded video from the "Testing the Environment" section, and the commands should be executed in the same directory where the video is located.

   1) 426x240 (no bitrate specified) - ONLY FOR PION-TO-PION
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=426:240 -g 30 426_240_n.ivf
      ```

   2) 854x480 (no bitrate specified) - ONLY FOR PION-TO-PION
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=854:480 -g 30 854_480_n.ivf
      ```

   3) 1280x720 (no bitrate specified) - ONLY FOR PION-TO-PION
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 1280_720_n.ivf
      ```

   4) 1920x1080 (no bitrate specified) - ONLY FOR PION-TO-PION
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1920:1080 -g 30 1920_1080_n.ivf
      ```

   5) 426x240 (2M bitrate specified) 
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=426:240 -g 30 -b:v 2M 426_240_2M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=426:240 -c:v libvpx -b:v 2M -g 30 426_240_2M.webm
      ```

   6) 854x480 (2M bitrate specified) 
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=854:480 -g 30 -b:v 2M 854_480_2M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=854:480 -c:v libvpx -b:v 2M -g 30 854_480_2M.webm
      ```

   7) 1280x720 (2M bitrate specified) 
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -g 30 -b:v 2M 1280_720_2M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1280:720 -c:v libvpx -b:v 2M -g 30 1280_720_2M.webm
      ```

   8) 1920x1080 (2M bitrate specified) 
      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1920:1080 -g 30 -b:v 2M 1920_1080_2M.ivf
      ```

      ```bash
      ffmpeg -i INPUT_FILE -vf scale=1920:1080 -c:v libvpx -b:v 2M -g 30 1920_1080_2M.webm
      ```

2. **Move the video to be tested to the correct location, depending on whether it will be used in a Pion-to-Pion or web-to-web connection. There is no need to move audio files for Pion-to-Pion connections, since, if all previous steps from the “Set up the Environment” section have been followed, an OGG file will already exist in the correct location. It does not need to be changed for performance testing purposes because Obscura does not use the audio stream to encapsulate data. `NAME_OF_FILE` should be the name of the video being copied and correspond to one of the encoded output videos from the previous step. `PATH` should be the path to the root directory of the Obscura artifact.**
   - For Pion-to-Pion connections:
      ```bash
      cp NAME_OF_FILE.ivf "PATH/Obscura---Artifact/Client/Media/output.ivf
      cp NAME_OF_FILE.ivf "PATH/Obscura---Artifact/Proxy/Media/output.ivf
      ```

   - For web-to-web connections:
      ```bash
      cp NAME_OF_FILE.webm "PATH/Obscura---Artifact/ClientBrowser/Node-Server/videos/output.webm
      cp NAME_OF_FILE.webm "PATH/Obscura---Artifact/Proxy-Web/videos/output.webm
      ```

3. **Set up the baseline network conditions for the client and proxy to the value "2" (RTT of 115 ms). To configure the network condition, the `entry.sh` file in the Proxy and Client folders for Pion, and in the ClientBrowser and Proxy-Web folders for web-based instances, must be modified as described in the previous experiment. Specifically, the `TEST_N` variable inside the script should be assigned the value "2".**


4. **To run the Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox connections, the evaluator must follow the "Testing the Environment" section of this document to run the artifact with all these connection types. Completely follow steps 1) and 2) to make the artifact run, but do not follow step 3) for the download test.** (we assume that the "Set up the environment" section of this document has already been completed. From this point onward, we also assume that the artifact is running with any of these three connection types, but only one at a time).


5. **In a different console, connect to the client virtual machine (the working directory of the console must be the same as the root directory of the artifact where the Vagrantfile is).**
   - For web-based instances:
      ```bash
      vagrant ssh client_web
      ```

   - For Pion-based instances:
      ```bash
      vagrant ssh client
      ```


6. **Inside the client virtual machine, run the following command to download a file 5 times from the HTTP server.**
   ```bash
   ./t.sh > r.txt
   ```


7. **The output will show, per line (5 lines), the download speed for each attempt in B/s inside the r.txt file (once it is finished). In our results, we report speeds in Mb/s or Kb/s. The evaluator can convert the results, calculate the average, and compare them to the results shown in Figure 8, which should be similar. We provide a script that performs this conversion (Kb/s) and calculates the average. It takes a text file as input and prints the result. The script can be run using the following command inside the client virtual machine (in a new console, after the five downloads).**
   ```bash
   python3 helper.py r.txt
   ```

8. **Repeat the process until all differently encoded videos are tested for all connection types.**
 
9. **After finishing this experiment, halt or destroy the VMs for a clean start, as described in the "Testing the Environment" section.**


#### Experiment 3: DTLS Fingerprintability for web-to-web Connections.
This experiment/discussion is related to main result 4, and we performed experimental testing to support it. However, we believe it is sufficient to discuss why this result holds instead of repeating the actual experiment, since the dataset contains sensitive data specific to the authors' devices used to collect DTLS fingerprints from Firefox and Chrome when using WebRTC-based applications (the dataset can be provided upon request if necessary). 

In the experiment described in the paper, we used a dataset that we collected with 50 handshakes per application per browser. The applications we used were Discord, Google Meet, Facebook Messenger, and an official sample WebRTC application used for demonstration purposes (we based our WebRTC streaming web application on this sample application). This gives a total of 100 handshakes per application, with 50 collected using Firefox and 50 with Chrome. Our claim is that the DTLS fingerprint of web-to-web connections will match at least one other DTLS fingerprint within the dataset. Since it is not possible to alter the DTLS handshake using a web application like the one we developed, in the worst case, the DTLS handshake from our web-to-web connections will be identical to that of the sample application for both the ClientHello and ServerHello messages, or to that of any other WebRTC application that uses the browser as both a client and a server. In the paper, we further show that the fingerprints are also identical to other applications in the dataset beyond the sample WebRTC application. However, for the purposes of our claim, we believe this is sufficient, as the potential collateral effect of blocking our fingerprints would be to block any other application that uses the default DTLS implementation of the browsers in use.


#### Experiment 4: Pion-to-Pion Resistance to the Media-Based Variant of Differential Degradation Attacks.
To perform this experiment, related to Main Result 5, follow the instructions below (it should take about 5 hours to complete):

1. **We implement the attack as described in its original paper. The basic idea is to drop a percentage of frames by targeting a number of packets per frame marked for dropping (since a single frame can be composed of multiple RTP packets). Thus, we define two scenarios: the best-case scenario and the worst-case scenario. In the best-case scenario, a single packet is dropped from each frame marked for dropping, and in the worst-case scenario, all packets belonging to each marked frame are dropped. We test both scenarios. For each scenario, we test two different video resolutions: `426×240` and `1920×1080`. For each resolution, we also use two different encodings: one where we set the bitrate in FFmpeg to 2M, and one where no target bitrate is enforced. Thus, it is possible to reuse the encoded videos from Experiment 2, specifically the following video files: `1920_1080_2M.ivf`, `1920_1080_n.ivf`, `426_240_2M.ivf`, and `426_240_n.ivf`.**


   **NOTE:** For functionality testing, a single video can be used for a single scenario to test the attack. However, by fully following this experiment as we describe it, it is possible to reproduce our results, either with the suggested 5 downloads for a less time-consuming version or the 100 downloads we performed in the paper.

2. **Copy the video to be tested to the respective folders (where `FILE_NAME` is one of the files defined above in step 1).**
   ```bash
   cp FILE_NAME.ivf "/Obscura---Artifact/Client/Media
   cp FILE_NAME.ivf "/Obscura---Artifact/Proxy/Media
   ```

3. **Inside both the Client and Proxy folders, there are two Python scripts called `DDA.py` and `DDA2.py`. `DDA.py` implements the attack for the best-case scenario, and `DDA2.py` implements the attack for the worst-case scenario. For the tests presented in the paper, we also varied the percentage of frames that were dropped.**

   We use the following percentages when testing videos with the 1920×1080 resolution, independently of the specified bitrate and scenario (best or worst).
   1. 0%
   2. 5%
   3. 15%
   4. 25%
   5. 45%

   Additionally, we use the following percentages when testing videos with the 426×240 resolution, independently of the specified bitrate and scenario (best or worst):
      1. 0%
      2. 20%
      3. 30%
      4. 50%
      5. 70%

   The reader/evaluator can also see this in Figures 10 and 11 of the paper, which provide a visualization of the test setup.
   
   To define the percentage of frames to drop, a variable called `FRAME_DROP_RATE` is available within the `DDA.py` (best-case script) and `DDA2.py` (worst-case script) scripts. The evaluator should assign a value between 0.0 and 1.0 that corresponds to the drop percentages defined above for each test. The same value should be assigned in both the Proxy and the Client scripts.
   
4. **The `entry.sh` file in both the Client and Proxy folders contains a variable called `DDA`. A value of `case_b` should be assigned when testing the system under the best‑case scenario in both the Client’s and Proxy’s `entry.sh` files, and a value of `case_w` should be assigned when testing the system under the worst‑case scenario in both files.**

   **Note:** To disable differential degradation attacks, the value `no` should be assigned to the `DDA` variable in the `entry.sh` file of both the Client and the Proxy.

5. **Once everything is set up, to run the Pion-to-Pion connection, the evaluator must follow the "Testing the Environment" section of this document to run the artifact. Follow step 1) completely to make the artifact run, but do not follow step 2) for the Web-to-Web connection setup or step 3) for the download test. (Start with the best-case scenario at 0%, increasing the percentage after each test and moving to the worst-case scenario once the best-case tests are complete).** (we assume that the "Set up the environment" section of this document has already been completed. From this point onward, we also assume that the artifact is running the Pion-to-Pion connection).

6. **In a different console, connect to the client virtual machine (the working directory of the console must be the same as the root directory of the artifact where the Vagrantfile is).**
   ```bash
   vagrant ssh client
   ```

7. **Inside the client virtual machine, run the following command to download a file 5 times from the HTTP server.**
   ```bash
   ./t.sh > r.txt
   ```

8. **The output will show, per line (5 lines), the download speed for each attempt in B/s inside the r.txt file (once it is finished). In our results, we report speeds in Mb/s or Kb/s. The evaluator can convert the results, calculate the average, and compare them to the results shown in Figures 10 and 11, which should be similar. We provide a script that performs this conversion (Kb/s) and calculates the average. It takes a text file as input and prints the result. The script can be run using the following command inside the client virtual machine (in a new console, after the five downloads).**
   ```bash
   python3 helper.py r.txt
   ```

8) **Repeat the process until all percentages of dropped frames have been tested for all four encoded videos and for both scenarios.**
 
   
#### Experiment 5: Testing Proxy Ephemerality
To test that proxies can go offline and that the connection is automatically restarted once a new proxy comes online, follow the steps below (related to Main Result 6) (estimated time: 30 minutes to 1 hour). Note: This experiment is not included in the paper; it is presented here to validate one of the features of Obscura.

1. **Run the Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox connections by following the "Testing the Environment" section of this document to run the artifact with all three types of connections. Completely follow steps 1) and 2) to make the artifact run, but do not follow step 3) for the download test.** (we assume that the "Set up the environment" section of this document has already been completed. From this point onward, we also assume that the artifact is running with any of these three connection types, but only one at a time).

2. **On the virtual machine running the proxy, cancel the execution of the `entry.sh` script by pressing Ctrl + C (or, if the virtual machine becomes unresponsive, by halting and restarting it). For Pion-to-Pion connections, you should see the following output on the Client virtual machine: "Peer Connection has gone to failed exiting." For Web-to-Web connections, no output should appear.**

3. **Once the proxy virtual machine has restarted or `entry.sh` is no longer running, start the `entry.sh` script again in the proxy virtual machine.**

4. **After some seconds (note that the Pion-to-Pion connection may take up to one minute), the client virtual machine should display a second “Connected” message if the artifact is running a Pion-to-Pion connection, or a second “READY” message if it is running a web-to-web connection. If this appears, it means the client successfully reconnected to the new proxy while maintaining the state of the connection between the client and the bridge.**

5. **Repeat this process for all connection types (i.e., Pion-to-Pion, Chrome-to-Chrome, and Firefox-to-Firefox).**
 

## Limitations 
Main Result 3 is not fully reproducible because the artifact does not include the videos used in the experiment, nor are they available online due to participant consent and copyright restrictions. However, the expected outcomes can be inferred from Main Result 2, which shows that videos with higher bitrates achieve higher throughput, while profiles such as coding videos naturally exhibit lower bitrates due to their static nature. We believe this does not hinder the reproducibility badge, as the underlying behavior remains consistent. Moreover, even if the exact videos cannot be reproduced, the experiment itself can be replicated by collecting videos representative of each profile and testing them with each connection type.
   

## Notes on Reusability 
When using the system as a Tor pluggable transport (tested only in the appendix), the client will be ready once the Tor print shows 100% on the screen. Additionally, the HTTP server must be publicly accessible over the Internet, since the inbound connection will come from the Tor exit relay, which is outside the control of our testing environment. The download scripts (i.e., `t.sh`) will also need to have their endpoints adjusted to the address they access (i.e., the public IP of the HTTP server). 

The folders containing the code for the components that differ when the artifact is running in Tor pluggable transport mode are:
- BridgeTor: contains the code for running the bridge as a Tor Bridge.
- ClientTor: contains the code for running the client as a Pion instance in Tor pluggable transport mode.
- ClientBrowserTor: contains the code for running the client as a Web instance in Tor pluggable transport mode.

(A virtual machine for each of these components is defined in the Vagrantfile; Apart from this, the setup process is the same for running the system).

We also provide a way to use self-generated animations instead of video files (tested in the appendix of the paper). The code for this is provided in the following folders:
- ProxyAnimation: for the proxy as a Pion instance with an animation running instead of a video.
- ClientAnimation: for the client as a Pion instance with an animation running instead of a video.

For web-based instances, the code for the self-generated animations is not organized into separate folders; instead, it resides within the Proxy-Web, ClientBrowser, and ClientBrowserTor folders. However, to run the animation code for web-based instances, the Python scripts inside the Selenium folder, available within each of the previously listed folders and which automate browser functionality, must have the following line modified (where `URL` represents the remainder of the URL already present in the unmodified line and should remain the same):

```bash  
url = "URL/video"
```

to

```bash
url = "URL/canvas"
```