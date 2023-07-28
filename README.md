<p align="center">
<img height="120" src="https://i.imgur.com/U5shZKs.png" />
</p>
<p align="center">
  <i>A Raspberry Pi & ws281x Led strip control panel, written in Golang</i>
   <br/>
  
  <br/>
  <a href="LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-green.svg" alt="333">
  </a>
  <a>
    <img src="https://img.shields.io/github/stars/keelus/horus" alt="dsa1">
  </a>
  <a>
    <img src="https://img.shields.io/github/downloads-pre/keelus/horus/latest/total" alt="asd">
  </a>
  <a>
    <img src="https://img.shields.io/badge/made%20with-%E2%98%95%EF%B8%8F%20coffee-yellow.svg" alt="asd">
  </a>
</p>




# Horus
Horus is a project I decided to make to be able to control the ws281B Led Strip that I installed on my house from anywhere (usually for my computer or phone). The led strip is attached to the Raspberry Pi's GPIO pins.

Horus has a well-made user interface (compatible with mobile devices) that lets the user log in, control the led strip & view your Pi's stats (CPU temperature, RAM usage, etc). 


## 🛠️ Tech Stack
**Client:** HTML, JavaScript/jQuery, CSS/Sass
<br>
**Server:** Gin, Golang


## ▶️ Demo
![Horus Demo](https://i.imgur.com/gMdqeiE.gif)


## ✨ Features
- ws281x led strip compatible software
- Save color and gradient presets on all 4 different modes available now:
  - Static color: Draws a color that will remain static.
  - Static gradient: Create your own gradient combining the colors you want.
  - Fading rainbow: A visually appealing rainbow that moves through your strip
  - Breathing color: A breathing/pulsating effect on the color you want
- Raspberry Pi live stats: CPU temperature, usage, RAM usage, Disk space & system uptime live in your browser.
- Light & dark interface color modes.
- Semi-customizable interface
  - Show or hide the features that you want or don't want to see!
  - Choose your desired session cookie lifetime.
- Logging system: Keep all the errors that could happen & changes you make to your Led & configuration logged.
- Restful API: Integrate Horus with your own scripts (e.g. change your led color)


## ️🚀 Installation
Clone the project
```bash
  git clone https://github.com/keelus/horus
```
or download the latest release. If you don't have GUI/screen and have to do it via terminal:
```bash
	wget https://github.com/keelus/horus/releases/download/vX.X.X/horus.zip # Replace X.X.X with the version of the release you want to install
```

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

Also, to prevent Linux from asking for sudoers password, or to prevent it from stopping `Horus` when disconnected from the `SSH` connection, please check the [RUN GUIDE](RUNGUIDE.md).

## 📦 Build it yourself
After cloning the repo and entering the project directory, install the dependencies:
```bash
  go mod tidy
```
then, go into the `cmd` folder:
```bash
  cd ./cmd
```
And run:
```bash
  go build -o horus
```
Now, you will be left with a `horus.sh`, which I recommend placing into the project directory to prevent issues from GO_PATH & relative path issues. Then, you can run it by:
```bash
  sudo ./horus
```

To prevent Linux from asking for sudoers password, or to prevent it from stopping `Horus` when disconnected from the `SSH` connection, please check the [RUN GUIDE](RUNGUIDE.md).


## 📸 Screenshots (big & small screen)
<p float="left">
  <img src="https://i.imgur.com/pA497dQ.gif" height="400"/>
  <img src="https://i.imgur.com/pRKPjEP.gif" height="400"/>
</p>


## 🤔 FAQ
#### Where should I connect my sw281X Led strip?
The led strip is connected via GPIO to the Raspberry Pi. You should connect the led strip data line (usually green) to pin 12 (GPIO 18) (as seen [here](https://i.imgur.com/nncVgoZ.png). It can vary where that pin is located depending on your Raspberry Pi model)

#### Why does Horus need sudo privileges?
To be able to control your ws281x Led strip, Horus needs access to the Raspberry Pi's GPIO, where the strip is connected. Also, the web server is initialized by default on port 80 (which is a privileged port on Linux).


## ⚖️ License
This project is open source under the terms of the [MIT License](https://github.com/keelus/horus/blob/main/LICENSE)


## 📬Feedback

If you have any feedback, please reach out to me at hugomoreda@hotmail.com

<br><br>

Made by [@keelus](https://github.com/keelus)