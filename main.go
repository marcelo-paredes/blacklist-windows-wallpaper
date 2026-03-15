package main

import (
	"fmt"
	"log"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const (
	SPI_SETWALLPAPER = 0x0014
)

func main() {

	user32 := syscall.NewLazyDLL("user32.dll")
	systemParametersInfo := user32.NewProc("SystemParametersInfoW")
	_ = systemParametersInfo.Find()

	k, err := registry.OpenKey(windows.HKEY_CURRENT_USER, "Control Panel\\Desktop", registry.NOTIFY|registry.QUERY_VALUE /*|registry.SET_VALUE*/)

	if err != nil {

		log.Fatalf("RegOpenKeyEx failed: %v", err)

	}
	keyHandle := windows.Handle(k)

	cleanupKey := func(key windows.Handle) {
		fmt.Printf("Cleanup Key")
		_ = windows.RegCloseKey(key)
	}

	defer cleanupKey(keyHandle)

	fileName, _, _ := k.GetStringValue("WallPaper")
	fmt.Printf("Current WallPaper: %s\n", fileName)

	event, err := windows.CreateEvent(nil, 1, 0, nil) // manualReset=1, initialState=0
	if err != nil {
		log.Fatalf("Failed to create event: %v", err)
	}

	cleanupEvent := func(event windows.Handle) {
		fmt.Printf("Cleanup Event")
		_ = windows.CloseHandle(event)
	}

	defer cleanupEvent(event)

	fmt.Println("RegOpenKeyEx succeeded")
	for {
		err = windows.RegNotifyChangeKeyValue(
			keyHandle,
			false,                                                             // bWatchSubtree: true to watch subkeys as well
			windows.REG_NOTIFY_CHANGE_NAME|windows.REG_NOTIFY_CHANGE_LAST_SET, // dwNotifyFilter: Notify on name or value changes
			event,                                                             // hEvent: The event to signal
			true,                                                              // fAsynchronous: true for asynchronous operation
		)
		if err != nil {
			log.Fatalf("Failed to set registry notification: %v", err)
		}

		waitRet, err := windows.WaitForSingleObject(event, windows.INFINITE)

		if err != nil {
			log.Fatalf("WaitForSingleObject failed: %v", err)
		}

		if waitRet == windows.WAIT_OBJECT_0 {
			fmt.Printf("Change detected at %s!\n", time.Now().Format(time.RFC3339))
			wallpaper, _, err := k.GetStringValue("WallPaper")
			if err != nil {
				log.Fatalf("Failed to get Wallpaper from registry: %v", err)
			}
			fmt.Printf("Wallpaper is %s\n", wallpaper)

			setWallPaper := func(fileName string) {
				winFileName, _ := windows.UTF16PtrFromString(fileName)
				_ = systemParametersInfo.Find()
				_, _, lastErr := systemParametersInfo.Call(
					uintptr(SPI_SETWALLPAPER),
					uintptr(0x0000),
					uintptr(unsafe.Pointer(winFileName)),
					uintptr(0x01|0x02), // write to profile AND win.ini
				)
				if lastErr != nil {
					fmt.Printf("Error: %s\n", lastErr)
				}
				//_ = k.SetStringValue("WallPaper", fileName)
			}

			if strings.EqualFold(wallpaper, "C:\\Temp\\downloadedWallpaper.jpg") {
				fmt.Println("resetting to " + fileName)
				setWallPaper(fileName)
				setWallPaper(fileName)

				_ = windows.ResetEvent(event)
			} else {
				fileName, _, _ = k.GetStringValue("WallPaper")
				fmt.Printf("Current WallPaper: %s\n", fileName)
			}
			// Reset the event manually since we created it as a manual reset event
		} else {
			// Handle timeout (not possible here due to windows.INFINITE) or other wait results
			fmt.Printf("Wait returned: %d\n", waitRet)
		}
	}
}
