package main

import (
	"fmt"
	"net"
	"net/http"
	"syscall"
	"time"
	"unsafe"
)

var (
	user32          = syscall.NewLazyDLL("user32.dll")
	procSendMessage = user32.NewProc("SendMessageW")
	procFindWindow  = user32.NewProc("FindWindowW")

	hostAwake       = false
	hostIgnoreCount = 0
)

const (
	SM_MOUSEPRESENT = 19
	WM_SYSCOMMAND   = 0x0112
	SC_MONITORPOWER = 0xF170

	MONITOR_ON  = -1
	MONITOR_OFF = 2
)

type HMONITOR syscall.Handle
type HDC syscall.Handle
type HWND syscall.Handle

func mouseConnected() bool {
	getSystemMetrics := user32.NewProc("GetSystemMetrics")
	ret, _, _ := getSystemMetrics.Call(uintptr(SM_MOUSEPRESENT))
	return ret != 0
}
func controlMonitor(power int) {
	if power == MONITOR_ON {
		mouse_event := user32.NewProc("mouse_event")
		const MOUSEEVENTF_MOVE = 0x0001
		mouse_event.Call(MOUSEEVENTF_MOVE, uintptr(int32(1)), uintptr(int32(0)), 0, 0)
		mouse_event.Call(MOUSEEVENTF_MOVE, uintptr(int32(10)), uintptr(int32(0)), 0, 0)
		return
	}

	cwc, _ := syscall.UTF16PtrFromString("ConsoleWindowClass")
	hwnd, _, _ := procFindWindow.Call(uintptr(unsafe.Pointer(cwc)), 0)
	if hwnd == 0 {
		hwnd, _, _ = procFindWindow.Call(0, 0)
	}

	_, _, err := procSendMessage.Call(hwnd,
		uintptr(WM_SYSCOMMAND),
		uintptr(SC_MONITORPOWER),
		uintptr(power))

	if err != nil && err.Error() != "The operation completed successfully." {
		fmt.Println("Failed to send message:", err)
	}

}

func hideConsole() {
	hwnd, _, _ := syscall.NewLazyDLL("kernel32.dll").
		NewProc("GetConsoleWindow").Call()
	if hwnd != 0 {
		syscall.NewLazyDLL("user32.dll").
			NewProc("ShowWindow").Call(hwnd, 0) // 0 = SW_HIDE
	}
}

func discoverListener(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "tinyrabbit-go monitor: %t", hostAwake)
}

func wakeListener(w http.ResponseWriter, r *http.Request) {
	hostAwake = true
	//to-do: store ip for sending wake requests
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		fmt.Println("Received http wake request from:", ip)
	} else {
		fmt.Println("Received http wake request")
	}
	hostIgnoreCount = 7
	controlMonitor(MONITOR_ON)
	fmt.Fprintf(w, "tinyrabbit-go: OK")
}

func rabbitWorker() {
	for {
		time.Sleep(1 * time.Second)
		isMouseConnected := mouseConnected()
		if isMouseConnected && !hostAwake {
			fmt.Println("Mouse is connected, turning on output")
			controlMonitor(MONITOR_ON)
			hostAwake = true
		} else if !isMouseConnected && hostAwake {
			if hostIgnoreCount < 1 {
				fmt.Println("Mouse disconnected, turning off output")
				controlMonitor(MONITOR_OFF)
				hostAwake = false
			}
		}
		if hostIgnoreCount > 0 {
			hostIgnoreCount--
		}
	}

}

func main() {
	hideConsole()
	go rabbitWorker()
	http.HandleFunc("/wake", wakeListener)
	http.HandleFunc("/discover", discoverListener)
	err := http.ListenAndServe("0.0.0.0:11812", nil)
	if err != nil {
		fmt.Println("failed to start http server")
	}
}
