# tinyrabbit-go
------------
Simple binary to help with KVM Setups in Windows.<br/>
For MacOS use [tinyrabbit](https://github.com/b1on1cdog/tinyrabbit).<br/>

# Purpose
Most monitors automatically switch to the active output<br/>
Using tinyrabbit along with a cheap USB switcher (aka: USB KM/KVM) we can skip using a KVM to switch the video<br/>
This also means a direct video connection from your computer to your Monitor, taking advantage of your full monitor specifications<br/>

# Features
- Turn off monitor when Mouse is disconnected<br/>
- Wake Monitor when Mouse is connected<br/>
- Wake Monitor when a HTTP Request to :11812/wake is received<br/>

# To-do:
- Add Linux support<br/>
- Auto-add binary to user login<br/>