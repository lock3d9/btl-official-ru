package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"

	bitlockerw "btl/funcs/bitlocker"
	checkdiskerrors "btl/funcs/checkdisk"
	checkflashff "btl/funcs/checkflashfake"
	"btl/funcs/checkheavyfiles"
	"btl/funcs/cleantemp"
	"btl/funcs/clearcachewindowsupdate"
	autostart "btl/funcs/controlautostart"
	"btl/funcs/createbackup"
	deleteonedriv "btl/funcs/deleteonedrive"
	"btl/funcs/deleteyarlk"
	"btl/funcs/disabledefender"
	sysinfo "btl/funcs/fullinformationpc"
	"btl/funcs/killtasks"
	"btl/funcs/offsmartscreen"
	"btl/funcs/offstickingbuttons"
	"btl/funcs/offupdatewindows"
	recoverysystemdism "btl/funcs/recoverysystem/dism"
	recoverysystem "btl/funcs/recoverysystem/sfc"
	"btl/funcs/showfilextandhidefiles"
	interfface "btl/menucli"

	"golang.org/x/sys/windows"
)

var (
	Version   = "00.1"
	BuildNum  = "0"
	BuildDate = "unknown"

	ntdll         = syscall.NewLazyDLL("ntdll.dll")
	rtlGetVersion = ntdll.NewProc("RtlGetVersion")
	user32        = syscall.NewLazyDLL("user32.dll")
	messageBox    = user32.NewProc("MessageBoxW")
)

type RTL_OSVERSIONINFOW struct {
	DwOSVersionInfoSize uint32
	DwMajorVersion      uint32
	DwMinorVersion      uint32
	DwBuildNumber       uint32
	DwPlatformId        uint32
	SzCSDVersion        [128]uint16
}

func CheckWindowsVersion() {
	var osInfo RTL_OSVERSIONINFOW
	osInfo.DwOSVersionInfoSize = uint32(unsafe.Sizeof(osInfo))

	status, _, _ := rtlGetVersion.Call(uintptr(unsafe.Pointer(&osInfo)))
	if status == 0 && osInfo.DwBuildNumber < 22000 {
		showErrorBox(
			"Ошибка совместимости",
			fmt.Sprintf("Ваша система (Build %d) не поддерживается.\n\nBackTooL работает только на Windows 11 и новее.", osInfo.DwBuildNumber),
		)
		os.Exit(1)
	}
}

func showErrorBox(title, message string) {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	messagePtr, _ := syscall.UTF16PtrFromString(message)

	// MB_OK (0x0) | MB_ICONERROR (0x10) | MB_TOPMOST (0x40000)
	messageBox.Call(0, uintptr(unsafe.Pointer(messagePtr)), uintptr(unsafe.Pointer(titlePtr)), 0x00040010)
}

func isAdmin() bool {
	var sid *windows.SID
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY, 2,
		windows.SECURITY_BUILTIN_DOMAIN_RID, windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0, &sid,
	)
	if err != nil {
		return false
	}
	defer windows.FreeSid(sid)

	token := windows.Token(0)
	member, err := token.IsMember(sid)
	return err == nil && member
}

func runAsAdmin() {
	exe, err := os.Executable()
	if err != nil {
		exe = os.Args[0]
	}
	cmd := exec.Command("powershell", "Start-Process", fmt.Sprintf("%q", exe), "-Verb", "RunAs")
	if err := cmd.Run(); err != nil {
		fmt.Println("[BTL] Не удалось запросить права администратора:", err)
	}
}

func forceWindowsTerminal() {
	if os.Getenv("WT_SESSION") == "" {
		if exe, err := os.Executable(); err == nil {
			cmd := exec.Command("wt", "new-tab", exe, strings.Join(os.Args[1:], " "))
			if cmd.Start() == nil {
				os.Exit(0)
			}
		}
	}
}

func main() {
	CheckWindowsVersion()
	forceWindowsTerminal()

	if !isAdmin() {
		fmt.Println("[BTL] Программе требуются права администратора. Делаем запрос...")
		runAsAdmin()
		os.Exit(0)
	}

	actions := map[string]func(){
		"Проверить диск на чтение, запись и скорость": checkdiskerrors.Ccheckdsk,
		"Полная информация о ПК":                      sysinfo.ShowSystemInfo,
		"Проверить флешку на муляж":                   checkflashff.H2testwCheck,
		"Убить зависшие приложения":                   killtasks.KillFrozenApps,
		"Анализ «тяжелых» файлов":                     checkheavyfiles.FindHeavyFiles,
		"Контроль автозагрузки":                       autostart.ManageAutostart,
		"[РЕКОМЕНДУЕМ] Создать точку восстановления":  createbackup.CreateRestorePoint,
		"Очистка временных файлов":                    cleantemp.CleanTempFiles,
		"Удалить встроенный Windows Defender":         disabledefender.DisableDefender,
		"Удалить OneDrive":                            deleteonedriv.RemoveOneDrive,
		"Показать скрытые файлы и папки":              showfilextandhidefiles.ShowHiddenFiles,
		"Показывать расширения файлов":                showfilextandhidefiles.ShowFileExtensions,
		//"Выключить размытие на экране блокировки",
		"Отключить SmartScreen":                 offsmartscreen.DisableSmartScreen,
		"Отключить залипание клавиш":            offstickingbuttons.DisableStickyKeys,
		"Убрать окончание «Ярлык»":              deleteyarlk.RemoveShortcutSuffix,
		"Восстановить систему через SFC":        recoverysystem.RunSFCFix,
		"Восстановить систему через DISM":       recoverysystemdism.RunDISMFix,
		"Переключить авто-шифрование BitLocker": bitlockerw.ToggleBitLockerAutoEncryption,
		"Очистить кэш Windows Update":           clearcachewindowsupdate.ClearWindowsUpdateCache,
		"Переключить Windows Update (Вкл/Выкл)": offupdatewindows.ToggleWindowsUpdate,
	}

	for {
		choice := interfface.Showmainmenu(Version, BuildNum, BuildDate)
		if choice == "Выход" {
			fmt.Println("[BTL] Выходим...")
			os.Exit(0)
		}

		if action, exists := actions[choice]; exists {
			action()
		} else {
			fmt.Println("[BTL] Неизвестное действие...")
		}

		fmt.Printf("[BTL] Версия: v%s (build #%s | %s)\n", Version, BuildNum, BuildDate)
	}
}
