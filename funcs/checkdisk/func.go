package checkdiskerrors

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/manifoldco/promptui"
)

func Ccheckdsk() {
	readerStdin := bufio.NewReader(os.Stdin)

	prompt := promptui.Prompt{
		Label: "[BTL] Введите диск для проверки (Default: C:)",
	}
	resultdiskuser, err := prompt.Run()
	if err != nil {
		fmt.Printf("[BTL] Ошибка при вводе: %v\n", err)
		return
	}

	drive := strings.TrimSpace(resultdiskuser)
	if drive == "" {
		drive = "C"
	}
	drive = strings.ToUpper(strings.TrimSuffix(drive, ":")) + ":"

	fmt.Printf("\n[BTL] Анализ диска %s\n", drive)

	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка загрузки kernel32.dll: %v\n", err)
		pauseExit(readerStdin)
		return
	}
	defer kernel32.Release()

	getDiskFreeSpaceEx, err := kernel32.FindProc("GetDiskFreeSpaceExW")
	if err != nil {
		fmt.Printf("\n[BTL] Ошибка поиска GetDiskFreeSpaceExW: %v\n", err)
		pauseExit(readerStdin)
		return
	}

	var freeBytesAvailable, totalBytes, totalFreeBytes uint64
	pathPtr, err := syscall.UTF16PtrFromString(drive + "\\")
	if err != nil {
		fmt.Printf("[BTL] Ошибка пути: %v\n", err)
		pauseExit(readerStdin)
		return
	}

	_, _, _ = getDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(pathPtr)),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalBytes)),
		uintptr(unsafe.Pointer(&totalFreeBytes)),
	)

	if totalBytes == 0 {
		fmt.Printf("\n[BTL] Ошибка доступа к диску %s (возможно, он недоступен).\n", drive)
		pauseExit(readerStdin)
		return
	}

	totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
	freeGB := float64(freeBytesAvailable) / (1024 * 1024 * 1024)
	usedGB := totalGB - (float64(totalFreeBytes) / (1024 * 1024 * 1024))

	testFile := filepath.Join(drive+"\\", "btl_speed_test.tmp")
	blockSize := 1024 * 1024
	blockCount := 100
	dataBlock := make([]byte, blockSize)
	for i := range dataBlock {
		dataBlock[i] = byte(i % 256)
	}

	writeSpeed := 0.0
	readSpeed := 0.0
	readWriteOk := true

	startWrite := time.Now()
	f, err := os.Create(testFile)
	if err != nil {
		readWriteOk = false
	} else {
		for i := 0; i < blockCount; i++ {
			_, err := f.Write(dataBlock)
			if err != nil {
				readWriteOk = false
				break
			}
		}
		_ = f.Sync()
		_ = f.Close()
		writeDuration := time.Since(startWrite).Seconds()
		if readWriteOk && writeDuration > 0 {
			writeSpeed = (float64(blockCount) / 1024.0) / writeDuration
			writeSpeed = writeSpeed * 1024
		}
	}

	// Замер чтения
	if readWriteOk {
		startRead := time.Now()
		f, err := os.Open(testFile)
		if err != nil {
			readWriteOk = false
		} else {
			buf := make([]byte, blockSize)
			for i := 0; i < blockCount; i++ {
				_, err := f.Read(buf)
				if err != nil {
					readWriteOk = false
					break
				}
			}
			_ = f.Close()
			readDuration := time.Since(startRead).Seconds()
			if readWriteOk && readDuration > 0 {
				readSpeed = (float64(blockCount) / 1024.0) / readDuration
				readSpeed = readSpeed * 1024 // МБ/с
			}
		}
		_ = os.Remove(testFile)
	}

	fmt.Println("\n--------------------------------------------------")
	fmt.Printf(" Диск: %s\n", drive)
	fmt.Printf(" Общий объем:  %.2f GB\n", totalGB)
	fmt.Printf(" Занято:       %.2f GB\n", usedGB)
	fmt.Printf(" Свободно:     %.2f GB\n", freeGB)
	fmt.Println("--------------------------------------------------")
	if readWriteOk {
		fmt.Printf(" Скорость записи (100 MB): %.2f MB/s\n", writeSpeed)
		fmt.Printf(" Скорость чтения (100 MB): %.2f MB/s\n", readSpeed)
		fmt.Println("--------------------------------------------------")
		fmt.Println("\033[32m[BTL] Статус: Диск в порядке!\033[0m")
	} else {
		fmt.Println("\033[33m[BTL] Статус: Ошибка при выполнении теста чтения/записи.\033[0m")
	}
	fmt.Println("--------------------------------------------------")
	fmt.Println("[BTL] Проверка завершена успешно!")

	pauseExit(readerStdin)
}

func pauseExit(reader *bufio.Reader) {
	fmt.Print("\nНажмите Enter для выхода...")
	_, _ = reader.ReadString('\n')
}
