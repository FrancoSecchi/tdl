package chat

import (
   "encoding/csv"
   "os"
   "io"
   "fmt"
   "strings"

   "path/filepath"
)


// getUsersFromCSV reads user data from a CSV file and returns a slice of users.
func getUsersFromCSV(filename string) (map[string]*User, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	var allUsers = make(map[string]*User)
	firstLine := true
	for {
		record, err := reader.Read()
		
        if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

        if firstLine {
            firstLine = false
            continue
        }

		// Crear un nuevo usuario a partir del registro CSV
		user := &User{
			name:       record[0],
			password:   record[1],
			registered: record[2] == "true",
		}
		allUsers[user.name] = user
	}

	return allUsers, nil
}

// appendUsersToCSV writes user data to a CSV file.
func appendUsersToCSV(users []*User, filename string) error {
    file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // Opcional: Si el archivo está vacío, escribe el encabezado.
    fileinfo, err := file.Stat()
    if err != nil {
        return err
    }

    if fileinfo.Size() == 0 {
        header := []string{"Name", "Password", "Registered"}
        if err := writer.Write(header); err != nil {
            return err
        }
    }

    // Escribe los datos de usuario
    for _, user := range users {
        if err := writer.Write(user.toCSVRecord()); err != nil {
            return err
        }
    }

    return nil
}

// writeChatHistory appends a message to a CSV file representing chat history.
func writeChatHistory(filename string, message []string, isGlobalChat bool) (bool, error) {
	chatsFolderPath := "chats"
	chatsFilePath := filepath.Join(chatsFolderPath, filename)
	chatFile, err := os.OpenFile(chatsFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error al abrir el archivo global_chat.csv:", err)
		return false, err
	}
	defer chatFile.Close()

	chatWriter := csv.NewWriter(chatFile)
	defer chatWriter.Flush()

	err = chatWriter.Write(message)
	if err != nil {
		fmt.Println("Error al escribir en el archivo " + filename + ": ", err)
		return false, err
	}
    
    return true, nil
}

// GetChatHistoryData reads chat history from a CSV file and returns a slice of ChatMessage.
func GetChatHistoryData(filePath string) ([]ChatMessage, error) {
	file, err := os.Open(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []ChatMessage{}, nil
		}
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var chatMessages []ChatMessage
 	var firstLine = true
	for _, record := range records {
        if firstLine && strings.Contains(strings.ToLower(filePath), "global") {
            firstLine = false
            continue
        }
		if len(record) == 3 {
			chatMessages = append(chatMessages, ChatMessage{
				User:    record[0],
				Message: record[1],
				Time:    record[2],
			})
		}
	}

	return chatMessages, nil
}
