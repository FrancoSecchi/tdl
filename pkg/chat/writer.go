package chat

import (
   "encoding/csv"
   "os"
   "io"
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
			registered: record[2] == "true", // Convertir "true" a true, "false" a false,
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
