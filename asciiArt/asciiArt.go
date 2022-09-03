package asciiArt

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func AsciiArt(banner, text string) (string, int) {
	check := checkBanner(banner)
	if check != http.StatusOK {
		return "", check
	}

	input := strings.Split(text, "\r\n")

	removeNewline(&input)

	var result string

	for _, s := range input {
		if s == "" {
			result = result + "\n"
			continue
		}

		tmp, err := PrintInput(s)
		if err != http.StatusOK {
			return "", err
		}

		result = result + tmp
	}

	return result, http.StatusOK
}

var Store [128][8]string // Переменная для хранения символов из файла

func checkBanner(filename string) int {
	if filename != "standard.txt" && filename != "thinkertoy.txt" && filename != "shadow.txt" {
		return http.StatusInternalServerError
	}

	if err := TxtFileCheck(filename); err != http.StatusOK {
		return err
	}

	f, err := os.Open("asciiArt/" + filename)
	if err != nil {
		return http.StatusNotFound
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	order := " !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~" // это порядок символов в файле, нужен чтобы правильно все считать в Store

	for _, r := range order {
		scanner.Scan() // Так как перед каждым символом в файле идет пустая строка, ее нужно пропускать, перед считыванием в Store
		for i := 0; i < 8; i++ {
			scanner.Scan()                    // здесь уже начинается считывание из файла в Store
			Store[int(r)][i] = scanner.Text() // Каждый символ помещается в ячейку под номером своего ascii значенения
		}
	}
	return http.StatusOK
}

// Выводит данную строку на консоль символами из файла
func PrintInput(s string) (string, int) {
	var result string

	for i := 0; i < 8; i++ {
		var tmp string

		for _, r := range s {
			if r < 0 || r > 127 || Store[int(r)][0] == "" {
				return "", http.StatusBadRequest
			}
			tmp = tmp + Store[int(r)][i]
		}

		result = result + tmp + "\n"
	}

	return result, http.StatusOK
}

// Проверяет файл на какие либо изменения, и выдает ошибку, если они были
func TxtFileCheck(fileName string) int {
	var hash string
	switch fileName {
	case "standard.txt":
		hash = "e194f1033442617ab8a78e1ca63a2061f5cc07a3f05ac226ed32eb9dfd22a6bf"
	case "shadow.txt":
		hash = "26b94d0b134b77e9fd23e0360bfd81740f80fb7f6541d1d8c5d85e73ee550f73"
	case "thinkertoy.txt":
		hash = "a57beec43fde6751ba1d30495b092658a064452f321e221d08c3ac34a9dc1294"
	}

	file, err := os.Open("asciiArt/" + fileName)
	if err != nil {
		return http.StatusNotFound
	}
	defer file.Close()
	buf := make([]byte, 30*1024)
	sha256 := sha256.New()
	for {
		n, err := file.Read(buf)
		if n > 0 {
			_, err := sha256.Write(buf[:n])
			if err != nil {
				log.Fatal(err)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Read %d bytes: %v", n, err)
			break
		}
	}
	sum := fmt.Sprintf("%x", sha256.Sum(nil))

	if sum == hash {
		return http.StatusOK
	} else {
		return http.StatusInternalServerError
	}
}

// Когда в вводной строке нет слов, а только '\n' либо вообще ничего, создается лишняя пустая строка, из-за которой на консоль выводится лишшняя новая линия. Эта функция убирает эту строку.
// "\n" -> "", ""
// "\nHello" -> "", "Hello"
func removeNewline(input *[]string) {
	nowords := true
	for _, s := range *input {
		if len(s) > 0 {
			nowords = false
		}
	}
	if nowords {
		*input = (*input)[1:]
	}
}
