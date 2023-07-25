<p align="center">
<img height="120" src="https://i.imgur.com/U5shZKs.png" />
</p>
<p align="center">
  <i>A Raspberry Pi & ws281x Led strip control panel, written in Golang</i>
   <br/>
  
  <br/>
  <a href="#Demo">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="333">
  </a>
  <a href="#Demo">
    <img src="https://img.shields.io/github/stars/keelus/horus" alt="dsa1">
  </a>
  <a href="#Demo">
    <img src="https://img.shields.io/github/downloads-pre/keelus/horus/latest/total" alt="asd">
  </a>
</p>


# Horus
Horus is a project I decided to make to be able to control the ws281B Led Strip that I installed on my house from anywhere (usually for my computer or phone). The led strip is attached to the Raspberry Pi, so, I the project focuses in that.

Horus has a well made user interface (compatible with mobile devices) that lets the user log in, and control the led strip & view your Pi's stats (cpu temperature, ram usage, etc). 


## 🛠️ Tech Stack
**Client:** HTML, JavaScript/jQuery, CSS/Sass
<br>
**Server:** Gin, Golang


## ▶️ Demo
![Horus demo GIF](https://via.placeholder.com/468x300?text=Horus+GIF)


## ✨ Features
- ws281x led strip compatible software
- Save color and gradient presets on all 4 different modes available now:
  - Static color: Draws a color that will remain static.
  - Static gradient: Create your own gradient combining the colors you want.
  - Fading rainbow: A visually appealing rainbow that moves throught your strip
  - Breathing color: A breathing/pulsating effect on the color you want
- Raspberry Pi live stats: CPU temperature, usage, RAM usage, Disk space & system uptime live in your browser.
- Light & dark interface color modes.
- Semi-customisable interface
  - Show or hide the features that you want or don't want to see!
  - Choose your desired session cookie lifetime.
- Logging system: Keep all the error that could happen & changes you make to your Led & configuration logged.
- Restful API: Integrate Horus with your own scripts (e.g. change your led color)
- Made with ❤️


## 👨‍💻 Authors
- [@keelus](https://www.github.com/keelus)


## 📚 Appendix
Any additional information goes here


## ️🚀 Installation
Clone the project
```bash
  git clone https://github.com/keelus/horus
```
or download the latest release
<br>

Then, go to the project directory
```bash
  cd horus
```

And simply run
```bash
  sudo ./horus
```
On the first time execution, you will be asked to enter a username and password, which will be used to log in.
<br>
When using a Raspberry Pi it's common to use SSH to connect to it. To be able to run horus, and keep it running after closing it, check (this)[#Run Horus in the background]

## 📦 Build it yourself
After cloning the repo and entering the project directory, install the dependencies:
```bash
  go mod tidy
```
then, go into `cmd` folder:
```bash
  cd ./cmd
```
And run:
```bash
  go build -o horus
```
Now, you will be left with a `horus.sh`, which I recommend placing into the project directory. Then, you can run it by:
```bash
  sudo ./horus
```

## 🏃‍♂️ [Run Horus in the background]
When using a Raspberry Pi, it's pretty common to connect to it via `SSH`. If you execute Horus, and close the `SSH` connection, Horus' process will be killed by Linux, because it was opened with a connection that was terminated. To prevent this, we have to run Horus in the background, via `nohup` command, which will be helpful. 

Also, keep in mind that Horus needs `sudo` privileges to access to GPIO Pins to be able to control the Led Strip, while also accessing the `HTTP Port 80` to be able to run the server. To prevent Linux asking for sudo password, while also making it run in the background, and run on the Raspberry Pi startup, we have to do the following!

### Make the script that opens Horus from it's directory folder
First, check where you cloned `Horus`, for example, in my case, it's in `~/horus` (`/home/YOUR_USERNAME/horus`). The first thing we want to do is to create a small `bash` script to run `Horus` easily (you can place this file wherever you want):
Let's call it, horusStartup.sh. I will place it in my home folder (`/home/YOUR_USERNAME/`):
```bash
  nano horusStartup.sh
```
Then, write this:
```bash
#!/bin/bash

cd /home/YOUR_USERNAME/horus   # Enters Horus directory [we need to be inside and then execute it, bcz GO_PATH]
sudo ./horus                   # Executes Horus script
```
Then we can save and close the file (if using `nano`, CTRL+X -> Y -> And press enter)
Once that done, we have to make that script we created an executable:
```bash
  chmod +x horusStartup.sh
```
### Make Linux execute that script in power on/startup
Now, we have to make Linux execute that script when we power on our Raspberry Pi. First, we create a systemd service file:
```bash
  sudo nano /etc/systemd/system/horusStartup.service
```
Then, we write the following:
```bash
[Unit]
Description=Horus Startup Script
After=network.target

[Service]
Type=simple
ExecStart=/home/YOUR_USERNAME/horusStartup.sh # Or wherever you placed that file we created earlier

[Install]
WantedBy=multi-user.target
```
Once saved, we have to enable it, so once we reboot, it will be executed:
```bash
  sudo systemctl enable horusStartup.service
```
### Prevent Linux from asking for a sudo password on background/startup
Now, here's the thing: Once the Raspberry Pi turns on, you will be asked for sudo privileges while executing that script, and if you have your Raspberry Pi not easily accessible (like I do), It won't execute Horus. To prevent that, we have to prevent that golang command with sudo privileges not to ask for a password. It's easy. First, open sudoers configuration file:
```bash
  sudo visudo
```
Then, go to the last line, and place this:
```bash
  YOUR_USERNAME ALL=(ALL) NOPASSWD: /usr/bin/nohup /home/YOUR_USERNAME/horusStartup.sh # Make the script executable without password
  YOUR_USERNAME ALL=(ALL) NOPASSWD: /home/YOUR_USERNAME/horus/horus                    # Make Horus executable without passwod
```
Save the file, and we should be ready to go. Now, Linux won't ask for the sudoers password while executing Horus. Now, you can try rebooting your Raspberry Pi. After waiting for a reasonable time to turn on and execute Horus, you should be able to access it via `http://YOUR_RASP_LOCAL_IP:80/` in your browser.

## ⚖️ License
This project is open source under the terms of the [MIT License](https://github.com/keelus/horus/blob/main/LICENSE)


## 📸 Screenshots
![App Screenshot](https://via.placeholder.com/468x300?text=App+Screenshot+Here)


## 🤔 FAQ
#### Why do Horus need sudo privileges?

To be able to control ws281x Led strip, Horus needs access to the Raspberry Pi's GPIO, where the strip is connected. Also, the web server is initialized by default on port 80 (which is a privileged port on linux).


## 📬Feedback

If you have any feedback, please reach out to me at hugomoreda@hotmail.com