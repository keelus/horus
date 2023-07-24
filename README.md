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

Install dependencies
```bash
  go mod tidy
```

Start horus
```bash
  sudo go run ./cmd
```

On the first time execution, you will be asked to enter a username and password, which will be used to log in.


## ⚖️ License
This project is open source under the terms of the [MIT License](https://github.com/keelus/horus/blob/main/LICENSE)


## 📸 Screenshots
![App Screenshot](https://via.placeholder.com/468x300?text=App+Screenshot+Here)


## 🤔 FAQ
#### Why do Horus need sudo privileges?

To be able to control ws281x Led strip, Horus needs access to the Raspberry Pi's GPIO, where the strip is connected. Also, the web server is initialized by default on port 80 (which is a privileged port on linux).


## 📬Feedback

If you have any feedback, please reach out to me at hugomoreda@hotmail.com