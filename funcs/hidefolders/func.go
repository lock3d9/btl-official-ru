package hidefolders

import (
	"fmt"
	"os/exec"

	"golang.org/x/sys/windows/registry"
)

var win11FolderGuids = []string{
	"{31C06258-0713-4170-8636-D5D317F8C5E9}",
	"{088e3905-0323-4b02-9826-5999d6092922}",
	"{24ad3114-04ec-497d-aa29-8756c802424b}",
	"{f8637988-a42e-4424-9bda-091219b1ed5d}",
	"{a0c69a99-21c8-4671-8703-7934162fcf1d}",
	"{f4246208-736e-417e-9db0-d5a3f313fe20}",
	"{B4BFCC3A-DB2C-424C-B029-7FE99A87C641}",
	"{0DB7E03F-FC29-4DC6-9020-FF41B59E513A}",
	"{d31a1100-ea5c-4558-8834-4280e22f074d}",
}

func HideExplorerLibraries() error {
	for _, guid := range win11FolderGuids {
		hidePropertyBag(guid)
		unpinFromCLSID(guid)
		removeDelegateFolder(guid)
	}

	_ = exec.Command("taskkill", "/F", "/IM", "explorer.exe").Run()
	_ = exec.Command("explorer.exe").Start()

	fmt.Println("\n\033[32m[УСПЕХ] Папки и библиотеки скрыты из Проводника Windows 11.\033[0m")
	return nil
}

func hidePropertyBag(guid string) {
	path := fmt.Sprintf(`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\FolderDescriptions\%s\PropertyBag`, guid)
	k, _, err := registry.CreateKey(registry.LOCAL_MACHINE, path, registry.ALL_ACCESS)
	if err == nil {
		_ = k.SetStringValue("ThisPCPolicy", "Hide")
		k.Close()
	}
}

func unpinFromCLSID(guid string) {
	path := fmt.Sprintf(`Software\Classes\CLSID\%s`, guid)
	k, _, err := registry.CreateKey(registry.CURRENT_USER, path, registry.SET_VALUE)
	if err == nil {
		_ = k.SetDWordValue("System.IsPinnedToNameSpaceTree", 0)
		k.Close()
	}
}

func removeDelegateFolder(guid string) {
	path := fmt.Sprintf(`SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\MyComputer\NameSpace\DelegateFolders\%s`, guid)
	_ = registry.DeleteKey(registry.LOCAL_MACHINE, path)
}
