package chat

import (
   "encoding/csv"
   "os"
)

// writeUsersToCSV writes user data to a CSV file.
func writeUsersToCSV(users []*User, filename string) error {
      file, err := os.Create(filename)
   if err != nil {
      return err
   }
   defer file.Close()

   writer := csv.NewWriter(file)
   defer writer.Flush()

   // Write header
   header := []string{"Name", "Password", "Registered"}
   if err := writer.Write(header); err != nil {
      return err
   }

   // Write user data
   for _, user := range users {
      if err := writer.Write(user.toCSVRecord()); err != nil {
         return err
      }
   }

   return nil
}
