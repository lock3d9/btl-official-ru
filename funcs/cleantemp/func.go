package cleantemp

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/eiannone/keyboard"
)

func CleanTempFiles() {
	fmt.Println("\n[BTL] Запуск очистки временных файлов...")

	tempPaths := []string{
		os.TempDir(),
		filepath.Join(os.Getenv("SystemRoot"), "Temp"),
		filepath.Join(os.Getenv("SystemRoot"), "Prefetch"),
	}

	var deletedCount int
	var freedBytes int64

	for _, path := range tempPaths {
		if path == "" {
			continue
		}

		entries, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			fullPath := filepath.Join(path, entry.Name())
			info, err := entry.Info()
			var size int64
			if err == nil {
				size = info.Size()
			}

			err = os.RemoveAll(fullPath)
			if err == nil {
				deletedCount++
				freedBytes += size
			}
		}
	}

	cmd := exec.Command("powershell", "-NoProfile", "-Command", "Clear-RecycleBin -Force -ErrorAction SilentlyContinue")
	_ = cmd.Run()

	freedMB := float64(freedBytes) / (1024 * 1024)

	fmt.Println("\n==================================================")
	fmt.Println("\033[32m[BTL] ОЧИСТКА ЗАВЕРШЕНА\033[0m")
	fmt.Println("==================================================")
	fmt.Printf("Удалено объектов: %d\n", deletedCount)
	fmt.Printf("Освобождено памяти: %.2f MB\n", freedMB)
	fmt.Println("Корзина успешно очищена.")
	fmt.Println("==================================================")

	fmt.Print("\nНажмите Enter для возврата в меню...")

	if err := keyboard.Open(); err == nil {
		defer keyboard.Close()
		for {
			_, key, err := keyboard.GetKey()
			if err != nil || key == keyboard.KeyEnter {
				break
			}
		}
	}
}
