package checkflashfake

import (
	"bufio"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/manifoldco/promptui"
)

type UsbDrive struct {
	Letter string `json:"Letter"`
	Model  string `json:"Model"`
	SizeGB string `json:"SizeGB"`
	FreeGB string `json:"FreeGB"`
}

func H2testwCheck() {
	readerStdin := bufio.NewReader(os.Stdin)

	fmt.Println("\n[BTL] Поиск подключенных USB-флешек...")
	drives := getRemovableUSBDisks()

	if len(drives) == 0 {
		fmt.Println("[BTL] Не найдено ни одной съемной USB-флешки!")
		pauseExit(readerStdin)
		return
	}

	var items []string
	driveMap := make(map[string]UsbDrive)
	for _, d := range drives {
		label := fmt.Sprintf("[%s:] %s (Объем: %s GB, Свободно: %s GB)", d.Letter, d.Model, d.SizeGB, d.FreeGB)
		items = append(items, label)
		driveMap[label] = d
	}
	items = append(items, "Отмена")

	prompt := promptui.Select{
		Label: "[BTL] Выберите флешку для проверки на муляжи",
		Items: items,
	}

	_, result, err := prompt.Run()
	if err != nil || result == "Отмена" {
		fmt.Println("[BTL] Проверка отменена.")
		pauseExit(readerStdin)
		return
	}

	selected := driveMap[result]
	drive := selected.Letter + ":"

	fmt.Printf("\n[BTL] Выбрана флешка: %s (%s)\n", drive, selected.Model)

	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		fmt.Printf("[BTL] Ошибка загрузки kernel32.dll: %v\n", err)
		pauseExit(readerStdin)
		return
	}
	defer kernel32.Release()

	getDiskFreeSpaceEx, err := kernel32.FindProc("GetDiskFreeSpaceExW")
	if err != nil {
		fmt.Printf("[BTL] Ошибка поиска GetDiskFreeSpaceExW: %v\n", err)
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

	var reserveBytes uint64 = 30 * 1024 * 1024 // 30 МБ резерв
	if freeBytesAvailable <= reserveBytes {
		fmt.Printf("[BTL] На флешке %s недостаточно свободного места.\n", drive)
		pauseExit(readerStdin)
		return
	}
	testBytesSize := freeBytesAvailable - reserveBytes
	testGB := float64(testBytesSize) / (1024 * 1024 * 1024)

	fmt.Printf("[BTL] Будет протестировано: %.2f GB свободного пространства\n", testGB)
	fmt.Println("[BTL] Внимание! Процесс записи и чтения может занять время.")

	confirmPrompt := promptui.Prompt{
		Label: "[BTL] Запустить тест? (y/n)",
	}
	confirm, err := confirmPrompt.Run()
	if err != nil || strings.ToLower(confirm) != "y" {
		fmt.Println("[BTL] Тест отменен.")
		pauseExit(readerStdin)
		return
	}

	testFile := filepath.Join(drive+"\\", "btl_h2test_filler.tmp")
	f, err := os.Create(testFile)
	if err != nil {
		fmt.Printf("[BTL] Ошибка создания файла: %v\n", err)
		pauseExit(readerStdin)
		return
	}

	blockSize := int64(2 * 1024 * 1024) // 2 МБ блок
	totalBlocks := testBytesSize / uint64(blockSize)

	fmt.Println("\n[BTL] Фаза 1/2: Запись данных...")
	startWrite := time.Now()

	writeOk := true
	var writtenBlocks uint64
	for i := uint64(0); i < totalBlocks; i++ {
		block := make([]byte, blockSize)
		binary.BigEndian.PutUint64(block[0:8], i)
		_, _ = rand.Read(block[8:])

		_, writeErr := f.Write(block)
		if writeErr != nil {
			fmt.Printf("\n[BTL] Ошибка записи на блоке %d (признак фейковой флешки!): %v\n", i, writeErr)
			writeOk = false
			break
		}
		writtenBlocks = i + 1

		if (i+1)%10 == 0 || (i+1) == totalBlocks {
			fmt.Printf("\r[BTL] Записано блоков: %d / %d (%.1f%%)", i+1, totalBlocks, float64(i+1)/float64(totalBlocks)*100)
		}
	}
	_ = f.Sync()
	_ = f.Close()
	fmt.Printf("\n[BTL] Запись завершена за %.2f сек.\n", time.Since(startWrite).Seconds())

	if writtenBlocks == 0 {
		_ = os.Remove(testFile)
		pauseExit(readerStdin)
		return
	}

	fmt.Println("\n[BTL] Фаза 2/2: Чтение и проверка целостности...")
	startRead := time.Now()

	fRead, err := os.Open(testFile)
	if err != nil {
		fmt.Printf("[BTL] Ошибка открытия файла: %v\n", err)
		_ = os.Remove(testFile)
		pauseExit(readerStdin)
		return
	}

	readOk := true
	corrupted := false
	for i := uint64(0); i < writtenBlocks; i++ {
		block := make([]byte, blockSize)
		_, readErr := fRead.Read(block)
		if readErr != nil {
			fmt.Printf("\n[BTL] Ошибка чтения на блоке %d: %v\n", i, readErr)
			readOk = false
			break
		}

		savedIndex := binary.BigEndian.Uint64(block[0:8])
		if savedIndex != i {
			fmt.Printf("\n[BTL] ОШИБКА! Муляж/Сбой: блок %d перезаписан (найден индекс %d).\n", i, savedIndex)
			corrupted = true
			break
		}

		if (i+1)%10 == 0 || (i+1) == writtenBlocks {
			fmt.Printf("\r[BTL] Проверено блоков: %d / %d (%.1f%%)", i+1, writtenBlocks, float64(i+1)/float64(writtenBlocks)*100)
		}
	}
	_ = fRead.Close()
	_ = os.Remove(testFile)
	fmt.Printf("\n[BTL] Проверка завершена за %.2f сек.\n", time.Since(startRead).Seconds())

	fmt.Println("\n--------------------------------------------------")
	if writeOk && readOk && !corrupted {
		fmt.Println("\033[32m[BTL] РЕЗУЛЬТАТ: Флешка исправна, реальный объем подтвержден!\033[0m")
	} else {
		fmt.Println("\033[31m[BTL] РЕЗУЛЬТАТ: ВНИМАНИЕ! Обнаружены ошибки памяти или признаки фейковой флешки.\033[0m")
	}
	fmt.Println("--------------------------------------------------")

	pauseExit(readerStdin)
}

func pauseExit(reader *bufio.Reader) {
	fmt.Print("\nНажмите Enter для выхода...")
	_, _ = reader.ReadString('\n')
}

func getRemovableUSBDisks() []UsbDrive {
	var drives []UsbDrive

	kernel32 := syscall.MustLoadDLL("kernel32.dll")
	getLogicalDrives := kernel32.MustFindProc("GetLogicalDrives")
	getDriveTypeW := kernel32.MustFindProc("GetDriveTypeW")

	ret, _, _ := getLogicalDrives.Call()
	bitmask := uint32(ret)

	for i := uint(0); i < 26; i++ {
		if (bitmask & (1 << i)) != 0 {
			letter := string(rune('A' + i))
			if letter == "C" {
				continue
			}

			pathPtr, _ := syscall.UTF16PtrFromString(letter + ":\\")
			driveType, _, _ := getDriveTypeW.Call(uintptr(unsafe.Pointer(pathPtr)))

			if driveType == 2 || driveType == 3 {
				var freeBytes, totalBytes, totalFreeBytes uint64
				getDiskFreeSpaceEx := kernel32.MustFindProc("GetDiskFreeSpaceExW")
				_, _, _ = getDiskFreeSpaceEx.Call(
					uintptr(unsafe.Pointer(pathPtr)),
					uintptr(unsafe.Pointer(&freeBytes)),
					uintptr(unsafe.Pointer(&totalBytes)),
					uintptr(unsafe.Pointer(&totalFreeBytes)),
				)

				if totalBytes > 0 {
					totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
					if totalGB < 2048 {
						freeGB := float64(freeBytes) / (1024 * 1024 * 1024)

						typeName := "USB/Съемный диск"
						if driveType == 3 {
							typeName = "Внешний диск"
						}

						drives = append(drives, UsbDrive{
							Letter: letter,
							Model:  typeName,
							SizeGB: fmt.Sprintf("%.2f", totalGB),
							FreeGB: fmt.Sprintf("%.2f", freeGB),
						})
					}
				}
			}
		}
	}
	return drives
}
