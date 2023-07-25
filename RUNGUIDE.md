## 🏃 Run Horus in the background
When using a Raspberry Pi, it's pretty common to connect to it via `SSH`. If you execute Horus, and close the `SSH` connection, Horus' process will be killed by Linux, because it was opened from a connection that was terminated. To prevent this, we have to run Horus in the background, via `nohup` command.

Also, keep in mind that Horus needs `sudo` privileges to access to GPIO Pins to be able to control the Led Strip, while also accessing the `HTTP Port 80` to be able to run the server. To prevent Linux asking for sudo password, while also making it run in the background, and run on the Raspberry Pi startup, you can follow this steps (in order):

### 📜 Make the script that opens Horus from it's directory folder
First, check where you cloned `Horus` repository (or where you placed the release folder `horus_vX.X.X`), for example, in my case, it's in my home folder (`/home/YOUR_USERNAME/horus`) (replace `YOUR_USERNAME` with your Linux username). The first thing we want to do is to create a small `bash` script to run `Horus` easily (you can place this file wherever you want).

Let's call it, `horusStartup.sh`. I will place it in my home folder (`/home/YOUR_USERNAME/horusStartup.sh`):
```bash
  nano /home/YOUR_USERNAME/horusStartup.sh
```
Then, write this:
```bash
#!/bin/bash

cd /home/YOUR_USERNAME/horus   # If downloaded a release, "/home/YOUR_USERNAME/horus_X.X.X"
sudo ./horus                   # Executes Horus script
```
Then we can save and close the file (if using `nano`, CTRL+X -> Y -> And press enter)
Once that done, we have to make that script we created an executable:
```bash
  chmod +x horusStartup.sh
```


### 🔌 Make Linux execute that script in power on/startup
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
Once saved, we have to enable the shortcut, so once we reboot, it will be executed:
```bash
  sudo systemctl enable horusStartup.service
```


### 🧙 Prevent Linux from asking for a sudo password on background/startup
Now, here's the problem; Once the Raspberry Pi turns on, you will be asked for sudo privileges while executing that script, and if you have your Raspberry Pi not easily accessible (like I do), It won't execute Horus. To prevent that, we have to prevent that golang command with sudo privileges not to ask for a password. It's easy. First, open sudoers configuration file:
```bash
  sudo visudo
```
Then, go to the last line, and place this:
```bash
  YOUR_USERNAME ALL=(ALL) NOPASSWD: /usr/bin/nohup /home/YOUR_USERNAME/horusStartup.sh # Make the script executable without password
  YOUR_USERNAME ALL=(ALL) NOPASSWD: /home/YOUR_USERNAME/horus/horus                    # Make Horus executable without passwod
```
Save the file, and we should be ready to go.

Now, Linux won't ask for the sudoers password while executing Horus. Now, you can try rebooting your Raspberry Pi. After waiting for a reasonable time to turn on and execute Horus, you should be able to access it via `http://YOUR_RASP_LOCAL_IP:80/` in your browser.

I hope that was clear. If you find any issues, feel free to contact me at hugomoreda@hotmail.com