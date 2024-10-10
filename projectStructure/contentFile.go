package projectstructure

import "os"

func CreateFileWithContent(filepath string, content string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return nil
	}
	return nil
}
