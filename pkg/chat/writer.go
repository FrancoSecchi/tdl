package chat

import (
   "encoding/csv"
   "os"
)

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
